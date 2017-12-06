package mq

import "encoding/json"

type PublishMessage struct {
	Header map[string]string `json:"header"`
	Ip string `json:"ip"`
	Body string `json:"body"`
	Method string `json:"method"`
	Args string `json:"args"`
}

// json encode publish message
func (pmg *PublishMessage) JsonEncode(publishMsg *PublishMessage) (string, error){
	b, err := json.Marshal(publishMsg)
	return string(b), err
}

// json decode publish message
func (pmg *PublishMessage) JsonDecode(publishMsg string) *PublishMessage {
	var p *PublishMessage
	json.Unmarshal([]byte(publishMsg), p)
	return p
}