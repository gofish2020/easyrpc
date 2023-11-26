package rpcmsg

import (
	"github.com/gofish2020/easyrpc/codec"
	"github.com/gofish2020/easyrpc/compress"
)

var Codecs = map[SerializeType]codec.Codec{
	Gob:  codec.GobCodec{},
	Json: codec.GobCodec{},
}

var Compressor = map[CompressType]compress.Compression{
	Snappy: compress.GetSnappyCompresser(),
	Zlib:   compress.GetZlibCompresser(),
	Lz4:    compress.GetLz4Compresser(),
}
