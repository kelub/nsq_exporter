package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

var endpoint = "http://127.0.0.1:4161/stats?format=json"

type statsResp struct {
}

var v interface{}

func Test_nsqdHTTP(t *testing.T) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		t.Error("Request error", err)
	}
	req.Header.Add("Accept", "application/vnd.nsq; version=1.0")

	resp, err := client.Do(req)
	if err != nil {
		t.Error("Response error", err)
	}
	body, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	err = json.Unmarshal(body, &v)
	if err != nil {
		t.Error("Unmarshal Response body error", err)

	}
	fmt.Println("",v)
}
