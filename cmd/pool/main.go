package main

import (
	"ProxyPool/pkg/crawler"
	"ProxyPool/pkg/settings"
	"ProxyPool/rpcserve"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		TimestampFormat: time.RFC3339,
	})
}


func main() {
	var (
		port string
		username string
		password string
		ss = settings.Init("config.ini")
		pool = crawler.NewPool(ss)
	)

	flag.StringVar(&port, "port", "8081", "listen port")
	flag.StringVar(&username, "username", "fuckpool", "pool username")
	flag.StringVar(&password, "password", "wtfpool", "pool password")
	flag.Parse()

	fmt.Println("Listen on", port)
	//启动爬虫
	go pool.Start()

	grpcServer := &rpcserve.Service{Port:"8082", Pool:pool}
	go grpcServer.RunGrpcServer()

	// http api
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
				key := pool.GetOneProxy()
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
			key := pool.GetOneProxy()

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

