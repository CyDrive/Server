package store

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/CyDrive/config"
	"github.com/CyDrive/consts"
	"github.com/CyDrive/model"
	"github.com/CyDrive/utils"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

type AccountStore interface {
	AddAccount(account *model.Account) error
	GetAccountByEmail(email string) (*model.Account, error)

	AddUsage(email string, usage int64) error
	ExpandCap(email string, newCap int64) error
}

// Store users in memory
// load from a json file
type MemStore struct {
	idGen           *utils.IdGenerator
	accountEmailMap map[string]*model.Account
	rwMutex         *sync.RWMutex // guard for accountEmailMap

	updatedFlag int32
}

func NewMemStore() *MemStore {
	store := MemStore{
		idGen:           utils.NewIdGenerator(),
		accountEmailMap: make(map[string]*model.Account),
		rwMutex:         &sync.RWMutex{},
		updatedFlag:     0,
	}

	data, err := ioutil.ReadFile(consts.MemAccountStoreJsonPath)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
	}

	accountArray := make([]*model.Account, 1)
	json.Unmarshal(data, &accountArray)
	for _, account := range accountArray {
		// Get the storage usage
		account.Usage, _ = utils.DirSize(account.DataDir)

		store.accountEmailMap[account.Name] = account
	}

	go store.persistThread()

	return &store
}

// required fields: Email, Password
// optional fields: Name, Cap
func (store *MemStore) AddAccount(account *model.Account) error {
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

func (store *MemStore) GetAccountByEmail(email string) (*model.Account, error) {
	store.rwMutex.RLock()
	defer store.rwMutex.RUnlock()

	return store.accountEmailMap[email], nil
}

func (store *MemStore) AddUsage(email string, usage int64) error {
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

func (store *MemStore) ExpandCap(email string, newCap int64) error {
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

func (store *MemStore) persistThread() {
	for {
		updatedNum := atomic.LoadInt32(&store.updatedFlag)
		if updatedNum > 0 {
			store.rwMutex.RLock()
			store.save()
			store.rwMutex.RUnlock()

			atomic.AddInt32(&store.updatedFlag, -updatedNum)
		}
	}
}

func (store *MemStore) save() {
	accountList := make([]*model.Account, 0, len(store.accountEmailMap))
	for _, account := range store.accountEmailMap {
		accountList = append(accountList, account)
	}

	accountListBytes, err := json.Marshal(accountList)
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
type RdbStore struct {
	db *gorm.DB
}

func NewRdbStore(config config.Config) *RdbStore {
	store := RdbStore{}
	store.db, _ = gorm.Open("mysql", config.PackDSN())
	return &store
}

func (store *RdbStore) GetAccountByEmail(email string) *model.Account {
	var account model.AccountORM

	if store.db.First(&account, "email = ?", email).RecordNotFound() {
		return nil
	}

	account.DataDir = filepath.Join(consts.UserDataDir, fmt.Sprint(account.Id))
	realAccount, _ := account.ToPB(context.Background())
	return &realAccount
}
