package pools

import (
	"github.com/jolestar/go-commons-pool"
	"wmqx/mq"
	"fmt"
	"wmqx/app"
	"github.com/streadway/amqp"
	"net"
	"time"
)

type RabbitMQ struct {
	ObjectPool *pool.ObjectPool
}

func NewRabbitMQPools() *RabbitMQ {
	return &RabbitMQ{}
}

// init rabbitMQ pools
func (rMQPool *RabbitMQ) Init(maxPools int) {
	p := pool.NewObjectPoolWithDefaultConfig(new(RabbitMQFactory))
	p.Config.MaxTotal = maxPools
	rMQPool.ObjectPool = p
}

// get a rabbitMQ obj
func (rMQPool *RabbitMQ) GetMQ() (*mq.RabbitMQ, error) {
	obj, err := rMQPool.ObjectPool.BorrowObject()
	if err != nil {
		return &mq.RabbitMQ{}, err
	}
	return obj.(*mq.RabbitMQ), nil
}

// recover a rabbitMQ obj
func (rMQPool *RabbitMQ) Recover(obj interface{}) {
	rMQPool.ObjectPool.ReturnObject(obj)
}

// pools factory
type RabbitMQFactory struct {

}

// make rabbitMQ object
func (f *RabbitMQFactory) MakeObject() (*pool.PooledObject, error) {

	username := app.Conf.GetString("rabbitmq.username")
	password := app.Conf.GetString("rabbitmq.password")
	host := app.Conf.GetString("rabbitmq.host")
	port := app.Conf.GetInt("rabbitmq.port")
	vHost := app.Conf.GetString("rabbitmq.vhost")
	heartBeat := app.Conf.GetInt("rabbitmq.heartbeat")
	connTimeout := app.Conf.GetInt("rabbitmq.connTimeout")

	uri := fmt.Sprintf("amqp://%s:%s@%s:%d%s", username, password, host, port, vHost)
	config := amqp.Config{
		Heartbeat: time.Duration(heartBeat) * time.Second,
		Dial: func(network, addr string) (net.Conn, error){
			c, err := net.DialTimeout(network, addr, time.Duration(connTimeout) * time.Second)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
	rq, err := mq.NewRabbitMQ(uri, config)
	if err != nil {
		return &pool.PooledObject{}, err
	}
	return pool.NewPooledObject(rq), nil
}

// do destroy
func (f *RabbitMQFactory) DestroyObject(object *pool.PooledObject) error {
	//do destroy
	object.Object.(*mq.RabbitMQ).Close()
	return nil
}

// do validate
func (f *RabbitMQFactory) ValidateObject(object *pool.PooledObject) bool {
	//do validate
	rabbitMq := object.Object.(*mq.RabbitMQ)
	conn := rabbitMq.Conn
	if conn == nil {
		return false
	}
	ch, err := conn.Channel()
	if err != nil || ch == nil {
		return false
	}
	ch.Close()
	return true
}

// do activate
func (f *RabbitMQFactory) ActivateObject(object *pool.PooledObject) error {
	//do activate
	return nil
}

// do passivate
func (f *RabbitMQFactory) PassivateObject(object *pool.PooledObject) error {
	//do passivate
	return nil
}
