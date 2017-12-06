package pools

import (
	"github.com/jolestar/go-commons-pool"
	"rmqc/mq"
)

type MQ struct {
	ObjectPool *pool.ObjectPool
}

func NewPoolsMQ() *MQ {
	return &MQ{}
}

// init mq pools
func (mqPool *MQ) Init(maxPools int) {
	p := pool.NewObjectPoolWithDefaultConfig(new(MQFactory))
	p.Config.MaxTotal = maxPools
	mqPool.ObjectPool = p
}

// get a mq obj
func (mqPool *MQ) GetMQ() (*mq.RabbitMQ, error) {
	obj, err := mqPool.ObjectPool.BorrowObject()
	if err != nil {
		return &mq.RabbitMQ{}, err
	}
	return obj.(*mq.RabbitMQ), nil
}

// recover a mq obj
func (mqPool *MQ) Recover(obj interface{}) {
	mqPool.ObjectPool.ReturnObject(obj)
}


type MQFactory struct {

}

// make rabbitMQ object
func (f *MQFactory) MakeObject() (*pool.PooledObject, error) {
	uri := "amqp://test:123456@192.168.30.131:5672/"
	rq, err := mq.NewRabbitMQ(uri)
	if err != nil {
		return &pool.PooledObject{}, err
	}
	return pool.NewPooledObject(rq), nil
}

// do destroy
func (f *MQFactory) DestroyObject(object *pool.PooledObject) error {
	//do destroy
	object.Object.(*mq.RabbitMQ).Close()
	return nil
}

// do validate
func (f *MQFactory) ValidateObject(object *pool.PooledObject) bool {
	//do validate
	return true
}

// do activate
func (f *MQFactory) ActivateObject(object *pool.PooledObject) error {
	//do activate
	return nil
}

// do passivate
func (f *MQFactory) PassivateObject(object *pool.PooledObject) error {
	//do passivate
	return nil
}
