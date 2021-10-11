package store

import (
	"time"

	"github.com/CyDrive/models"
)

type MessageStore interface {
	GetMessageByTime(userId int32, count int32, time time.Time) []*models.Message
	PutMessage(message *models.Message)
}

type MessageMemStore struct {
	messageMap map[int32][]*models.Message
}

func (store MessageMemStore) GetMessageByTime(userId int32, count int32, time time.Time) []*models.Message {
	messageList, ok := store.messageMap[userId]
	if !ok {
		return []*models.Message{}
	}
	left := 0
	right := len(messageList) - 1
	if messageList[left].SendedAt.AsTime().After(time) {
		return []*models.Message{}
	}
	for left < right {
		mid := (left + right) / 2
		if time.After(messageList[mid+1].SendedAt.AsTime()) {
			left = mid + 1
		} else {
			right = mid
		}
	}
	if left-int(count)+1 < 0 {
		return messageList[0:left]
	} else {
		return messageList[left-int(count)+1 : left]
	}
}

func (store MessageMemStore) PutMessage(message *models.Message) {
	receiverId := message.Receiver
	_, ok := store.messageMap[receiverId]
	if !ok {
		store.messageMap[receiverId] = []*models.Message{}
	}
	store.messageMap[receiverId] = append(store.messageMap[receiverId], message)
}

var messageStore MessageStore = MessageMemStore{}

func GetMsgMemStore() MessageStore { return messageStore }
