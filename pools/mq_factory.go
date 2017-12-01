package pools

import (
	"github.com/jolestar/go-commons-pool"
	"rmqc/mq"
)

type MQFactory struct {

}

// make rabbitMQ object
func (f *MQFactory) MakeObject() (*pool.PooledObject, error) {
	uri := "amqp://guest:guest@localhost:5672/"
	rq, _ := mq.NewRabbitMQ(uri)
	return pool.NewPooledObject(rq), nil
}

// do destroy
func (f *MQFactory) DestroyObject(object *pool.PooledObject) error {
	//do destroy
	object.Object.(mq.RabbitMQ).Close()
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
