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

func (l *Queue) SetSize(newsize int32) bool {
	l.rwm.Lock()
	defer l.rwm.Unlock()
	if int(newsize) > l.size {
		l.size = int(newsize)
		return true
	}
	return false
}

func (l *Queue) Push(s interface{}) {
	defer l.rwm.Unlock()

	l.rwm.Lock()
	if l.size == 0 {
		return
	}
	if l.list.Len() == l.size {
		l.list.Remove(l.list.Back())
	}
	l.list.PushFront(element{timestamp: time.Now(), data: s})
}

func (l *Queue) GetElementsAfter(t time.Time) <-chan interface{} {
	out := make(chan interface{})

	defer l.rwm.RUnlock()
	l.rwm.RLock()

	go func() {
		defer close(out)

		lastElm := l.list.Back().Value.(element)
		if lastElm.timestamp.Unix() <= t.Unix() {
			for e := l.list.Front(); e != nil; e = e.Next() {
				elm := e.Value.(element)
				if elm.timestamp.Before(t) {
					return
				}
				out <- elm.data
			}
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
