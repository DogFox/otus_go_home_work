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
	Value interface{} // значение
	Next  *ListItem   // следующий элемент
	Prev  *ListItem   // предыдущий элемент
}

type list struct {
	// items []ListItem
	len   int
	first *ListItem
	last  *ListItem
}

func (list *list) Len() int {
	return list.len //len(list.items)
}
func (list *list) IncreaseCount() {
	list.len++
}
func (list *list) DecreaseCount() {
	list.len--
}
func (list *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
		Next:  list.first,
		Prev:  nil,
	}
	if list.len == 0 {
		list.first = newItem
		list.last = newItem
	}
	list.IncreaseCount()
	list.first = newItem
	// list.items = append(list.items, *newItem)

	return newItem
}
func (list *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  list.last,
	}
	list.IncreaseCount()
	list.last = newItem
	// list.items = append(list.items, *newItem)

	return newItem
}

func (list *list) Front() *ListItem {
	return list.first
}

func (list *list) Back() *ListItem {
	return list.last
}

// удалить элемент
func (list *list) Remove(deleted *ListItem) {
	deleted.Prev.Next = deleted.Next
	deleted.Next.Prev = deleted.Prev
	list.DecreaseCount()
}

// переместить элемент в начало
func (list *list) MoveToFront(moved *ListItem) {
	moved.Prev.Next = moved.Next
	moved.Next.Prev = moved.Prev
	list.first.Prev = moved
	moved.Prev = nil
}

func NewList() List {
	return new(list)
}
