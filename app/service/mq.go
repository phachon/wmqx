package service

import (
	"rmqc/container"
	"errors"
)

func NewMQ() *MqService {
	return &MqService{}
}

type MqService struct {

}

// declare a Exchange
func (s *MqService) DeclareExchange(name string, mode string, durable bool) error {
	rabbitMq, err := container.Ctx.RabbitMQPools.GetMQ()
	if err != nil {
		return errors.New("rabbitmq pools faild: " + err.Error())
	}
	defer container.Ctx.RabbitMQPools.Recover(rabbitMq)

	// declare exchange
	err = rabbitMq.DeclareExchange(name, mode, durable)
	if err != nil {
		return errors.New("Declare exchange "+name+" faild: "+err.Error())
	}

	return nil
}

// delete a exchange
func (s *MqService) DeleteExchange(name string) error {
	rabbitMq, err := container.Ctx.RabbitMQPools.GetMQ()
	if err != nil {
		return errors.New("rabbitmq pools faild: " + err.Error())
	}
	defer container.Ctx.RabbitMQPools.Recover(rabbitMq)

	consumers := container.Ctx.QMessage.GetConsumersByMessageName(name)
	if len(consumers) > 0 {
		//unbind queue exchange
		for _, consumer := range consumers {
			consumerKey := container.Ctx.GetConsumerKey(name, consumer.ID)
			routeKey := consumer.RouteKey
			err := rabbitMq.UnBindQueueToExchange(consumerKey, name, routeKey)
			if err != nil {
				return errors.New("Unbind exchange "+name+" consumer id "+consumer.ID+" faild: "+err.Error())
			}
			// stop consumer
			container.Worker.SendConsumerSign(container.Consumer_Action_Delete, consumerKey)
		}
	}
	// delete exchange
	err = rabbitMq.DeleteExchange(name)
	if err != nil {
		return errors.New("Delete exchange "+name+" faild: "+err.Error())
	}

	return nil
}

// declare a consumer
func (s *MqService) DeclareConsumer(consumerId string, messageName string, consumerRouteKey string, isUpdate bool) error {
	rabbitMq, err := container.Ctx.RabbitMQPools.GetMQ()
	if err != nil {
		return errors.New("Rabbitmq pools faild: " + err.Error())
	}
	defer container.Ctx.RabbitMQPools.Recover(rabbitMq)

	message, err := container.Ctx.QMessage.GetMessageByName(messageName)
	if err != nil {
		return err
	}

	consumerKey := container.Ctx.GetConsumerKey(messageName, consumerId)

	// declare consumer
	err = rabbitMq.DeclareConsumer(consumerKey, message.Durable, messageName, consumerRouteKey)
	if err != nil {
		return errors.New("Declare queue faild: "+err.Error())
	}

	if isUpdate == true {
		container.Worker.SendConsumerSign(container.Consumer_Action_Update, consumerKey)
	}else {
		container.Worker.SendConsumerSign(container.Consumer_Action_Insert, consumerKey)
	}
	return nil
}
