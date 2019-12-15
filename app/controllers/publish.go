package controllers

import (
	"github.com/phachon/wmqx/app"
	"github.com/phachon/wmqx/app/service"
	"github.com/phachon/wmqx/container"
	"github.com/phachon/wmqx/message"
	"github.com/valyala/fasthttp"
	"strings"
)

type PublishController struct {
	BaseController
}

// return PublishController
func NewPublishController() *PublishController {
	return &PublishController{}
}

// publish message
func (this *PublishController) Publish(ctx *fasthttp.RequestCtx) {

	messageTokenHeader := app.Conf.GetString("publish.messageTokenHeader")
	messageRouteKey := app.Conf.GetString("publish.messageRouteKeyHeader")
	queryString := string(ctx.QueryArgs().QueryString())
	exchangeName := ctx.UserValue("name").(string)
	messageToken := string(ctx.Request.Header.Peek(messageTokenHeader))
	routeKey := string(ctx.Request.Header.Peek(messageRouteKey))
	method := strings.ToLower(string(ctx.Request.Header.Method()))
	bodyByte := ctx.Request.Body()
	ip := ctx.RemoteIP().String()

	qMessage, err := container.Ctx.QMessage.GetMessageByName(exchangeName)
	if err != nil {
		this.jsonError(ctx, err.Error(), nil)
		return
	}
	if qMessage.IsNeedToken && qMessage.Token != messageToken {
		this.jsonError(ctx, "message token error", nil)
		return
	}

	// ignore header
	headerMap := make(map[string]string)
	ignores := app.Conf.GetStringSlice("publish.ignoreHeaders")
	ctx.Request.Header.VisitAll(func(k, v []byte) {
		var found = false
		for _, ignore := range ignores {
			k1 := strings.ToLower(strings.TrimSpace(string(k)))
			k2 := strings.ToLower(strings.TrimSpace(ignore))
			if k1 == k2 {
				found = true
				break
			}
		}
		if !found {
			headerMap[strings.TrimSpace(string(k))] = string(v)
		}
	})

	publishMsg := &message.PublishMessage{
		Header: headerMap,
		Ip:     ip,
		Body:   string(bodyByte),
		Method: method,
		Args:   queryString,
	}
	publishJson, err := publishMsg.Encode()
	if err != nil {
		this.jsonError(ctx, err.Error(), nil)
		return
	}

	err = service.MQ.Publish(publishJson, exchangeName, messageToken, routeKey)
	if err != nil {
		app.Log.Errorf("message %s publish failed, %s", exchangeName, err.Error())
		this.jsonError(ctx, err.Error(), nil)
		return
	}

	app.Log.Infof("message %s publish message %s success!", exchangeName, publishMsg.EncodeOriginalString())
	this.jsonSuccess(ctx, "success", nil)
}
