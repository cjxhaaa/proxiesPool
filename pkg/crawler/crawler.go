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
		cch := make(chan interface{}, 30)
		stop := make(chan bool)
		for _, crawler := range p.crawlers {
			go func(crawler Crawler) {
				defer func() {
					close(cch)
					if r := recover(); r != nil {
						logrus.Errorf("Recovered a panic %s", r)
						stop<-true
					}
				}()
				ps, err := crawler.GetProxy()
				if err != nil {
					panic(err)
				}
				p.Address = append(p.Address, ps...)
				for _, p := range ps {
					cch <- p
				}
			}(crawler)
			go func() {
				var pp interface{}
				ok := true
				for {
					select {
					case pp, ok = <-cch:
						if ok {
							ch <- pp
						} else {
							break
						}
					case <- stop:
						ok = false
						break
					}
					if !ok {
						break
					}
				}
			}()
		}

		time.Sleep(15*time.Second)
	}

}

type Crawler interface {
	GetProxy() ([]*proxies.Proxy, error)
}

