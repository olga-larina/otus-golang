package hw04lrucache

// Двусвязный список.
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
	Value interface{} // значение
	Next  *ListItem   // следующий элемент
	Prev  *ListItem   // предыдущий элемент
}

type list struct {
	head *ListItem // 1-ый элемент
	tail *ListItem // последний элемент
	size int       // количество элементов
}

func NewList() List {
	return new(list)
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
		l.tail = item
	} else {
		item.Next = l.head
		l.head.Prev = item
	}
	l.head = item
	l.size++
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	if l.tail == nil {
		l.head = item
	} else {
		item.Prev = l.tail
		l.tail.Next = item
	}
	l.tail = item
	l.size++
	return item
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	if i == l.head {
		l.head = i.Next
	}
	if i == l.tail {
		l.tail = i.Prev
	}
	l.size--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.head {
		return
	}
	i.Prev.Next = i.Next
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.tail = i.Prev
	}
	i.Prev = nil
	i.Next = l.head
	l.head.Prev = i
	l.head = i
}
