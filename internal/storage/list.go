package storage

/*
type LimitedStack


type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	Remove(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	head     *ListItem
	tail     *ListItem
	size     int32
	max_size int32
}

func (l *list) Len() int {
	return (l.size)
}

func (l *list) Front() *ListItem {
	return (l.head)
}

func (l *list) Back() *ListItem {
	return (l.tail)
}

func (l *list) PushFront(v interface{}) *ListItem {
	elm := &ListItem{v, nil, nil}

	if l.size == 0 {
		l.tail = elm
	} else {
		l.head.Prev = elm
		elm.Next = l.head
	}

	l.head = elm
	l.size++

	return (l.head)
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}

	elmPrev := i.Prev
	elmNext := i.Next

	if elmPrev != nil {
		elmPrev.Next = elmNext
		if elmPrev.Next == nil {
			l.tail = elmPrev
		}
	}

	if elmNext != nil {
		elmNext.Prev = elmPrev
		if elmNext.Prev == nil {
			l.head = elmNext
		}
	}

	l.size--
	if l.size == 0 {
		l.head = nil
		l.tail = nil
	}
}

func NewList(len int32) List {
	return &list{}
}
*/
