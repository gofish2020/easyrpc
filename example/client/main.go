package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/gofish2020/easyrpc/example/client/user"
	"github.com/gofish2020/easyrpc/rpcclient"
)

func main() {
	gob.RegisterName("info", user.Info{})
	option := rpcclient.DefaultOption
	client := rpcclient.NewRPCClient(option)

	err := client.Connect("127.0.0.1:6060")
	if err != nil {
		log.Println(err)
		return
	}

	var sayHello func(s string) (string, error)
	result, err := client.Call(context.Background(), "User.SayHello", &sayHello, "hello client")
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(result)
	// 支持并行调用
	for i := 0; i < 10; i++ {
		//time.Sleep(1 * time.Second)
		go func(idx int) {
			str, err := sayHello("hello client!!!!!")
			fmt.Println(idx, str, err)

		}(i)
	}
	//time.Sleep(10 * time.Second)

	var GetUserInfoById func(id int) (user.Info, error)
	result, err = client.Call(context.Background(), "User.GetUserInfoById", &GetUserInfoById, 1)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(result)
	info, err := GetUserInfoById(1)
	fmt.Println(info, err)

}
