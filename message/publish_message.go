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
func (pmg *PublishMessage) JsonEncode(publishMsg *PublishMessage) (string, error){
	// body base64 encode
	body := base64.StdEncoding.EncodeToString(publishMsg.BodyByte)
	publishMsg.Body = body
	b, err := json.Marshal(publishMsg)
	return string(b), err
}

// json decode publish message
func (pmg *PublishMessage) JsonDecode(publishMsg string) (p PublishMessage, err error) {
	json.Unmarshal([]byte(publishMsg), &p)
	// body base64 decode
	requestBody, err := base64.StdEncoding.DecodeString(p.Body)
	if err != nil {
		return
	}
	p.Body = string(requestBody)
	return p, nil
}