package mq

import (
	"runtime"
	"github.com/streadway/amqp"
	"os"
)

// consumer worker

type ConsumerWorker struct {
	ConsumerKey string
	SignChan chan string
	DeliveryChan <-chan amqp.Delivery
}

func NewConsumerWork(consumerKey string, delivery <-chan amqp.Delivery) *ConsumerWorker {
	return &ConsumerWorker{
		ConsumerKey: consumerKey,
		SignChan: make(chan string, 1),
		DeliveryChan: delivery,
	}
}

// star a consumer worker
func (cw *ConsumerWorker) Start() {

	go func() {
		//for {
			select {
			case deliveryChan := <-cw.DeliveryChan:
				//todo handle message
				file, _ := os.OpenFile("consumer.log", os.O_RDWR|os.O_APPEND, 0766)
				file.Write(deliveryChan.Body)
				file.Close()
				deliveryChan.Nack(false, true)
			case str := <-cw.SignChan:
				if str == "stop" {
					// go exit
					runtime.Goexit()
				}
			}
		//}
	}()
}

// stop a consumer worker
func (cw *ConsumerWorker) Stop() {
	cw.SignChan <- "stop"
}

// restart a consumer worker
func (cw *ConsumerWorker) Restart() {

}

// status
func (cw *ConsumerWorker) Status() {

}