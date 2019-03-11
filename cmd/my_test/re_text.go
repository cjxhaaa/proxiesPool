package main

import (
	"ProxyPool/pkg/config"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main()  {
	dirPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	//mm := regexp.MustCompile("\\[(?P<header>[^]]+)]")
	//
	//fmt.Println(tools.ReCompileGroup(mm,"[cjxh]"))
	con := config.New(dirPath+"/config.ini")
	fmt.Println(con.GetOptionValue("aa","aa"))

}
