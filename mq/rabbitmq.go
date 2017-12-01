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
func NewRabbitMQ(uri string) (*RabbitMQ) {
	return &RabbitMQ{
		Uri:uri,
	}
}

// get mq connect
func (rq *RabbitMQ) GetConnect() (*amqp.Connection, error) {
	conn, err := amqp.Dial(rq.Uri)
	return conn, err
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
	autoDelete := true // When there is no consumer, the server can delete the Exchange
	internal := false // Exchanges declared as `internal` do not accept accept publishings.
	noWait := false  // When noWait is true, declare without waiting for a confirmation from the server.

	err := rq.Channel.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, nil)
	if err != nil {
		return err
	}
	return nil
}

func (rq *RabbitMQ) DeclareQueue() {

}

// close
func (rq *RabbitMQ) Close() {
	rq.Conn.Close()
}