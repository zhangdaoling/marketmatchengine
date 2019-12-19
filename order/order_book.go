package order

import (
	"fmt"
	"github.com/zhangdaoling/marketmatchengine/common"
)

type OrderBook struct {
	Index              uint64   `json:"index"`
	MatchOrderID       uint64   `json:"match_order_id"`
	MatchTime          uint64   `json:"match_time"`
	MatchPrice         uint64   `json:match_price`
	BuyQuotationSlice  []uint64 `json:"buy_quotations"`
	SellQuotationSlice []uint64 `json:"sell_quotations"`
}

type BookSlice struct {
	Data []uint64
}

func NewQuotation(cap int) (q *BookSlice) {
	q = &BookSlice{make([]uint64, 0, 2*cap)}
	return
}

func (q BookSlice) String() string {
	return fmt.Sprint("OrderBook \n")
}

//index is the position
func (q *BookSlice) BinarySearch(isBuy bool, price uint64) (isExist bool, index int, amount uint64) {
	length := len(q.Data) / 2
	if isBuy {
		index = binarySearch(length, func(i int) bool { return q.Data[2*i] >= price })
	} else {
		index = binarySearch(length, func(i int) bool { return q.Data[2*i] <= price })
	}
	if index < length && q.Data[2*index] == price { //find the price
		return true, index, q.Data[2*index+1]
	}
	// not found
	return false, index, amount
}

//BookSlice must the same with isBuy
func (q *BookSlice) Insert(isBuy bool, price uint64, amount uint64) {
	isExist, index, _ := q.BinarySearch(isBuy, price)
	length := len(q.Data) / 2
	if isExist {
		q.Data[2*index+1] += amount
	} else if index == length {
		q.Data = append(q.Data, price, amount)
	} else {
		q.Data = append(q.Data, 0, 0)
		copy(q.Data[2*index+2:2*length+2], q.Data[2*index:2*length])
		q.Data[2*index] = price
		q.Data[2*index+1] = amount
	}
	return
}

//use search before subAmount
//BookSlice must the same with isBuy
//amount <= sliceAmount; index < length
func (q *BookSlice) SubAmount(index int, amount uint64) {
	length := len(q.Data) / 2
	q.Data[2*index+1] -= amount
	if q.Data[2*index+1] == 0 { //rebuild slice
		if index != length-1 {
			copy(q.Data[2*index:2*length-2], q.Data[2*index+2:2*length])
		}
		q.Data = q.Data[0 : 2*length-2]
	}
	return
}

func (q *BookSlice) Serialize(zero *common.ZeroCopySink) {
	zero = common.NewZeroCopySink(nil, 2*len(q.Data))
	for _, v := range q.Data {
		zero.WriteUint64(v)
	}
}

//copy from sort.BinarySearch(), no test
func binarySearch(n int, f func(int) bool) int {
	var i, j, h int
	i = 0
	j = n
	for i < j {
		h = int(uint(i+j) >> 1)
		if !f(h) {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}
