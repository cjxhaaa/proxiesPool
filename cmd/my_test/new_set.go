package main

import (
	"ProxyPool/pkg/new_set"
	"fmt"
)

func main() {
	mm := new_set.New()
	mm.Set("117.241.98.149:35789",100)
	mm.Set("117.241.98.149:357",100)
	mm.Set("117.241.98.149:3",100)
	fmt.Println(mm.Length())
	fmt.Println(mm.GetScore("117.241.98.149:35789"))
	fmt.Println(mm.GetScore("117.241.98.149:357"))
	//fmt.Println(mm.GetRank("123"))
	//fmt.Println(mm.GetRank("122"))
	//fmt.Println(mm.GetRank("12222.2"))
	//mm.Delete("122")
	fmt.Println(mm.GetRank("117.241.98.149:35789",true))
	fmt.Println(mm.GetRank("117.241.98.149:357",true))
	fmt.Println(mm.GetRank("117.241.98.149:3",true))


	mm.Increase("117.241.98.149:3",-1)
	fmt.Println(mm.GetRank("117.241.98.149:35789",true))
	fmt.Println(mm.GetRank("117.241.98.149:357",true))
	fmt.Println(mm.GetRank("117.241.98.149:3",true))
}
