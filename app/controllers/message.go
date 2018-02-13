package controllers

import (
	"rmqc/container"
	"github.com/valyala/fasthttp"
	"rmqc/message"
	"rmqc/app/service"
	"rmqc/app"
)

type MessageController struct {
	BaseController
}

// return MessageController
func NewMessageController() *MessageController {
	return &MessageController{}
}

// add a message
func (this *MessageController) Add(ctx *fasthttp.RequestCtx) {
	r := this.AccessToken(ctx)
	if r != true {
		this.jsonError(ctx, "token error", nil)
		return
	}

	name := this.GetCtxString(ctx, "name")
	comment := this.GetCtxString(ctx, "comment")
	durable := this.GetCtxBool(ctx, "durable")
	isNeedToken := this.GetCtxBool(ctx, "is_need_token")
	mode := this.GetCtxString(ctx, "mode")
	token := this.GetCtxString(ctx, "token")

	if name == "" || comment == ""{
		this.jsonError(ctx, "param require!", nil)
		return
	}
	if (mode != "fanout") && (mode != "topic") && (mode != "direct") {
		this.jsonError(ctx, "param error!", nil)
		return
	}

	// check message is exists
	ok := container.Ctx.QMessage.IsExistsMessage(name)
	if ok == true {
		this.jsonError(ctx, "message "+name+" is exist", nil)
		return
	}

	msg := &message.Message{
		Consumers  : []*message.Consumer{},
		Durable     : durable,
		IsNeedToken : isNeedToken,
		Mode        : mode,
		Name        : name,
		Token       : token,
		Comment     : comment,
	}

	err := service.NewMQ().DeclareExchange(name, mode, durable)
	if err != nil {
		app.Log.Error("Add message "+name+" failed: "+err.Error())
		this.jsonError(ctx, "add message failed: "+err.Error(), nil)
		return
	}

	err = container.Ctx.QMessage.AddMessage(msg)
	if err != nil {
		app.Log.Error("Add message "+name+" failed: "+err.Error())
		this.jsonError(ctx, "add message failed"+err.Error(), nil)
		return
	}

	app.Log.Info("Add message "+name+" success!")
	this.jsonSuccess(ctx, "success", nil)
}

// update a message
func (this *MessageController) Update(ctx *fasthttp.RequestCtx) {
	r := this.AccessToken(ctx)
	if r != true {
		this.jsonError(ctx, "token error", nil)
		return
	}

	name := this.GetCtxString(ctx, "name")
	comment := this.GetCtxString(ctx, "comment")
	durable := this.GetCtxBool(ctx, "durable")
	isNeedToken := this.GetCtxBool(ctx, "is_need_token")
	mode := this.GetCtxString(ctx, "mode")
	token := this.GetCtxString(ctx, "token")

	if name == "" || comment == "" {
		this.jsonError(ctx, "param require!", nil)
		return
	}
	if (mode != "fanout") && (mode != "topic") && (mode != "direct") {
		this.jsonError(ctx, "param error!", nil)
		return
	}

	// check message is exists
	ok := container.Ctx.QMessage.IsExistsMessage(name)
	if ok == false {
		this.jsonError(ctx, "message "+name+" not exist", nil)
		return
	}

	msg := &message.Message{
		Durable     : durable,
		IsNeedToken : isNeedToken,
		Mode        : mode,
		Name        : name,
		Token       : token,
		Comment     : comment,
	}

	err := service.NewMQ().DeclareExchange(name, mode, durable)
	if err != nil {
		app.Log.Error("Update message "+name+" failed: "+err.Error())
		this.jsonError(ctx, "update message failed: "+err.Error(), nil)
		return
	}

	err = container.Ctx.QMessage.UpdateMessageByName(name, msg)
	if err != nil {
		app.Log.Error("Update message "+name+" failed: "+err.Error())
		this.jsonError(ctx, "update message failed: "+err.Error(), nil)
		return
	}

	app.Log.Info("Update message "+name+" success!")

	this.jsonSuccess(ctx, "success", nil)
}

// delete a message
func (this *MessageController) Delete(ctx *fasthttp.RequestCtx) {
	r := this.AccessToken(ctx)
	if r != true {
		this.jsonError(ctx, "token error", nil)
		return
	}

	name := this.GetCtxString(ctx, "name")
	if name == "" {
		this.jsonError(ctx, "param require!", nil)
		return
	}

	// check message is exists
	ok := container.Ctx.QMessage.IsExistsMessage(name)
	if ok == false {
		this.jsonError(ctx, "message "+name+" not exist", nil)
		return
	}

	err := service.NewMQ().DeleteExchange(name)
	if err != nil {
		app.Log.Error("Delete message "+name+" failed: "+err.Error())
		this.jsonError(ctx, "delete message failed: "+err.Error(), nil)
		return
	}

	err = container.Ctx.QMessage.DeleteMessageByName(name)
	if err != nil {
		app.Log.Error("Delete message "+name+" failed: "+err.Error())
		this.jsonError(ctx, "delete message failed: "+err.Error(), nil)
		return
	}

	app.Log.Info("Delete message "+name+" success!")

	this.jsonSuccess(ctx, "success", nil)
}

// get message all consumer status
func (this *MessageController) Status(ctx *fasthttp.RequestCtx) {

	r := this.AccessToken(ctx)
	if r != true {
		this.jsonError(ctx, "token error", nil)
		return
	}
	name := this.GetCtxString(ctx, "name")
	if name == "" {
		this.jsonError(ctx, "param require!", nil)
		return
	}

	// check message is exists
	ok := container.Ctx.QMessage.IsExistsMessage(name)
	if ok == false {
		this.jsonError(ctx, "message "+name+" not exist", nil)
		return
	}

	data := []map[string]interface{}{}
	consumers := container.Ctx.QMessage.GetConsumersByMessageName(name)

	consumerProcess := container.Ctx.ConsumerProcess.ProcessMessages
	if len(consumers) > 0 {
		for _, consumer := range consumers {
			item := map[string]interface{}{
				"name": name,
				"consumer_id": consumer.ID,
				"status": 0,
				"last_time": 0,
			}
			consumerKey := container.Ctx.GetConsumerKey(name, consumer.ID)
			for _, process := range consumerProcess {
				if process.Key == consumerKey {
					item["status"] = 1
					item["last_time"] = process.LastTime
				}
			}
			data = append(data, item)
		}
	}

	this.jsonSuccess(ctx, "success", data)
}

// get all message list
func (this *MessageController) List(ctx *fasthttp.RequestCtx) {
	r := this.AccessToken(ctx)
	if r != true {
		this.jsonError(ctx, "token error", nil)
		return
	}

	messages := container.Ctx.QMessage.GetMessages()

	this.jsonSuccess(ctx, "success", messages)
}

// get message by name
func (this *MessageController) GetMessageByName(ctx *fasthttp.RequestCtx) {
	r := this.AccessToken(ctx)
	if r != true {
		this.jsonError(ctx, "token error", nil)
		return
	}

	name := this.GetCtxString(ctx, "name")
	if name == "" {
		this.jsonError(ctx, "param require!", nil)
		return
	}
	msg, err := container.Ctx.QMessage.GetMessageByName(name)
	if err != nil {
		this.jsonError(ctx, err.Error(), nil)
		return
	}

	this.jsonSuccess(ctx, "success", msg)
}

// get consumer by name
func (this *MessageController) GetConsumerByName(ctx *fasthttp.RequestCtx) {
	r := this.AccessToken(ctx)
	if r != true {
		this.jsonError(ctx, "token error", nil)
		return
	}

	name := this.GetCtxString(ctx, "name")
	if name == "" {
		this.jsonError(ctx, "param require!", nil)
		return
	}
	consumers := container.Ctx.QMessage.GetConsumersByMessageName(name)

	this.jsonSuccess(ctx, "success", consumers)
}