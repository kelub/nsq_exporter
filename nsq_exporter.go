package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"runtime"
)

type NsqCollector struct {
	nsqinfoDesc *prometheus.Desc

	// metrics to describe and collect
	//metrics memStatsMetrics
}

// memStatsMetrics provide description, value, and value type for memstat metrics.
type memStatsMetrics []struct {
	desc    *prometheus.Desc
	eval    func(*runtime.MemStats) float64
	valType prometheus.ValueType
}

func NewNsqCollector(opts NsqOpts) (*NsqCollector, error) {
	return &NsqCollector{
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
	ch <- prometheus.MustNewConstMetric(c.nsqinfoDesc, prometheus.GaugeValue, 1)
}
