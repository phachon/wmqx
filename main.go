package main

import (
	"wmqx/app"
	"github.com/valyala/fasthttp"
	"os"
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
		app.Log.Errorf("Init QMessage Error: %s", err.Error())
		os.Exit(1)
	}
	container.Ctx.QMessage = qm

	app.Log.Info("Init QMessage "+recordType+ " success!")
}

// init Ctx RabbitMq pools and check rabbitMQ conn
func initRabbitMQPools() {
	poolNumber := app.Conf.GetInt("rabbitmq.poolNumber")
	container.Ctx.SetRabbitMQPools(poolNumber)
	mq, err := container.Ctx.RabbitMQPools.GetMQ()
	defer container.Ctx.RabbitMQPools.Recover(mq)
	if err != nil {
		app.Log.Errorf("Init Rabbitmq pools falied: %s", err.Error())
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
				app.Log.Errorf("strat api server crash: %v", e)
			}
		}()
		apiListen := app.Conf.GetString("listen.api")
		app.Log.Info("Api server start listen: "+apiListen +"!")
		err := fasthttp.ListenAndServe(apiListen, NewRouter().Api().Handler)
		if err != nil {
			app.Log.Errorf("Api server listen failed: %s", err.Error())
		}
	}()
}

// start publish server
func startPublishServer()  {

	publishListen := app.Conf.GetString("listen.publish")
	app.Log.Info("Publish Server start listen: "+publishListen+"!")
	err := fasthttp.ListenAndServe(publishListen, NewRouter().Publish().Handler)
	if err != nil {
		app.Log.Errorf("Publish Server listen failed: %s", err.Error())
		os.Exit(1)
	}
}