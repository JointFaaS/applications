package app

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/JointFaaS/Client-go/client"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	host1 string = "http://"
	host2 string = "http://"
	pushgateway string = ""
)

type imgBody struct {
	Width int `json:"width"`
	Height int `json:"height"`
	Img     string `json:"img"`
}

func Testing() {
	c1, _ := client.NewClient(client.Config{
		Host: host1,
	})
	c2, _ := client.NewClient(client.Config{
		Host: host2,
	})
	p := InitMetrics(pushgateway, "client")
	f, err := os.Open("test.img")
	if err != nil {
		panic(err)
	}

	fb, _ := ioutil.ReadAll(f)
	body := imgBody{
		Width: 50,
		Height: 50,
		Img: base64.StdEncoding.EncodeToString(fb),
	}
	bodyBytes := bytes.NewBuffer(nil)
	err = json.NewEncoder(bodyBytes).Encode(body)
	if err != nil {
		panic(err)
	}
	r, err := c1.FcInvoke(&client.FcInvokeInput{
		FuncName: "picture",
		Args: bodyBytes.Bytes(),
		EnableNative: "true",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(r.RespBody)
	price1 := 2.0
	price2 := 3.0
	for i := 0; i < 10; i++ {
		go func ()  {
			for j := 0; j < 1000; j++ {
				start := time.Now()
				if price1 < price2 {
					c1.FcInvoke(&client.FcInvokeInput{
						FuncName: "picture",
						Args: bodyBytes.Bytes(),
						EnableNative: "true",
					})
				} else {
					c2.FcInvoke(&client.FcInvokeInput{
						FuncName: "picture",
						Args: bodyBytes.Bytes(),
						EnableNative: "true",
					})
				}
				cost := time.Since(start)
				processed_time.With(prometheus.Labels{"funcName": "picture"}).Observe(cost.Seconds())
			}
		}()
	}

	for j := 0; j < 10; j++ {
		price.With(prometheus.Labels{"cloud": "aliyun"}).Set(price1)
		price.With(prometheus.Labels{"cloud": "aws"}).Set(price2)
		p.Push()
		if j > 4 {
			price2 = 1.0
		}
		time.Sleep(time.Duration(time.Second * 5))
	}
}