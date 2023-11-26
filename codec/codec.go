/*
purpose: 序列化和反序列化
author: nash
date: 2023/11/23
*/
package codec

type Codec interface {
	Encode(i interface{}) ([]byte, error)
	Decode(data []byte, i interface{}) error
}





