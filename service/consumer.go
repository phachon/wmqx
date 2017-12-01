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

	exchangeName := this.GetCtxString(ctx, "Name")
	comment := this.GetCtxString(ctx, "comment")
	checkCode := this.GetCtxBool(ctx, "check_code")
	code := this.GetCtxFloat64(ctx, "code")
	routeKey := this.GetCtxString(ctx,"route_key")
	timeout:= this.GetCtxFloat64(ctx, "timeout")
	url := this.GetCtxString(ctx, "url")

	if exchangeName == "" || code == 0 || timeout == 0 || url == "" {
		this.jsonError(ctx, "param require!", nil)
		return
	}

	uuId, _ := uuid.NewV4()
	// add consumer to QMessage
	consumer := &mq.Consumer{
		ID: uuId.String(),
		URL: url,
		RouteKey: routeKey,
		Timeout: timeout,
		Code: code,
		CheckCode: checkCode,
		Comment: comment,
	}

	err := container.Ctx.QMessage.AddConsumer(exchangeName, consumer)
}

// update a consumer
func (consumerService *ConsumerService) Update(ctx *fasthttp.RequestCtx) {

}

// delete a consumer
func (consumerService *ConsumerService) Delete(ctx *fasthttp.RequestCtx) {

}

// get consumer status
func (consumerService *ConsumerService) Status(ctx *fasthttp.RequestCtx) {

}