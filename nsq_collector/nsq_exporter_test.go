package nsq_collector

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

type statsResp struct {
}

var nsqlookupdaddr = "http://127.0.0.1:4161"

func Test_nsqdHTTP(t *testing.T) {
	var v interface{}
	endpoint := ""
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
	fmt.Println("", v)
}

func Test_GETV1(t *testing.T) {
	endpoint := fmt.Sprintf("%s/nodes", nsqlookupdaddr)
	var resp respType
	fmt.Println("endpoint:", endpoint)
	c := &Client{
		c: &http.Client{},
	}
	err := c.GETV1(endpoint, &resp)
	if err != nil {
		t.Error("GETV1 error", err)
	}
	fmt.Print("DATA", resp)
	for _, producers := range resp.Producers {
		fmt.Println("producers", producers)
	}
}
