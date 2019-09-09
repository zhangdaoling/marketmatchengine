package queue

import (
	"github.com/zhangdaoling/marketmatchengine/common"
	"gotest.tools/assert"
	"testing"
)

type order struct {
	Number int
	ID     uint32
}

func (o *order) Compare(i interface{}) int {
	o1 := i.(*order)
	if o.Number > o1.Number {
		return 1
	} else if o.Number == o1.Number {
		return 0
	} else {
		return -1
	}
}

func (o *order) Key() uint32 {
	return o.ID
}

func (o *order) Serialize() (zero *common.ZeroCopySink) {
	return
}

func TestInsert(t *testing.T) {
	data := []*order{
		&order{3, 1},
		&order{2, 2},
		&order{5, 3},
		&order{1, 4},
		&order{4, 5},
		&order{7, 6},
		&order{6, 7}}
	q := NewPriorityList()
	for _, v := range data {
		q.Insert(v)
	}
	v := q.First()
	a := v.(*order)
	assert.Equal(t, 7, a.Number, a)
	for i := 0; i < 7; i++ {
		v := q.Pop()
		a := v.(*order)
		assert.Equal(t, 7-i, a.Number, a)
	}
	return
}
