package rpcmsg

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"

	"github.com/gofish2020/easyrpc/utils"
)

const (
	DATA_LEN uint32 = 4
)

// RPCMsg: 一个完整的数据包 header + body
type RPCMsg struct {
	Header
	Seq int64 //请求编号
	// uint32 表示长度
	ObjectName string
	// uint32 表示长度
	MethodName string
	// uint32 表示长度
	Payload []byte
}

func NewRPCMsg() *RPCMsg {
	rpcMsg := RPCMsg{
		Header: NewHeader(),
	}
	rpcMsg.Header[0] = magicNumber
	return &rpcMsg
}

// SendMsg 发送消息
func (t *RPCMsg) SendMsg(w io.Writer) error {
	var err error
	//******
	_, err = w.Write(t.Header[:]) // 1.发送header头 5字节
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.BigEndian, uint64(t.Seq)) // 8字节
	if err != nil {
		return err
	}
	//******
	totalLen := DATA_LEN + uint32(len(t.ObjectName)) + DATA_LEN + uint32(len(t.MethodName)) + DATA_LEN + uint32(len(t.Payload))
	err = binary.Write(w, binary.BigEndian, uint32(totalLen)) // 2.写入总长度 4字节
	if err != nil {
		return err
	}
	//******
	err = binary.Write(w, binary.BigEndian, uint32(len(t.ObjectName))) // 3.写入 ObjectName 长度
	if err != nil {
		return err
	}
	_, err = w.Write(utils.String2Bytes(t.ObjectName)) // 4.写入 ObjectName
	if err != nil {
		return err
	}

	//******
	err = binary.Write(w, binary.BigEndian, uint32(len(t.MethodName))) // 5.写入 MethodName 长度
	if err != nil {
		return err
	}
	_, err = w.Write(utils.String2Bytes(t.MethodName)) // 6.写入 MethodName
	if err != nil {
		return err
	}
	//******
	err = binary.Write(w, binary.BigEndian, uint32(len(t.Payload))) // 7.写入 Payload 长度
	if err != nil {
		return err
	}
	_, err = w.Write(t.Payload) // 8.写入 Payload

	return err
}

// RecvMsg 接收消息
func (t *RPCMsg) RecvMsg(r io.Reader) error {

	var err error
	//1. 读取header数据
	_, err = io.ReadFull(r, t.Header[:])
	if err != nil {
		return err
	}
	if !t.Header.CheckMagicNumber() {
		return fmt.Errorf("magic number error: %v", t.Header[0])
	}
	seqByte := make([]byte, 8)
	_, err = io.ReadFull(r, seqByte)
	if err != nil {
		return err
	}
	t.Seq = int64(binary.BigEndian.Uint64(seqByte))

	//2. 读取总长度
	totalByte := make([]byte, 4)
	_, err = io.ReadFull(r, totalByte)
	if err != nil {
		return err
	}
	totalLen := binary.BigEndian.Uint32(totalByte)
	//3. 读取全部数据
	data := make([]byte, totalLen)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return err
	}

	left, right := uint32(0), DATA_LEN
	//4. 获取ObjectName
	objectNameLen := binary.BigEndian.Uint32(data[left:right])

	left = right
	right = left + objectNameLen
	t.ObjectName = utils.Bytes2String(data[left:right])

	//5 .获取 MethodName
	left = right
	right = left + DATA_LEN
	methodNameLen := binary.BigEndian.Uint32(data[left:right])

	left = right
	right = left + methodNameLen
	t.MethodName = utils.Bytes2String(data[left:right])

	// 6. 获取 Payload

	left = right
	right = left + DATA_LEN

	payLoadLen := binary.BigEndian.Uint32(data[left:right])

	left = right
	right = left + payLoadLen

	t.Payload = data[left:right]

	return err
}

type RPCMsgConfig struct {
	MsgTypeConf       MsgType
	CompressTypeConf  CompressType
	SerializeTypeConf SerializeType
	VersionConf       byte
	ObjectName        string
	MethodName        string
	Seq               int64
}

func SendTo(w io.Writer, payload []byte, msgConfig RPCMsgConfig) error {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	msg := NewRPCMsg()
	msg.SetMsgType(msgConfig.MsgTypeConf)
	msg.SetCompressType(msgConfig.CompressTypeConf)
	msg.SetSerializeType(msgConfig.SerializeTypeConf)
	msg.SetVersion(msgConfig.VersionConf)
	msg.Seq = msgConfig.Seq
	msg.ObjectName = msgConfig.ObjectName
	msg.MethodName = msgConfig.MethodName
	msg.Payload = payload
	return msg.SendMsg(w)
}

func RecvFrom(r io.Reader) (*RPCMsg, error) {
	msg := NewRPCMsg()
	err := msg.RecvMsg(r)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
