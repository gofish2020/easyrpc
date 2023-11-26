package rpcclient

import (
	"github.com/gofish2020/easyrpc/rpcmsg"
	"github.com/gofish2020/easyrpc/utils"
)

type waitMsg struct {
	done chan struct{}
	msg  *rpcmsg.RPCMsg
	seq  int64
}

func newWaitMsg() *waitMsg {
	return &waitMsg{
		done: make(chan struct{}, 1),
		msg:  nil,
		seq:  utils.CreateGUID(),
	}
}

func (t *waitMsg) GetSeq() int64 {
	return t.seq
}

func (t *waitMsg) Wait() *rpcmsg.RPCMsg {
	<-t.done
	return t.msg
}

func (t *waitMsg) Ready(msg *rpcmsg.RPCMsg) {
	t.msg = msg
	t.done <- struct{}{}
}
