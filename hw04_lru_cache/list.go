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
	len   int
	first *ListItem
	last  *ListItem
}

func (list *list) Len() int {
	return list.len //len(list.items)
}
func (list *list) increaseCount() {
	list.len++
}
func (list *list) decreaseCount() {
	list.len--
}
func (list *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
		Next:  list.first,
		Prev:  nil,
	}
	if list.len == 0 {
		list.last = newItem
	} else {
		list.first.Prev = newItem
	}
	list.first = newItem

	list.increaseCount()

	return newItem
}
func (list *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  list.last,
	}
	if list.len == 0 {
		list.first = newItem
	} else {
		list.last.Next = newItem
	}
	list.last = newItem

	list.increaseCount()

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
	wasDeleted := false
	if deleted.Prev != nil {
		wasDeleted = true
		deleted.Prev.Next = deleted.Next
	}
	if deleted.Next != nil {
		wasDeleted = true
		deleted.Next.Prev = deleted.Prev
	}
	// если флаг не поднимали - то нам передали какую то фигню не из списка
	if wasDeleted {
		list.decreaseCount()
	}
}

// переместить элемент в начало
func (list *list) MoveToFront(moved *ListItem) {
	// нет смысла вперед двигать то что спереди
	if moved.Prev == nil {
		return
	}

	// если последний взяли - то обновим лист
	if moved.Next != nil {
		moved.Next.Prev = moved.Prev
	} else {
		list.last = moved.Prev
	}

	moved.Prev.Next = moved.Next
	moved.Next = list.first
	list.first.Prev = moved
	list.first = moved
	moved.Prev = nil
}

func NewList() List {
	return new(list)
}
