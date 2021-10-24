package store

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/CyDrive/config"
	"github.com/CyDrive/consts"
	"github.com/CyDrive/models"
	"github.com/CyDrive/utils"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

type AccountStore interface {
	AddAccount(account *models.Account) error
	GetAccountByEmail(email string) (*models.Account, error)

	AddUsage(email string, usage int64) error
	ExpandCap(email string, newCap int64) error
}

// Store users in memory
// load from a json file
type AccountStoreMem struct {
	idGen           *utils.IdGenerator
	accountEmailMap map[string]*models.Account
	rwMutex         *sync.RWMutex // guard for accountEmailMap

	updatedFlag int32
}

func NewMemStore() *AccountStoreMem {
	store := AccountStoreMem{
		idGen:           utils.NewIdGenerator(),
		accountEmailMap: make(map[string]*models.Account),
		rwMutex:         &sync.RWMutex{},
		updatedFlag:     0,
	}

	data, err := ioutil.ReadFile(consts.MemAccountStoreJsonPath)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
	}

	accountArray := models.AccountList{}
	utils.GetJsonDecoder().Unmarshal(data, &accountArray)
	for _, account := range accountArray.AccountList {
		// Get the storage usage
		account.Usage, _ = utils.DirSize(utils.GetAccountDataDir(account.Id))

		store.idGen.Ref(account.Id)
		store.accountEmailMap[account.Email] = account
	}

	go store.persistThread()

	return &store
}

// required fields: Email, Password
// optional fields: Name, Cap
func (store *AccountStoreMem) AddAccount(account *models.Account) error {
	store.rwMutex.RLock()
	_, ok := store.accountEmailMap[account.Email]
	if ok {
		store.rwMutex.RUnlock()
		return fmt.Errorf("email %v has been registered", account.Email)
	}
	store.rwMutex.RUnlock()

	store.rwMutex.Lock()
	defer store.rwMutex.Unlock()

	account.Id = store.idGen.NextAndRef()
	store.accountEmailMap[account.Email] = account
	store.updatedFlag++

	return nil
}

func (store *AccountStoreMem) GetAccountByEmail(email string) (*models.Account, error) {
	store.rwMutex.RLock()
	defer store.rwMutex.RUnlock()

	return store.accountEmailMap[email], nil
}

func (store *AccountStoreMem) AddUsage(email string, usage int64) error {
	store.rwMutex.RLock()
	defer store.rwMutex.RUnlock()

	account, err := store.GetAccountByEmail(email)
	if err != nil {
		return err
	}

	if usage != 0 {
		atomic.AddInt32(&store.updatedFlag, 1)
	}
	atomic.AddInt64(&account.Usage, usage)

	return nil
}

func (store *AccountStoreMem) ExpandCap(email string, newCap int64) error {
	store.rwMutex.RLock()
	defer store.rwMutex.RUnlock()

	account, err := store.GetAccountByEmail(email)
	if err != nil {
		return err
	}

	for {
		old := atomic.LoadInt64(&account.Cap)
		if atomic.CompareAndSwapInt64(&account.Cap, old, newCap) {
			if old != newCap {
				atomic.AddInt32(&store.updatedFlag, 1)
				break
			}
		}
	}

	return nil
}

func (store *AccountStoreMem) persistThread() {
	for {
		updatedNum := atomic.LoadInt32(&store.updatedFlag)
		if updatedNum > 0 {
			store.rwMutex.RLock()
			store.save()
			store.rwMutex.RUnlock()

			atomic.AddInt32(&store.updatedFlag, -updatedNum)
		}

		time.Sleep(time.Second)
	}
}

func (store *AccountStoreMem) save() {

	accountList := models.AccountList{
		AccountList: make([]*models.Account, 0, len(store.accountEmailMap)),
	}
	for _, account := range store.accountEmailMap {
		accountList.AccountList = append(accountList.AccountList, account)
	}

	accountListBytes, err := utils.GetJsonEncoder().Marshal(&accountList)
	if err != nil {
		log.Error(err)
		return
	}

	err = ioutil.WriteFile(consts.MemAccountStoreJsonPath, accountListBytes, 0666)
	if err != nil {
		log.Error(err)
		return
	}
}

// Store users in a relational db
type AccountStoreRdb struct {
	db                *gorm.DB
	accountUsageCache sync.Map
}

func NewRdbStore(config config.Config) *AccountStoreRdb {
	store := AccountStoreRdb{}
	store.db, _ = gorm.Open("mysql", config.PackDSN())

	go store.MonitorUsageCache(5)
	return &store
}

func (store *AccountStoreRdb) AddAccount(account *models.Account) error {
	accountOrm, err := account.ToORM(context.Background())
	if err != nil {
		return err
	}

	if store.db.Create(accountOrm).Error != nil {
		return fmt.Errorf("email %v has been registered", account.Email)
	}

	return nil
}

func (store *AccountStoreRdb) MonitorUsageCache(delay int64) error {
	for {
		store.accountUsageCache.Range(func(key, value interface{}) bool {
			email := key.(string)
			usage := value.(int64)

			store.accountUsageCache.Delete(email)
			store.db.Model(models.AccountORM{}).Where("email = ?", email).UpdateColumn("Usage", gorm.Expr("usage + ?", usage))

			return true
		})

		time.Sleep(time.Duration(delay) * time.Second)
	}
}

func (store *AccountStoreRdb) GetAccountByEmail(email string) *models.Account {
	var account models.AccountORM

	if store.db.First(&account, "email = ?", email).RecordNotFound() {
		return nil
	}

	realAccount, _ := account.ToPB(context.Background())

	value, ok := store.accountUsageCache.Load(email)
	usage := value.(int64)

	if ok {
		realAccount.Usage += usage
	}

	return &realAccount
}

func (store *AccountStoreRdb) AddUsage(email string, usage int64) error {
	value, ok := store.accountUsageCache.Load(email)
	oldUsage := value.(int64)

	if ok {
		store.accountUsageCache.Store(email, oldUsage+usage)
	} else {
		store.accountUsageCache.Store(email, usage)
	}

	return nil
}

func (store *AccountStoreRdb) ExpandCap(email string, newCap int64) error {
	err := store.db.Model(models.AccountORM{}).Where("email = ?", email).Update("Cap", newCap).Error

	if err != nil {
		return err
	}

	return nil
}
