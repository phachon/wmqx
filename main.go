package main

import (
	"wmqx/app"
	"github.com/valyala/fasthttp"
	"os"
	"fmt"
	"wmqx/container"
	"wmqx/message"
)

// WMQX RabbitMQ Callback

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

	recordType := app.Conf.GetString("message.record_type")
	filename := app.Conf.GetString("message.filename")
	jsonBeautify := app.Conf.GetBool("message.jsonBeautify")

	fileConfig := &message.RecordFileConfig{
		Filename: filename,
		JsonBeautify: jsonBeautify,
	}
	qm, err := message.NewQMessage(recordType, message.NewRecordConfigFile(fileConfig))
	if err != nil {
		app.Log.Error(err.Error())
		os.Exit(1)
	}
	container.Ctx.QMessage = qm

	app.Log.Info("Init QMessage "+recordType+ " success!")
}

// init Ctx RabbitMq pools and check rabbitMQ conn
func initRabbitMQPools() {
	container.Ctx.SetRabbitMQPools(20)
	mq, err := container.Ctx.RabbitMQPools.GetMQ()
	defer container.Ctx.RabbitMQPools.Recover(mq)
	if err != nil {
		app.Log.Info("Init Rabbitmq pools falied: "+err.Error())
		os.Exit(1)
	}
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