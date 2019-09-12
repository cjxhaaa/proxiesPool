package crawler

import (
	"ProxyPool/pkg/crawler/ip3366"
	"ProxyPool/pkg/crawler/proxies"
	"ProxyPool/pkg/crawler/qingting"
	"ProxyPool/pkg/settings"
	"github.com/sirupsen/logrus"
	"time"
)



type Pool struct {
	crawlers   []Crawler
	address    []*proxies.Proxy
	set        *Seter
}

func NewPool(ss *settings.Settings) *Pool {
	pool := &Pool{}
	for _, st := range ss.ProxySetting {
		pool.register(&st)
	}
	pool.initSet(ss)
	return pool
}


func (p *Pool) register(st *settings.ProxyParams) {
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

func (p *Pool) initSet(ss *settings.Settings) {
	p.set = NewSet(ss)
	go p.set.Run()
}

func (p *Pool) Start() {

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
				p.address = append(p.address, ps...)
				for _, proxy := range ps {
					 p.set.add <- proxy.Address
				}
			}(crawler)
		}

		time.Sleep(15*time.Second)
	}

}

func (p *Pool) GetOneProxy() string {
	for i := 0 ; i < int(p.set.SortSet.Length()); i++ {
		key,score := p.set.SortSet.GetDataByRank(0,true)

		if score < p.set.PassScore {
			p.set.pop <- key
			continue
		}
		p.set.decline <- key
		return key
	}
	return ""
}

type Crawler interface {
	GetProxy() ([]*proxies.Proxy, error)
}

