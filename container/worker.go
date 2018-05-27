package container

import (
	"wmqx/app"
	"wmqx/message"
	"time"
	"strconv"
)

var Worker = NewWorker()

const Consumer_Action_Insert = "insert"
const Consumer_Action_Update = "update"
const Consumer_Action_Delete = "delete"
const Consumer_Action_Status = "status"

func NewWorker() *worker {
	return &worker{
		ConsumerWorkChan: make(chan ConsumerWorker, 100),
	}
}

type worker struct {
	ConsumerWorkChan chan ConsumerWorker
}

type ConsumerWorker struct {
	Action string
	ConsumerKey string
}

// send consumer sign
func (w *worker) SendConsumerSign(action string, consumerKey string) {
	w.ConsumerWorkChan <- ConsumerWorker{
		Action: action,
		ConsumerKey: consumerKey,
	}
}

// consumer main process worker
func (w *worker) Consumer() {
	go func() {
		defer func() {
			e := recover()
			if e != nil {
				app.Log.Errorf("consumer worker crash, %v", e)
			}
		}()

		for {
			select {
			case consumerWorker := <-w.ConsumerWorkChan:
				app.Log.Info("Consumer worker receive "+consumerWorker.Action+" action, consumerKey: "+consumerWorker.ConsumerKey)
				// insert consumer
				if consumerWorker.Action == Consumer_Action_Insert {
					err := Ctx.ConsumerProcess.AddProcess(consumerWorker.ConsumerKey)
					if err != nil {
						app.Log.Error(err.Error())
						continue
					}
					cp, err := Ctx.ConsumerProcess.GetProcessMessage(consumerWorker.ConsumerKey)
					if err != nil {
						app.Log.Error(err.Error())
						continue
					}
					w.startConsumerProcess(cp)
				}
				// update consumer
				if consumerWorker.Action == Consumer_Action_Update {
					err := Ctx.ConsumerProcess.StopProcessByKey(consumerWorker.ConsumerKey)
					if err != nil {
						app.Log.Error(err.Error())
						continue
					}
					w.SendConsumerSign(Consumer_Action_Insert, consumerWorker.ConsumerKey)
				}
				// delete consumer
				if consumerWorker.Action == Consumer_Action_Delete {
					err := Ctx.ConsumerProcess.StopProcessByKey(consumerWorker.ConsumerKey)
					if err != nil {
						app.Log.Error(err.Error())
					}
				}
				// get consumer status
				if consumerWorker.Action == Consumer_Action_Status {

				}
			}
		}
	}()
}

// start a consumer process
func (w *worker) startConsumerProcess(processMessage *message.ConsumerProcessMessage) {

	go func(processMessage *message.ConsumerProcessMessage) {
		rabbitMq, _ := Ctx.RabbitMQPools.GetMQ()
		channel, _ := rabbitMq.Conn.Channel()

		defer func() {
			channel.Close()
			Ctx.RabbitMQPools.Recover(rabbitMq)
			// ack consumer process exit
			processMessage.ExitAck<-true
			e := recover()
			if e != nil {
				app.Log.Errorf("Consumer %s consum process crash, %v", processMessage.Key, e)
				// retry consumer
				w.SendConsumerSign(Consumer_Action_Insert, processMessage.Key)
				app.Log.Infof("Consumer %s retry start", processMessage.Key)
			}
		}()

		autoAck := false
		exclusive := false
		noLocal := false
		noWait := false
		delivery, _ := channel.Consume(processMessage.Key, "", autoAck, exclusive, noLocal, noWait, nil)

		app.Log.Info("Consumer "+processMessage.Key+" process start, wait message...")
		// update last_time
		Ctx.ConsumerProcess.UpdateProcessByKey(processMessage.Key, time.Now().Unix())
		failRetryTime := app.Conf.GetInt("consume.failRetryTime")
		for {
			select {
			case d := <-delivery:
				// update last_time
				Ctx.ConsumerProcess.UpdateProcessByKey(processMessage.Key, time.Now().Unix())
				// decode publish message
				publishMessage, err := message.NewPublishMessage().JsonDecode(string(d.Body))
				if err != nil {
					app.Log.Error("Consumer "+processMessage.Key+" decode publish message error: "+err.Error())
					time.Sleep(time.Duration(failRetryTime) * time.Second)
					d.Nack(false, true)
					continue
				}
				app.Log.Info("Consumer "+processMessage.Key+" receive message body: "+string(publishMessage.Body))
				// request consumer url
				resBody, code, err := Ctx.RequestConsumerUrl(processMessage.Key, publishMessage)
				if err != nil {
					app.Log.Error("Consumer "+processMessage.Key+" request url failed: "+err.Error())
					time.Sleep(time.Duration(failRetryTime) * time.Second)
					d.Nack(false, true)
					continue
				}
				app.Log.Info("Consumer "+processMessage.Key+" request url success, response code: "+strconv.Itoa(code)+", body: "+resBody)
				d.Ack(false)
			case sign := <-processMessage.SignalChan:
				app.Log.Info("Consumer "+processMessage.Key+" receive stop sign")
				if sign == message.Consumer_Sign_Stop {
					app.Log.Info("Consumer "+processMessage.Key+" process exit!")
					return
				}
			}
		}
	}(processMessage)
}