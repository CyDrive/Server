package master

import (
	"encoding/gob"
	"net"
	"time"

	"github.com/CyDrive/config"
	"github.com/CyDrive/consts"
	"github.com/CyDrive/master/envs"
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

func GetEnv() envs.Env {
	return master.env
}

type NodeManagerServer struct {
	rpc.UnimplementedManageServer
}

type Master struct {
	nodeManagerServer *NodeManagerServer

	fileTransferor *network.FileTransferor

	env          envs.Env
	accountStore store.AccountStore
}

func NewMaster(config config.Config) *Master {
	var (
		env          envs.Env
		accountStore store.AccountStore
	)

	if config.EnvType == consts.EnvTypeLocal {
		env = envs.NewLocalEnv()
	}
	if config.AccountStoreType == consts.AccountStoreTypeMem {
		accountStore = store.NewMemStore()
	}

	if env == nil || accountStore == nil {
		panic("error when initialize")
	}

	master = &Master{
		nodeManagerServer: &NodeManagerServer{},
		fileTransferor:    network.NewFileTransferor(env),
		env:               env,
		accountStore:      accountStore,
	}

	return master
}

func (m *Master) Start() {
	// HTTP services
	log.Info("start http services...")
	router := gin.Default()
	memStore := memstore.NewStore([]byte("ProjectMili"))
	router.Use(sessions.SessionsMany([]string{"account"}, memStore))
	router.Use(LoginAuth(router))
	// router.Use(SetFileInfo())

	router.POST("/register", RegisterHandle)
	router.POST("/login", LoginHandle)

	router.GET("/list/*path", ListHandle)

	// router.GET("/file_info/*path", GetFileInfoHandle)
	// router.PUT("/file_info/*path", PutFileInfoHandle)

	router.GET("/file/*path", DownloadHandle)
	router.PUT("/file/*path", UploadHandle)
	router.DELETE("/file/*path", DeleteHandle)

	go router.Run(consts.HttpListenPortStr)

	// RPC services
	log.Info("start rpc services...")
	listen, err := net.Listen("tcp", consts.RpcListenPortStr)
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	rpc.RegisterManageServer(grpcServer, m.nodeManagerServer)
	go grpcServer.Serve(listen)

	// Start FileTransferManager
	log.Info("start file transfer manager...")
	m.fileTransferor.Listen()
}
