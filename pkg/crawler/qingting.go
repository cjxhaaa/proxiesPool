package crawler

import (
	"ProxyPool/pkg/requests"
	"ProxyPool/pkg/settings"
	"encoding/json"
	"fmt"
	"log"
)


type qingTingProxyJson struct {
	Host      string  `json:"host"`
	Port      string  `json:"port"`
	Country   string  `json:"country_cn"`
	Province  string  `json:"province_cn"`
	City      string  `json:"city_cn"`
}


func (proxies *Proxies) GetQingTingProxy(st settings.ProxyParams) {
	start_url := fmt.Sprintf("https://proxy.horocn.com/api/proxies?order_id=%s&num=%d&format=json&line_separator=unix&can_repeat=no",st.OrderID,st.Num)

	client,_ := requests.InitClient()
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
	if err != err {
		return
	}

	for _, qt := range qtJson {
		ip := qt.Host
		port := qt.Port
		proxy := fmt.Sprintf("%s:%s",ip,port)
		proxies.Address = append(proxies.Address, Proxy{ip,port,proxy})
	}
}
