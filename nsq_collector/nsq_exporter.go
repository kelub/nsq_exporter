package nsq_collector

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"nsq_exporter/structs"
)

const (
	namespace = "nsq"
)

type NsqCollector struct {
	nsqlookupdAddr string   //
	nsqdAddr       []string //nsqd http address
	client         *Client

	*nsqDesc
}

type Client struct {
	c *http.Client
}

type nsqDesc struct {
	memoryDesc   *prometheus.Desc
	topicsDesc   *prometheus.Desc
	channelsDesc *prometheus.Desc
}

func memstatNamespace(s string) string {
	return fmt.Sprintf("nsq_stats_%s", s)
}

func NewNSQCollector(opts structs.NsqOpts) (*NsqCollector, error) {
	nsqDesc := &nsqDesc{
		memoryDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "memory"),
			"memory",
			[]string{"node", "memory"}, nil,
		),
		topicsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "topics"),
			"topic",
			[]string{"node", "topic_name", "stat"}, nil,
		),
		channelsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "channel"),
			"topic",
			[]string{"node", "topic_name", "channel_name", "stat"}, nil,
		),
	}
	return &NsqCollector{
		nsqlookupdAddr: opts.NsqlookupHttpAddr,
		nsqdAddr:       opts.NsqdHttpAddr,
		client: &Client{
			&http.Client{},
		},
		nsqDesc: nsqDesc,
	}, nil
}

func (c *NsqCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.memoryDesc
	ch <- c.topicsDesc
	ch <- c.channelsDesc
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

	for i := 0; i < len(c.nsqdAddr); i++ {
		nsqdAddr := c.nsqdAddr[i]
		endpointStats := fmt.Sprintf("%s/%s", nsqdAddr, "stats?format=json")
		if err := c.client.GETV1(endpointStats, &nodeStats); err != nil {
			logrus.Error("get ", nsqdAddr, err)
			return
		}
		logrus.Infof("nsqd %s", nodeStats)
		ch <- prometheus.MustNewConstMetric(c.memoryDesc, prometheus.GaugeValue, float64(nodeStats.Memory.Gc_pause_usec_95), nsqdAddr, "gc_pause_usec_95")
		ch <- prometheus.MustNewConstMetric(c.memoryDesc, prometheus.GaugeValue, float64(nodeStats.Memory.Gc_pause_usec_99), nsqdAddr, "gc_pause_usec_99")
		ch <- prometheus.MustNewConstMetric(c.memoryDesc, prometheus.GaugeValue, float64(nodeStats.Memory.Gc_pause_usec_100), nsqdAddr, "gc_pause_usec_100")
		ch <- prometheus.MustNewConstMetric(c.memoryDesc, prometheus.GaugeValue, float64(nodeStats.Memory.Gc_total_runs), nsqdAddr, "gc_total_runs")
		ch <- prometheus.MustNewConstMetric(c.memoryDesc, prometheus.GaugeValue, float64(nodeStats.Memory.Heap_idle_bytes), nsqdAddr, "heap_idle_bytes")
		ch <- prometheus.MustNewConstMetric(c.memoryDesc, prometheus.GaugeValue, float64(nodeStats.Memory.Heap_in_use_bytes), nsqdAddr, "heap_in_use_bytes")
		ch <- prometheus.MustNewConstMetric(c.memoryDesc, prometheus.GaugeValue, float64(nodeStats.Memory.Heap_objects), nsqdAddr, "heap_objects")
		ch <- prometheus.MustNewConstMetric(c.memoryDesc, prometheus.GaugeValue, float64(nodeStats.Memory.Heap_released_bytes), nsqdAddr, "heap_released_bytes")
		ch <- prometheus.MustNewConstMetric(c.memoryDesc, prometheus.GaugeValue, float64(nodeStats.Memory.Next_gc_bytes), nsqdAddr, "next_gc_bytes")

		for j := 0; j < len(nodeStats.Topics); j++ {
			nodeTopics := nodeStats.Topics[j]
			ch <- prometheus.MustNewConstMetric(c.topicsDesc, prometheus.GaugeValue, float64(nodeTopics.MessageCount), nsqdAddr, nodeTopics.Name, "MessageCount")
			ch <- prometheus.MustNewConstMetric(c.topicsDesc, prometheus.GaugeValue, float64(nodeTopics.BackendDepth), nsqdAddr, nodeTopics.Name, "BackendDepth")
			ch <- prometheus.MustNewConstMetric(c.topicsDesc, prometheus.GaugeValue, float64(nodeTopics.Depth), nsqdAddr, nodeTopics.Name, "Depth")
			for k := 0; k < len(nodeTopics.Channels); k++ {
				nodeTopicChannel := nodeTopics.Channels[k]
				ch <- prometheus.MustNewConstMetric(c.channelsDesc, prometheus.GaugeValue, float64(nodeTopicChannel.Depth), nsqdAddr, nodeTopics.Name, nodeTopicChannel.Name, "Depth")
				ch <- prometheus.MustNewConstMetric(c.channelsDesc, prometheus.GaugeValue, float64(nodeTopicChannel.BackendDepth), nsqdAddr, nodeTopics.Name, nodeTopicChannel.Name, "BackendDepth")

				ch <- prometheus.MustNewConstMetric(c.channelsDesc, prometheus.GaugeValue, float64(nodeTopicChannel.DeferredCount), nsqdAddr, nodeTopics.Name, nodeTopicChannel.Name, "DeferredCount")
				ch <- prometheus.MustNewConstMetric(c.channelsDesc, prometheus.GaugeValue, float64(nodeTopicChannel.RequeueCount), nsqdAddr, nodeTopics.Name, nodeTopicChannel.Name, "RequeueCount")
				ch <- prometheus.MustNewConstMetric(c.channelsDesc, prometheus.GaugeValue, float64(nodeTopicChannel.InFlightCount), nsqdAddr, nodeTopics.Name, nodeTopicChannel.Name, "in_flight_count")
				ch <- prometheus.MustNewConstMetric(c.channelsDesc, prometheus.GaugeValue, float64(nodeTopicChannel.MessageCount), nsqdAddr, nodeTopics.Name, nodeTopicChannel.Name, "message_count")
				ch <- prometheus.MustNewConstMetric(c.channelsDesc, prometheus.GaugeValue, float64(nodeTopicChannel.TimeoutCount), nsqdAddr, nodeTopics.Name, nodeTopicChannel.Name, "timeout_count")
			}
		}
	}

}
