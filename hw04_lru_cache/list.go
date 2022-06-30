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
	len   int
}

// NewList creates default initialized list object.
func NewList() List {
	return new(list)
}

// Len returns length of the list.
func (ls *list) Len() int {
	return ls.len
}

// Front returns the first element of the list.
func (ls *list) Front() *ListItem {
	if ls.len == 0 {
		return nil
	}
	return ls.front
}

// Back returns the last element of the list.
func (ls *list) Back() *ListItem {
	if ls.len == 0 {
		return nil
	}
	return ls.back
}

// PushFront pushes the new element with the value v at the front of the list.
func (ls *list) PushFront(v interface{}) *ListItem {
	i := new(ListItem)
	i.Value = v
	i.Next = ls.front
	i.Prev, i.Next = nil, ls.front
	if ls.front != nil {
		ls.front.Prev = i
	}
	ls.front = i
	if ls.len == 0 {
		ls.back = i
	}
	ls.len++
	return i
}

// PushBack pushes the new element with the value v at the back of the list.
func (ls *list) PushBack(v interface{}) *ListItem {
	i := new(ListItem)
	i.Value = v
	i.Prev, i.Next = ls.back, nil
	if ls.back != nil {
		ls.back.Next = i
	}
	ls.back = i
	if ls.len == 0 {
		ls.front = i
	}
	ls.len++
	return i
}

// Remove removes the element from the list.
func (ls *list) Remove(i *ListItem) {
	ls.pickOut(i)

	i.Prev, i.Next = nil, nil
	ls.len--
}

// MoveToFront moves the element at the front of the list.
func (ls *list) MoveToFront(i *ListItem) {
	ls.pickOut(i)
	i.Prev, i.Next = nil, ls.front

	if ls.front != nil {
		ls.front.Prev = i
	}
	ls.front = i
	if ls.len == 0 {
		ls.back = i
	}
}

// pickOut picks out the element from the list without destroying en element for for following moving.
func (ls *list) pickOut(i *ListItem) {
	if ls.Len() <= 1 {
		ls.front, ls.back = nil, nil
		return
	}

	if i == ls.front {
		ls.front = ls.front.Next
		if ls.front != nil {
			ls.front.Prev = nil
		}
	} else {
		i.Prev.Next = i.Next
	}
	if i == ls.back {
		ls.back = ls.back.Prev
		if ls.back != nil {
			ls.back.Next = nil
		}
	} else {
		i.Next.Prev = i.Prev
	}
}
