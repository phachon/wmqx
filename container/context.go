package container

import (
	"wmqx/message"
	"wmqx/pools"
	"errors"
	"strings"
)

var Ctx = NewContext()

func NewContext() *Context {
	return &Context{
		QMessage:&message.QMessage{},
		RabbitMQPools: &pools.RabbitMQ{},
		ConsumerProcess: message.NewConsumerProcess(),
	}
}

type Context struct {

	// QMessage
	QMessage *message.QMessage

	// RabbitMQ pools
	RabbitMQPools *pools.RabbitMQ

	// Consumer Process
	ConsumerProcess *message.ConsumerProcess
}

func (ctx *Context) SetRabbitMQPools(n int)  {
	ctx.RabbitMQPools = pools.NewRabbitMQPools()
	ctx.RabbitMQPools.Init(n)
}

func (ctx *Context) GetConsumerKey(messageName string, consumerId string) string {
	return messageName +"_"+ consumerId
}

func (ctx *Context) SplitConsumerKey(consumerKey string) (messageName string, consumerId string){
	d := strings.Split(consumerKey, "_")
	if len(d) == 2 {
		return d[0], d[1]
	}else {
		return "", d[0]
	}
}

// init RabbitMQ all exchanges
func (ctx *Context) InitExchanges() error {
	rabbitMq, err := ctx.RabbitMQPools.GetMQ()
	if err != nil {
		return errors.New("rabbitmq pools faild: " + err.Error())
	}
	defer ctx.RabbitMQPools.Recover(rabbitMq)

	ctx.QMessage.Lock.Lock()
	defer ctx.QMessage.Lock.Unlock()

	messages := ctx.QMessage.Messages
	for _, msg := range messages {
		// declare exchange
		err := rabbitMq.DeclareExchange(msg.Name, msg.Mode, msg.Durable)
		if err != nil {
			return errors.New("Declare exchange faild: "+err.Error())
		}
		// declare queue
		for _, consumer := range msg.Consumers {
			consumerKey := ctx.GetConsumerKey(msg.Name, consumer.ID)
			err := rabbitMq.DeclareQueue(consumerKey, msg.Durable)
			if err != nil {
				return errors.New("Declare queue faild: "+err.Error())
			}
			// bind queue to exchange
			err = rabbitMq.BindQueueToExchange(consumerKey, msg.Name, consumer.RouteKey)
			if err != nil {
				return errors.New("Bind queue exchange fail: "+err.Error())
			}
			Worker.SendConsumerSign(Consumer_Action_Insert, consumerKey)
		}
	}
	return nil
}