package app

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)


var (
	processed_time = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "manager",
			Subsystem: "requests",
			Name: "fn_request_duration",
			Help: "Duration of processed requests.",
		},
		[]string{"funcName"},
	)
)

func InitMetrics(url string, job string) (*push.Pusher) {
	p := push.New(url, job)
	p.Collector(processed_time)
	return p
}