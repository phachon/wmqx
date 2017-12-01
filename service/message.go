package service

import (
	"rmqc/mq"
	"rmqc/container"
	"github.com/valyala/fasthttp"
)

type MessageService struct {
	BaseService
}

// return MessageService
func NewMessageService() *MessageService {
	return &MessageService{}
}

// add a message
func (this *MessageService) Add(ctx *fasthttp.RequestCtx) {
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

	if name == "" || comment == "" || durable == "" {
		this.jsonError(ctx, "param require!", nil)
		return
	}

	if (mode != "fanout") && (mode != "topic") && (mode != "direct") {
		this.jsonError(ctx, "param error!", nil)
		return
	}

	message := &mq.Message{
		Consumers  : []*mq.Consumer{},
		Durable     : durable,
		IsNeedToken : isNeedToken,
		Mode        : mode,
		Name        : name,
		Token       : token,
		Comment     : comment,
	}

	err := container.Ctx.QMessage.AddMessage(message)
	if err != nil {
		this.jsonError(ctx, err.Error(), nil)
		return
	}
	container.Ctx.ResetQMessage()

	this.jsonSuccess(ctx, "ok", nil)
}

// update a message
func (this *MessageService) Update(ctx *fasthttp.RequestCtx) {
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

	if name == "" || comment == "" || durable == "" {
		this.jsonError(ctx, "param require!", nil)
		return
	}

	if (mode != "fanout") && (mode != "topic") && (mode != "direct") {
		this.jsonError(ctx, "param error!", nil)
		return
	}

	message := &mq.Message{
		Durable     : durable,
		IsNeedToken : isNeedToken,
		Mode        : mode,
		Name        : name,
		Token       : token,
		Comment     : comment,
	}

	err := container.Ctx.QMessage.UpdateMessageByName(name, message)
	if err != nil {
		this.jsonError(ctx, err.Error(), nil)
		return
	}
	container.Ctx.ResetQMessage()

	this.jsonSuccess(ctx, "ok", nil)
}

// delete a message
func (this *MessageService) Delete(ctx *fasthttp.RequestCtx) {
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

	container.Ctx.QMessage.DeleteMessageByName(name)
	container.Ctx.ResetQMessage()

	this.jsonSuccess(ctx, "ok", nil)
}

// get message status
func (messageService *MessageService) Status(ctx *fasthttp.RequestCtx) {

}