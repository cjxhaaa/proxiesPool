package main

import (
	"ProxyPool/pkg/crawler"
	settings2 "ProxyPool/pkg/settings"
	"flag"
	"fmt"
	"github.com/cjxhaaa/skiplist"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"time"
)

type params struct {
	Num  string   `form:"num"`
}


func main() {
	var (
		port string
		username string
		password string
	)

	flag.StringVar(&port, "port", "8080", "listen port")
	flag.StringVar(&username, "username", "fuckpool", "pool username")
	flag.StringVar(&password, "password", "wtfpool", "pool password")
	flag.Parse()

	fmt.Println("Listen on", port)

	ss := settings2.Init("config.ini")
	if ss.Server.Port != "" {
		port = ss.Server.Port
	}

	//代理排名
	proxiesRank := skipList.New()

	q_queue := make(chan interface{},300)

	//代理获取
	go func() {
		for true {
			pp := crawler.Proxies{}
			for _, st := range ss.ProxySetting {
				pp.GetProxy(st)
			}
			fmt.Println("======================")

			qq := pp.Address


			for _,q := range qq{
				q_queue <- q
			}
			time.Sleep(10*time.Second)
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
				proxiesRank.Set(f.Address,ss.InitScore)  // 存入有序set中
			}
		}
	}()


	engine := gin.New()
	engine.Use(gin.Logger(),gin.Recovery())
	engine.Use(AuthMiddle)
	engine.GET("/ping", func(c *gin.Context) {
		var ppp params
		err := c.ShouldBindQuery(&ppp)
		if err == nil {
			var proxies []string
			num,err := strconv.Atoi(ppp.Num)
			if err != nil {
				log.Fatal(err)
			}
			for i:=0 ; i < num; i++ {
				key,_ := checkProxyAvaliable(proxiesRank,ss)
				proxiesRank.Increase(key,ss.ScoreInterval)

				if key == "" || checkExist(proxies,key) {
					continue
				}
				proxies = append(proxies,key)

			}
			if proxies != nil {
				if len(proxies) < num {
					c.JSON(200, gin.H{
						"message":"not enough proxies",
						"proxies":proxies,
					})
				} else {
					c.JSON(200, gin.H{
						"message":"success",
						"proxies":proxies,
					})
				}
			} else {
				c.JSON(200, gin.H{
					"message":"no proxies",
					"proxies":proxies,
				})
			}

		} else {
			key,score := checkProxyAvaliable(proxiesRank,ss)  //代理超时每次主动监测

			if key == "" && score == 0 {  //没取到合适的删除

				c.JSON(200, gin.H{
					"message":"all proxies time out",
				})
			} else {
				proxiesRank.Increase(key,ss.ScoreInterval)
				c.JSON(200, gin.H{
					"message":key,
					"score":score,
				})
			}
		}





	})
	engine.Run("0.0.0.0:"+port)

	//not_stop := make(chan bool)
	//<- not_stop
}

func checkProxyAvaliable(proxies *skipList.SortedSet,setting *settings2.Settings) (string, float64){
	//检测是否可用是被动的，此处可优化
	for i := 0 ; i < int(proxies.Length()); i++ {
		key,score := proxies.GetDataByRank(0,true)
		time_stamp := proxies.GetTimeStamp(key)
		if !check_time(time_stamp,setting.Timeout) || score < setting.PassScore{ //超时或者分值低删除
			proxies.Delete(key)
			continue
		}
		return key, score
	}
	return "", 0
}


func check_time(time_stamp interface{},timeout int) bool {
	now_stamp := time.Now().Unix()
	if stamp,ok := time_stamp.(int64);ok {
		if int(now_stamp - stamp) > timeout {
			return false
		} else {
			return true
		}
	}
	return false
}

func AuthMiddle(c *gin.Context) {
	if ua := c.Request.Header.Get("User-Agent"); ua != "fuckPool" {
		c.String(403, "bye~~")
		c.Abort()
	}
}

func checkExist(ll []string, l string) bool {
	for _,x := range ll {
		if x == l {
			return true
		}
	}
	return false
}