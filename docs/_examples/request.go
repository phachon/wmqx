package main

import (
	"strings"
	"net/url"
	"net/http"
	"io/ioutil"
)

func httpGet(queryUrl string, queryValues map[string]string, headerValues map[string]string) (body []byte, code int, err error) {

	if !strings.Contains(queryUrl, "?") {
		queryUrl += "?"
	}
	queryString := ""
	for queryKey, queryValue := range queryValues {
		queryString = queryString + "&" + queryKey + "=" + url.QueryEscape(queryValue)
	}
	queryString = strings.Replace(queryString, "&", "", 1)
	queryUrl += queryString

	req, err := http.NewRequest("GET", queryUrl, nil)
	if err != nil {
		return
	}
	if (headerValues != nil) && (len(headerValues) > 0) {
		for key, value := range headerValues {
			req.Header.Set(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	code = resp.StatusCode
	defer resp.Body.Close()

	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return bodyByte, code, nil
}

// http post request
func httpPost(queryUrl string, queryValues map[string]string, headerValues map[string]string) (body []byte, code int, err error) {
	if !strings.Contains(queryUrl, "?") {
		queryUrl += "?"
	}
	queryString := ""
	for queryKey, queryValue := range queryValues {
		queryString = queryString + "&" + queryKey + "=" + url.QueryEscape(queryValue)
	}
	queryString = strings.Replace(queryString, "&", "", 1)

	req, err := http.NewRequest("POST", queryUrl, strings.NewReader(queryString))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if (headerValues != nil) && (len(headerValues) > 0) {
		for key, value := range headerValues {
			req.Header.Set(key, value)
		}
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	code = resp.StatusCode
	defer resp.Body.Close()

	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return bodyByte, code, nil
}
