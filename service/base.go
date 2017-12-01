package service

import (
	"github.com/valyala/fasthttp"
	"encoding/json"
	"rmqc/container"
	"strconv"
)

type BaseService struct {

}

type JsonResult struct {
	Code int `json:"code"`
	Message interface{} `json:"message"`
	Data interface{} `json:"data"`
}

// get request content text string
func (baseService *BaseService) GetCtxString(ctx *fasthttp.RequestCtx, key string) string {
	return string(ctx.QueryArgs().Peek(key))
}

// get request content text bool
func (baseService *BaseService) GetCtxBool(ctx *fasthttp.RequestCtx, key string) bool {
	str := string(ctx.QueryArgs().Peek(key))
	if str == "1" {
		return true
	}
	return false
}

// get request content text int
func (baseService *BaseService) GetCtxInt(ctx *fasthttp.RequestCtx, key string) int {
	str := string(ctx.QueryArgs().Peek(key))
	i, _ := strconv.Atoi(str)
	return i
}

// get request content text float64
func (baseService *BaseService) GetCtxFloat64(ctx *fasthttp.RequestCtx, key string) float64 {
	str := string(ctx.QueryArgs().Peek(key))
	i, _ := strconv.Atoi(str)
	return float64(i)
}


// access token
func (baseService *BaseService) AccessToken(ctx *fasthttp.RequestCtx) bool {
	token := ctx.QueryArgs().Peek("api_token")
	apiToken := container.Ctx.Config.GetString("api.token")
	if token == nil || string(token) != apiToken {
		return false
	}
	return true
}

// return json error
func (baseService *BaseService) jsonError(ctx *fasthttp.RequestCtx, message interface{}, data interface{}) {
	baseService.jsonResult(ctx, 0, message, data)
}

// return json success
func (baseService *BaseService) jsonSuccess(ctx *fasthttp.RequestCtx, message interface{}, data interface{}) {
	baseService.jsonResult(ctx, 1, message, data)
}

// return json result
func (baseService *BaseService) jsonResult(ctx *fasthttp.RequestCtx, code int, message interface{}, data interface{}) {
	if message == nil {
		message = ""
	}
	if data == nil {
		data = map[string]string{}
	}

	res := JsonResult {
		Code:    code,
		Message: message,
		Data:    data,
	}

	jsonByte, err := json.Marshal(res)
	if err != nil {
		ctx.Write([]byte(err.Error()))
	} else {
		ctx.Write(jsonByte)
	}
}