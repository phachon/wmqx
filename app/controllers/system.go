package controllers

import (
	"github.com/valyala/fasthttp"
	"rmqc/app/service"
	"rmqc/app"
	"rmqc/container"
)

type SystemController struct {
	BaseController
}

// return SystemController
func NewSystemController() *SystemController {
	return &SystemController{}
}

// reload rmqc
func (this *SystemController) Reload(ctx *fasthttp.RequestCtx) {
	r := this.AccessToken(ctx)
	if r != true {
		this.jsonError(ctx, "token error", nil)
		return
	}

	app.Log.Info("Start reload all exchange")
	err := service.NewMQ().ReloadExchanges()
	if err != nil {
		app.Log.Error("Reload error: "+err.Error())
		this.jsonError(ctx, "reload error: "+err.Error(), nil)
		return
	}
	app.Log.Info("Reload all exchange success!")

	this.jsonSuccess(ctx, "success", nil)
}

// restart rmqc
func (this *SystemController) Restart(ctx *fasthttp.RequestCtx) {
	r := this.AccessToken(ctx)
	if r != true {
		this.jsonError(ctx, "token error", nil)
		return
	}

	app.Log.Info("Start restart")

	service.NewMQ().StopAllConsumer()

	err := container.Ctx.InitExchanges()
	if err != nil {
		app.Log.Error("Restart error: "+err.Error())
		this.jsonError(ctx, "reload error: "+err.Error(), nil)
		return
	}
	app.Log.Info("Restart init exchange success!")

	app.Log.Info("Restart consumer success!")

	this.jsonSuccess(ctx, "success", nil)
}