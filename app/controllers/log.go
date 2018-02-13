package controllers

import "github.com/valyala/fasthttp"

type LogController struct {
	BaseController
}

// return LogController
func NewLogController() *LogController {
	return &LogController{}
}

// log index
func (this *LogController) Index(ctx *fasthttp.RequestCtx) {

}

// log file
func (this *LogController) File(ctx *fasthttp.RequestCtx) {

}

// log list
func (this *LogController) List(ctx *fasthttp.RequestCtx) {

}


