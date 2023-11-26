package compress

import "github.com/golang/snappy"

type Snappy struct {
}

func GetSnappyCompresser() Snappy {
	return Snappy{}
}

func (t Snappy) Compress(data []byte) ([]byte, error) {
	return snappy.Encode(nil, data), nil
}

func (t Snappy) UnCompress(data []byte) ([]byte, error) {
	return snappy.Decode(nil, data)
}
