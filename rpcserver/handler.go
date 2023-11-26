package rpcserver

import "reflect"

type Handler interface {
	Handle(string, []interface{}) ([]interface{}, error)
}

type RPCHandler struct {
	object reflect.Value
}

func (handler *RPCHandler) Handle(methodName string, params []interface{}) ([]interface{}, error) {
	argsIn := make([]reflect.Value, len(params))
	for i := range params {
		argsIn[i] = reflect.ValueOf(params[i])
	}

	method := handler.object.MethodByName(methodName)
	argsOut := method.Call(argsIn)

	result := make([]interface{}, len(argsOut))
	for i := range argsOut {
		result[i] = argsOut[i].Interface()
	}

	var err error
	if _, ok := result[len(result)-1].(error); ok {
		err = result[len(result)-1].(error)
	}
	return result, err
}
