package master

import (
	"net"

	"github.com/CyDrive/config"
	"github.com/CyDrive/consts"
	"github.com/CyDrive/master/env"
	"github.com/CyDrive/master/store"
	"github.com/CyDrive/rpc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

var (
	master *Master
)

func GetMaster() *Master {
	return master
}

func GetFileTransferManager() *FileTransferManager {
	return master.fileTransferManager
}

func GetAccountStore() store.AccountStore {
	return master.accountStore
}

func GetEnv() env.Env {
	return master.env
}

type NodeManagerServer struct {
	rpc.UnimplementedManageServer
}

type Master struct {
	nodeManagerServer *NodeManagerServer

	fileTransferManager *FileTransferManager

	env          env.Env
	accountStore store.AccountStore
}

func NewMaster(config config.Config, env env.Env, accountStore store.AccountStore) *Master {
	return nil
}

func (m *Master) Start() {
	// HTTP services
	router := gin.Default()
	memStore := memstore.NewStore([]byte("ProjectMili"))
	router.Use(sessions.SessionsMany([]string{"user"}, memStore))
	router.Use(LoginAuth(router))
	// router.Use(SetFileInfo())

	router.POST("/login", LoginHandle)
	router.GET("/list/*path", ListHandle)

	router.GET("/file_info/*path", GetFileInfoHandle)
	// router.PUT("/file_info/*path", PutFileInfoHandle)

	router.GET("/file/*path", DownloadHandle)
	router.PUT("/file/*path", UploadHandle)

	go router.Run(consts.HttpListenPortStr)

	// RPC services
	listen, err := net.Listen("tcp", consts.RpcListenPortStr)
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	rpc.RegisterManageServer(grpcServer, m.nodeManagerServer)
	go grpcServer.Serve(listen)

	// Start FileTransferManager
	go m.fileTransferManager.Listen()
}
