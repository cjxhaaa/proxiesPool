package main

import (
	"ProxyPool/pkg/crawler"
	"fmt"
	"github.com/cjxhaaa/skiplist"
	"github.com/gin-gonic/gin"
	"time"
)


func check_time(time_stamp interface{}) bool {
	now_stamp := time.Now().Unix()
	if stamp,ok := time_stamp.(int64);ok {
		if now_stamp - stamp > 600 {
			return false
		} else {
			return true
		}
	}
	return false
}


func main() {
	//代理排名
	proxiesRank := skipList.New()

	q_queue := make(chan interface{},300)

	//代理获取
	go func() {
		for true {
			pp := crawler.Proxies{}
			pp.GetProxy()

			fmt.Println("======================")

			qq := pp.Address


			for _,q := range qq{
				q_queue <- q
			}
			time.Sleep(500*time.Second)
			fmt.Println("等待5秒")
		}

	}()

	//代理提取
	go func() {
		for true {
			var qqq interface{}
			qqq = <- q_queue
			//fmt.Println(qqq.Ip)
			fmt.Println(qqq)
			if f, ok := qqq.(crawler.Proxy);ok {
				proxiesRank.Set(f.Address,100)  // 存入有序set中
			}
		}
	}()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		key,score := proxiesRank.GetDataByRank(0,true)
		time_stamp := proxiesRank.GetTimeStamp(key)

		if !check_time(time_stamp) {  //超时删除

			c.JSON(200, gin.H{
				"message":"time out",
			})
		} else {
			if score < 80 {     //分值低删除
				fmt.Println(proxiesRank.Delete(key))
				c.JSON(200, gin.H{
					"message":"no good proxies",
					"score":score,
				})
			} else {
				proxiesRank.Increase(key,-1)
				c.JSON(200, gin.H{
					"message":key,
					"score":score,
				})
			}
		}



	})
	r.Run()

	//not_stop := make(chan bool)
	//<- not_stop
}