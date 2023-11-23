package codec

import (
	"bytes"
	"encoding/gob"
)

type GobCodec struct {
}

func (t GobCodec) Encode(i interface{}) ([]byte, error) {
	// buffer 空内存流
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	// 编码后的结果保存到buffer中
	if err := encoder.Encode(i); err != nil {
		return nil, err
	}
	// 返回内存流中的数据
	return buffer.Bytes(), nil
}

func (t GobCodec) Decode(data []byte, i interface{}) error {
	// buffer 内存流（初始内容为data）
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	// 解码buffer中的数据
	return decoder.Decode(i)
}
