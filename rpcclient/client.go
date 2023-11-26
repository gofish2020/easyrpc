package rpcclient

import (
	"context"
	"fmt"
	"log"
	"net"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gofish2020/easyrpc/rpcmsg"
)

type Client interface {
	Connect(addr string) error
	Call()
	Close()
}

func NewRPCClient(option Option) *RPCClient {
	return &RPCClient{
		option:         option,
		waiting:        make(map[int64]*waitMsg),
		serverShutdown: 0,
		clientClose:    0,
	}
}

type RPCClient struct {
	conn   net.Conn
	option Option
	addr   string

	mutex sync.Mutex

	mu      sync.RWMutex
	waiting map[int64]*waitMsg

	serverShutdown int32
	clientClose    int32
}

func (client *RPCClient) addWaitMsg(wMsg *waitMsg) {
	client.mu.Lock()
	defer client.mu.Unlock()
	client.waiting[wMsg.GetSeq()] = wMsg
}

func (client *RPCClient) removeWaitMsg(seq int64) *waitMsg {
	client.mu.Lock()
	defer client.mu.Unlock()
	wMsg := client.waiting[seq]
	delete(client.waiting, seq)
	return wMsg
}

func (client *RPCClient) removeAllWaitMsg() {
	atomic.CompareAndSwapInt32(&client.serverShutdown, 0, 1)
	client.mu.Lock()
	defer client.mu.Unlock()
	for _, waitMsg := range client.waiting {
		waitMsg.Ready(nil)
	}
	client.waiting = make(map[int64]*waitMsg)
}

func (client *RPCClient) loopWaitMsg() {
	for {
		resMsg, err := rpcmsg.RecvFrom(client.conn)
		if err != nil {
			break
		}
		wMsg := client.removeWaitMsg(resMsg.Seq)
		if wMsg != nil { // 说明这个序列号，不存在
			wMsg.Ready(resMsg)
		}
	}
	client.removeAllWaitMsg()
}
func (client *RPCClient) Connect(addr string) error {
	conn, err := net.DialTimeout(client.option.Network, addr, client.option.ConnectTimeout)
	if err != nil {
		atomic.CompareAndSwapInt32(&client.clientClose, 0, 1)
		atomic.CompareAndSwapInt32(&client.serverShutdown, 0, 1)
		return err
	}
	client.conn = conn
	client.addr = addr
	atomic.CompareAndSwapInt32(&client.clientClose, 1, 0)
	atomic.CompareAndSwapInt32(&client.serverShutdown, 1, 0)
	go client.loopWaitMsg()
	return nil
}

func (client *RPCClient) Close() {
	atomic.CompareAndSwapInt32(&client.clientClose, 0, 1)
	if client.conn != nil {
		client.conn.Close()
	}
}

// servicePath
func (client *RPCClient) Call(ctx context.Context, servicePath string, stub interface{}, params ...interface{}) (interface{}, error) {
	serviceInfo := strings.Split(servicePath, ".")
	if len(serviceInfo) != 2 {
		return nil, fmt.Errorf("servicePath format is splitted by point ObjectXXX.MethodXXX")
	}

	// stub 函数指针
	funcValue := reflect.ValueOf(stub).Elem()

	fn := func(args []reflect.Value) (results []reflect.Value) {

		numOut := funcValue.Type().NumOut()

		errorHandler := func(err error) []reflect.Value {
			argsOut := make([]reflect.Value, numOut)
			// 前面的[0,len(argsout)-1)设置为0值
			for i := 0; i < len(argsOut)-1; i++ {
				argsOut[i] = reflect.Zero(funcValue.Type().Out(i))
			}
			argsOut[len(argsOut)-1] = reflect.ValueOf(err)
			return argsOut
		}

		if atomic.LoadInt32(&client.clientClose) == 1 {
			return errorHandler(ErrClient)
		}
		if atomic.LoadInt32(&client.serverShutdown) == 1 {
			return errorHandler(ErrServer)
		}

		// 入参
		argsIn := make([]interface{}, 0, len(args))
		for _, arg := range args {
			argsIn = append(argsIn, arg.Interface())
		}
		// 序列化器
		codeTool := rpcmsg.Codecs[client.option.SerializeType]
		encodeRes, err := codeTool.Encode(argsIn)
		if err != nil {
			log.Printf("encode err:%+v\n", err)
			return errorHandler(err)
		}
		// 压缩器
		compressor := rpcmsg.Compressor[client.option.CompressType]
		payload, err := compressor.Compress(encodeRes)
		if err != nil {
			log.Printf("compress err:%+v\n", err)
			return errorHandler(err)
		}

		waitMsg := newWaitMsg()
		// 发送请求
		conf := rpcmsg.RPCMsgConfig{
			MsgTypeConf:       rpcmsg.Request,
			CompressTypeConf:  client.option.CompressType,
			SerializeTypeConf: client.option.SerializeType,
			VersionConf:       client.option.Version,
			ObjectName:        serviceInfo[0],
			MethodName:        serviceInfo[1],
			Seq:               waitMsg.GetSeq(),
		}
		// 设置读取超时时间
		if client.option.WriteTimeout != 0 {
			client.conn.SetWriteDeadline(time.Now().Add(client.option.WriteTimeout))
		}
		client.addWaitMsg(waitMsg)
		client.mutex.Lock()
		err = rpcmsg.SendTo(client.conn, payload, conf)
		client.mutex.Unlock()
		if err != nil {
			client.removeWaitMsg(waitMsg.GetSeq())
			return errorHandler(err)
		}

		//log.Println("send to server success")

		// 获取请求

		resMsg := waitMsg.Wait()
		if resMsg == nil {
			return errorHandler(ErrServer)
		}
		// 解压缩
		compressRes, err := compressor.UnCompress(resMsg.Payload)
		if err != nil {
			return errorHandler(err)
		}

		argsOut := make([]interface{}, 0)
		// 反序列化
		err = codeTool.Decode(compressRes, &argsOut)
		if err != nil {
			return errorHandler(err)
		}
		//log.Println("decode success")

		// 出参
		if len(argsOut) != numOut {
			argsOut = make([]interface{}, numOut)
			log.Println("出参个数和实际需要不一致")
		}
		// 利用 argsOut填充result
		results = make([]reflect.Value, numOut)

		for i := 0; i < len(results); i++ {
			if argsOut[i] == nil {
				results[i] = reflect.Zero(funcValue.Type().Out(i))
			} else {
				results[i] = reflect.ValueOf(argsOut[i])
			}
		}

		return results
	}
	// 相当于修改了stub指向的函数为fn
	funcValue.Set(reflect.MakeFunc(funcValue.Type(), fn))

	// 执行 funcValue函数

	if len(params) != funcValue.Type().NumIn() {
		return nil, ErrParam
	}

	in := make([]reflect.Value, len(params))
	for idx, param := range params {
		in[idx] = reflect.ValueOf(param)
	}

	result := funcValue.Call(in)
	return result, nil
}
