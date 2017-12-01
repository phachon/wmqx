package main

import (
	"github.com/valyala/fasthttp"
	"rmqc/container"
)

type Wmq struct {
}

// return wmq
func NewWmq() *Wmq {
	return &Wmq{}
}

// run
func (wmq *Wmq) Run()  {
	wmq.beforeInit()
	go wmq.startApiServer()
	go wmq.startPublishServer()
	select {

	}
}

// before init config
func (wmq *Wmq) beforeInit() {
	container.Ctx.InitConfig()
	container.Ctx.InitQMessage()
}

// start api server
func (wmq *Wmq) startApiServer() {
	fasthttp.ListenAndServe("127.0.0.1:3302", NewRouter().Api().Handler)
}

// start publish server
func (wmq *Wmq) startPublishServer()  {
	fasthttp.ListenAndServe("127.0.0.1::3303", NewRouter().Api().Handler)
}