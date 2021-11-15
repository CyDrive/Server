package store

import (
	"time"

	"github.com/CyDrive/models"
)

type MessageStore interface {
	GetMessagesByTime(accountId int32, recverId string, count int32, time time.Time) []*models.Message
	SaveMessage(accountId int32, message *models.Message)
}

type DeviceMessageMap = map[string][]*models.Message // deviceId -> []*Message
type AccountMessageMap = map[int32]DeviceMessageMap  // accountId -> DeviceMessageMap
type MessageStoreMem struct {
	messageMap AccountMessageMap // map: device_id -> []*Message, inclue the messages the device sended and received
}

func NewMessageStoreMem() *MessageStoreMem {
	return &MessageStoreMem{
		messageMap: AccountMessageMap{},
	}
}

func (store MessageStoreMem) GetMessagesByTime(accountId int32, deviceId string, count int32, time time.Time) []*models.Message {
	deviceMsgMap, ok := store.messageMap[accountId]
	if !ok {
		return []*models.Message{}
	}
	messages, ok := deviceMsgMap[deviceId]
	if !ok {
		return []*models.Message{}
	}

	left := 0
	right := len(messages) - 1
	if messages[left].SendedAt.AsTime().After(time) {
		return []*models.Message{}
	}

	for left < right {
		mid := (left + right) / 2
		if time.After(messages[mid+1].SendedAt.AsTime()) {
			left = mid + 1
		} else {
			right = mid
		}
	}

	if left-int(count)+1 < 0 {
		return messages[0:left]
	} else {
		return messages[left-int(count)+1 : left]
	}
}

func (store MessageStoreMem) SaveMessage(accountId int32, message *models.Message) {
	senderId := message.Sender
	receiverId := message.Receiver

	// create slice if miss
	deviceMsgMap, ok := store.messageMap[accountId]
	if !ok {
		deviceMsgMap = make(DeviceMessageMap)
		store.messageMap[accountId] = deviceMsgMap
	}

	_, ok = deviceMsgMap[senderId]
	if !ok {
		deviceMsgMap[senderId] = []*models.Message{}
	}
	_, ok = deviceMsgMap[receiverId]
	if !ok {
		deviceMsgMap[receiverId] = []*models.Message{}
	}

	deviceMsgMap[senderId] = append(deviceMsgMap[senderId], message)
	deviceMsgMap[receiverId] = append(deviceMsgMap[receiverId], message)
}
