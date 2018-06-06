package main

import (
	"fmt"
	"encoding/json"
	"errors"
)

var (
	messageAddPath = "/message/add"
)

func addMessage(message map[string]string) (err error) {

	requestUrl := fmt.Sprintf("%s%s", managerUri, messageAddPath)

	// set token header
	headerValue := map[string]string{
		tokenHeaderName: token,
	}

	body, code, err := httpPost(requestUrl, message, headerValue)
	if err != nil {
		return
	}
	if len(body) == 0 {
		return errors.New(fmt.Sprintf("request wmqx failed, httpStatus: %d", code))
	}
	v := map[string]interface{}{}
	if json.Unmarshal(body, &v) != nil {
		return
	}
	if v["code"].(float64) == 0 {
		return errors.New(fmt.Sprintf(v["message"].(string)))
	}

	return nil
}

// update message..

// delete message..

// message status..