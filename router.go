package main

import (
	"github.com/buaazp/fasthttprouter"
	"rmqc/service"
)

type Router struct {

}

// return router
func NewRouter() *Router {
	return &Router{}
}

// api manager server router
func (r *Router) Api() *fasthttprouter.Router {
	router := fasthttprouter.New()

	// message router
	messageService := service.NewMessageService()
	router.GET("/message/add", messageService.Add)
	router.GET("/message/update", messageService.Update)
	router.GET("/message/delete", messageService.Delete)
	router.GET("/message/status", messageService.Status)

	// consumer router
	consumerService := service.NewConsumerService()
	router.GET("/consumer/add", consumerService.Add)
	router.GET("/consumer/update", consumerService.Update)
	router.GET("/consumer/delete", consumerService.Delete)
	router.GET("/consumer/status", consumerService.Status)

	// system router
	systemService := service.NewSystemService()
	router.GET("/reload", systemService.Reload)
	router.GET("/restart", systemService.Restart)
	router.GET("/config", systemService.Config)

	// log router
	logService := service.NewLogService()
	router.GET("/log", logService.Index)
	router.GET("/log/file", logService.File)
	router.GET("/log/list", logService.List)

	return router
}

// publish server router
func (r *Router) Publish() *fasthttprouter.Router {
	router := fasthttprouter.New()

	// publish router
	publishService := service.NewPublishService()
	router.GET("/publish/:name", publishService.Publish)
	router.POST("/publish/:name", publishService.Publish)

	return router
}