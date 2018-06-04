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
	Body string `json:"body"`
	Method string `json:"method"`
	Args string `json:"args"`
}

// json encode publish message
func (pmg *PublishMessage) Encode() (string, error) {
	// body base64 encode
	body := base64.StdEncoding.EncodeToString([]byte(pmg.Body))
	pmg.Body = body
	b, err := json.Marshal(pmg)
	return string(b), err
}

// json decode publish message
func (pmg *PublishMessage) Decode(publishMsg string) (*PublishMessage, error) {
	json.Unmarshal([]byte(publishMsg), pmg)
	// body base64 decode
	requestBody, err := base64.StdEncoding.DecodeString(pmg.Body)
	if err != nil {
		return pmg, err
	}
	pmg.Body = string(requestBody)
	return pmg, nil
}

// publish message Body original encode
func (pmg *PublishMessage) EncodeOriginalString() (string) {
	// body base64 decode
	requestBody, err := base64.StdEncoding.DecodeString(pmg.Body)
	if err != nil {
		return ""
	}
	pmg.Body = string(requestBody)
	b, _ := json.Marshal(pmg)
	return string(b)
}

// publish message Body no base64 encode
func (pmg *PublishMessage) OriginalString() (string) {
	b, _ := json.Marshal(pmg)
	return string(b)
}