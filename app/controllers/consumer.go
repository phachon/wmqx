package controllers

import (
	"github.com/valyala/fasthttp"
	"github.com/nu7hatch/gouuid"
	"rmqc/container"
	"rmqc/message"
	"rmqc/app/service"
	"rmqc/app"
)

type ConsumerController struct {
	BaseController
}

// return ConsumerController
func NewConsumerController() *ConsumerController {
	return &ConsumerController{}
}

// add a consumer
func (this *ConsumerController) Add(ctx *fasthttp.RequestCtx) {
	r := this.AccessToken(ctx)
	if r != true {
		this.jsonError(ctx, "token error", nil)
		return
	}

	exchangeName := this.GetCtxString(ctx, "name")
	comment := this.GetCtxString(ctx, "comment")
	checkCode := this.GetCtxBool(ctx, "check_code")
	code := this.GetCtxFloat64(ctx, "code")
	routeKey := this.GetCtxString(ctx,"route_key")
	timeout:= this.GetCtxFloat64(ctx, "timeout")
	url := this.GetCtxString(ctx, "url")

	if exchangeName == "" || timeout == 0 || url == "" {
		this.jsonError(ctx, "param require!", nil)
		return
	}
	if checkCode == true && code == 0 {
		this.jsonError(ctx, "param code require!", nil)
		return
	}

	// check message name is exist
	ok := container.Ctx.QMessage.IsExistsMessage(exchangeName)
	if ok == false {
		this.jsonError(ctx, "message name not exists", nil)
		return
	}

	uuId, _ := uuid.NewV4()
	consumer := &message.Consumer{
		ID: uuId.String(),
		URL: url,
		RouteKey: routeKey,
		Timeout: timeout,
		Code: code,
		CheckCode: checkCode,
		Comment: comment,
	}

	// declare queue and bind consumer to exchange
	err := service.NewMQ().DeclareConsumer(consumer.ID, exchangeName, routeKey, false)
	if err != nil {
		app.Log.Error("Add Consumer faild: "+err.Error())
		this.jsonError(ctx, "add consumer faild: "+err.Error(), nil)
		return
	}

	// add a consumer to QMessage
	err = container.Ctx.QMessage.AddConsumer(exchangeName, consumer)
	if err != nil {
		app.Log.Error("Add Consumer faild: "+err.Error())
		this.jsonError(ctx, "add consumer faild: "+err.Error(), nil)
		return
	}

	app.Log.Info("add consumer success, message: "+exchangeName+" consumer_id: "+consumer.ID)
	this.jsonSuccess(ctx, "success", nil)
}

// update a consumer
func (this *ConsumerController) Update(ctx *fasthttp.RequestCtx) {
	r := this.AccessToken(ctx)
	if r != true {
		this.jsonError(ctx, "token error", nil)
		return
	}
	
	consumerId := this.GetCtxString(ctx, "consumer_id")
	exchangeName := this.GetCtxString(ctx, "name")
	comment := this.GetCtxString(ctx, "comment")
	checkCode := this.GetCtxBool(ctx, "check_code")
	code := this.GetCtxFloat64(ctx, "code")
	routeKey := this.GetCtxString(ctx,"route_key")
	timeout:= this.GetCtxFloat64(ctx, "timeout")
	url := this.GetCtxString(ctx, "url")
	
	if consumerId == "" || exchangeName == "" || timeout == 0 || url == "" {
		this.jsonError(ctx, "param require!", nil)
		return
	}
	if checkCode == true && code == 0 {
		this.jsonError(ctx, "param code require!", nil)
		return
	}

	// check message and consumerId is exist
	ok := container.Ctx.QMessage.IsExistsMessageAndConsumerId(exchangeName, consumerId)
	if ok == false {
		this.jsonError(ctx, "message name or conusmerId not exists", nil)
		return
	}

	consumer := &message.Consumer{
		ID: consumerId,
		URL: url,
		RouteKey: routeKey,
		Timeout: timeout,
		Code: code,
		CheckCode: checkCode,
		Comment: comment,
	}

	// declare queue and bind consumer to exchange
	err := service.NewMQ().DeclareConsumer(consumer.ID, exchangeName, routeKey, true)
	if err != nil {
		app.Log.Error("Update message "+exchangeName+" consumer "+consumerId+" faild: "+err.Error())
		this.jsonError(ctx, "update consumer faild: "+err.Error(), nil)
		return
	}

	// update a consumer to QMessage
	err = container.Ctx.QMessage.UpdateConsumerByName(exchangeName, consumer)
	if err != nil {
		app.Log.Error("Update message "+exchangeName+" consumer "+consumerId+" faild: "+err.Error())
		this.jsonError(ctx, "udpate consumer failed: "+err.Error(), nil)
		return
	}

	app.Log.Error("Update message "+exchangeName+" consumer "+consumerId+" success")
	this.jsonSuccess(ctx, "ok", nil)
}

// delete a consumer
func (this *ConsumerController) Delete(ctx *fasthttp.RequestCtx) {
	r := this.AccessToken(ctx)
	if r != true {
		this.jsonError(ctx, "token error", nil)
		return
	}

	consumerId := this.GetCtxString(ctx, "consumer_id")
	exchangeName := this.GetCtxString(ctx, "name")

	if consumerId == "" || exchangeName == "" {
		this.jsonError(ctx, "param require!", nil)
		return
	}

	// check message and consumerId is exist
	ok := container.Ctx.QMessage.IsExistsMessageAndConsumerId(exchangeName, consumerId)
	if ok == false {
		this.jsonError(ctx, "message name or conusmerId not exists", nil)
		return
	}

	// No need to be deleted queue, queue no consumer is auto delete
	consumerKey := container.Ctx.GetConsumerKey(exchangeName, consumerId)
	container.Worker.SendConsumerSign(container.Consumer_Action_Delete, consumerKey)
	
	 // delete a consumer to QMessage
	err := container.Ctx.QMessage.DeleteConsumerByNameAndId(exchangeName, consumerId)
	if err != nil {
		app.Log.Error("Delete message "+exchangeName+" consumer "+consumerId+" faild: "+err.Error())
		this.jsonError(ctx, "delete consumer failed: "+err.Error(), nil)
		return
	}

	app.Log.Info("Delete message "+exchangeName+" consumer "+consumerId+" success")
	this.jsonSuccess(ctx, "ok", nil)
}

// get consumer status
func (this *ConsumerController) Status(ctx *fasthttp.RequestCtx) {

	r := this.AccessToken(ctx)
	if r != true {
		this.jsonError(ctx, "token error", nil)
		return
	}

	name := this.GetCtxString(ctx, "name")
	consumerId := this.GetCtxString(ctx, "consumer_id")
	if name == "" {
		this.jsonError(ctx, "param require!", nil)
		return
	}
	if consumerId == "" {
		this.jsonError(ctx, "param require!", nil)
		return
	}

	// check message and consumerId is exists
	ok := container.Ctx.QMessage.IsExistsMessageAndConsumerId(name, consumerId)
	if ok == false {
		this.jsonError(ctx, "message or consumerId not exist", nil)
		return
	}

	consumer, err := container.Ctx.QMessage.GetConsumerById(name, consumerId)
	if err != nil {
		this.jsonError(ctx, "get status failed:"+err.Error(), nil)
		return
	}
	data := map[string]interface{}{
		"name": name,
		"consumer_id": consumerId,
		"status": 0,
		"last_time": 0,
	}

	consumerProcess := container.Ctx.ConsumerProcess.ProcessMessages
	consumerKey := container.Ctx.GetConsumerKey(name, consumer.ID)
	for _, process := range consumerProcess {
		if process.Key == consumerKey {
			data["status"] = 1
			data["last_time"] = process.LastTime
		}
	}

	this.jsonSuccess(ctx, "success", data)
}

// get consumer by id
func (this *ConsumerController) GetConsumerById(ctx *fasthttp.RequestCtx) {

	r := this.AccessToken(ctx)
	if r != true {
		this.jsonError(ctx, "token error", nil)
		return
	}

	name := this.GetCtxString(ctx, "name")
	consumerId := this.GetCtxString(ctx, "consumer_id")
	if name == "" {
		this.jsonError(ctx, "param require!", nil)
		return
	}
	if consumerId == "" {
		this.jsonError(ctx, "param require!", nil)
		return
	}

	consumer, err := container.Ctx.QMessage.GetConsumerById(name, consumerId)
	if err != nil {
		this.jsonError(ctx, err.Error(), nil)
		return
	}

	this.jsonSuccess(ctx, "success", consumer)
}
