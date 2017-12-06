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
	if err := container.Ctx.InitQMessage(); err != nil {
		panic(err.Error())
	}
	container.Ctx.InitMQPools()
	if err := container.Ctx.InitRabbitMQ(); err != nil {
		panic(err.Error())
	}
	container.Ctx.InitConsumer()
}

// start api server
func (wmq *Wmq) startApiServer() {
	err := fasthttp.ListenAndServe("127.0.0.1:3302", NewRouter().Api().Handler)
	if err != nil {
		panic("start publish server fail:"+err.Error())
	}
}

// start publish server
func (wmq *Wmq) startPublishServer()  {
	err := fasthttp.ListenAndServe("127.0.0.1:3303", NewRouter().Publish().Handler)
	if err != nil {
		panic("start publish server fail:"+err.Error())
	}
}