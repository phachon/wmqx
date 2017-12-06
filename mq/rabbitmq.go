package mq

import (
	"github.com/streadway/amqp"
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

	err := channel.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, nil)
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

	autoDelete := true // When there is no consumer, the server can delete the Exchange
	exclusive := false //Exclusive queues are only accessible by the connection that declares them and will be deleted when the connection closes.
	noWait := false // When noWait is true, declare without waiting for a confirmation from the server.
	_, err := channel.QueueDeclare(name, durable, autoDelete, exclusive, noWait, nil)
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

// publish message
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

// close
func (rq *RabbitMQ) Close() {
	rq.Conn.Close()
}