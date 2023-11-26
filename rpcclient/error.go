package rpcclient

import "errors"

var ErrParam = errors.New("param not adapted")

var ErrClient = errors.New("client disconnection ")

var ErrServer = errors.New("server inner error")
