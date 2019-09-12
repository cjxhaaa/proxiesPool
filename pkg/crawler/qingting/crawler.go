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
const tooManyRequest = "请求频率过快"

type Crawler struct {
	Setting  *settings.ProxyParams
	cli      *http.Client
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

	fmt.Printf("GET: %s\n",start_url)
	options := requests.Options{
		Method:"GET",
		URL:start_url,
		Retry: 1,
	}
	response, err := requests.Request(options)
	if err != nil {
		return nil, err
	}

	var qtJson []qingTingProxyJson

	err = json.Unmarshal(response.Bytes,&qtJson)
	if err != nil {
		var msgJson qingTingMsgJson
		if err = json.Unmarshal(response.Bytes, &msgJson); err != nil {
			switch msgJson.Msg {
			case tooManyRequest:
				log.Println(tooManyRequest)
				time.Sleep(5*time.Second)
				return nil, errors.New(tooManyRequest)
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
