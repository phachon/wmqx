package message

import (
	"errors"
	"encoding/json"
	"bytes"
	"wmqx/utils"
)

const QMessage_Record_Type = "file"

type RecordFile struct {
	config *RecordFileConfig
}

type RecordFileConfig struct {
	Filename string
	JsonBeautify bool
}

func NewRecordFile() QMessageRecord {
	return &RecordFile{}
}

// init file record
func (r *RecordFile) Init(config *RecordConfig) error {
	if config.File.Filename == "" {
		return errors.New("QMessage record file config Filename not is empty!")
	}

	// check file is exists
	ok, _ := utils.NewFile().PathExists(config.File.Filename)
	if ok == false {
		err := utils.NewFile().CreateFile(config.File.Filename)
		if err != nil {
			return err
		}
	}

	r.config = config.File
	return nil
}

// write messages to file
func (r *RecordFile) Write(messages []*Message) error {

	messageByte, err := json.Marshal(messages)
	if r.config.JsonBeautify == true {
		var out bytes.Buffer
		err = json.Indent(&out, messageByte, "", "\t")
		if err != nil {
			return err
		}
		messageByte = out.Bytes()
	}
	err = utils.NewFile().WriteFile(r.config.Filename, string(messageByte))
	return err
}

// read file
func (r *RecordFile) Read() (messages []*Message, err error) {
	data, err := utils.NewFile().ReadAll(r.config.Filename)
	if err != nil {
		return
	}
	json.Unmarshal([]byte(data), &messages)
	return
}

// rewrite file empty
func (r *RecordFile) Clean() error {
	err := utils.NewFile().WriteFile(r.config.Filename, "[]")
	return err
}

func init()  {
	Register(QMessage_Record_Type, NewRecordFile)
}