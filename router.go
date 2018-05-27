package main

import (
	"github.com/buaazp/fasthttprouter"
	"wmqx/app/controllers"
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
	messageController := controllers.NewMessageController()
	router.POST("/message/add", messageController.Add)
	router.POST("/message/update", messageController.Update)
	router.GET("/message/delete", messageController.Delete)
	router.GET("/message/status", messageController.Status)
	router.GET("/message/list", messageController.List)
	router.GET("/message/getMessageByName", messageController.GetMessageByName)
	router.GET("/message/getConsumersByName", messageController.GetConsumersByName)

	// consumer router
	consumerController := controllers.NewConsumerController()
	router.POST("/consumer/add", consumerController.Add)
	router.POST("/consumer/update", consumerController.Update)
	router.GET("/consumer/delete", consumerController.Delete)
	router.GET("/consumer/status", consumerController.Status)
	router.GET("/consumer/getConsumerById", consumerController.GetConsumerById)

	// system router
	systemController := controllers.NewSystemController()
	router.GET("/system/reload", systemController.Reload)
	router.GET("/system/restart", systemController.Restart)

	// log router
	logController := controllers.NewLogController()
	router.GET("/log/index", logController.Index)
	router.GET("/log/search", logController.Search)
	router.GET("/log/list", logController.List)
	router.GET("/log/download", logController.Download)

	return router
}

// publish server router
func (r *Router) Publish() *fasthttprouter.Router {
	router := fasthttprouter.New()

	// publish router
	publishController := controllers.NewPublishController()
	router.GET("/publish/:name", publishController.Publish)
	router.POST("/publish/:name", publishController.Publish)

	return router
}