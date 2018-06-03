package message

import (
	"encoding/json"
	"encoding/base64"
)

func NewPublishMessage() *PublishMessage {
	return &PublishMessage{}
}

type PublishMessage struct {
	Header map[string]string `json:"header"`
	Ip string `json:"ip"`
	BodyByte []byte
	Body string `json:"body"`
	Method string `json:"method"`
	Args string `json:"args"`
}

// json encode publish message
func (pmg *PublishMessage) JsonEncode() (string, error){
	// body base64 encode
	body := base64.StdEncoding.EncodeToString(pmg.BodyByte)
	pmg.Body = body
	b, err := json.Marshal(pmg)
	return string(b), err
}

// json decode publish message
func (pmg *PublishMessage) JsonDecode(publishMsg string) (*PublishMessage, error) {
	json.Unmarshal([]byte(publishMsg), pmg)
	// body base64 decode
	requestBody, err := base64.StdEncoding.DecodeString(pmg.Body)
	if err != nil {
		return pmg, err
	}
	pmg.Body = string(requestBody)
	return pmg, nil
}