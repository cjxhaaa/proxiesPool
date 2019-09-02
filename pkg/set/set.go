package set

import (
	"ProxyPool/pkg/crawler/proxies"
	"ProxyPool/pkg/settings"
	"fmt"
	"github.com/cjxhaaa/skiplist"
	"sync"
	"time"
)

type Seter struct {
	InitScore      float64
	PassScore      float64
	ScoreInterval  float64
	Timeout        int
	SortSet        *skipList.SortedSet
	sync.RWMutex
}

func NewSet(st *settings.Settings) *Seter {
	return &Seter{
		InitScore: st.InitScore,
		PassScore: st.PassScore,
		ScoreInterval: st.ScoreInterval,
		Timeout: st.Timeout,
		SortSet: skipList.New(),
	}
}

func (s *Seter) AddProxy(key string) {
	s.Lock()
	defer s.Unlock()
	s.SortSet.Set(key, s.InitScore)
}

func (s *Seter) AddProxies(ch <-chan interface{}) {
	qqq := <- ch
	if f, ok := qqq.(*proxies.Proxy);ok {
		fmt.Println(f.Address)
		s.AddProxy(f.Address) // 存入有序set中
	}
}

func (s *Seter) GetOneProxy() string {


	for i := 0 ; i < int(s.SortSet.Length()); i++ {
		s.RLock()
		key,score := s.SortSet.GetDataByRank(0,true)
		s.RUnlock()

		if score < s.PassScore {
			s.Lock()
			s.SortSet.Delete(key)
			s.Unlock()
			continue
		}

		s.Lock()
		s.SortSet.Increase(key, s.ScoreInterval)
		s.Unlock()

		return key
	}
	return ""
}

func (s *Seter) CheckIpsValid() {
	s.Lock()
	defer s.Unlock()
	for i := int64(0) ; i < int64(s.SortSet.Length()); i++ {
		key,score := s.SortSet.GetDataByRank(i,true)

		time_stamp := s.SortSet.GetTimeStamp(key)

		if !check_time(time_stamp,s.Timeout) || score < s.PassScore {
			s.SortSet.Delete(key)
			continue
		}
	}

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


