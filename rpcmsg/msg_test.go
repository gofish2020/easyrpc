package rpcmsg

import (
	"bytes"
	"testing"

	"github.com/gofish2020/easyrpc/codec"
	"github.com/stretchr/testify/assert"
)

func TestMsg(t *testing.T) {
	msg := NewRPCMsg()
	msg.SetMsgType(Request)
	msg.SetVersion(0x11)
	msg.SetCompressType(Zlib)
	msg.SetSerializeType(Json)
	msg.ObjectName = "UserService"
	msg.MethodName = "GetUserIds"
	json := codec.JsonCodec{}
	payload, _ := json.Encode(map[string]interface{}{
		"nash": 1,
		"yu":   "fdsf",
	})

	compressor := Compressor[msg.CompressType()]
	compressPayload, _ := compressor.Compress(payload)
	msg.Payload = compressPayload

	var buf bytes.Buffer
	msg.SendMsg(&buf)

	msg2 := NewRPCMsg()
	err := msg2.RecvMsg(&buf)
	t.Log("msg2.RecvMsg err", err)
	t.Log("MsgType", msg2.Header.MsgType())
	t.Log("Version", msg2.Header.Version())
	t.Log("CompressType", msg2.Header.CompressType())
	t.Log("SerializeType", msg2.Header.SerializeType())
	t.Log("MagicNumber", msg2.MagicNumber())
	t.Log("MethodName", msg2.MethodName)
	t.Log("ServiceName", msg2.ObjectName)

	unCompressPayload, _ := compressor.UnCompress(msg2.Payload)

	assert.Equal(t, unCompressPayload, payload)

	m := make(map[string]interface{})
	err = json.Decode(unCompressPayload, &m)
	t.Log(m, err)
}
