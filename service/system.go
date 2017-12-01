package service

import "github.com/valyala/fasthttp"

type SystemService struct {
	BaseService
}

// return SystemService
func NewSystemService() *SystemService {
	return &SystemService{}
}

// reload wmq
func (systemService *SystemService) Reload(ctx *fasthttp.RequestCtx) {

}

// restart wmq
func (systemService *SystemService) Restart(ctx *fasthttp.RequestCtx) {

}

// config wmq
func (systemService *SystemService) Config(ctx *fasthttp.RequestCtx) {

}
