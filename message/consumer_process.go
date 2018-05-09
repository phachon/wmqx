package message

import (
	"sync"
	"errors"
	"time"
	"wmqx/app"
)

func NewConsumerProcess() *ConsumerProcess {
	return &ConsumerProcess{}
}

// consumer process
type ConsumerProcess struct {
	lock sync.RWMutex
	ProcessMessages []*ConsumerProcessMessage
}

// consumer process message
type ConsumerProcessMessage struct {
	Key string // key
	LastTime int64 //last time
	SignalChan chan string //signal chan
	ExitAck chan bool // consumer exit ack
}

const Consumer_Sign_Stop = "stop"

// consumer process is exist
func (cp *ConsumerProcess) ProcessIsExist(consumerKey string) bool {
	isExist := false
	for _, process := range cp.ProcessMessages {
		cp.lock.RLock()
		if process.Key == consumerKey {
			cp.lock.RUnlock()
			isExist = true
			break
		}
		cp.lock.RUnlock()
	}
	return isExist
}

// get consumer process message by consumer key
func (cp *ConsumerProcess) GetProcessMessage(consumerKey string) (processMsg *ConsumerProcessMessage, err error) {
	ok := cp.ProcessIsExist(consumerKey)
	if ok == false {
		return processMsg, errors.New("consumer process not exists!")
	}

	for _, process := range cp.ProcessMessages {
		cp.lock.RLock()
		if process.Key == consumerKey {
			cp.lock.RUnlock()
			return process, nil
		}
		cp.lock.RUnlock()
	}
	return
}

// add consumer process
func (cp *ConsumerProcess) AddProcess(consumerKey string) error {

	ok := cp.ProcessIsExist(consumerKey)
	if ok == true {
		return errors.New("consumer process is exists!")
	}

	process := &ConsumerProcessMessage{
		Key : consumerKey,
		LastTime: time.Now().Unix(),
		SignalChan: make(chan string, 1),
		ExitAck: make(chan bool, 1),
	}

	cp.lock.Lock()
	cp.ProcessMessages = append(cp.ProcessMessages, process)
	cp.lock.Unlock()
	return nil
}

// update consumer process by consumer key
func (cp *ConsumerProcess) UpdateProcessByKey(consumerKey string, lastTime int64) error {

	ok := cp.ProcessIsExist(consumerKey)
	if ok == false {
		return errors.New("consumer process not exists!")
	}
	for _, process := range cp.ProcessMessages {
		cp.lock.RLock()
		if process.Key == consumerKey {
			cp.lock.RUnlock()
			cp.lock.Lock()
			process.LastTime = lastTime
			cp.lock.Unlock()
			break
		}
		cp.lock.RUnlock()
	}
	return nil
}

// delete consumer process by consumer key
func (cp *ConsumerProcess) DeleteProcessByKey(consumerKey string) error {

	processes := []*ConsumerProcessMessage{}
	for _, process := range cp.ProcessMessages {
		cp.lock.RLock()
		if process.Key != consumerKey {
			cp.lock.RUnlock()
			cp.lock.Lock()
			processes = append(processes, process)
			cp.lock.Unlock()
		}
		cp.lock.RUnlock()
	}
	cp.ProcessMessages = processes
	return nil
}

// stop a consumer process by consumer key
func (cp *ConsumerProcess) StopProcessByKey(consumerKey string) error {

	process, err := cp.GetProcessMessage(consumerKey)
	if err != nil {
		return err
	}

	process.SignalChan<-Consumer_Sign_Stop

	ok := <-process.ExitAck
	if ok == true {
		app.Log.Info("consumer "+consumerKey+" process ack exit!")
		cp.DeleteProcessByKey(consumerKey)
	}
	return nil
}
