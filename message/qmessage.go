package message

import (
	"sync"
	"errors"
)

type Message struct {
	Consumers   []*Consumer `json:"consumers"`
	Durable     bool `json:"durable"`
	IsNeedToken bool `json:"is_need_token"`
	Mode        string `json:"mode"`
	Name        string `json:"name"`
	Token       string `json:"token"`
	Comment     string `json:"comment"`
}

type Consumer struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	RouteKey  string `json:"route_key"`
	Timeout   float64 `json:"timeout"`
	Code      float64 `json:"code"`
	CheckCode bool `json:"check_code"`
	Comment   string `json:"comment"`
}

type QMessageRecordFunc func() QMessageRecord

type QMessage struct {
	Lock *sync.RWMutex
	Messages []*Message
	record QMessageRecord
}

type QMessageRecord interface {
	Init(QMessageRecordConfig) error
	Write([]*Message) error
	Read() ([]*Message, error)
	Clean() error
}

type QMessageRecordConfig interface {
	Name() string
}

var records = make(map[string]QMessageRecordFunc)

// Register QMessage record
func Register(recordType string, record QMessageRecordFunc)  {
	if records[recordType] != nil {
		panic("wmqx: QMessage record type "+ recordType +" already registered!")
	}
	if record == nil {
		panic("wmqx: QMessage record type "+ recordType +" is nil!")
	}

	records[recordType] = record
}

// New QMessage
func NewQMessage(recordTypeName string, config QMessageRecordConfig) (qm *QMessage, err error) {

	recordType, ok := records[recordTypeName]
	if ok == false {
		return qm, errors.New("QMessage record type "+ recordTypeName +" not support!")
	}
	recordFun := recordType()
	err = recordFun.Init(config)
	if err != nil {
		return
	}
	qm = &QMessage{
		Lock: &sync.RWMutex{},
		Messages: []*Message{},
		record: recordFun,
	}
	// load record to Messages
	err = qm.LoadRecord()
	if err != nil {
		return
	}
	return
}

// check message name is exists
func (qm *QMessage) IsExistsMessage(name string) bool {

	isExists := false
	for _, message := range qm.Messages {
		qm.Lock.RLock()
		if message.Name == name {
			isExists = true
			qm.Lock.RUnlock()
			break
		}
		qm.Lock.RUnlock()
	}
	return isExists
}

// add a message
func (qm *QMessage) AddMessage(messageValue *Message) error {

	if qm.IsExistsMessage(messageValue.Name) {
		return errors.New("message is exist! ")
	}
	qm.Lock.Lock()
	qm.Messages = append(qm.Messages, messageValue)
	err := qm.record.Write(qm.Messages)
	qm.Lock.Unlock()

	return err
}

// update a message by name
func (qm *QMessage) UpdateMessageByName(name string, messageValue *Message) error {

	if name != messageValue.Name {
		return errors.New("message name error! ")
	}
	isExist := false

	qm.Lock.Lock()
	defer qm.Lock.Unlock()
	for _, message := range qm.Messages {
		if message.Name == name {
			message.Durable = messageValue.Durable
			message.IsNeedToken = messageValue.IsNeedToken
			message.Mode = messageValue.Mode
			message.Token = messageValue.Token
			message.Comment = messageValue.Comment
			isExist = true
			break
		}
	}
	if isExist == true {
		err := qm.record.Write(qm.Messages)
		return err
	}else {
		return errors.New("message not exist! ")
	}
}

// delete a message by name
func (qm *QMessage) DeleteMessageByName(name string) error {

	messages := []*Message{}
	for _, message := range qm.Messages {
		qm.Lock.RLock()
		if message.Name != name {
			messages = append(messages, message)
		}
		qm.Lock.RUnlock()
	}
	qm.Lock.Lock()
	qm.Messages = messages
	err := qm.record.Write(qm.Messages)
	qm.Lock.Unlock()
	return err
}

// get message by name
func (qm *QMessage) GetMessageByName(name string) (*Message, error) {
	for _, message := range qm.Messages {
		qm.Lock.RLock()
		if message.Name == name {
			qm.Lock.RUnlock()
			return message, nil
		}
		qm.Lock.RUnlock()
	}
	return &Message{}, errors.New("message not exist! ")
}

// get all messages
func (qm *QMessage) GetMessages() (messages []*Message) {
	qm.Lock.RLock()
	messages = qm.Messages
	qm.Lock.RUnlock()
	return messages
}

// delete all message
func (qm *QMessage) ClearMessages() error {
	qm.Lock.Lock()
	qm.Messages = []*Message{}
	err := qm.record.Write(qm.Messages)
	qm.Lock.Unlock()
	return err
}

// check message and consumerId is exist
func (qm *QMessage) IsExistsMessageAndConsumerId(messageName string, consumerId string) bool {

	isExist := false
	isExist = qm.IsExistsMessage(messageName)
	if isExist == false {
		return isExist
	}
	for _, message := range qm.Messages {
		qm.Lock.RLock()
		for _, consumer := range message.Consumers {
			if consumer.ID == consumerId {
				isExist = true
				break
			}
		}
		qm.Lock.RUnlock()
	}
	return isExist
}

// add a consumer for message
func (qm *QMessage) AddConsumer(name string, consumerValue *Consumer) error {
	qm.Lock.Lock()
	defer qm.Lock.Unlock()
	for _, message := range qm.Messages {
		if message.Name == name {
			for _, consumer := range message.Consumers {
				if consumer.ID == consumerValue.ID {
					return errors.New("consumer id is exist! ")
				}
			}
			message.Consumers = append(message.Consumers, consumerValue)
			err := qm.record.Write(qm.Messages)
			return err
		}
	}
	return errors.New("message not exist!")
}

// get consumers by message name
func (qm *QMessage) GetConsumersByMessageName(name string) []*Consumer {
	for _, message := range qm.Messages {
		qm.Lock.RLock()
		if name == message.Name {
			qm.Lock.RUnlock()
			return message.Consumers
		}
		qm.Lock.RUnlock()
	}
	return make([]*Consumer, 0)
}

// get consumer by message name and consumer id
func (qm *QMessage) GetConsumerById(name string, id string) (*Consumer, error) {
	consumers := qm.GetConsumersByMessageName(name)
	if len(consumers) == 0 {
		return &Consumer{}, errors.New("consumer not exist!")
	}
	for _, consumer := range consumers {
		qm.Lock.RLock()
		if consumer.ID == id {
			qm.Lock.RUnlock()
			return consumer, nil
		}
		qm.Lock.RUnlock()
	}
	return &Consumer{}, errors.New("consumer not exist!")
}

// update consumer by message name and consumer id
func (qm *QMessage) UpdateConsumerByName(name string, consumerVal *Consumer) error {
	qm.Lock.Lock()
	defer qm.Lock.Unlock()
	message, err := qm.GetMessageByName(name)
	if err != nil {
		return err
	}
	isExist := false
	for _, consumer := range message.Consumers {
		if consumer.ID == consumerVal.ID {
			consumer.URL = consumerVal.URL
			consumer.RouteKey = consumerVal.RouteKey
			consumer.Timeout = consumerVal.Timeout
			consumer.Code = consumerVal.Code
			consumer.CheckCode = consumerVal.CheckCode
			consumer.Comment = consumerVal.Comment
			isExist = true
			break
		}
	}
	if isExist == true {
		err := qm.record.Write(qm.Messages)
		return err
	}else {
		return errors.New("consumer id not exist!")
	}
}

// delete consumer by message name and consumer id
func (qm *QMessage) DeleteConsumerByNameAndId(name string, consumerId string) error {
	qm.Lock.Lock()
	defer qm.Lock.Unlock()
	message, err := qm.GetMessageByName(name)
	if err != nil {
		return err
	}
	consumers := []*Consumer{}
	for _, consumer := range message.Consumers {
		if consumer.ID != consumerId {
			consumers = append(consumers, consumer)
		}
	}
	message.Consumers = consumers
	err = qm.record.Write(qm.Messages)
	return err
}

// update messages record
func (qm *QMessage) UpdateRecord() (err error) {
	qm.Lock.Lock()
	err = qm.record.Write(qm.Messages)
	qm.Lock.Unlock()
	return
}

// load record messages to QMessage
func (qm *QMessage) LoadRecord() (err error) {
	qm.Lock.Lock()
	messages, err := qm.record.Read()
	if err != nil {
		qm.Lock.Unlock()
		return err
	}
	qm.Messages = messages
	qm.Lock.Unlock()
	return
}