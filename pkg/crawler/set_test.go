package crawler

import (
	"ProxyPool/pkg/settings"
	"fmt"
	"testing"
	"time"
)


func TestSeter_GetOneProxy(t *testing.T) {
	st := &settings.Settings{
		Timeout: 30,
		ScoreInterval: -1,
		PassScore: 90,
		InitScore: 100,
	}
	set := NewSet(st)
	set.Run()

	set.add <- "1.1.1.1"
	set.add <- "2.2.2.2"
	set.add <- "3.3.3.3"
	set.add <- "4.4.4.4"
	set.add <- "5.5.5.5"

	fmt.Println(set.SortSet.Length())

	key := set.GetOneProxy()
	fmt.Println("proxy: ", key)

	score, ok  := set.SortSet.GetScore(key)
	if ok {
		fmt.Printf("score: %f", score)
	} else {
		fmt.Printf("key score %s不存在", key)
	}
	for i := 0; i < 300; i++ {
		fmt.Println(set.GetOneProxy())
	}
	fmt.Println(set.SortSet.Length())


}

func TestSeter_CheckIpsValid(t *testing.T) {
	st := &settings.Settings{
		Timeout: 10,
		ScoreInterval: -1,
		PassScore: 90,
		InitScore: 100,
	}
	set := NewSet(st)
	go func() {
		for {
			set.CheckIpsValid()
			time.Sleep(1*time.Second)
		}

	}()
	set.add <- "1.1.1.1"
	fmt.Println()
	time.Sleep(2*time.Second)
	set.add <- "2.2.2.2"
	time.Sleep(2*time.Second)
	set.add <- "3.3.3.3"
	time.Sleep(2*time.Second)
	set.add <- "4.4.4.4"
	time.Sleep(2*time.Second)
	set.add <- "5.5.5.5"
	fmt.Println(set.SortSet.Length())
	time.Sleep(3*time.Second)
	fmt.Println(set.SortSet.Length())
	time.Sleep(2*time.Second)
	fmt.Println(set.SortSet.Length())
	time.Sleep(2*time.Second)
	fmt.Println(set.SortSet.Length())
	time.Sleep(2*time.Second)
	fmt.Println(set.SortSet.Length())
	time.Sleep(2*time.Second)
	fmt.Println(set.SortSet.Length())
	time.Sleep(3*time.Second)
	fmt.Println("up: ",set.GetOneProxy())

}
