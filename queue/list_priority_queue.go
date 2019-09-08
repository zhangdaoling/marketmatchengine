//no lock, not go
package queue

import (
	"github.com/zhangdaoling/marketmatchengine/common"
)

var List_Queue_Name = "list_priority_queue"

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
	len    uint32
	search map[uint32]*Element
}

func NewPriorityList() (p *PriorityList) {
	p = &PriorityList{}
	p.search = make(map[uint32]*Element, 100)
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

//for PriorityQueue interface
func (p *PriorityList) Insert(item Item) (i Item) {
	if item == nil {
		return nil
	}
	if _, ok := p.search[item.Key()]; ok {
		return nil
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

//for PriorityQueue interface
func (p *PriorityList) Cancel(key uint32) (item Item) {
	if element, ok := p.search[key]; ok {
		delete(p.search, key)
		p.remove(element)
		return element.Value
	}
	return nil
}

//for PriorityQueue interface
func (p *PriorityList) First() (item Item) {
	e := p.root.next
	if e == nil {
		return nil
	}
	return e.Value
}

//for PriorityQueue interface
func (p *PriorityList) Pop() (item Item) {
	e := p.root.next
	if e == nil {
		return nil
	}
	e = p.remove(e)
	return e.Value
}

//for PriorityQueue interface
func (p *PriorityList) Len() uint32 {
	return p.len
}

//for PriorityQueue interface
func (p *PriorityList) Serialize() (zero *common.ZeroCopySink) {
	zero = common.NewZeroCopySink(nil)
	zero.WriteString(List_Queue_Name)
	zero.WriteUint32(p.len)
	for element := p.root.Next(); element != nil; element = element.Next() {
		data := element.Value.Serialize()
		zero.WriteVarBytes(data.Bytes())
	}
	return
}

func (p *PriorityList) remove(e *Element) *Element {
	e.prev.next = e.next
	if e.next != nil {
		e.next.prev = e.prev
	}
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
