package controllers

import (
	"github.com/valyala/fasthttp"
	"rmqc/app/service"
	"rmqc/app"
	"rmqc/container"
	"time"
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

	app.Log.Info("Reload start")
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

	app.Log.Info("Restart start ")

	service.NewMQ().StopAllConsumer()

	// wait all consumer stop
	for {
		if len(container.Ctx.ConsumerProcess.ProcessMessages) != 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		app.Log.Info("Restart stop all consumer success!")
		break
	}

	err := container.Ctx.InitExchanges()
	if err != nil {
		app.Log.Error("Restart error: "+err.Error())
		this.jsonError(ctx, "reload error: "+err.Error(), nil)
		return
	}
	app.Log.Info("Restart init exchange success!")

	app.Log.Info("Restart all consumer success!")

	this.jsonSuccess(ctx, "success", nil)
}