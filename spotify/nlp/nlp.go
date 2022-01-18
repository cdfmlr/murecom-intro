package nlp

import (
	"container/heap"
	"log"
	"regexp"
	"sort"
	"spotifyplaylist/priority"
	"strings"
	"sync"
)

//type item struct {
//	Word  string
//	Count int
//}
//
//type WordsHeap []item
//
//func (h WordsHeap) Len() int { return len(h) }
//
//func (h WordsHeap) Less(i, j int) bool {
//	// We want Pop to give us the highest, not lowest, Priority so we use greater than here.
//	return h[i].Count > h[j].Count
//}
//
//func (h WordsHeap) Swap(i, j int) {
//	h[i], h[j] = h[j], h[i]
//}
//
//func (h *WordsHeap) Push(x interface{}) {
//	// Push and Pop use pointer receivers because they modify the slice's length,
//	// not just its contents.
//	*h = append(*h, x.(item))
//}
//
//func (h *WordsHeap) Pop() interface{} {
//	old := *h
//	n := len(old)
//	x := old[n-1]
//	*h = old[0 : n-1]
//	return x
//}

////////////////
// Interface  //
////////////////

// WordsCounter 词频统计
//
// 目前只有英文的，推荐用 EnglishWordsCounterPQ，这个性能比较好：
//     $  go test -bench=.  -benchmem
//     goos: darwin
//     goarch: amd64
//     pkg: spotifyplaylist/nlp
//     cpu: Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz
//     BenchmarkEnglishWordCounterPQ-4    33    35152096 ns/op    7236090 B/op    57066 allocs/op
//     BenchmarkEnglishWordCounterQS-4    28    39783905 ns/op    7251620 B/op    57171 allocs/op
type WordsCounter interface {
	// AddSentence 从句子中分词，统计
	AddSentence(s string)
	// PopMostCommon 获取频次最高的词
	PopMostCommon() string
}

/////////////////////////////////////////////
// counter based on priority queue (heap)  //
/////////////////////////////////////////////

// EnglishWordsCounterPQ 分词，统计频数，英文适用，基于堆（优先队列）
type EnglishWordsCounterPQ struct {
	Map      map[string]*priority.Item
	Priority *priority.Queue

	sync.Mutex
}

func NewEnglishWordCounterPQ() WordsCounter {
	mp := make(map[string]*priority.Item)
	pq := make(priority.Queue, 0)
	heap.Init(&pq)
	c := &EnglishWordsCounterPQ{
		Map:      mp,
		Priority: &pq,
	}

	return c
}

// onlyWords removes all characters except letters and spaces
func onlyWords(s string) string {
	// only want letters and numbers
	reg, err := regexp.Compile("[^a-zA-Z0-9 ]+")
	if err != nil {
		log.Fatal(err)
	}
	s = reg.ReplaceAllString(s, "")

	return s
}

func (c *EnglishWordsCounterPQ) AddSentence(s string) {
	c.Lock()
	defer c.Unlock()

	s = strings.ToLower(onlyWords(s))

	words := strings.Fields(s)
	for _, w := range words {
		if item, ok := c.Map[w]; ok { // exist
			c.Priority.UpdatePriority(item, item.Priority+1)
		} else { // new
			item = c.Priority.PushValue(w)
			c.Map[w] = item
		}
	}
}

func (c *EnglishWordsCounterPQ) PopMostCommon() string {
	c.Lock()
	defer c.Unlock()

	word := c.Priority.PopValue()
	delete(c.Map, word)
	return word
}

////////////////////////////
// counter based on sort  //
////////////////////////////

// EnglishWordsCounterQS 是基于顺序表和快速排序的 counter
// Benchmark 时空性能不如 PQ
type EnglishWordsCounterQS struct {
	Map  map[string]int // [word]: idx_in_List
	List wordsList

	sync.Mutex
}

func NewEnglishWordCounterQS() WordsCounter {
	return &EnglishWordsCounterQS{Map: map[string]int{}}
}

type item struct {
	value string
	count int
}

type wordsList []*item

func (w wordsList) Len() int {
	return len(w)
}

func (w wordsList) Less(i, j int) bool {
	return w[i].count > w[j].count
}

func (w wordsList) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

func (c *EnglishWordsCounterQS) AddSentence(s string) {
	c.Lock()
	defer c.Unlock()

	s = strings.ToLower(onlyWords(s))
	words := strings.Fields(s)
	for _, w := range words {
		if i, ok := c.Map[w]; ok { // exist
			c.List[i].count++
		} else { // new
			c.List = append(c.List, &item{
				value: w,
				count: 1,
			})
			c.Map[w] = len(c.List) - 1
		}
	}
}

func (c *EnglishWordsCounterQS) PopMostCommon() string {
	c.Lock()
	defer c.Unlock()

	sort.Sort(c.List)
	// XXX: 现在的性能瓶颈应该就在这个地方，如果可以提升，有望在时间上超越 PQ
	for idx, item := range c.List {
		c.Map[item.value] = idx
	}

	w := c.List[0]
	c.List = c.List[1:]

	delete(c.Map, w.value)
	return w.value
}
