// Package priority implements a PriorityQueue
//
// From: https://pkg.go.dev/container/heap@go1.17.6#example-package-PriorityQueue
package priority

import "container/heap"

// An Item is something we manage in a Priority queue.
type Item struct {
	Value    string // The Value of the item; arbitrary.
	Priority int    // The Priority of the item in the queue.
	// The Index is needed by Update and is maintained by the Priority.Interface methods.
	Index int // The Index of the item in the Priority.
}

// A Queue implements Priority.Interface and holds Items.
//
// A Queue with some items, adds and manipulates an item,
// and then removes the items in Priority order.
//
// Example:
//        // Some items and their priorities.
//        items := map[string]int{
//            "banana": 3, "apple": 2, "pear": 4,
//        }
//
//        // Create a Priority queue, put the items in it, and
//        // establish the Priority queue (Priority) invariants.
//        pq := make(Queue, len(items))
//        i := 0
//        for Value, Priority := range items {
//            pq[i] = &Item{
//                Value:    Value,
//                Priority: Priority,
//                Index:    i,
//            }
//            i++
//        }
//        Priority.Init(&pq)
//
//        // Insert a new item and then modify its Priority.
//        item := &Item{
//            Value:    "orange",
//            Priority: 1,
//        }
//        Priority.Push(&pq, item)
//        pq.Update(item, item.Value, 5)
//
//        // Take the items out; they arrive in decreasing Priority order.
//        for pq.Len() > 0 {
//            item := Priority.Pop(&pq).(*Item)
//            fmt.Printf("%.2d:%s ", item.Priority, item.Value)
//        }
//
// From: https://pkg.go.dev/container/heap@go1.17.6#example-package-PriorityQueue
type Queue []*Item

func (pq Queue) Len() int { return len(pq) }

func (pq Queue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, Priority so we use greater than here.
	return pq[i].Priority > pq[j].Priority
}

func (pq Queue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *Queue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.Index = n
	*pq = append(*pq, item)
}

// Update modifies the Priority and Value of an Item in the queue.
func (pq *Queue) Update(item *Item, value string, priority int) {
	item.Value = value
	item.Priority = priority
	heap.Fix(pq, item.Index)
}

func (pq *Queue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *Queue) PushValue(value string) *Item {
	i := &Item{
		Value:    value,
		Priority: 0,
		Index:    0,
	}
	heap.Push(pq, i)
	return i
}

func (pq *Queue) PopValue() string {
	return heap.Pop(pq).(*Item).Value
}

func (pq *Queue) UpdatePriority(item *Item, priority int) {
	pq.Update(item, item.Value, priority)
}
