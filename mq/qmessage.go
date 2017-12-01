package mq

import (
	"sync"
	"errors"
	"encoding/json"
	"io"
	"bytes"
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

type QMessage struct {
	Lock *sync.Mutex
	Messages []*Message
	MessageChan chan int
}

func NewQMessage() *QMessage {
	return &QMessage{
		Lock: &sync.Mutex{},
		Messages: []*Message{},
		MessageChan: make(chan int, 1),
	}
}

// add a message
func (qm *QMessage) AddMessage(messageValue *Message) error {
	qm.Lock.Lock()
	defer qm.Lock.Unlock()
	for _, message := range qm.Messages {
		if message.Name == messageValue.Name {
			return errors.New("message is exist!")
		}
	}
	qm.Messages = append(qm.Messages, messageValue)
	return nil
}

// update a message by name
func (qm *QMessage) UpdateMessageByName(name string, messageValue *Message) error {
	qm.Lock.Lock()
	defer qm.Lock.Unlock()
	if name != messageValue.Name {
		return errors.New("message name error!")
	}
	for _, message := range qm.Messages {
		if message.Name == name {
			message.Durable = messageValue.Durable
			message.IsNeedToken = messageValue.IsNeedToken
			message.Mode = messageValue.Mode
			message.Token = messageValue.Token
			message.Comment = messageValue.Comment
			return nil
		}
	}
	return errors.New("message not exist!")
}

// delete a message by name
func (qm *QMessage) DeleteMessageByName(name string) error {
	qm.Lock.Lock()
	defer qm.Lock.Unlock()
	messages := []*Message{}
	for _, message := range qm.Messages {
		if message.Name != name {
			messages = append(messages, message)
		}
	}
	qm.Messages = messages
	return nil
}

// get message by name
func (qm *QMessage) GetMessageByName(name string) *Message {
	for _, message := range qm.Messages {
		if message.Name == name {
			return message
		}
	}
	return &Message{}
}

// get all messages
func (qm *QMessage) GetMessages() []*Message {
	return qm.Messages
}

// delete all message
func (qm *QMessage) ClearMessages() {
	qm.Lock.Lock()
	defer qm.Lock.Unlock()
	qm.Messages = []*Message{}
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
				if consumer.URL == consumerValue.URL {
					return errors.New("consumer url is exist! ")
				}
			}
			message.Consumers = append(message.Consumers, consumerValue)
			return nil
		}
	}
	return errors.New("message not exist!")
}

// get consumers by message name
func (qm *QMessage) GetConsumersByMessageName(name string) []*Consumer {
	for _, message := range qm.Messages {
		if name == message.Name {
			return message.Consumers
		}
	}
	return []*Consumer{}
}

// get consumers by message and consumer id
func (qm *QMessage) GetConsumerById(name string, id string) *Consumer {
	consumers := qm.GetConsumersByMessageName(name)
	if len(consumers) == 0 {
		return &Consumer{}
	}
	for _, consumer := range consumers {
		if consumer.ID == id {
			return consumer
		}
	}
	return &Consumer{}
}

// write to ...
func (qm *QMessage) WriteTo(write io.Writer, isFormat bool) error {
	qm.Lock.Lock()
	defer qm.Lock.Unlock()
	messages := qm.Messages
	messageByte, err := json.Marshal(messages)
	if isFormat {
		var out bytes.Buffer
		err = json.Indent(&out, messageByte, "", "\t")
		if err != nil {
			return err
		}
		out.WriteTo(write)
		return nil
	}
	_, err = write.Write(messageByte)
	if err != nil {
		return err
	}
	return nil
}

// read from ...
func (qm *QMessage) ReadFrom(reader io.Reader) error {
	qm.Lock.Lock()
	defer qm.Lock.Unlock()

	buf := make([]byte, 2048)

	n, err := reader.Read(buf)
	if err != nil {
		return err
	}
	data := buf[:n]
	var messages []*Message
	json.Unmarshal(data, &messages)
	qm.Messages = messages
	return nil
}