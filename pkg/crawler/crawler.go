package crawler

import (
	"ProxyPool/pkg/requests"
	"fmt"
	"log"
)

type Proxy struct {
	Ip       string
	Port     string
	Address  string
}

type Proxies struct {
	Address  []Proxy
}


func (proxies *Proxies) GetProxy() {
	start_url := "http://www.ip3366.net/free/?stype=3&page=%d"
	var urls []string
	for i := 0; i < 4; i++ {
		urls = append(urls, fmt.Sprintf(start_url,i))
	}

	client,_ := requests.InitClient()
	//uu,_ := url.Parse("http://www.66ip.cn/0.html")

	//var cookies []*http.Cookie
	//
	//cookies = append(cookies,&http.Cookie{Name:"yd_cookie",Value:"a11dfc0b-b13a-4d10a45ee2d8b235160063231fc1ce86d2a8"})
	//cookies = append(cookies,&http.Cookie{Name:"_ydclearance",Value:"3add658201d8db08ce273b98-395d-40bd-81a3-45ab5a276907-1548836277"})
	//
	//client.Jar.SetCookies(uu,cookies)
	for _, _url := range urls {
		fmt.Printf("GET: %s\n",_url)
		options := requests.Options{
			Client:client,
			Method:"GET",
			URL:_url,
		}
		response, err := requests.Request(options)
		if err != nil {
			log.Fatal(err)
		}

		nodes, err := response.Selector.GetNodes(`//div[@id="list"]/table/tbody/tr`)

		for _, node := range nodes {
			ip,_ := node.Get(`./td[1]`)
			port,_ := node.Get(`./td[2]`)
			proxy := fmt.Sprintf("%s:%s",ip,port)
			fmt.Println(proxy)
			proxies.Address = append(proxies.Address, Proxy{ip.String(),port.String(), proxy})
		}

	}
}