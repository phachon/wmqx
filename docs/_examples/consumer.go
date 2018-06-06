package main

import (
	"fmt"
	"encoding/json"
	"errors"
)

var (
	consumerAddPath = "/consumer/add"
)

func addConsumer(consumer map[string]string) (err error) {

	requestUrl := fmt.Sprintf("%s%s", managerUri, consumerAddPath)

	// set token header
	headerValue := map[string]string{
		tokenHeaderName: token,
	}

	body, code, err := httpPost(requestUrl, consumer, headerValue)
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