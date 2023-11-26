package main

import (
	"context"
	"encoding/gob"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofish2020/easyrpc/example/server/user"
	"github.com/gofish2020/easyrpc/rpcserver"
)

func main() {
	gob.RegisterName("info", user.Info{})
	option := rpcserver.Option{
		Ip:           "127.0.0.1",
		Port:         6060,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	server := rpcserver.NewRPCServer(option)
	server.RegisterByName("User", &user.UserService{})

	go server.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Printf("start shutdown server")
	// 优雅关闭服务
	server.Shutdown()

	select {
	case <-ctx.Done():
		log.Println("server close timeout")
	default:
	}

	log.Println("server exiting")
}
