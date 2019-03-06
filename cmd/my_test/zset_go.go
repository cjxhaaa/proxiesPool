package main

import (
	"ProxyPool/pkg/mset"
	"fmt"
)

func main() {
	s := mset.New()
	s.Set(100,"5555","test5")
	s.Set(100,"2222", "test1")
	s.Set(100,"3333", "test1")
	fmt.Println(s.Length())
	s.Set(101,"5555", "test1")
	fmt.Println(s.Length())
	key22,score22, extra22 := s.GetDataByRank(0,true)

	fmt.Println(key22)
	fmt.Println(score22)
	fmt.Println(extra22)
	//s.Set(100,"2222","test2")
	//s.Set(100,"3333","test3")

	//rank, score, extra := s.GetRank("1111",true)
	//fmt.Println(rank)
	//fmt.Println(score)
	//fmt.Println(extra)
	//fmt.Println("==================")
	//
	////key,score, extra := s.GetDataByRank(1,false)
	//s.Increase("1111",-50)
	//rank2, score2, extra2 := s.GetRank("1111",true)
	//fmt.Println(rank2)
	//fmt.Println(score2)
	//fmt.Println(extra2)
	//fmt.Println("==================")
	////fmt.Println(s)
	//s.Increase("2222",-10)
	//rank3, score3, extra3 := s.GetRank("1111",true)
	//fmt.Println(rank3)
	//fmt.Println(score3)
	//fmt.Println(extra3)

	//fmt.Println("length:",s.Length())
	//
	//key,score, extra := s.GetDataByRank(0,true)
	//fmt.Println("length:",s.Length())
	//fmt.Println(key)
	//fmt.Println(score)
	//fmt.Println(extra)
	//fmt.Println("==================")
	//s.Increase("1111",1)
	//fmt.Println("==================")
	//fmt.Println("length:",s.Length())
	//key22,score22, extra22 := s.GetDataByRank(0,true)
	//
	//fmt.Println(key22)
	//fmt.Println(score22)
	//fmt.Println(extra22)
	//fmt.Println("==================")
	////s.Increase("1111",1)
	//s.Set(101,"1111", "test1")
	//fmt.Println("==================")
	//key22,score22, extra22 = s.GetDataByRank(0,true)
	//
	//fmt.Println(key22)
	//fmt.Println(score22)
	//fmt.Println(extra22)

	//fmt.Println("length:",s.Length())
	//rank3, score3, extra3 := s.GetRank(key22,true)
	//fmt.Println(rank3)
	//fmt.Println(score3)
	//fmt.Println(extra3)

}
