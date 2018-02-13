package mq

import (
	"github.com/streadway/amqp"
	"errors"
)

type RabbitMQ struct {
	Uri string //amqp://guest:guest@localhost:5672/
	Conn *amqp.Connection
	Channel *amqp.Channel
}

// return a rabbitMQ object
func NewRabbitMQ(uri string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return &RabbitMQ{}, err
	}
	return &RabbitMQ{
		Uri:uri,
		Conn:conn,
	}, nil
}

// get mq connect
func (rq *RabbitMQ) GetConnect() *amqp.Connection {
	return rq.Conn
}

// get mq Channel
func (rq *RabbitMQ) GetChannel() (*amqp.Channel, error) {
	mqChannel, err := rq.Conn.Channel()
	return mqChannel, err
}

// declare exchange
// params:
// name : exchange name
// kind : exchange type (fanout, topic, direct)
// durable: true or false save exchange when the server is restarted
func (rq *RabbitMQ) DeclareExchange(name string, kind string, durable bool) error {
	channel, _:= rq.Conn.Channel()
	defer channel.Close()

	autoDelete := true // When there is no consumer, the server can delete the Exchange
	internal := false // Exchanges declared as `internal` do not accept accept publishings.
	noWait := false  // When noWait is true, declare without waiting for a confirmation from the server.
	//
	maxRetryCount := 1
	retryCount := 0
	RETRY:
	err := channel.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, nil)
	if err != nil {
		// delete exchange
		rq.DeleteExchange(name)
		retryCount++
		if retryCount > maxRetryCount {
			return err
		}
		goto RETRY
	}
	return nil
}

// delete a exchange
// params:
// name : exchange name
func (rq *RabbitMQ) DeleteExchange(name string) error {
	channel, _:= rq.Conn.Channel()
	defer channel.Close()

	ifUnused := true // When ifUnused is true, the server will only delete the exchange if it has no queue bindings.
	noWait := false // When noWait is true, declare without waiting for a confirmation from the server.
	err := channel.ExchangeDelete(name, ifUnused, noWait)
	if err != nil {
		return err
	}
	return nil
}

// declare queue
// params:
// name : queue name
// durable : durable
func (rq *RabbitMQ) DeclareQueue(name string, durable bool) error {
	channel, _:= rq.Conn.Channel()
	defer channel.Close()

	autoDelete := true // When there is no consumer, the server can delete the queue
	exclusive := false //Exclusive queues are only accessible by the connection that declares them and will be deleted when the connection closes.
	noWait := false // When noWait is true, declare without waiting for a confirmation from the server.

	maxRetryCount := 1
	retryCount := 0
	RETRY:
	_, err := channel.QueueDeclare(name, durable, autoDelete, exclusive, noWait, nil)
	if err != nil {
		rq.DeleteQueue(name)
		retryCount++
		if retryCount > maxRetryCount {
			return err
		}
		goto RETRY
	}
	return nil
}

// delete queue
// params:
// name : queue name
func (rq *RabbitMQ) DeleteQueue(name string) error {
	channel, _:= rq.Conn.Channel()
	defer channel.Close()

	ifUnused := false // When ifUnused is true, the queue will not be deleted if there are any consumers on the queue.
	ifEmpty := false // When ifEmpty is true, the queue will not be deleted if there are any messages remaining on the queue.
	noWait := false // When noWait is true, the queue will be deleted without waiting for a response from the server.
	_, err := channel.QueueDelete(name, ifUnused, ifEmpty, noWait)
	if err != nil {
		return err
	}

	return nil

}

// bind queue to exchange
// params:
// consumerKey : message name+"-"consumerID
// messageName : message name
// routeKey : consumer routeKey
func (rq *RabbitMQ) BindQueueToExchange(consumerKey string, messageName string, routeKey string) error {
	channel, _:= rq.Conn.Channel()
	defer channel.Close()

	noWait := false
	err := channel.QueueBind(consumerKey, routeKey, messageName, noWait, nil)
	return err
}

// unbind queue to exchange
// params:
// consumerKey : message name+"-"consumerID
// messageName : message name
// routeKey : consumer routeKey
func (rq *RabbitMQ) UnBindQueueToExchange(consumerKey string, messageName string, routeKey string) error {
	channel, _:= rq.Conn.Channel()
	defer channel.Close()

	err := channel.QueueUnbind(consumerKey, routeKey, messageName, nil)
	return err
}

// publish message
// params:
// exchange : message name
// routeKey : route_key
// body : publish body
func (rq *RabbitMQ) Publish(exchange string, routeKey string, body string) error {
	channel, _:= rq.Conn.Channel()
	defer channel.Close()

	mandatory := false // todo
	immediate := false // todo
	msg := amqp.Publishing{
		Body: []byte(body),
	}

	return channel.Publish(exchange, routeKey, mandatory, immediate, msg)
}

// declare consumer
func (rq *RabbitMQ) DeclareConsumer(consumerKey string, durable bool, messageName string, consumerRouteKey string) (err error) {

	// declare queue
	err = rq.DeclareQueue(consumerKey, durable)
	if err != nil {
		return errors.New("Declare queue faild: "+err.Error())
	}
	// bind queue to exchange
	err = rq.BindQueueToExchange(consumerKey, messageName, consumerRouteKey)
	if err != nil {
		return errors.New("bind queue exchange fail: "+err.Error())
	}
	return
}

// consume
func (rq *RabbitMQ) Consume(queue string, consumer string) (<-chan amqp.Delivery,error) {
	channel, _ := rq.Conn.Channel()
	defer channel.Close()

	autoAck := false // is auto ack
	exclusive := false //
	noLocal := false
	noWait := false

	delivery, err := channel.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, nil)

	// todo consumer
	return delivery, err
}

// close conn
func (rq *RabbitMQ) Close() {
	rq.Conn.Close()
}