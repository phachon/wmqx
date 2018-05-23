package message

import (
	"testing"
)

func GetQMessage() (qm *QMessage, err error){
	fileConfig := &RecordFileConfig{
		Filename: "../message.json",
		JsonBeautify: true,
	}
	qm, err = NewQMessage("file", NewRecordConfigFile(fileConfig))
	if err != nil {
		return
	}
	return
}

func TestNewQMessage(t *testing.T) {
	_, err := GetQMessage()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestQMessage_AddMessage(t *testing.T) {
	qm, err := GetQMessage()
	if err != nil {
		t.Error(err)
	}

	msg := &Message{
		Consumers: []*Consumer{},
		Durable: true,
		IsNeedToken: true,
		Mode: "topic",
		Name: "ada",
		Token: "this is token",
		Comment: "this is comment",
	}
	err = qm.AddMessage(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestQMessage_GetMessageByName(t *testing.T) {

	qm, err := GetQMessage()
	if err != nil {
		t.Error(err)
	}
	messages, err := qm.GetMessageByName("aa")
	if err != nil {
		t.Error(err)
	}
	t.Log(messages.Comment)
}

func TestQMessage_UpdateRecord(t *testing.T) {
	qm, err := GetQMessage()
	if err != nil {
		t.Error(err)
	}

	msg := &Message{
		Consumers: []*Consumer{},
		Durable: true,
		IsNeedToken: true,
		Mode: "topic",
		Name: "PPPPPADADADASD",
		Token: "this is token",
		Comment: "this is comment",
	}
	err = qm.AddMessage(msg)
	if err != nil {
		t.Error(err)
	}

	qm.UpdateRecord()
}

func TestQMessage_AddConsumer(t *testing.T) {
	qm, err := GetQMessage()
	if err != nil {
		t.Error(err)
	}

	consum := &Consumer{
		ID: "b0a38e68-8b36-4d85-6610-4e10425e4ada8",
		URL: "htp://127.0.0.1:8099/index.php",
		RouteKey: "test routekey",
		Timeout: 5000,
		Code: 500,
		CheckCode: false,
		Comment: "this is a test comment",
	}
	err = qm.AddConsumer("aa", consum)
	if err != nil {
		t.Error(err)
	}
}