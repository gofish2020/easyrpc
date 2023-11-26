package rpcserver

import (
	"reflect"
	"time"
)

type Server interface {
	Register(string, interface{})
	Run()
	Shutdown()
}

// 服务启动配置参数
type Option struct {
	Ip           string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var DefaultOption = Option{
	ReadTimeout:  5 * time.Second,
	WriteTimeout: 5 * time.Second,
}

func NewRPCServer(option Option) *RPCServer {
	return &RPCServer{
		listener: NewRPCListener(option),
		option:   option,
	}
}

type RPCServer struct {
	listener Listener
	option   Option
}

func (server *RPCServer) Register(obj interface{}) {
	objectName := reflect.Indirect(reflect.ValueOf(obj)).Type().Name()
	server.RegisterByName(objectName, obj)

}

func (server *RPCServer) RegisterByName(objectName string, obj interface{}) {
	server.listener.SetHandler(objectName, &RPCHandler{object: reflect.ValueOf(obj)})
}

func (server *RPCServer) Run() {
	go server.listener.Run()
}

func (server *RPCServer) Shutdown() {
	if server.listener != nil {
		server.listener.Shutdown()
	}

}
