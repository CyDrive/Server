package managers

import (
	"sync"

	"github.com/CyDrive/consts"
	"github.com/CyDrive/master/store"
	"github.com/CyDrive/models"
	"github.com/CyDrive/utils"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (mgr *MessageManager) GetHub(accountId int32) *Hub {
	hubI, ok := mgr.hubMap.LoadOrStore(accountId, NewHub(accountId, mgr.messageStore))
	hub := hubI.(*Hub)

	// it's the first connection of this account
	// create a new hub to serve
	// need to start a goroutine to deliver messages
	if !ok {
		log.Infof("create new hub for accountId=%+v", accountId)
		go hub.deliverMessage()
	}

	return hub
}

func (mgr *MessageManager) GetMessageStore() store.MessageStore {
	return mgr.messageStore
}

type Hub struct {
	// all the 3 queues will be handled in only 1 goroutine
	// so we don't need a lock to protect connMap
	accountId       int32
	messageQueue    chan *models.Message
	registerQueue   chan *MessageConn
	unregisterQueue chan string

	connMap map[string]*MessageConn // map: deviceId -> MessageConn

	// store
	messageStore store.MessageStore
}

func NewHub(accountId int32, messageStore store.MessageStore) *Hub {
	hub := Hub{
		accountId:       accountId,
		messageQueue:    make(chan *models.Message, 10),
		registerQueue:   make(chan *MessageConn, 3),
		unregisterQueue: make(chan string, 1),
		connMap:         map[string]*MessageConn{},

		messageStore: messageStore,
	}

	return &hub
}

func (hub *Hub) Register(conn *MessageConn) {
	log.Infof("register new conn, deviceId=%+v", conn.DeviceId)
	hub.registerQueue <- conn
	conn.PushQueue <- &models.Message{
		Id:         0,
		Sender:     "",
		SenderName: "CyberBot",
		Receiver:   "",
		Type:       consts.MessageType_Text,
		Content:    "it's a test message",
		SendedAt:   timestamppb.Now(),
		Expire:     30,
	}
	go conn.SendMessageHandle()
	go conn.PushMessage()
}

func (hub *Hub) Unregister(deviceId string) {
	log.Infof("unregister conn, deviceId=%+v", deviceId)
	hub.unregisterQueue <- deviceId
}

func (hub *Hub) HandleMessage(message *models.Message) {
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
			log.Infof("storing new message=%+v...", message)

			// it's a broadcast message
			if message.Receiver == "" {
				for id, conn := range hub.connMap {
					if id != message.Sender {
						filledMessage := &models.Message{}
						*filledMessage = *message
						filledMessage.Receiver = id
						hub.messageStore.SaveMessage(hub.accountId, filledMessage)
						conn.PushQueue <- filledMessage
					}
				}
			} else {
				hub.messageStore.SaveMessage(hub.accountId, message)
				conn, ok := hub.connMap[message.Receiver]
				if ok {
					conn.PushQueue <- message
				}
			}
		}
	}
}

type MessageConn struct {
	Hub        *Hub
	DeviceId   string
	DeviceName string // todo: fill the field when establishing connection
	Conn       *websocket.Conn
	PushQueue  chan *models.Message
}

func NewMessageConn(hub *Hub, deviceId string, conn *websocket.Conn) *MessageConn {
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
			log.Errorf("failed to read message, err=%+v, will close this connection", err)
			return
		}

		log.Infof("client sends message with type=%+v", msgType)

		switch msgType {
		case websocket.TextMessage:
			var message models.Message
			err = utils.GetJsonDecoder().Unmarshal(messageBytes, &message)
			if err != nil {
				log.Errorf("failed to unmarshal the message, messageBytes=%+v", string(messageBytes))
				return
			}

			log.Infof("client sends message=%+v", message)

			conn.Hub.HandleMessage(&message)
		}
	}
}

func (conn *MessageConn) PushMessage() {
	defer func() {
		close(conn.PushQueue)
	}()

	for message := range conn.PushQueue {
		log.Infof("push message=%+v", message)

		messageBytes, err := utils.GetJsonEncoder().Marshal(message)
		if err != nil {
			log.Errorf("failed to marshal message, message=%+v, err=%+v", message, err)
			return
		}

		err = conn.Conn.WriteMessage(websocket.TextMessage,
			messageBytes)

		if err != nil {
			log.Errorf("failed to push message, message=%+v, err=%+v, will retry later", message, err)
			conn.PushQueue <- message
		}
	}
}
