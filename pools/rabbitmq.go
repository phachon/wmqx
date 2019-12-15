package pools

import (
	"context"
	"fmt"
	"github.com/jolestar/go-commons-pool"
	"github.com/phachon/wmqx/app"
	"github.com/phachon/wmqx/mq"
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
	ctx := context.Background()
	p := pool.NewObjectPoolWithDefaultConfig(ctx, new(RabbitMQFactory))
	p.Config.MaxTotal = maxPools
	rMQPool.ObjectPool = p
}

// get a rabbitMQ obj
func (rMQPool *RabbitMQ) GetMQ() (*mq.RabbitMQ, error) {
	ctx := context.Background()
	obj, err := rMQPool.ObjectPool.BorrowObject(ctx)
	if err != nil {
		return &mq.RabbitMQ{}, err
	}
	return obj.(*mq.RabbitMQ), nil
}

// recover a rabbitMQ obj
func (rMQPool *RabbitMQ) Recover(obj interface{}) {
	ctx := context.Background()
	rMQPool.ObjectPool.ReturnObject(ctx, obj)
}

// pools factory
type RabbitMQFactory struct {
}

// make rabbitMQ object
func (f *RabbitMQFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {

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
		Dial: func(network, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(network, addr, time.Duration(connTimeout)*time.Second)
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
func (f *RabbitMQFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	//do destroy
	object.Object.(*mq.RabbitMQ).Close()
	return nil
}

// do validate
func (f *RabbitMQFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
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
func (f *RabbitMQFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	//do activate
	return nil
}

// do passivate
func (f *RabbitMQFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	//do passivate
	return nil
}
