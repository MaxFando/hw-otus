package hw04lrucache

type List interface {
	Len() int                          // длина списка
	Front() *ListItem                  // первый элемент списка
	Back() *ListItem                   // последний элемент списка
	PushFront(v interface{}) *ListItem // добавить значение в начало
	PushBack(v interface{}) *ListItem  // добавить значение в конец
	Remove(i *ListItem)                // удалить элемент
	MoveToFront(i *ListItem)           // переместить элемент в начало
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

var _ List = (*list)(nil)

type list struct {
	size int
	head *ListItem
	tail *ListItem
}

func (l *list) Len() int {
	return l.size
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	if l.head == nil {
		l.head = item
		l.tail = item
	} else {
		item.Next = l.head
		l.head.Prev = item
		l.head = item
	}
	l.size++
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	if l.tail == nil {
		l.head = item
		l.tail = item
	} else {
		item.Prev = l.tail
		l.tail.Next = item
		l.tail = item
	}
	l.size++
	return item
}

func (l *list) Remove(item *ListItem) {
	if item.Prev != nil {
		item.Prev.Next = item.Next
	} else {
		l.head = item.Next
	}
	if item.Next != nil {
		item.Next.Prev = item.Prev
	} else {
		l.tail = item.Prev
	}
	l.size--
}

func (l *list) MoveToFront(item *ListItem) {
	if item.Prev == nil {
		return
	}
	l.Remove(item)
	l.PushFront(item.Value)
}

func NewList() List {
	return new(list)
}
