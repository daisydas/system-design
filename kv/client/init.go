package main

import (
	"fmt"
	"sync"
	"system-design/kv/data"
)

var (
	handlerInstance *Handler
	once            sync.Once
)

func Initialize() *Handler {
	once.Do(func() {
		ops := data.NewOperation()
		ops.CreateTable()
		handlerInstance = &Handler{ops}
		fmt.Println("initializing")
	})
	return handlerInstance

}
