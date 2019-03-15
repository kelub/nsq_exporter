package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
)

//func main() {
//	cmd.Execute()
//}

type NsqOpts struct {
	listenAddr string

	nsqdHttpAddr string
}

var nsqOpts = NsqOpts{}

func initCommandLine() {
	kingpin.Flag(
		"listenAddr",
		"address to nsq exporter listen, default :9101",
	).Default(":9101").StringVar(&nsqOpts.listenAddr)
	kingpin.Flag(
		"nsqdHttpAddr",
		"nsqd http addr , default :4151",
	).Default(":9101").StringVar(&nsqOpts.nsqdHttpAddr)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	fmt.Println(nsqOpts)
}

func main() {
	initCommandLine()
	nsqcollector, err := NewNsqCollector(nsqOpts)
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
	logrus.Infof("Listening on %s", nsqOpts.listenAddr)
	if err := http.ListenAndServe(nsqOpts.listenAddr, nil); err != nil {
		logrus.WithError(err).Fatalln("启动失败", err)
	}
}
