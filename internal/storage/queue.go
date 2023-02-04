package storage

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type element struct {
	timestamp time.Time
	data      interface{}
}

type Queue struct {
	rwm  sync.RWMutex
	list *list.List
	size int
}

func New(size int) *Queue {
	return &Queue{rwm: sync.RWMutex{}, list: list.New(), size: size}
}

func (l *Queue) Push(s interface{}) {
	elem := element{timestamp: time.Now(), data: s}
	// defer l.rwm.Unlock() defer is high costs

	l.rwm.Lock()

	if l.list.Len() == l.size {
		l.list.Remove(l.list.Back())
	}

	l.list.PushFront(elem)

	l.rwm.Unlock()
}

func (l *Queue) GetElementsAfter(t time.Time) <-chan interface{} {
	out := make(chan interface{})

	defer l.rwm.RUnlock()
	l.rwm.RLock()

	go func() {
		defer close(out)

		lastElm := l.list.Back().Value.(element)
		if lastElm.timestamp.Unix() != t.Unix() && lastElm.timestamp.After(t) {
			return
		}

		for e := l.list.Front(); e != nil; e = e.Next() {
			elm := e.Value.(element)
			if elm.timestamp.Before(t) {
				return
			}
			out <- elm.data
		}
	}()

	return out
}

func (l *Queue) Print() {
	defer l.rwm.RUnlock()
	l.rwm.RLock()

	for e := l.list.Front(); e != nil; e = e.Next() {
		fmt.Printf("Elem:%v, %T\n", e.Value, e.Value)
	}
}

/*
func (l *List) PrintAvg(d time.Duration) {
	t := time.Now().Add(d * -1)
	defer l.rwm.RUnlock()
	l.rwm.RLock()

	lastElm := l.list.Back().Value.(sensors.Interface)
	if lastElm.GetTimestamp().Unix() != t.Unix() && lastElm.GetTimestamp().After(t) {
		fmt.Println("too early")
		return
	}

	for e := l.list.Front(); e != nil; e = e.Next() {
		if e.Value.(sensors.Interface).GetTimestamp().Before(t) {
			break
		}
		fmt.Printf("Sensors:%v, %T\n", e.Value, e.Value)
	}
}

func (l *List) GetSensorsAfter(t time.Time) <-chan sensors.Interface {
	out := make(chan sensors.Interface)

	defer l.rwm.RUnlock()
	l.rwm.RLock()

	go func() {
		defer close(out)

		lastElm := l.list.Back().Value.(sensors.Interface)
		if lastElm.GetTimestamp().Unix() != t.Unix() && lastElm.GetTimestamp().After(t) {
			return
		}

		for e := l.list.Front(); e != nil; e = e.Next() {
			if e.Value.(sensors.Interface).GetTimestamp().Before(t) {
				return
			}
			out <- e.Value.(sensors.Interface)
		}
	}()

	return out
}

func (l *List) Avg(in <-chan sensors.Interface) <-chan sensors.Interface {
	out := make(chan sensors.Interface)

	go func() {
		count := 0
		var a sensors.Interface
		for v := range in {
			count++
			a = v.Sum(&a)
		}
		a = a.Div(count)
		out <- a
		close(out)
	}()

	return out
}
*/
