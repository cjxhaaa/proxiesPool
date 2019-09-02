package qingting

import (
	"ProxyPool/pkg/settings"
	"fmt"
	"log"
	"testing"
)

func TestCrawler_GetProxy(t *testing.T) {
	var st = &settings.ProxyParams{"qingting","CQLF1627412170278709", 3}
	var crawer = &Crawler{Setting:st}
	pp, err := crawer.GetProxy()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pp[0])
}
