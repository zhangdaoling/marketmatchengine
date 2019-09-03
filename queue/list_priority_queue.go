package queue

import "errors"

type Element struct {
	next, prev *Element
	list       *PriorityList
	Value      Item
}

func (e *Element) Next() *Element {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

func (e *Element) Prev() *Element {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

type PriorityList struct {
	root   Element
	len    int
	search map[int32]*Element
}

func (p *PriorityList) Insert(item Item) (err error) {
	if item == nil {
		return
	}
	if _, ok := p.search[item.Key()]; ok {
		return errors.New("same ID")
	}
	element := &p.root
	for {
		next := element.Next()
		if next == nil {
			p.insertAfter(item, element)
			return
		}
		if item.Compare(next.Value) >= 0 {
			p.insertAfter(item, element)
			return
		}
		element = next
	}
	return
}

func (p *PriorityList) Cancel(key int32) (item Item) {
	if element, ok := p.search[key]; ok {
		delete(p.search, key)
		p.remove(element)
		return element.Value
	}
	return nil
}

func (p *PriorityList) First() (item Item) {
	e := p.root.next
	if e == nil {
		return nil
	}
	return e.Value
}

func (p *PriorityList) Pop() (item Item) {
	e := p.root.next
	if e == nil {
		return nil
	}
	e = p.remove(e)
	return e.Value
}

func NewPriorityList() (p *PriorityList) {
	p = &PriorityList{}
	p.search = make(map[int32]*Element, 100)
	p.init()
	return p
}

func (p *PriorityList) init() *PriorityList {
	p.root.next = &p.root
	p.root.prev = &p.root
	p.root.list = p
	p.len = 0
	return p
}

func (p *PriorityList) remove(e *Element) *Element {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil
	e.prev = nil
	e.list = nil
	p.len--
	delete(p.search, e.Value.Key())
	return e
}
func (p *PriorityList) insertBefore(i Item, mark *Element) *Element {
	if mark.list != p {
		return nil
	}
	return p.insertValue(i, mark.prev)
}

func (p *PriorityList) insertAfter(i Item, mark *Element) *Element {
	if mark.list != p {
		return nil
	}
	return p.insertValue(i, mark)
}

func (p *PriorityList) insertValue(i Item, at *Element) *Element {
	return p.insert(&Element{Value: i}, at)
}

func (p *PriorityList) insert(e, at *Element) *Element {
	n := at.next
	at.next = e
	e.prev = at
	e.next = n
	n.prev = e
	e.list = p
	p.len++
	p.search[e.Value.Key()] = e
	return e
}
