package crawler

import (
	"ProxyPool/pkg/crawler/ip3366"
	"ProxyPool/pkg/crawler/proxies"
	"ProxyPool/pkg/crawler/qingting"
	"ProxyPool/pkg/settings"
	"github.com/sirupsen/logrus"
	"time"
)



type Parser struct {
	crawlers   []Crawler
	proxies.Proxies
}


func (p *Parser) Register(st *settings.ProxyParams) {
	var cr Crawler
	switch st.ProxyName {
	case ip3366.Name:
		cr = ip3366.NewCrawler(st)
	case qingting.Name:
		cr = qingting.NewCrawler(st)
	}
	if cr != nil {
		p.crawlers = append(p.crawlers, cr)
	}
}

func (p *Parser) Start(ch chan<- interface{}) {
	for {
		for _, crawler := range p.crawlers {

			go func(crawler Crawler) {
				defer func() {
					if r := recover(); r != nil {
						logrus.Errorf("Recovered a panic %s", r)
					}
				}()
				ps, err := crawler.GetProxy()
				if err != nil {
					panic(err)
				}
				p.Address = append(p.Address, ps...)
				for _, p := range ps {
					ch <- p
				}
			}(crawler)
		}
		time.Sleep(15*time.Second)
	}

}

type Crawler interface {
	GetProxy() ([]*proxies.Proxy, error)
}

