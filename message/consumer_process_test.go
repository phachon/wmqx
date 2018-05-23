package message

import (
	"testing"
	"time"
)

func TestNewConsumerProcess(t *testing.T) {

	NewConsumerProcess()
}

func TestConsumerProcess_AddProcess(t *testing.T) {

	cp := NewConsumerProcess()

	err := cp.AddProcess("key1")
	if err != nil {
		t.Error(err.Error())
	}

	err = cp.AddProcess("key2")
	if err != nil {
		t.Error(err.Error())
	}

	if len(cp.ProcessMessages) != 2 {
		t.Error("add error")
	}
}

func TestConsumerProcess_ProcessIsExist(t *testing.T) {

	cp := NewConsumerProcess()

	err := cp.AddProcess("key1")
	if err != nil {
		t.Error(err.Error())
	}

	ok := cp.ProcessIsExist("key1")
	if ok == false {
		t.Error("faild")
	}
}

func TestConsumerProcess_GetProcessMessage(t *testing.T) {
	cp := NewConsumerProcess()

	err := cp.AddProcess("key1")
	if err != nil {
		t.Error(err.Error())
	}

	process, err := cp.GetProcessMessage("key1")
	if err != nil {
		t.Error(err.Error())
	}
	if process.key != "key1" {
		t.Error("faild")
	}
}

func TestConsumerProcess_UpdateProcessByKey(t *testing.T) {
	cp := NewConsumerProcess()

	err := cp.AddProcess("key1")
	if err != nil {
		t.Error(err.Error())
	}

	time.Sleep(2*time.Second)
	updateTime := time.Now().Unix()
	err = cp.UpdateProcessByKey("key1", updateTime)
	if err != nil {
		t.Error(err.Error())
	}
	process, err := cp.GetProcessMessage("key1")
	if err != nil {
		t.Error(err.Error())
	}
	if updateTime != process.lastTime {
		t.Error("faild")
	}
}

func TestConsumerProcess_DeleteProcessByKey(t *testing.T) {
	cp := NewConsumerProcess()

	err := cp.AddProcess("key1")
	if err != nil {
		t.Error(err.Error())
	}

	err = cp.DeleteProcessByKey("key1")
	if err != nil {
		t.Error(err.Error())
	}
	ok := cp.ProcessIsExist("key1")
	if ok != false {
		t.Error("faild")
	}
}

func TestConsumerProcess_StopProcessByKey(t *testing.T) {
	cp := NewConsumerProcess()

	err := cp.AddProcess("key1")
	if err != nil {
		t.Error(err.Error())
	}

	err = cp.StopProcessByKey("key1")
	if err != nil {
		t.Error(err.Error())
	}
}