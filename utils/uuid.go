package utils

import (
	"github.com/zheng-ji/goSnowFlake"
)

var iw *goSnowFlake.IdWorker

func CreateGUID() int64 {
	if iw == nil {
		iw, _ = goSnowFlake.NewIdWorker(1)
	}

	if id, err := iw.NextId(); err != nil {
		return 0
	} else {
		return id
	}
}
