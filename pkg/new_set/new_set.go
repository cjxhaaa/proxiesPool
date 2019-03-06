package new_set

import (
	"math/rand"
	"sync"
	"time"
)

const (
	skipListMaxLevel = 32
	skipListP = 0.25
)

type (

	skipListLevel struct {
		forward    *skipListNode     //指向当前level下一个node
		span       uint64            //跨度，和下个node的距离
	}

	skipListNode struct {
		key      string              //key值
		score    float64             //节点分支
		backward *skipListNode       //指向上一个node
		level    []*skipListLevel    //该节点层
	}

	skipList struct {
		header   *skipListNode       // 头部结点
		tail     *skipListNode       // 尾部节点
		length   int64               // 节点个数
		maxLevel int16               // 最大的节点层数
	}

	obj struct {
		key       string
		score     float64
		timestamp int64
	}

	SortedSet struct {
		dict sync.Map
		sl   *skipList
	}

)


func skipListCreateNode(maxLevel int16, score float64, key string) *skipListNode {
	// 创建一个节点
	n := &skipListNode{
		score: score,
		key: key,
		level: make([]*skipListLevel, maxLevel),
	}

	for i := range n.level {
		n.level[i] = new(skipListLevel)
	}
	return n
}

func skipListInit() *skipList {
	// 跳表初始化，会创建一个头节点，层为1，分值0，key暂时为空字符串
	return &skipList{
		maxLevel: 1,
		header: skipListCreateNode(skipListMaxLevel, 0, ""),
	}
}

func randomLevel() int16 {
	// 随机层数
	level := int16(1)
	for float32(rand.Int31()&0xFFFF) < (skipListP * 0xFFFF) {
		level++
	}
	if level < skipListMaxLevel {
		return level
	}
	return skipListMaxLevel
}

func (sl *skipList) skipListInsert(key string, score float64) *skipListNode {
	// 跳表插入，先按层依次比较score,key值，找到正确的插入位置
	update := make([]*skipListNode, skipListMaxLevel)
	rank := make([]uint64, skipListMaxLevel)
	x := sl.header
	for i := sl.maxLevel - 1; i >=0; i-- {
		if i == sl.maxLevel - 1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}
		if x.level[i] != nil {
			for x.level[i].forward != nil && (x.level[i].forward.score < score ||
				(x.level[i].forward.score == score && x.level[i].forward.key < key)) {
					rank[i] += x.level[i].span //记录一下位置
					x = x.level[i].forward     //查看下一个节点
			}
		}
		update[i] = x
	}

	level := randomLevel()

	//当层数超过当前最大层数时，当前节点的span变成了和头节点的距离
	if level > sl.maxLevel {
		for i := sl.maxLevel; i < level; i++ {
			rank[i] = 0
			update[i] = sl.header  //把超过maxlevel的层用头结点补全
			update[i].level[i].span = uint64(sl.length)
		}
		sl.maxLevel = level
	}

	x = skipListCreateNode(level, score, key)
	for i := int16(0); i < level; i++ {
		x.level[i].forward = update[i].level[i].forward
		update[i].level[i].forward = x

		x.level[i].span = update[i].level[i].span - (rank[0] - rank[i])
		update[i].level[i].span = (rank[0] - rank[i]) + 1
	}

	for i := level; i < sl.maxLevel; i++ {
		update[i].level[i].span++
	}

	// 插入位置是头结点
	if update[0] == sl.header {
		x.backward = nil
	} else {
		x.backward = update[0]
	}

	// 插入位置是尾节点
	if x.level[0].forward != nil {
		x.level[0].forward.backward = x
	} else {
		sl.tail = x
	}
	sl.length++
	return x
}

func (sl *skipList) skipListDeleteNode(x *skipListNode, update []*skipListNode) {
	for i := int16(0); i < sl.maxLevel; i++ {
		if update[i].level[i].forward == x {
			update[i].level[i].span += x.level[i].span - 1
			update[i].level[i].forward = x.level[i].forward
		} else {
			update[i].level[i].span--
		}
	}

	//插入位置是尾节点
	if x.level[0].forward != nil {
		x.level[0].forward.backward = x.backward
	} else {
		sl.tail = x.backward
	}

	//删除元素后如果头结点最高层尾节点为nil，说明当前层可弃用
	for sl.maxLevel > 1 && sl.header.level[sl.maxLevel-1].forward == nil {
		sl.maxLevel--
	}

	sl.length--
}

func (sl *skipList) skipListDelete(key string, score float64) int {
	update := make([]*skipListNode, skipListMaxLevel)  //更新的位置
	x := sl.header
	for i := sl.maxLevel - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			(x.level[i].forward.score < score ||
				(x.level[i].forward.score == score && x.level[i].forward.key < key)) {
					x = x.level[i].forward
		}
		update[i] = x
	}

	x = x.level[0].forward
	if x != nil && score == x.score && x.key == key {
		sl.skipListDeleteNode(x,update)
		return 1
	}
	return 0 //没找到
}

func (sl *skipList) skipListGetRank(key string, score float64) int64 {
	rank := uint64(0)
	x := sl.header
	for i := sl.maxLevel - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			(x.level[i].forward.score < score || (x.level[i].forward.score == score && x.level[i].forward.key <= key)) {
			rank += x.level[i].span
			x = x.level[i].forward
		}

		if x.key == key {
			return int64(rank)
		}
	}
	return 0
}

func (sl *skipList) skipListGetElementByRank(rank uint64) *skipListNode {
	traversed := uint64(0)
	x := sl.header
	for i := sl.maxLevel - 1; i >= 0; i-- {
		for x.level[i].forward != nil && (traversed+x.level[i].span) <= rank {
			traversed += x.level[i].span
			x = x.level[i].forward
		}
		if traversed == rank {
			return x
		}
	}
	return nil
}


/*
---------------------------------------------------------------------------------分割线
*/

// 初始化
func New() *SortedSet {
	s := &SortedSet{
		dict: sync.Map{},
		sl:   skipListInit(),
	}
	return s
}

// 返回set长度
func (s *SortedSet) Length() int64 {
	return s.sl.length
}

func (s *SortedSet) Set(key string, score float64) {
	v, ook := s.dict.LoadOrStore(key,&obj{key:key,score:score})
	if vv,ok := v.(*obj);ok {
		if ook {
			if score != vv.score {
				s.sl.skipListDelete(key, vv.score)
				s.sl.skipListInsert(key, score)
				s.dict.Store(key,&obj{key:key,score:score,timestamp:vv.timestamp})
			}
		} else {
			s.sl.skipListInsert(key, score)
			s.dict.Store(key,&obj{key:key,score:score,timestamp:time.Now().Unix()})
		}
	}


}

func (s *SortedSet) Delete(key string) (ok bool) {
	v, ok := s.dict.Load(key)
	if ok {
		if vv,ok := v.(*obj);ok {
			s.sl.skipListDelete(key, vv.score)
			s.dict.Delete(key)
			return true
		}
	}
	return false
}

//根据key查询score,没有暂时返回0
func (s *SortedSet) GetScore(key string) (score float64, ok bool) {
	v, ok := s.dict.Load(key)
	if ok {
		if vv,ok := v.(*obj);ok {
			return vv.score, true
		}
	}
	return 0, false
}

func (s *SortedSet) GetRank(key string, reverse bool) (rank int64, score float64) {
	v, ok := s.dict.Load(key)
	if !ok {
		return -1,0 //没找到返回位置-1，分值0
	}
	if vv,ok := v.(*obj);ok {
		r := s.sl.skipListGetRank(key, vv.score)

		if reverse {
			r = s.sl.length - r
		} else {
			r--
		}
		return int64(r), vv.score
	}
	return -1,0
}

func (s *SortedSet) GetDataByRank(rank int64,reverse bool) (key string, score float64) {
	if rank < 0 || rank > s.sl.length {
		return "", 0
	}
	if reverse {
		rank = s.sl.length - rank
	} else {
		rank++
	}

	n := s.sl.skipListGetElementByRank(uint64(rank))
	if n == nil {
		return "", 0
	}
	v, ok := s.dict.Load(n.key)
	if !ok || v == nil {
		return "", 0
	}
	if vv,ok := v.(*obj);ok {
		return vv.key, vv.score
	}
	return "", 0
}

func (s *SortedSet) Increase(key string, descore float64) {
	_, score := s.GetRank(key,false)
	score += descore
	s.Set(key,score)
}

func (s *SortedSet) GetTimeStamp(key string) int64 {
	v, ok := s.dict.Load(key)
	if !ok {
		return 0
	}
	if vv,ok := v.(*obj);ok {
		return vv.timestamp
	}
	return 0
}



