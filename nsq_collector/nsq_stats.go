package nsq_collector

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"net/http"
)

type Producer struct {
	RemoteAddresses  []string `json:"remote_addresses"`
	RemoteAddress    string   `json:"remote_address"`
	Hostname         string   `json:"hostname"`
	BroadcastAddress string   `json:"broadcast_address"`
	TCPPort          int      `json:"tcp_port"`
	HTTPPort         int      `json:"http_port"`
	Version          string   `json:"version"`
	//VersionObj       semver.Version `json:"-"`
	Topics    []string `json:"topics"`
	OutOfDate bool     `json:"out_of_date"`
}

type respType struct {
	Producers []*Producer `json:"producers"`
}

type nodestatsResponse struct {
	StatusCode int `json:"status_code"`
	StatusText string `json:"status_text"`
	//Data       stats  `json:"data"`
	Memory		memory `json:"memory"`
}

type memory struct {
	Gc_pause_usec_95 string `json:"gc_pause_usec_95"`
	Gc_pause_usec_99 string `json:"gc_pause_usec_99"`
	Gc_pause_usec_100 string `json:"gc_pause_usec_100"`
	Gc_total_runs string `json:"gc_total_runs"`
	Heap_idle_bytes string `json:"heap_idle_bytes"`
	Heap_in_use_bytes string `json:"heap_in_use_bytes"`
	Heap_objects string `json:"heap_objects"`
	Heap_released_bytes string `json:"heap_released_bytes"`
	Hext_gc_bytes string `json:"next_gc_bytes"`
}

type stats struct {
	Version   string   `json:"version"`
	Health    string   `json:"health"`
	StartTime int64    `json:"start_time"`
	Topics    []*topic `json:"topics"`
}

// see https://github.com/nsqio/nsq/blob/master/nsqd/stats.go
type topic struct {
	Name         string     `json:"topic_name"`
	Paused       bool       `json:"paused"`
	Depth        int64      `json:"depth"`
	BackendDepth int64      `json:"backend_depth"`
	MessageCount uint64     `json:"message_count"`
	E2eLatency   e2elatency `json:"e2e_processing_latency"`
	Channels     []*channel `json:"channels"`
}

type channel struct {
	Name          string     `json:"channel_name"`
	Paused        bool       `json:"paused"`
	Depth         int64      `json:"depth"`
	BackendDepth  int64      `json:"backend_depth"`
	MessageCount  uint64     `json:"message_count"`
	InFlightCount int        `json:"in_flight_count"`
	DeferredCount int        `json:"deferred_count"`
	RequeueCount  uint64     `json:"requeue_count"`
	TimeoutCount  uint64     `json:"timeout_count"`
	E2eLatency    e2elatency `json:"e2e_processing_latency"`
	Clients       []*client  `json:"clients"`
}

type client struct {
	ID            string `json:"client_id"`
	Hostname      string `json:"hostname"`
	Version       string `json:"version"`
	RemoteAddress string `json:"remote_address"`
	State         int32  `json:"state"`
	FinishCount   uint64 `json:"finish_count"`
	MessageCount  uint64 `json:"message_count"`
	ReadyCount    int64  `json:"ready_count"`
	InFlightCount int64  `json:"in_flight_count"`
	RequeueCount  uint64 `json:"requeue_count"`
	ConnectTime   int64  `json:"connect_ts"`
	SampleRate    int32  `json:"sample_rate"`
	Deflate       bool   `json:"deflate"`
	Snappy        bool   `json:"snappy"`
	TLS           bool   `json:"tls"`
}

type e2elatency struct {
	Count       int                  `json:"count"`
	Percentiles []map[string]float64 `json:"percentiles"`
}

func (c *Client) GETV1(endpoint string, respv interface{}) error {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		logrus.Error("", err)
		return err
	}
	req.Header.Add("Accept", "application/vnd.nsq; version=1.0")
	resp, err := c.c.Do(req)
	if err != nil {
		logrus.Error("Request Error", err)
		return err
	}
	//body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&respv); err != nil {
		logrus.Error("Unmarshal Error", err)
		return err
	}
	//err = json.Unmarshal(body, &respv)
	//if err != nil {
	//	logrus.Error("Unmarshal Error", err)
	//	return err
	//}
	return nil
}
