package master

import (
	"encoding/gob"
	"net"
	"time"

	"github.com/CyDrive/config"
	"github.com/CyDrive/consts"
	"github.com/CyDrive/master/envs"
	"github.com/CyDrive/master/managers"
	"github.com/CyDrive/master/store"
	"github.com/CyDrive/models"
	"github.com/CyDrive/network"
	"github.com/CyDrive/rpc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func init() {
	gob.Register(&models.Account{})
	gob.Register(time.Time{})
}

var (
	master *Master
)

func GetMaster() *Master {
	return master
}

func GetFileTransferor() *network.FileTransferor {
	return master.fileTransferor
}

func GetAccountStore() store.AccountStore {
	return master.accountStore
}

func GetMessageStore() store.MessageStore {
	return master.messageStore
}

func GetShareStore() store.ShareStore {
	return master.shareStore
}

func GetEnv() envs.Env {
	return master.env
}
func GetNodeManager() *managers.NodeManager {
	return master.nodeManager
}

func GetMessageManager() *managers.MessageManager {
	return master.messageManager
}

type Master struct {
	// node
	nodeManager      *managers.NodeManager
	nodeManageServer *NodeManageServer

	// file
	fileTransferor *network.FileTransferor

	// message
	messageManager *managers.MessageManager

	env          envs.Env
	accountStore store.AccountStore
	messageStore store.MessageStore
	shareStore   store.ShareStore
}

func NewMaster(config config.Config) *Master {
	var (
		env          envs.Env
		accountStore store.AccountStore
		messageStore store.MessageStore
		shareStore   store.ShareStore
	)

	if config.EnvType == consts.EnvTypeLocal {
		env = envs.NewLocalEnv()
	}
	if config.AccountStoreType == consts.AccountStoreTypeMem {
		accountStore = store.NewMemStore()
	}
	switch config.MessageStoreType {
	case consts.MessageStoreTypeMem:
		messageStore = store.NewMessageStoreMem()
	}

	shareStore = store.NewShareStoreMem()

	if env == nil || accountStore == nil {
		panic("error when initialize")
	}

	master = &Master{
		nodeManager:      managers.NewNodeManager(),
		nodeManageServer: &NodeManageServer{},

		fileTransferor: network.NewFileTransferor(env),

		messageManager: managers.NewMessageManager(messageStore),

		env:          env,
		accountStore: accountStore,
		messageStore: messageStore,
		shareStore:   shareStore,
	}

	return master
}

func (m *Master) Start() {
	// HTTP services
	log.Info("start http services...")
	router := gin.Default()
	memStore := memstore.NewStore([]byte("ProjectMili"))
	router.Use(sessions.SessionsMany([]string{"account"}, memStore))
	router.Use(SetRequestId(router))
	router.Use(Log(router))
	router.Use(LoginAuth(router))
	// router.Use(SetFileInfo())

	router.POST("/register", RegisterHandle)
	router.POST("/login", LoginHandle)

	router.GET("/list/*path", ListHandle)

	// storage service
	router.GET("/file/*path", DownloadHandle)
	router.PUT("/file/*path", UploadHandle)
	router.DELETE("/file/*path", DeleteHandle)

	// message service
	router.GET("/message_service", ConnectMessageServiceHandle)
	router.GET("/message", GetMessageHandle)

	// share
	router.POST("/share", ShareHandle)
	router.GET("/share/*uri", GetSharedFileHandle)

	go router.Run(consts.HttpListenPortStr)

	// RPC services
	log.Info("start rpc services...")
	listen, err := net.Listen("tcp", consts.RpcListenPortStr)
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	rpc.RegisterManageServer(grpcServer, m.nodeManageServer)
	go grpcServer.Serve(listen)

	// Start FileTransferManager
	log.Info("start file transfer managers...")
	m.fileTransferor.Listen()
}
