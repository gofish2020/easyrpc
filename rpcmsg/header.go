/*
purpose: 定义网络传输的数据包
*/
package rpcmsg

const (
	magicNumber byte = 0xFF // 魔法数
	version     byte = 0x01 //协议版本
	HEADER_LEN  int  = 5    // 固定5字节
)

// 消息类型
type MsgType byte

const (
	Request MsgType = iota
	Response
)

// 压缩类型
type CompressType byte

const (
	None CompressType = iota
	Gzip
	Snappy
)

// 序列化类型
type SerializeType byte

const (
	Gob SerializeType = iota
	Json
)

func NewHeader() Header {

	return Header([HEADER_LEN]byte{})
}

// ********数据包头格式： 【魔法数 协议版本 消息类型 压缩类型 序列化类型】*******
type Header [HEADER_LEN]byte

// 魔法数
func (t *Header) CheckMagicNumber() bool {
	return t[0] == magicNumber
}

func (t *Header) MagicNumber() byte {
	return t[0]
}

// 协议版本
func (t *Header) Version() byte {
	return t[1]
}

func (t *Header) SetVersion(version byte) {
	t[1] = version
}

// 消息类型
func (t *Header) MsgType() MsgType {
	return MsgType(t[2])
}

func (t *Header) SetMsgType(msgType MsgType) {
	t[2] = byte(msgType)
}

// 压缩类型
func (t *Header) CompressType() CompressType {
	return CompressType(t[3])
}

func (t *Header) SetCompressType(compressType CompressType) {
	t[3] = byte(compressType)
}

//序列化类型

func (t *Header) SerializeType() SerializeType {
	return SerializeType(t[4])
}

func (t *Header) SetSerializeType(serializeType SerializeType) {
	t[4] = byte(serializeType)
}
