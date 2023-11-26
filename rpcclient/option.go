package rpcclient

import (
	"time"

	"github.com/gofish2020/easyrpc/rpcmsg"
)

type Option struct {
	Network        string
	Retries        int
	FailMode       FailMode
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	SerializeType  rpcmsg.SerializeType
	CompressType   rpcmsg.CompressType
	Version        byte
}

var DefaultOption = Option{
	Network:        "tcp",
	Retries:        3,
	FailMode:       Failover,
	ConnectTimeout: 5 * time.Second,
	ReadTimeout:    1 * time.Second,
	WriteTimeout:   1 * time.Second,
	SerializeType:  rpcmsg.Gob,
	CompressType:   rpcmsg.Zlib,
	Version:        rpcmsg.Version,
}
