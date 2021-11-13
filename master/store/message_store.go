package store

import (
	"time"

	"github.com/CyDrive/models"
)

type MessageStore interface {
	GetMessagesByTime(recverId string, count int32, time time.Time) []*models.Message
	SaveMessage(message *models.Message)
}

type MessageStoreMem struct {
	messageMap map[string][]*models.Message
}

func NewMessageStoreMem() *MessageStoreMem {
	return &MessageStoreMem{
		messageMap: map[string][]*models.Message{},
	}
}

func (store MessageStoreMem) GetMessagesByTime(recverId string, count int32, time time.Time) []*models.Message {
	messageList, ok := store.messageMap[recverId]
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

func (store MessageStoreMem) SaveMessage(message *models.Message) {
	receiverId := message.Receiver
	_, ok := store.messageMap[receiverId]
	if !ok {
		store.messageMap[receiverId] = []*models.Message{}
	}
	store.messageMap[receiverId] = append(store.messageMap[receiverId], message)
}
