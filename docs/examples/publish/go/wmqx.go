package main

import (
	"fmt"
	"net/http"
	"strings"
	"log"
	"io/ioutil"
	"os"
)

func main()  {

	wmqxPublishUri := "http://127.0.0.1:3303/publish/"
	method := "post"

	// message info
	messageName := "ada"
	messageTokenHeader := "WMQX_MESSAGE_TOKEN"
	messageToken := "this is tokenssss"
	messageRouteKeyHeader := "WMQX_MESSAGE_ROUTEKEY"
	routeKey := "test222"

	// header
	headerValues := map[string]string{
		messageTokenHeader: messageToken,
		messageRouteKeyHeader: routeKey,
	}

	url := fmt.Sprintf("%s%s", wmqxPublishUri, messageName)
	data := "name=wmqx&func=publish"

	var req *http.Request
	var err error
	if method == "get" {
		if !strings.Contains(url, "?") {
			url += "?"
		}
		req, err = http.NewRequest("GET", url+data, nil)
	}else {
		req, err = http.NewRequest("POST", url, strings.NewReader(data))
	}
	if err != nil {
		log.Println("error : "+err.Error())
		os.Exit(100)
	}

	// set http header
	if len(headerValues) > 0 {
		for key, value := range headerValues {
			req.Header.Set(key, value)
		}
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("error : "+err.Error())
		os.Exit(100)
	}
	code := resp.StatusCode
	defer resp.Body.Close()
	if code != 200 {
		log.Printf("request error status: %d", code)
	}else {
		bodyByte, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(bodyByte))
	}
}