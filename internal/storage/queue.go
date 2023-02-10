package storage

import (
	"container/list"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/FRiniZ/system-stats-daemon/internal/services/sensors/common"
)

type element struct {
	timestamp time.Time
	data      common.Sensor
}

type Queue struct {
	rwm  sync.RWMutex
	list *list.List
	size int
}

func New(size int) *Queue {
	return &Queue{rwm: sync.RWMutex{}, list: list.New(), size: size}
}

func (l *Queue) SetSize(owner string, newsize int32) {
	l.rwm.Lock()
	defer l.rwm.Unlock()
	if int(newsize) > l.size {
		l.size = int(newsize)
		log.Printf("[%s] Changed size of queue to:%d", owner, newsize)
	}
}

func (l *Queue) Push(s common.Sensor, t time.Time) {
	defer l.rwm.Unlock()

	l.rwm.Lock()
	if l.size == 0 {
		return
	}
	if l.list.Len() == l.size {
		l.list.Remove(l.list.Back())
	}
	l.list.PushFront(element{timestamp: t, data: s})
}

func (l *Queue) GetElementsAfter(t time.Time) <-chan common.Sensor {
	out := make(chan common.Sensor)

	defer l.rwm.RUnlock()
	l.rwm.RLock()

	go func() {
		defer close(out)
		for e := l.list.Front(); e != nil; e = e.Next() {
			elm := e.Value.(element)
			if t.After(elm.timestamp) {
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
