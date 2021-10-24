package managers

import (
	"sync"

	"github.com/CyDrive/master/store"
	"github.com/CyDrive/models"
	"github.com/CyDrive/utils"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type MessageManager struct {
	hubMap       *sync.Map // map: accountId -> Hub
	messageStore store.MessageStore
}

func NewMessageManager(messageStore store.MessageStore) *MessageManager {
	return &MessageManager{
		hubMap:       &sync.Map{},
		messageStore: messageStore,
	}
}

func (mgr *MessageManager) GetHub(accountId int32) (*Hub, bool) {
	hubI, ok := mgr.hubMap.LoadOrStore(accountId, NewHub())
	hub := hubI.(*Hub)

	// it's the first connection of this account
	// create a new hub to serve
	// need to start a goroutine to deliver messages
	if !ok {
		go hub.deliverMessage()
	}

	return hub, ok
}

func (mgr *MessageManager) GetMessageStore() store.MessageStore {
	return mgr.messageStore
}

type Hub struct {
	// all the 3 queues will be handled in only 1 goroutine
	// so we don't need a lock to protect connMap
	messageQueue    chan *models.Message
	registerQueue   chan *MessageConn
	unregisterQueue chan int32

	connMap map[int32]*MessageConn
}

func NewHub() *Hub {
	hub := Hub{
		messageQueue:    make(chan *models.Message, 10),
		registerQueue:   make(chan *MessageConn, 3),
		unregisterQueue: make(chan int32, 1),
		connMap:         map[int32]*MessageConn{},
	}

	return &hub
}

func (hub *Hub) Register(conn *MessageConn) {
	hub.registerQueue <- conn
}

func (hub *Hub) Unregister(deviceId int32) {
	hub.unregisterQueue <- deviceId
}

func (hub *Hub) PushMessage(message *models.Message) {
	hub.messageQueue <- message
}

func (hub *Hub) deliverMessage() {
	for {
		select {
		case conn := <-hub.registerQueue:
			hub.connMap[conn.DeviceId] = conn

		case deviceId := <-hub.unregisterQueue:
			delete(hub.connMap, deviceId)

		case message := <-hub.messageQueue:
			// it's a broadcast message
			if message.Receiver <= 0 {
				for id, conn := range hub.connMap {
					if id != message.Sender {
						conn.PushQueue <- message
					}
				}
			} else {
				conn, ok := hub.connMap[message.Receiver]
				if ok {
					conn.PushQueue <- message
				}
			}
		}
	}
}

type MessageConn struct {
	Hub       *Hub
	DeviceId  int32
	Conn      *websocket.Conn
	PushQueue chan *models.Message
}

func NewMessageConn(hub *Hub, deviceId int32, conn *websocket.Conn) *MessageConn {
	return &MessageConn{
		Hub:       hub,
		DeviceId:  deviceId,
		Conn:      conn,
		PushQueue: make(chan *models.Message, 10),
	}
}

func (conn *MessageConn) SendMessageHandle() {
	defer func() {
		conn.Hub.Unregister(conn.DeviceId)
		conn.Conn.Close()
	}()

	for {
		msgType, messageBytes, err := conn.Conn.ReadMessage()
		if err != nil {
			log.Errorf("failed to read message, err=%+v", err)
			return
		}

		switch msgType {
		case websocket.BinaryMessage:
			var message models.Message
			err = utils.GetJsonDecoder().Unmarshal(messageBytes, &message)
			if err != nil {
				log.Errorf("failed to unmarshal the message, messageBytes=%+v", string(messageBytes))
				return
			}

			conn.Hub.PushMessage(&message)
		}
	}
}

func (conn *MessageConn) PushMessage() {
	defer func() {
		conn.Hub.Unregister(conn.DeviceId)
		close(conn.PushQueue)
		conn.Conn.Close()
	}()

	for message := range conn.PushQueue {
		messageBytes, err := utils.GetJsonEncoder().Marshal(message)
		if err != nil {
			log.Errorf("failed to marshal message, message=%+v, err=%+v", message, err)
			return
		}

		err = conn.Conn.WriteMessage(websocket.BinaryMessage,
			messageBytes)

		if err != nil {
			log.Errorf("failed to push message, message=%+v, err=%+v", message, err)
			return
		}
	}
}
