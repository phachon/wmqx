package service

import "github.com/valyala/fasthttp"

type LogService struct {
	BaseService
}

// return LogService
func NewLogService() *LogService {
	return &LogService{}
}

// log index
func (logService *LogService) Index(ctx *fasthttp.RequestCtx) {

}

// log file
func (logService *LogService) File(ctx *fasthttp.RequestCtx) {

}

// log list
func (logService *LogService) List(ctx *fasthttp.RequestCtx) {

}


