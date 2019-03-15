package nsq_collector

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"runtime"
	"nsq_exporter/structs"
)

type NsqCollector struct {
	nsqlookupdAddr string   //
	nsqdAddr       []string //nsqd http address
	client         *Client

	nsqinfoDesc *prometheus.Desc

	// metrics to describe and collect
	//metrics memStatsMetrics
}

type Client struct {
	c *http.Client
}

// memStatsMetrics provide description, value, and value type for memstat metrics.
type memStatsMetrics []struct {
	desc    *prometheus.Desc
	eval    func(*runtime.MemStats) float64
	valType prometheus.ValueType
}

func NewNSQCollector(opts structs.NsqOpts) (*NsqCollector, error) {
	return &NsqCollector{
		nsqlookupdAddr: "127.0.0.1:4161",
		nsqdAddr:       []string{"127.0.0.1:4151", "127.0.0.1:5151"},
		client: &Client{
			&http.Client{},
		},
		nsqinfoDesc: prometheus.NewDesc(
			"nsq_info",
			"nsq version",
			nil, prometheus.Labels{"version": runtime.Version()},
		),
	}, nil
}

func (c *NsqCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.nsqinfoDesc
}

func (c *NsqCollector) Collect(ch chan<- prometheus.Metric) {
	var nsqlookupdNodes respType
	endpointNodes := fmt.Sprintf("%s/%s", c.nsqdAddr[0], "nodes")
	if err := c.client.GETV1(endpointNodes, &nsqlookupdNodes); err != nil {
		logrus.Error("get ", err)
		return
	}

	var nodeStats nodestatsResponse
	endpointStats := fmt.Sprintf("%s/%s", c.nsqlookupdAddr, "stats")
	if err := c.client.GETV1(endpointStats, &nodeStats); err != nil {
		logrus.Error("get ", err)
		return
	}

	ch <- prometheus.MustNewConstMetric(c.nsqinfoDesc, prometheus.GaugeValue, float64(nodeStats.StatusCode))
}
