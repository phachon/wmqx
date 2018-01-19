package pools

import (
	"github.com/jolestar/go-commons-pool"
	"rmqc/mq"
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
	uri := "amqp://test:123456@192.168.30.131:5672/"
	rq, err := mq.NewRabbitMQ(uri)
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
