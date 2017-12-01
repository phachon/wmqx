package pools

import "github.com/jolestar/go-commons-pool"
type PoolsFactory struct {

}

// init mq pools
func (pools *PoolsFactory) InitMQ() {
	p := pool.NewObjectPoolWithDefaultConfig(new(MQFactory))
	p.Config.MaxTotal = 20
	obj, _ := p.BorrowObject()
	p.ReturnObject(obj)
}

// init mq channel pools
func (pools *PoolsFactory) InitMQChannel() {
	p := pool.NewObjectPoolWithDefaultConfig(new(MQChannelFactory))
	p.Config.MaxTotal = 20
	obj, _ := p.BorrowObject()
	p.ReturnObject(obj)
}
