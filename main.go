package main

import (
	"rmqc/app"
	"github.com/valyala/fasthttp"
	"os"
	"fmt"
	"rmqc/container"
	"rmqc/message"
)

// RMQC RabbitMQ Callback

func main()  {
	initQMessage()
	initRabbitMQPools()
	startConsumerWorker()
	initExchange()
	startApiServer()
	startPublishServer()
}

// init Ctx QMessage
func initQMessage() {
	fileConfig := &message.RecordFileConfig{
		Filename: "message.json",
		JsonBeautify: true,
	}
	qm, err := message.NewQMessage("file", message.NewRecordConfigFile(fileConfig))
	if err != nil {
		app.Log.Error(err.Error())
		os.Exit(1)
	}
	container.Ctx.QMessage = qm

	app.Log.Info("Init QMessage file success!")
}

// init Ctx RabbitMq pools
func initRabbitMQPools() {
	container.Ctx.SetRabbitMQPools(20)

	app.Log.Info("Init Rabbitmq pools success!")
}

// start consumer worker
func startConsumerWorker() {
	container.Worker.Consumer()
	app.Log.Info("Rabbitmq consumer worker start!")
}

// init RabbitMq exchange
func initExchange() {
	container.Ctx.InitExchanges()

	app.Log.Info("Init Rabbitmq exchange success!")
}

// start api server
func startApiServer() {
	go func() {
		defer func() {
			e := recover()
			if e != nil {
				fmt.Printf("go runtime error: %v", e)
			}
		}()
		apiListen := app.Conf.GetString("listen.api")
		app.Log.Info("Api server listen start: "+apiListen +"!")
		err := fasthttp.ListenAndServe(apiListen, NewRouter().Api().Handler)
		if err != nil {
			app.Log.Error("Api server listen failed: "+err.Error())
		}
	}()
}

// start publish server
func startPublishServer()  {

	publishListen := app.Conf.GetString("listen.publish")
	app.Log.Info("Publish Server listen start: "+publishListen+"!")
	err := fasthttp.ListenAndServe(publishListen, NewRouter().Publish().Handler)
	if err != nil {
		app.Log.Error("Publish Server listen failed: "+err.Error())
		os.Exit(1)
	}
}