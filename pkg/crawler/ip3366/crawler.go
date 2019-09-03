package ip3366

import (
	"ProxyPool/pkg/crawler/proxies"
	"ProxyPool/pkg/settings"
	"fmt"
	"github.com/cjxhaaa/myGoTools/requests"
	"log"
	"net/http"
)

const Name = "ip3366"

type Crawler struct {
	Setting  *settings.ProxyParams
	cli      *http.Client
	proxies.Proxies
}

func NewCrawler(st *settings.ProxyParams) *Crawler {
	return &Crawler{Setting:st}
}


func (c *Crawler)GetProxy() ([]*proxies.Proxy, error)  {
	start_url := "http://www.ip3366.net/free/?stype=3&page=%d"
	var urls []string
	for i := 0; i < c.Setting.Num; i++ {
		urls = append(urls, fmt.Sprintf(start_url,i))
	}


	var p = []*proxies.Proxy{}

	for _, _url := range urls {
		fmt.Printf("GET: %s\n",_url)
		options := requests.Options{
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
			p = append(p, &proxies.Proxy{ip.String(),port.String(), proxy})
		}
	}
	return p, nil
}
