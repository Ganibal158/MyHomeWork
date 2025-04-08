package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	front *ListItem
	back  *ListItem
	kolv  int
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.kolv
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	knot := &ListItem{}
	knot.Value = v

	if l.kolv == 0 {
		l.front = knot
		l.back = knot
	} else {
		knot.Next = l.front
		l.front.Prev = knot
		l.front = knot
	}
	l.kolv++
	return knot
}

func (l *list) PushBack(v interface{}) *ListItem {
	knot := &ListItem{}
	knot.Value = v
	if l.kolv == 0 {
		l.front = knot
		l.back = knot
	} else {
		knot.Prev = l.back
		l.back.Next = knot
		l.back = knot
	}
	l.kolv++
	return knot
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.front = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.back = i.Prev
	}
	l.kolv--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil || i == l.front {
		return
	} else if l.back == i {
		i.Prev.Next = nil
		l.back = i.Prev
		l.front.Prev = i
		i.Prev = nil
		i.Next = l.front
		l.front = i
	} else {
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
		i.Next = l.front
		i.Prev = nil
		l.front.Prev = i
		l.front = i
	}
}
