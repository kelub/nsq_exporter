package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"nsq_exporter/nsq_collector"
	"nsq_exporter/structs"
	"strings"
)

//func main() {
//	cmd.Execute()
//}

var nsqOpts structs.NsqOpts

func initCommandLine() {
	nsqOpts = structs.NsqOpts{
		ListenAddr:        "",
		NsqlookupHttpAddr: "",
		NsqdHttpAddr:      []string{},
	}
	nsqdHttpAddrs := "http://127.0.0.1:4151"
	kingpin.Flag(
		"listenAddr",
		"address to nsq exporter listen, default :9101",
	).Default(":9101").StringVar(&nsqOpts.ListenAddr)
	kingpin.Flag(
		"nsqlookupHttpAddr",
		"nsqlookup http addr , default http://127.0.0.1:4161",
	).Default("http://127.0.0.1:4161").StringVar(&nsqOpts.NsqlookupHttpAddr)
	kingpin.Flag(
		"nsqdAddrs",
		"nsqdAddrs http addr , default http://127.0.0.1:4151",
	).Default("http://127.0.0.1:4151").StringVar(&nsqdHttpAddrs)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	nsqOpts.NsqdHttpAddr = strings.Split(nsqdHttpAddrs, ",")
}

func main() {
	initCommandLine()
	nsqcollector, err := nsq_collector.NewNSQCollector(nsqOpts)
	if err != nil {
		logrus.WithError(err).Fatalln("创建 NsqExporter 失败")
	}
	prometheus.MustRegister(nsqcollector)

	http.Handle("/metrics", prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Nsq Exporter</title></head>
			<body>
			<h1>Node Exporter</h1>
			<p><a href="` + "/metrics" + `">Metrics</a></p>
			</body>
			</html>`))
	})
	logrus.Infof("Listening on %s", nsqOpts.ListenAddr)
	if err := http.ListenAndServe(nsqOpts.ListenAddr, nil); err != nil {
		logrus.WithError(err).Fatalln("启动失败", err)
	}
}
