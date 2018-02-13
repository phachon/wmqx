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
func (logController *LogController) Index(ctx *fasthttp.RequestCtx) {

}

// log file
func (logController *LogController) File(ctx *fasthttp.RequestCtx) {

}

// log list
func (logController *LogController) List(ctx *fasthttp.RequestCtx) {

}


