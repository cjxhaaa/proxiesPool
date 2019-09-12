package crawler

import (
	"ProxyPool/pkg/settings"
	"github.com/cjxhaaa/myGoTools/skipList"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Seter struct {
	InitScore      float64
	PassScore      float64
	ScoreInterval  float64
	Timeout        int
	SortSet        *skipList.SortedSet
	sync.Mutex
	add            chan string
	pop            chan string
	decline        chan string
}

func NewSet(st *settings.Settings) *Seter {
	return &Seter{
		InitScore: st.InitScore,
		PassScore: st.PassScore,
		ScoreInterval: st.ScoreInterval,
		Timeout: st.Timeout,
		SortSet: skipList.New(),
		add:      make(chan string),   // 这里使用无缓冲保证操作的原子性
		pop:      make(chan string),
		decline : make(chan string),
	}
}

func (s *Seter) Run() {
	// 增删操作
	go func() {
		for {
			select {
			case key := <- s.add:
				s.SortSet.Set(key, s.InitScore)

			case key := <- s.pop:
				s.SortSet.Delete(key)

			case key := <- s.decline:
				s.SortSet.Increase(key, s.ScoreInterval)
			}
		}
	}()

	// 检查过期
	go func() {
		for {
			s.CheckIpsValid()
			logrus.Infof("pool剩余代理数%d", s.SortSet.Length())
			time.Sleep(2 * time.Second)
		}
	}()
}

func (s *Seter) CheckIpsValid() {
	s.Lock()
	defer s.Unlock()
	for i := int64(0) ; i < int64(s.SortSet.Length()); i++ {
		key,score := s.SortSet.GetDataByRank(i,true)

		time_stamp := s.SortSet.GetTimeStamp(key)

		if !check_time(time_stamp,s.Timeout) || score < s.PassScore {
			s.pop <- key
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


