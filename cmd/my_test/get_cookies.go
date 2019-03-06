package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Cookie struct {
	Domain  string    `json:"domain"`
	Name    string    `json:"name"`
	Value   string    `json:"value"`
}

func main() {
	//fmt.Println(os.Args[1])
	//oo_cookies := []byte(os.Args[1])
	oo_cookies := []byte(`[{"domain": ".michaelkors.com", "name": "_abck", "value": "E8792DD1F6B8BD71B7EFDFA249E85F5ACCEDDDED990A00006569625CEC394621~-1~9LgVErnDpgOR3NPnCgIxyYpVOiXMqIKuU2CMqwQuPmE=~-1~-1"}]`)
	//fmt.Println(oo_cookies)
	var cookies []Cookie
	if err := json.Unmarshal(oo_cookies, &cookies); err != nil {
		log.Fatalf("Json unmashaling failed: %s", err)
	}
	fmt.Println(cookies)
	for _,i := range cookies {
		fmt.Println(i.Name)
	}
}
