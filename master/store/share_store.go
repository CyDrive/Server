package store

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/CyDrive/utils"
	log "github.com/sirupsen/logrus"
)

type ShareLink struct {
	Uri             string
	FilePath        string
	From            int32
	To              int32
	Password        string
	LeftAccessCount int32
	Expire          int32 // in minutes
	CreatedAt       time.Time
}

type ShareStore interface {
	CreateShareLink(link *ShareLink) error
	CheckPermission(uri string, accountId int32, password string) bool
}

type ShareStoreMem struct {
	linkGroupMap *sync.Map // map: uri -> *ShareLinkGroup
	linkCount    uint64
	linkLength   int
}

type ShareLinkGroup struct {
	Uri             string
	FilePath        string
	From            int32
	To              []int32
	Password        string
	LeftAccessCount int32 // enable when `LeftAccessCount` >= 0
	Expire          int32 // in minutes, enable when `Expire` > 0
	CreatedAt       time.Time
}

func NewShareLinkGroup(link *ShareLink) *ShareLinkGroup {
	linkGroup := &ShareLinkGroup{
		Uri:             link.Uri,
		FilePath:        link.FilePath,
		From:            link.From,
		To:              make([]int32, 0, 1),
		Password:        link.Uri,
		LeftAccessCount: link.LeftAccessCount,
		Expire:          link.Expire,
	}

	if link.To > 0 {
		linkGroup.To = append(linkGroup.To, link.To)
	}

	return linkGroup
}

func (linkGroup *ShareLinkGroup) IsExpired() bool {
	return linkGroup.Expire > 0 &&
		linkGroup.CreatedAt.Add(time.Duration(linkGroup.Expire)*time.Minute).Before(time.Now())
}

func NewShareStoreMem() *ShareStoreMem {
	store := &ShareStoreMem{
		linkGroupMap: &sync.Map{},
		linkCount:    0,
		linkLength:   4,
	}

	go store.gcMaintenance()

	return store
}

// required fields: FilePath, From
// optional fields: To, Password, LeftAccessCount, Expire
func (store *ShareStoreMem) CreateShareLink(link *ShareLink, accountIds ...int32) error {
	linkGroup := NewShareLinkGroup(link)

	for {
		link.Uri = utils.GenRandomString(store.linkLength)

		linkGroupI, ok := store.linkGroupMap.Load(link.Uri)
		if !ok {
			break
		}

		old := linkGroupI.(*ShareLinkGroup)
		if old.IsExpired() {
			break
		}
	}

	linkGroup.Uri = link.Uri
	linkGroup.To = append(linkGroup.To, accountIds...)
	linkGroup.CreatedAt = time.Now()

	store.linkGroupMap.Store(link.Uri, linkGroup)
	return nil
}

func (store *ShareStoreMem) CheckPermission(uri string, accountId int32, password string) error {
	linkGroup, ok := store.GetShareLinkGroup(uri)
	if !ok {
		return fmt.Errorf("no such share-link, may be expired")
	}

	if linkGroup.IsExpired() {
		go store.linkGroupMap.Delete(uri)
		return fmt.Errorf("this share-link has expired")
	}

	if password != linkGroup.Password {
		return fmt.Errorf("wrong password")
	}

	if len(linkGroup.To) > 0 {
		hasPermission := false
		for _, id := range linkGroup.To {
			if id == accountId {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			return fmt.Errorf("no access")
		}
	}

	for {
		leftAccessCount := atomic.LoadInt32(&linkGroup.LeftAccessCount)

		if leftAccessCount < 0 {
			break
		} else if leftAccessCount == 0 {
			go store.linkGroupMap.Delete(uri)
			return fmt.Errorf("reach access limit, this share-link is invalid")
		} else {
			if atomic.CompareAndSwapInt32(&linkGroup.LeftAccessCount, leftAccessCount, leftAccessCount-1) {
				break
			}
		}
	}

	return nil
}

func (store *ShareStoreMem) GetShareLinkGroup(uri string) (*ShareLinkGroup, bool) {
	linkGroupI, ok := store.linkGroupMap.Load(uri)
	if !ok {
		return nil, false
	}

	linkGroup := linkGroupI.(*ShareLinkGroup)
	return linkGroup, true
}

func (store *ShareStoreMem) gcMaintenance() {
	removeLinkGroups := make([]string, 0)
	for {
		store.linkGroupMap.Range(func(key, value interface{}) bool {
			uri := key.(string)
			linkGroup := value.(*ShareLinkGroup)
			if linkGroup.IsExpired() {
				removeLinkGroups = append(removeLinkGroups, uri)
			}

			return true
		})

		if len(removeLinkGroups) > 0 {
			log.Infof("remove share-links: %+v", removeLinkGroups)
			for _, uri := range removeLinkGroups {
				store.linkGroupMap.Delete(uri)
			}
		}

		time.Sleep(5 * time.Second)
	}
}
