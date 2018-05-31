package service

import (
	"wmqx/container"
	"errors"
)

var MQ = NewMQ()

func NewMQ() *MqService {
	return &MqService{}
}

type MqService struct {

}

// reload RabbitMQ all exchanges
func (s *MqService) ReloadExchanges() error {
	rabbitMq, err := container.Ctx.RabbitMQPools.GetMQ()
	if err != nil {
		return errors.New("rabbitmq pools failed: " + err.Error())
	}
	defer container.Ctx.RabbitMQPools.Recover(rabbitMq)

	container.Ctx.QMessage.Lock.Lock()
	defer container.Ctx.QMessage.Lock.Unlock()

	messages := container.Ctx.QMessage.Messages
	for _, msg := range messages {
		// declare exchange
		err := rabbitMq.DeclareExchange(msg.Name, msg.Mode, msg.Durable)
		if err != nil {
			return errors.New("Declare exchange failed: "+err.Error())
		}
		// declare queue
		for _, consumer := range msg.Consumers {
			consumerKey := container.Ctx.GetConsumerKey(msg.Name, consumer.ID)
			err := rabbitMq.DeclareQueue(consumerKey, msg.Durable)
			if err != nil {
				return errors.New("Declare queue failed: "+err.Error())
			}
			// bind queue to exchange
			err = rabbitMq.BindQueueToExchange(consumerKey, msg.Name, consumer.RouteKey)
			if err != nil {
				return errors.New("Bind queue exchange fail: "+err.Error())
			}
		}
	}
	return nil
}

// reload one exchange
func (s *MqService) ReloadExchange(messageName string) error {
	rabbitMq, err := container.Ctx.RabbitMQPools.GetMQ()
	if err != nil {
		return errors.New("rabbitmq pools failed: " + err.Error())
	}
	defer container.Ctx.RabbitMQPools.Recover(rabbitMq)

	container.Ctx.QMessage.Lock.Lock()
	defer container.Ctx.QMessage.Lock.Unlock()

	message, err := container.Ctx.QMessage.GetMessageByName(messageName)
	if err != nil {
		return err
	}
	// declare exchange
	err = rabbitMq.DeclareExchange(message.Name, message.Mode, message.Durable)
	if err != nil {
		return errors.New("Declare exchange failed: "+err.Error())
	}
	// declare queue
	for _, consumer := range message.Consumers {
		consumerKey := container.Ctx.GetConsumerKey(message.Name, consumer.ID)
		err := rabbitMq.DeclareQueue(consumerKey, message.Durable)
		if err != nil {
			return errors.New("Declare queue failed: "+err.Error())
		}
		// bind queue to exchange
		err = rabbitMq.BindQueueToExchange(consumerKey, message.Name, consumer.RouteKey)
		if err != nil {
			return errors.New("Bind queue exchange fail: "+err.Error())
		}
		// restart consumer process
		container.Worker.SendConsumerSign(container.Consumer_Action_Insert, consumerKey)
	}
	
	return nil
}

// declare a Exchange
func (s *MqService) DeclareExchange(name string, mode string, durable bool) error {
	rabbitMq, err := container.Ctx.RabbitMQPools.GetMQ()
	if err != nil {
		return errors.New("rabbitmq pools failed: " + err.Error())
	}
	defer container.Ctx.RabbitMQPools.Recover(rabbitMq)

	// declare exchange
	err = rabbitMq.DeclareExchange(name, mode, durable)
	if err != nil {
		return errors.New("Declare exchange "+name+" failed: "+err.Error())
	}

	return nil
}

// delete a exchange
func (s *MqService) DeleteExchange(name string) error {
	rabbitMq, err := container.Ctx.RabbitMQPools.GetMQ()
	if err != nil {
		return errors.New("rabbitmq pools failed: " + err.Error())
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
				return errors.New("Unbind exchange "+name+" consumer id "+consumer.ID+" failed: "+err.Error())
			}
			// stop consumer
			container.Worker.SendConsumerSign(container.Consumer_Action_Delete, consumerKey)
		}
	}
	// delete exchange
	err = rabbitMq.DeleteExchange(name)
	if err != nil {
		return errors.New("Delete exchange "+name+" failed: "+err.Error())
	}

	return nil
}

// declare a consumer
func (s *MqService) DeclareConsumer(consumerId string, messageName string, consumerRouteKey string) error {
	rabbitMq, err := container.Ctx.RabbitMQPools.GetMQ()
	if err != nil {
		return errors.New("Rabbitmq pools failed: " + err.Error())
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
		return errors.New("Declare queue failed: "+err.Error())
	}

	container.Worker.SendConsumerSign(container.Consumer_Action_Insert, consumerKey)
	return nil
}

func (s *MqService ) UnbindStopConsumer(consumerId string, messageName string, consumerRouteKey string) error {

	consumerKey := container.Ctx.GetConsumerKey(messageName, consumerId)

	// stop consumer
	err := container.Ctx.ConsumerProcess.StopProcessByKey(consumerKey)
	if err != nil {
		return errors.New("Stop "+consumerKey+" process failed: "+err.Error())
	}

	rabbitMq, err := container.Ctx.RabbitMQPools.GetMQ()
	if err != nil {
		return errors.New("rabbitmq pools failed: " + err.Error())
	}
	defer container.Ctx.RabbitMQPools.Recover(rabbitMq)

	// unbind queue to exchange
	err = rabbitMq.UnBindQueueToExchange(consumerKey, messageName, consumerRouteKey)
	if err != nil {
		return errors.New("Unbind exchange "+messageName+" consumer id "+consumerId+" failed: "+err.Error())
	}
	// delete queue
	err = rabbitMq.DeleteQueue(consumerKey)
	if err != nil {
		return errors.New("Delete exchange "+messageName+" consumer id "+consumerId+" failed: "+err.Error())
	}
	return nil
}

// stop all consumer
func (s *MqService) StopAllConsumer()  {

	messages := container.Ctx.QMessage.GetMessages()
	if len(messages) > 0 {
		for _, message := range messages {
			messageName := message.Name
			consumers := message.Consumers
			if len(consumers) > 0 {
				for _, consumer := range consumers {
					consumerKey := container.Ctx.GetConsumerKey(messageName, consumer.ID)
					// stop consumer
					container.Worker.SendConsumerSign(container.Consumer_Action_Delete, consumerKey)
				}
			}
		}
	}
}

// publish message to mq
func (s *MqService) Publish(body string, exchangeName string, token string, routeKey string) (err error) {
	qMessage, err := container.Ctx.QMessage.GetMessageByName(exchangeName)
	if err != nil {
		return
	}
	if qMessage.IsNeedToken && qMessage.Token != token {
		return errors.New("token error")
	}

	rabbitMq, err := container.Ctx.RabbitMQPools.GetMQ()
	if err != nil {
		return errors.New("rabbitmq pools failed: "+err.Error())
	}
	defer container.Ctx.RabbitMQPools.Recover(rabbitMq)

	return rabbitMq.Publish(exchangeName, routeKey, body)
}