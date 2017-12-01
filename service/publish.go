package service

import "github.com/valyala/fasthttp"

type PublishService struct {
	BaseService
}

// return PublishService
func NewPublishService() *PublishService {
	return &PublishService{}
}

func (publishService *PublishService) Publish(ctx *fasthttp.RequestCtx) {

}
