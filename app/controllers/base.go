package controllers

import (
	"github.com/valyala/fasthttp"
	"encoding/json"
	"strconv"
	"wmqx/app"
)

type BaseController struct {

}

type JsonResult struct {
	Code int `json:"code"`
	Message interface{} `json:"message"`
	Data interface{} `json:"data"`
}

// get request content text string
func (baseService *BaseController) GetCtxString(ctx *fasthttp.RequestCtx, key string) string {
	return string(ctx.FormValue(key))
}

// get request content text bool
func (baseService *BaseController) GetCtxBool(ctx *fasthttp.RequestCtx, key string) bool {
	str := string(ctx.FormValue(key))
	if str == "1" {
		return true
	}
	return false
}

// get request content text int
func (baseService *BaseController) GetCtxInt(ctx *fasthttp.RequestCtx, key string) int {
	str := string(ctx.FormValue(key))
	i, _ := strconv.Atoi(str)
	return i
}

// get request content text float64
func (baseService *BaseController) GetCtxFloat64(ctx *fasthttp.RequestCtx, key string) float64 {
	str := string(ctx.FormValue(key))
	i, _ := strconv.Atoi(str)
	return float64(i)
}

// access token
func (baseService *BaseController) AccessToken(ctx *fasthttp.RequestCtx) bool {
	tokenHeaderName := app.Conf.GetString("api.tokenHeaderName")
	token := ctx.Request.Header.Peek(tokenHeaderName)
	apiToken := app.Conf.GetString("api.token")
	if token == nil || string(token) != apiToken {
		return false
	}
	return true
}

// return json error
func (baseService *BaseController) jsonError(ctx *fasthttp.RequestCtx, message interface{}, data interface{}) {
	baseService.jsonResult(ctx, 0, message, data)
}

// return json success
func (baseService *BaseController) jsonSuccess(ctx *fasthttp.RequestCtx, message interface{}, data interface{}) {
	baseService.jsonResult(ctx, 1, message, data)
}

// return json result
func (baseService *BaseController) jsonResult(ctx *fasthttp.RequestCtx, code int, message interface{}, data interface{}) {
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