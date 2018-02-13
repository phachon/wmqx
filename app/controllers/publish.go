package controllers

import (
	"github.com/valyala/fasthttp"
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

	//queryString := string(ctx.QueryArgs().QueryString())
	//exchangeName := ctx.UserValue("name").(string)
	//token := string(ctx.Request.Header.Peek("Token"))
	//routeKey := string(ctx.Request.Header.Peek("RouteKey"))
	//method := strings.ToLower(string(ctx.Request.Header.Method()))
	//bodyByte := ctx.Request.Body()
	//ip := string([]byte(ctx.RemoteIP()))
	//
	//qMessage, err := container.Ctx.QMessage.GetMessageByName(exchangeName)
	//if err != nil {
	//	this.jsonError(ctx, err.Error(), nil)
	//	return
	//}
	//if qMessage.IsNeedToken && qMessage.Token != token {
	//	this.jsonError(ctx, "Token error", nil)
	//	return
	//}
	//
	//// ignore header
	//headerMap := make(map[string]string)
	//ignores := container.Ctx.Config.GetStringSlice("publish.IgnoreHeaders")
	//ctx.Request.Header.VisitAll(func(k, v []byte) {
	//	var found = false
	//	for _, ignore := range ignores {
	//		k1 := strings.ToLower(strings.TrimSpace(string(k)))
	//		k2 := strings.ToLower(strings.TrimSpace(ignore))
	//		if k1 == k2 {
	//			found = true
	//			break
	//		}
	//	}
	//	if !found {
	//		headerMap[strings.TrimSpace(string(k))] = string(v)
	//	}
	//})
	//
	//body := base64.StdEncoding.EncodeToString(bodyByte)
	//
	//publishMsg := &message.PublishMessage{
	//	Header:headerMap,
	//	Ip: ip,
	//	Body: body,
	//	Method: method,
	//	Args: queryString,
	//}
	//
	//publishJson, err := publishMsg.JsonEncode(publishMsg)
	//if err != nil {
	//	this.jsonError(ctx, err.Error(), nil)
	//	return
	//}
	//
	//err = container.Ctx.Publish(publishJson, exchangeName, token, routeKey)
	//if err != nil {
	//	this.jsonError(ctx, err.Error(), nil)
	//	return
	//}

	this.jsonSuccess(ctx, "success", 1)
}
