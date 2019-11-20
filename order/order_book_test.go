package order

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuyInsert(t *testing.T) {
	//test sort buy_order
	var isBuy, isExist bool
	isBuy = true
	data := []uint64{5, 50, 10, 1, 100, 1000, 3, 1, 50, 5, 9, 90, 2, 2, 90, 900}
	expectData := []uint64{2, 2, 3, 1, 5, 50, 9, 90, 10, 1, 50, 5, 90, 900, 100, 1000}
	q := NewQuotation(len(data) / 2)
	for i := 0; i < len(data)/2; i++ {
		price := data[2*i]
		amount := data[2*i+1]
		q.Insert(isBuy, price, amount)
	}
	equalArray(t, expectData, q.Data)

	var price, amount uint64
	var index int

	//test search buy index
	price = 1
	amount = 10
	isExist, index, _ = q.BinarySearch(isBuy, price)
	assert.Equal(t, false, isExist)
	assert.Equal(t, 0, index)
	q.Insert(isBuy, price, amount)
	expectData = []uint64{1, 10, 2, 2, 3, 1, 5, 50, 9, 90, 10, 1, 50, 5, 90, 900, 100, 1000}
	equalArray(t, expectData, q.Data)

	price = 8
	amount = 100
	isExist, index, _ = q.BinarySearch(isBuy, price)
	assert.Equal(t, false, isExist)
	assert.Equal(t, 4, index)
	q.Insert(true, price, amount)
	expectData = []uint64{1, 10, 2, 2, 3, 1, 5, 50, 8, 100, 9, 90, 10, 1, 50, 5, 90, 900, 100, 1000}
	equalArray(t, expectData, q.Data)

	price = 1001
	amount = 100
	isExist, index, _ = q.BinarySearch(isBuy, price)
	assert.Equal(t, false, isExist)
	assert.Equal(t, 10, index)
	q.Insert(true, price, amount)
	expectData = []uint64{1, 10, 2, 2, 3, 1, 5, 50, 8, 100, 9, 90, 10, 1, 50, 5, 90, 900, 100, 1000, 1001, 100}
	equalArray(t, expectData, q.Data)

	price = 5
	amount = 200
	isExist, index, _ = q.BinarySearch(isBuy, price)
	assert.Equal(t, true, isExist)
	assert.Equal(t, 3, index)
	q.Insert(true, price, amount)
	expectData = []uint64{1, 10, 2, 2, 3, 1, 5, 250, 8, 100, 9, 90, 10, 1, 50, 5, 90, 900, 100, 1000, 1001, 100}
	equalArray(t, expectData, q.Data)

	price = 100
	amount = 100
	isExist, index, _ = q.BinarySearch(isBuy, price)
	assert.Equal(t, true, isExist)
	assert.Equal(t, 9, index)
	q.Insert(true, price, amount)
	expectData = []uint64{1, 10, 2, 2, 3, 1, 5, 250, 8, 100, 9, 90, 10, 1, 50, 5, 90, 900, 100, 1100, 1001, 100}
	equalArray(t, expectData, q.Data)
}

func TestSellInsert(t *testing.T) {
	//test sort buy_order
	var isBuy, isExist bool
	isBuy = false
	data := []uint64{100, 1000, 40, 4, 50, 500, 700, 70, 11, 110, 2, 10, 9, 90, 3, 30, 55, 5, 80, 800}
	expectData := []uint64{700, 70, 100, 1000, 80, 800, 55, 5, 50, 500, 40, 4, 11, 110, 9, 90, 3, 30, 2, 10}
	q := NewQuotation(len(data) / 2)
	for i := 0; i < len(data)/2; i++ {
		price := data[2*i]
		amount := data[2*i+1]
		q.Insert(false, price, amount)
	}
	equalArray(t, expectData, q.Data)

	var price, amount uint64
	var index int

	//test search buy index
	price = 1
	amount = 100
	isExist, index, _ = q.BinarySearch(isBuy, price)
	assert.Equal(t, false, isExist)
	assert.Equal(t, 10, index)
	q.Insert(isBuy, price, amount)
	expectData = []uint64{700, 70, 100, 1000, 80, 800, 55, 5, 50, 500, 40, 4, 11, 110, 9, 90, 3, 30, 2, 10, 1, 100}
	equalArray(t, expectData, q.Data)

	price = 35
	amount = 100
	isExist, index, _ = q.BinarySearch(isBuy, price)
	assert.Equal(t, false, isExist)
	assert.Equal(t, 6, index)
	q.Insert(isBuy, price, amount)
	expectData = []uint64{700, 70, 100, 1000, 80, 800, 55, 5, 50, 500, 40, 4, 35, 100, 11, 110, 9, 90, 3, 30, 2, 10, 1, 100}
	equalArray(t, expectData, q.Data)

	price = 1001
	amount = 10
	isExist, index, _ = q.BinarySearch(isBuy, price)
	assert.Equal(t, false, isExist)
	assert.Equal(t, 0, index)
	q.Insert(isBuy, price, amount)
	expectData = []uint64{1001, 10, 700, 70, 100, 1000, 80, 800, 55, 5, 50, 500, 40, 4, 35, 100, 11, 110, 9, 90, 3, 30, 2, 10, 1, 100}
	equalArray(t, expectData, q.Data)

	price = 80
	amount = 1000
	isExist, index, _ = q.BinarySearch(isBuy, price)
	assert.Equal(t, true, isExist)
	assert.Equal(t, 3, index)
	q.Insert(isBuy, price, amount)
	expectData = []uint64{1001, 10, 700, 70, 100, 1000, 80, 1800, 55, 5, 50, 500, 40, 4, 35, 100, 11, 110, 9, 90, 3, 30, 2, 10, 1, 100}
	equalArray(t, expectData, q.Data)

	price = 1
	amount = 1000
	isExist, index, _ = q.BinarySearch(isBuy, price)
	assert.Equal(t, true, isExist)
	assert.Equal(t, 12, index)
	q.Insert(isBuy, price, amount)
	expectData = []uint64{1001, 10, 700, 70, 100, 1000, 80, 1800, 55, 5, 50, 500, 40, 4, 35, 100, 11, 110, 9, 90, 3, 30, 2, 10, 1, 1100}
	equalArray(t, expectData, q.Data)
}

func equalArray(t *testing.T, b1 []uint64, b2 []uint64) {
	assert.Equal(t, len(b1), len(b2), "b1 len: %d, b2 len: %d", len(b1), len(b2))
	for i := 0; i < len(b1); i++ {
		assert.Equal(t, b1[i], b2[i])
	}
}

func TestSubAmount(t *testing.T) {
	var expectData []uint64
	data := []uint64{2, 2, 3, 1, 5, 50, 9, 90, 10, 1, 50, 5, 90, 900, 100, 1000}
	q := NewQuotation(len(data) / 2)
	q.Data = data

	q.SubAmount(3, 40)
	expectData = []uint64{2, 2, 3, 1, 5, 50, 9, 50, 10, 1, 50, 5, 90, 900, 100, 1000}
	equalArray(t, expectData, q.Data)

	q.SubAmount(3, 50)
	expectData = []uint64{2, 2, 3, 1, 5, 50, 10, 1, 50, 5, 90, 900, 100, 1000}
	equalArray(t, expectData, q.Data)

	q.SubAmount(0, 2)
	expectData = []uint64{3, 1, 5, 50, 10, 1, 50, 5, 90, 900, 100, 1000}
	equalArray(t, expectData, q.Data)

	q.SubAmount(5, 1000)
	expectData = []uint64{3, 1, 5, 50, 10, 1, 50, 5, 90, 900}
	equalArray(t, expectData, q.Data)
}

func TestBinarySerach(t *testing.T) {
	a := []int{1, 3, 6, 10, 15, 21, 28, 36, 45, 55}
	x := 57
	my := binarySearch(len(a), func(i int) bool { return a[i] >= x })
	t.Logf("my=%d", my)
	i := sort.Search(len(a), func(i int) bool { return a[i] >= x })
	if i < len(a) && a[i] == x {
		t.Logf("found %d at index %d in %v\n", x, i, a)
	} else {
		t.Logf("%d not found in %v, i: %d\n", x, a, i)
	}
}
