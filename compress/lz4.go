package compress

import (
	"bytes"
	"io"

	"github.com/pierrec/lz4/v4"
)

type Lz4 struct {
}

func (t Lz4) Compress(data []byte) ([]byte, error) {

	zout := bytes.NewBuffer([]byte{})

	zw := lz4.NewWriter(zout)
	if err := zw.Apply([]lz4.Option{lz4.CompressionLevelOption(lz4.Level1)}...); err != nil {
		return nil, err
	}

	_, err := io.Copy(zw, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	err = zw.Close()
	if err != nil {
		return nil, err
	}
	return zout.Bytes(), nil
}

func (t Lz4) UnCompress(data []byte) ([]byte, error) {

	zin := bytes.NewReader(data)
	zr := lz4.NewReader(zin)
	zout := new(bytes.Buffer)

	_, err := io.Copy(zout, zr)
	if err != nil {
		return nil, err
	}

	return zout.Bytes(), nil
}

func GetLz4Compresser() Lz4 {
	return Lz4{}
}
