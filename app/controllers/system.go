package controllers

import "github.com/valyala/fasthttp"

type SystemController struct {
	BaseController
}

// return SystemController
func NewSystemController() *SystemController {
	return &SystemController{}
}

// reload wmq
func (systemController *SystemController) Reload(ctx *fasthttp.RequestCtx) {

}

// restart wmq
func (systemController *SystemController) Restart(ctx *fasthttp.RequestCtx) {

}

// config wmq
func (systemController *SystemController) Config(ctx *fasthttp.RequestCtx) {

}
