package app

import (
	"time"

	"github.com/JointFaaS/Client-go/client"
	"github.com/prometheus/client_golang/prometheus"
)

func Testing() {
	c1, _ := client.NewClient(client.Config{
		Host: "http://",
	})
	c2, _ := client.NewClient(client.Config{
		Host: "http://",
	})
	p := InitMetrics("xx", "client")
	for i := 0; i < 10; i++ {
		go func ()  {
			for j := 0; j < 1000; j++ {
				limit := (i + 1) * 100
				start := time.Now()
				if j < limit {
					c1.FcInvoke(&client.FcInvokeInput{
						FuncName: "hello",
						Args: make([]byte, 0),
						EnableNative: "True",
					})
				} else {
					c2.FcInvoke(&client.FcInvokeInput{
						FuncName: "hello",
						Args: make([]byte, 0),
						EnableNative: "True",
					})
				}
				cost := time.Since(start)
				processed_time.With(prometheus.Labels{"funcName": "hello"}).Observe(cost.Seconds())
			}
		}()
	}

	for {
		p.Push()
		time.Sleep(time.Duration(time.Second * 5))
	}
}