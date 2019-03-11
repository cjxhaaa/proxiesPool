package crawler

import "ProxyPool/pkg/settings"

type Proxy struct {
	Ip       string
	Port     string
	Address  string
}

type Proxies struct {
	Address  []Proxy
}

func (proxies *Proxies)GetProxy(st settings.ProxyParams) {
	if st.ProxyName == "ip3366_proxy" {
		proxies.GetIp3366Proxy(st)
	} else if st.ProxyName == "qingting_proxy" {
		proxies.GetQingTingProxy(st)
	}
}