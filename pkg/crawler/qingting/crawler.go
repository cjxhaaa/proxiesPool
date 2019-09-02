package qingting

import (
	"ProxyPool/pkg/crawler/proxies"
	"ProxyPool/pkg/settings"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cjxhaaa/myGoTools/requests"
	"log"
	"net/http"
	"time"
)

const Name = "qingting"

type Crawler struct {
	Setting  *settings.ProxyParams
	cli      *http.Client
	proxies.Proxies
}

func NewCrawler(st *settings.ProxyParams) *Crawler {
	return &Crawler{Setting:st}
}


type qingTingProxyJson struct {
	Host      string  `json:"host"`
	Port      string  `json:"port"`
	Country   string  `json:"country_cn"`
	Province  string  `json:"province_cn"`
	City      string  `json:"city_cn"`

}

type qingTingMsgJson struct {
	Code      string  `json:"code"`
	Msg       string  `json:"msg"`
}


func (c *Crawler)GetProxy() ([]*proxies.Proxy, error) {
	start_url := fmt.Sprintf("https://proxy.horocn.com/api/proxies?order_id=%s&num=%d&format=json&line_separator=unix&can_repeat=no",c.Setting.OrderID,c.Setting.Num)

	client := requests.InitClient()
	fmt.Printf("GET: %s\n",start_url)
	options := requests.Options{
		Client:client,
		Method:"GET",
		URL:start_url,
	}
	response, err := requests.Request(options)
	if err != nil {
		log.Fatal(err)
	}

	var qtJson []qingTingProxyJson

	err = json.Unmarshal(response.Bytes,&qtJson)
	if err != nil {
		var msgJson qingTingMsgJson
		if err = json.Unmarshal(response.Bytes, &msgJson); err != nil {
			switch msgJson.Msg {
			case "请求频率过快":
				log.Println("请求频率过快")
				time.Sleep(5*time.Second)
				return nil, nil
			}
			return nil, errors.New(msgJson.Msg)
		}
		return nil, errors.New(string(response.Bytes))
	}


	var p = []*proxies.Proxy{}

	for _, qt := range qtJson {
		ip := qt.Host
		port := qt.Port
		proxy := fmt.Sprintf("%s:%s",ip,port)
		p = append(p, &proxies.Proxy{ip,port,proxy})
	}
	return p, nil
}
