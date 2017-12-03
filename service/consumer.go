package service

import (
	"github.com/valyala/fasthttp"
	"github.com/nu7hatch/gouuid"
	"rmqc/mq"
	"rmqc/container"
)

type ConsumerService struct {
	BaseService
}

// return ConsumerService
func NewConsumerService() *ConsumerService {
	return &ConsumerService{}
}

// add a consumer
func (this *ConsumerService) Add(ctx *fasthttp.RequestCtx) {
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

	uuId, _ := uuid.NewV4()
	consumer := &mq.Consumer{
		ID: uuId.String(),
		URL: url,
		RouteKey: routeKey,
		Timeout: timeout,
		Code: code,
		CheckCode: checkCode,
		Comment: comment,
	}
	
	// add a consumer to QMessage
	err := container.Ctx.QMessage.AddConsumer(exchangeName, consumer)
	if err != nil {
		this.jsonError(ctx, err.Error(), nil)
		return
	}
	
	container.Ctx.ResetQMessage()
	
	this.jsonSuccess(ctx, "ok", nil)
}

// update a consumer
func (this *ConsumerService) Update(ctx *fasthttp.RequestCtx) {
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
	
	consumer := &mq.Consumer{
		ID: consumerId,
		URL: url,
		RouteKey: routeKey,
		Timeout: timeout,
		Code: code,
		CheckCode: checkCode,
		Comment: comment,
	}
	
	// update a consumer to QMessage
	err := container.Ctx.QMessage.UpdateConsumerByName(exchangeName, consumer)
	if err != nil {
		this.jsonError(ctx, err.Error(), nil)
		return
	}
	
	container.Ctx.ResetQMessage()
	
	this.jsonSuccess(ctx, "ok", nil)
}

// delete a consumer
func (this *ConsumerService) Delete(ctx *fasthttp.RequestCtx) {
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
	
	// delete a consumer to QMessage
	container.Ctx.QMessage.DeleteConsumerByNameAndId(exchangeName, consumerId)
	
	container.Ctx.ResetQMessage()
	
	this.jsonSuccess(ctx, "ok", nil)
}

// get consumer status
func (this *ConsumerService) Status(ctx *fasthttp.RequestCtx) {

}