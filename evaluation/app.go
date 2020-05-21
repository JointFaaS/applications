package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/JointFaaS/Client-go/client"
)

var (
	host1 = ""
	host2 = ""
)

type imgBody struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Img    string `json:"img"`
}

var rootCmd = &cobra.Command{
	Use:   "hcloud-tester",
	Short: "simple tester",
	PreRun: func(cmd *cobra.Command, args []string) {

	},
}

var priceCmd = &cobra.Command{
	Use:   "price",
	Short: "price test",
	Run: func(cmd *cobra.Command, args []string) {
		c1, _ := client.NewClient(client.Config{
			Host: host1,
		})
		c2, _ := client.NewClient(client.Config{
			Host: host2,
		})
		f, err := os.Open("test.img")
		if err != nil {
			panic(err)
		}

		fb, _ := ioutil.ReadAll(f)
		body := imgBody{
			Width:  50,
			Height: 50,
			Img:    base64.StdEncoding.EncodeToString(fb),
		}
		bodyBytes := bytes.NewBuffer(nil)
		err = json.NewEncoder(bodyBytes).Encode(body)
		if err != nil {
			panic(err)
		}
		r, err := c1.FcInvoke(&client.FcInvokeInput{
			FuncName:     "picture",
			Args:         bodyBytes.Bytes(),
			EnableNative: "true",
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(r.RespBody)
		price1 := 2.0
		price2 := 3.0
		timeSlice := make([]time.Duration, 500)
		for i := 0; i < 10; i++ {
			if price1 < price2 {
				price1 = price2 + 2
			} else {
				price2 = price1 + 2
			}
			for j := 0; j < 50; j++ {
				start := time.Now()
				if price1 < price2 {
					c1.FcInvoke(&client.FcInvokeInput{
						FuncName:     "picture",
						Args:         bodyBytes.Bytes(),
						EnableNative: "true",
					})
				} else {
					c2.FcInvoke(&client.FcInvokeInput{
						FuncName:     "picture",
						Args:         bodyBytes.Bytes(),
						EnableNative: "true",
					})
				}
				cost := time.Since(start)
				timeSlice = append(timeSlice, cost)
			}
		}

		for i := 0; i < 500; i++ {
			log.Println(int(timeSlice[i]))
		}
	},
}

var latencyCmd = &cobra.Command{
	Use:   "latency",
	Short: "latency test",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func rootInit() {
	priceCmd.Flags().StringVarP(&host1, "host1", "1", "", "target manager 1")
	priceCmd.Flags().StringVarP(&host2, "host2", "2", "", "target manager 2")
	rootCmd.AddCommand(priceCmd, latencyCmd)
}

func main() {
	rootInit()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
