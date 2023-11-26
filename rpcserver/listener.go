package rpcserver

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
	"time"

	"github.com/gofish2020/easyrpc/rpcmsg"
)

type Listener interface {
	Run()
	Shutdown()
	SetHandler(string, Handler)
}

func NewRPCListener(option Option) *RPCListener {
	return &RPCListener{
		Ip:        option.Ip,
		Port:      option.Port,
		option:    option,
		Handlers:  make(map[string]Handler),
		shutdown:  0,
		running:   0,
		closechan: make(chan struct{}),
	}
}

type RPCListener struct {
	Ip       string
	Port     int
	option   Option
	Handlers map[string]Handler

	l net.Listener

	running  int32 // 运行中的连接
	shutdown int32 // 服务关闭标识

	closechan chan struct{} // 监听关闭
}

func (listen *RPCListener) isShutDonw() bool {
	return atomic.LoadInt32(&listen.shutdown) == 1
}
func (listen *RPCListener) Run() {
	addr := fmt.Sprintf("%s:%d", listen.Ip, listen.Port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	listen.l = l

	log.Printf("server listen on %s\n", addr)

	go listen.acceptConn()

}

// 监听处理
func (listen *RPCListener) acceptConn() {
	for {
		conn, err := listen.l.Accept()
		if err != nil {
			select {
			case <-listen.closechan:
				return
			default:
			}
			if e, ok := err.(net.Error); ok && e.Timeout() {
				time.Sleep(2 * time.Microsecond)
				continue
			}
			log.Printf("accept() err:%+v\n", err)
			return
		}

		go listen.handleConn(conn)
	}
}

// 客户端连接处理
func (listen *RPCListener) handleConn(conn net.Conn) {
	// 如果服务正在关闭中...新连接进来自动关闭
	if listen.isShutDonw() {
		conn.Close()
		return
	}

	log.Printf("new client connection come in %s\n", conn.RemoteAddr().String())
	// 避免 panic
	defer func() {
		if err := recover(); err != nil {
			log.Printf("addr %s panic err :%+v\n", conn.RemoteAddr().String(), err)
		}
		conn.Close()

	}()
	// 记录处理中的连接个数
	atomic.AddInt32(&listen.running, 1)
	defer func() {
		atomic.AddInt32(&listen.running, -1)
	}()
	// 服务度是否关闭
	for !listen.isShutDonw() {

		// 读超时时间
		// if listen.option.ReadTimeout != 0 {
		// 	conn.SetReadDeadline(time.Now().Add(listen.option.ReadTimeout))
		// }

		// 从连接冲接收一个完整的数据包
		msg, err := rpcmsg.RecvFrom(conn)
		if err != nil {
			log.Printf("receive msg error:%+v\n", err)
			return
		}

		startTime := time.Now()
		// 压缩器
		compressor := rpcmsg.Compressor[msg.Header.CompressType()]
		payload, err := compressor.UnCompress(msg.Payload)
		if err != nil {
			log.Printf("uncompress msg error; %+v\n", err)
			return
		}
		// 序列化器
		codeTool := rpcmsg.Codecs[msg.Header.SerializeType()]

		// 入参解码
		argsIn := make([]interface{}, 0)
		err = codeTool.Decode(payload, &argsIn)
		if err != nil {
			log.Printf("decode msg error; %+v\n", err)
			return
		}
		// 并行读 Handlers是安全的
		handler, ok := listen.Handlers[msg.ObjectName]
		if !ok {
			log.Printf("%s is't registered!\n", msg.ObjectName)
			return
		}
		// 执行对象的具体方法
		result, err := handler.Handle(msg.MethodName, argsIn)
		if err != nil {
			log.Printf("%s.%s func exec error(可忽略错误)\n", msg.ObjectName, msg.MethodName)
		}
		// 编码结果
		encodeRes, err := codeTool.Encode(result)
		if err != nil {
			log.Printf("encode msg error:%+v\n", err)
			return
		}
		// 压缩结果
		compressRes, err := compressor.Compress(encodeRes)
		if err != nil {
			log.Printf("compress msg error:%+v\n", err)
			return
		}

		// 写超时时间

		if listen.option.WriteTimeout != 0 {
			conn.SetWriteDeadline(time.Now().Add(listen.option.WriteTimeout))
		}
		config := rpcmsg.RPCMsgConfig{
			MsgTypeConf:       rpcmsg.Response,
			CompressTypeConf:  msg.CompressType(),
			SerializeTypeConf: msg.SerializeType(),
			VersionConf:       msg.Version(),
			Seq:               msg.Seq,
			ObjectName:        "",
			MethodName:        "",
		}
		// 将结果返回给客户端
		err = rpcmsg.SendTo(conn, compressRes, config)
		if err != nil {
			log.Printf("send msg error:%+v\n", err)
			return
		}

		log.Printf("%s.%s total runtime %d ms\n", msg.ObjectName, msg.MethodName, time.Since(startTime).Milliseconds())
	}
}
func (listen *RPCListener) Shutdown() {
	// 设置关闭标识
	atomic.CompareAndSwapInt32(&listen.shutdown, 0, 1)
	// 关闭监听
	listen.closeChan()
	if listen.l != nil {
		listen.l.Close()
	}
	// 说明还有连接在处理（自旋锁）
	for atomic.LoadInt32(&listen.running) != 0 {
		//log.Printf("还有 %d task running\n", listen.running)
	}
	log.Printf("server shutdown success!!!\n")

}

func (listen *RPCListener) closeChan() {
	select {
	case <-listen.closechan:
	default:
		close(listen.closechan)
	}
}
func (listen *RPCListener) SetHandler(objectName string, handler Handler) {
	if _, ok := listen.Handlers[objectName]; ok {
		log.Printf("objectName %s is registered!\n", objectName)
		return
	}
	listen.Handlers[objectName] = handler
}
