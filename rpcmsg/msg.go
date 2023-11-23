package rpcmsg

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/gofish2020/easyrpc/utils"
)

const (
	DATA_LEN uint32 = 4
)

// RPCMsg: 一个完整的数据包 header + body【
type RPCMsg struct {
	Header
	// uint32 表示长度
	ServiceName string
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
	_, err = w.Write(t.Header[:]) // 1.发送header头
	if err != nil {
		return err
	}
	//******
	totalLen := DATA_LEN + uint32(len(t.ServiceName)) + DATA_LEN + uint32(len(t.MethodName)) + DATA_LEN + uint32(len(t.Payload))
	err = binary.Write(w, binary.BigEndian, uint32(totalLen)) // 2.写入总长度
	if err != nil {
		return err
	}
	//******
	err = binary.Write(w, binary.BigEndian, uint32(len(t.ServiceName))) // 3.写入 ServiceName 长度
	if err != nil {
		return err
	}
	_, err = w.Write(utils.String2Bytes(t.ServiceName)) // 4.写入 ServiceName
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
	//4. 获取ServiceName
	serviceNameLen := binary.BigEndian.Uint32(data[left:right])

	left = right
	right = left + serviceNameLen
	t.ServiceName = utils.Bytes2String(data[left:right])

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
