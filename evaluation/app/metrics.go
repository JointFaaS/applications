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

	price = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "manager",
			Subsystem: "requests",
			Name: "price_of_request",
		},
		[]string{"cloud"},
	)
)

func InitMetrics(url string, job string) (*push.Pusher) {
	p := push.New(url, job)
	p.Collector(processed_time)
	p.Collector(price)
	return p
}