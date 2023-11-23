package utils

import "testing"

func TestString2Bytes(t *testing.T) {
	s := "easyrpc is 很好"
	data := String2Bytes(s)
	t.Log(Bytes2String(data))
}
