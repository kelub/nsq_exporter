package nsq_collector

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"nsq_exporter/structs"
	"strconv"
)

type NsqCollector struct {
	nsqlookupdAddr string   //
	nsqdAddr       []string //nsqd http address
	client         *Client

	nsqinfoDesc *prometheus.Desc

	// metrics to describe and collect
	metrics memStatsMetrics
}

type Client struct {
	c *http.Client
}

// memStatsMetrics provide description, value, and value type for memstat metrics.
type memStatsMetrics []struct {
	desc    *prometheus.Desc
	eval    func(*nodestatsResponse) float64
	valType prometheus.ValueType
}

func memstatNamespace(s string) string {
	return fmt.Sprintf("nsq_stats_%s", s)
}

func NewNSQCollector(opts structs.NsqOpts) (*NsqCollector, error) {
	return &NsqCollector{
		nsqlookupdAddr: "http://127.0.0.1:4161",
		nsqdAddr:       []string{"http://127.0.0.1:4151", "http://127.0.0.1:5151"},
		client: &Client{
			&http.Client{},
		},
		nsqinfoDesc: prometheus.NewDesc(
			"nsq_info",
			"nsq version",
			nil, nil,
		),
		metrics: memStatsMetrics{
			{
				desc: prometheus.NewDesc(
					memstatNamespace("memory_heap_objects"),
					"memory heap objects.",
					nil, nil,
				),
				eval:    func(nodeStats *nodestatsResponse) float64 { return nodeStats.Memory.Heap_objects},
				valType: prometheus.GaugeValue,
			},
		},
	}, nil
}

func (c *NsqCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.nsqinfoDesc
}

func (c *NsqCollector) Collect(ch chan<- prometheus.Metric) {
	var nsqlookupdNodes respType
	endpointNodes := fmt.Sprintf("%s/%s", c.nsqlookupdAddr, "nodes?format=json")
	if err := c.client.GETV1(endpointNodes, &nsqlookupdNodes); err != nil {
		logrus.Error("get  ", nsqlookupdNodes, err)
		return
	}
	logrus.Infof("nsqlookupdNodes %s", nsqlookupdNodes)
	var nodeStats nodestatsResponse
	endpointStats := fmt.Sprintf("%s/%s", c.nsqdAddr[0], "stats?format=json")
	if err := c.client.GETV1(endpointStats, &nodeStats); err != nil {
		logrus.Error("get ", c.nsqdAddr[0], err)
		return
	}
	logrus.Infof("nsqd %s", nodeStats)
	heap_objects,_ := strconv.ParseFloat(nodeStats.Memory.Heap_objects, 64)
	ch <- prometheus.MustNewConstMetric(c.nsqinfoDesc, prometheus.GaugeValue, heap_objects)
}
