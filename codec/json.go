package codec

import (
	"bytes"
	"encoding/json"
)

/*
purpose: json序列化和反序列化
*/
type JsonCodec struct {
}

func (t JsonCodec) Encode(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

func (t JsonCodec) Decode(data []byte, i interface{}) error {
	decode := json.NewDecoder(bytes.NewBuffer(data))
	decode.UseNumber()
	return decode.Decode(i)
	//return json.Unmarshal(data, i)
}
