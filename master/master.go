package master

import (
	"net"

	"../config"
	"../consts"
	"../env"
	rpc "../rpc"
	"./handlers"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type Master struct {
	config config.Config
	env    env.Env
	rpc.UnimplementedManageServer
}

func NewMaster(config config.Config, env env.Env) *Master {
	return nil
}

func (m *Master) Start() {
	// HTTP services
	router := gin.Default()
	router.Use(sessions.SessionsMany([]string{"user"}, memStore))
	router.Use(handlers.LoginAuth(router))
	// router.Use(SetFileInfo())

	router.POST("/login", handlers.LoginHandle)
	router.GET("/list/*path", handlers.ListHandle)

	router.GET("/file_info/*path", handlers.GetFileInfoHandle)
	router.PUT("/file_info/*path", handlers.PutFileInfoHandle)

	router.GET("/file/*path", handlers.DownloadHandle)
	router.PUT("/file/*path", handlers.UploadHandle)

	go router.Run(consts.HttpListenPortStr)

	// RPC services
	listen, err := net.Listen("tcp", consts.RpcListenPortStr)
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	rpc.RegisterManageServer(grpcServer, m)
	grpcServer.Serve(listen)
}
