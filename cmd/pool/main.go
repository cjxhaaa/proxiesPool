package main

import (
	"ProxyPool/pkg/crawler"
	"ProxyPool/pkg/set"
	"ProxyPool/pkg/settings"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

var (
	parser crawler.Parser
	ss *settings.Settings
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		TimestampFormat: time.RFC3339,
	})
	ss = settings.Init("config.ini")
	for _, st := range ss.ProxySetting {
		fmt.Println(st)
		parser.Register(&st)
	}
}


func main() {
	var (
		port string
		username string
		password string
	)

	flag.StringVar(&port, "port", "8081", "listen port")
	flag.StringVar(&username, "username", "fuckpool", "pool username")
	flag.StringVar(&password, "password", "wtfpool", "pool password")
	flag.Parse()

	fmt.Println("Listen on", port)



	//代理排名
	proxiesRank := set.NewSet(ss)

	q_queue := make(chan interface{},300)
	fmt.Println(parser)
	//代理获取
	go parser.Start(q_queue)

	//代理提取
	go func() {
		for {
			proxiesRank.AddProxies(q_queue)
		}
	}()

	//检测代理时效性
	go func() {
		for {
			proxiesRank.CheckIpsValid()
			logrus.Infof("pool剩余代理数%d", proxiesRank.SortSet.Length())
			time.Sleep(2 * time.Second)
		}
	}()


	engine := gin.New()
	engine.Use(gin.Logger(),gin.Recovery())
	engine.Use(AuthMiddle)
	engine.GET("/ping", func(c *gin.Context) {
		var ppp struct {
			Num  string   `form:"num"`
		}
		err := c.ShouldBindQuery(&ppp)
		if err == nil {
			var ps []string
			num,err := strconv.Atoi(ppp.Num)
			if err != nil {
				panic(err)
			}
			for i:=0 ; i < num; i++ {
				key := proxiesRank.GetOneProxy()
				if key == "" || checkExist(ps,key) {
					continue
				}
				ps = append(ps,key)

			}
			if ps != nil {
				if len(ps) < num {
					c.JSON(200, gin.H{
						"message":"not enough proxies",
						"proxies":ps,
					})
				} else {
					c.JSON(200, gin.H{
						"message":"success",
						"proxies":ps,
					})
				}
			} else {
				c.JSON(200, gin.H{
					"message":"no proxies",
					"proxies":ps,
				})
			}

		} else {
			key := proxiesRank.GetOneProxy()

			if key == "" {  //没取到合适的删除

				c.JSON(200, gin.H{
					"message":"all proxies time out",
				})
			} else {
				c.JSON(200, gin.H{
					"message":"success",
					"proxies": key,
				})
			}
		}





	})
	engine.Run("0.0.0.0:"+port)

	//not_stop := make(chan bool)
	//<- not_stop
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

