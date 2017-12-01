package pools

import "github.com/jolestar/go-commons-pool"

type MQChannelFactory struct {

}

func (f *MQChannelFactory) MakeObject() (*pool.PooledObject, error) {
	return pool.NewPooledObject(&MQChannelFactory{}), nil
}

func (f *MQChannelFactory) DestroyObject(object *pool.PooledObject) error {
	//do destroy
	return nil
}

func (f *MQChannelFactory) ValidateObject(object *pool.PooledObject) bool {
	//do validate
	return true
}

func (f *MQChannelFactory) ActivateObject(object *pool.PooledObject) error {
	//do activate
	return nil
}

func (f *MQChannelFactory) PassivateObject(object *pool.PooledObject) error {
	//do passivate
	return nil
}
