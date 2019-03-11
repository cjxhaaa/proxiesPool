package settings

import (
	"ProxyPool/pkg/config"
	"ProxyPool/pkg/secret"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type Settings struct {
	Server          serverInfo
	Timeout         int
	ScoreInterval   float64
	PassScore       float64
	InitScore       float64
	ProxySetting    []ProxyParams
}

type ProxyParams struct {
	ProxyName  string
	OrderID    string
	Num        int
}

type serverInfo struct {
	Port      string
	Username  string
	Password  string
}


func Init(filename string) *Settings {
	var server serverInfo
	var pp []ProxyParams

	execPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	settings := config.New(filepath.Join(execPath,filename))

	if settings.HasSection("server") {
		port := settings.GetOptionValue("server","port")
		if port == "" {
			port = "8080"
		}
		server.Port = port

		username := settings.GetOptionValue("server","username")
		if username == "" {
			username = "fuckpool"
		}
		server.Username = username

		password := settings.GetOptionValue("server","password")
		if password == "" {
			password = "wtfpool"
		}
		server.Password = password
	}



	timeout,err := strconv.Atoi(settings.GetOptionValue("settings","timeout"))
	if err != nil {
		timeout = 60
	}

	score_interval, err := strconv.ParseFloat(settings.GetOptionValue("settings","score_interval"), 64)
	if err != nil {
		score_interval = -1
	}

	pass_score, err := strconv.ParseFloat(settings.GetOptionValue("settings","pass_score"), 64)
	if err != nil {
		pass_score = 80
	}

	init_score, err := strconv.ParseFloat(settings.GetOptionValue("settings","init_score"), 64)
	if err != nil {
		init_score = 100
	}

	if settings.HasSection("qingting_proxy") &&  settings.GetOptionBool("qingting_proxy","enable"){
		orderId := secret.AesDecrypt(settings.GetOptionValue("qingting_proxy","order_id"),"1234567812345678")
		if orderId == "" {
			log.Fatal("please check you order")
		}

		num, err := strconv.Atoi(settings.GetOptionValue("qingting_proxy","num"))
		if err != nil {
			num = 10
		}
		pp = append(pp,ProxyParams{"qingting_proxy",orderId,num})

	}

	if settings.HasSection("ip3366_proxy") &&  settings.GetOptionBool("ip3366_proxy","enable") {
		num, err := strconv.Atoi(settings.GetOptionValue("qingting_proxy","num"))
		if err != nil {
			num = 10
		}
		pp = append(pp,ProxyParams{"ip3366_proxy","",num})
	}

	return &Settings{server,timeout,score_interval,pass_score,init_score,pp}
}
