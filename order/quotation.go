package order

import (
	"fmt"
	"github.com/zhangdaoling/marketmatchengine/common"
)

type Quotation struct {
	Time               uint64
	BuyQuotationSlice  []uint64 `json:"buy_quotations"`
	SellQuotationSlice []uint64 `json:"sell_quotations"`
}

type QuotationSlice struct {
	Data []uint64
}

func NewQuotation(cap int) (q *QuotationSlice) {
	q = &QuotationSlice{make([]uint64, 0, 2*cap)}
	return
}

func (q QuotationSlice) String() string {
	return fmt.Sprint("Quotation \n")
}

func (q *QuotationSlice) Insert(isBuy bool, price uint64, amount uint64) {
	var index int
	length := len(q.Data) / 2
	if isBuy {
		index = binarySearch(length, func(i int) bool { return q.Data[2*i] >= price })
	} else {
		index = binarySearch(length, func(i int) bool { return q.Data[2*i] <= price })
	}

	if index < length && q.Data[2*index] == price { //find the price
		q.Data[2*index+1] += amount
	} else if index == length { // not found price, add to the tail
		q.Data = append(q.Data, price, amount)
	} else { //not found price, rebuild slice
		q.Data = append(q.Data, 0, 0)
		copy(q.Data[2*index+2:2*length+2], q.Data[2*index:2*length])
		q.Data[2*index] = price
		q.Data[2*index+1] = amount
	}
	return
}

func (q *QuotationSlice) SubAmount(price uint64, amount uint64, isBuy bool) (isExist bool, err error) {
	if len(q.Data) == 0 {
		return false, nil
	}
	var index int
	length := len(q.Data) / 2
	if isBuy {
		index = binarySearch(length, func(i int) bool { return q.Data[2*i] >= price })
	} else {
		index = binarySearch(length, func(i int) bool { return q.Data[2*i] <= price })
	}
	if index < length && q.Data[2*index] == price { //find the price
		if q.Data[2*index+1] < amount {
			return false, common.ErrQuotationAmount
		}
		q.Data[2*index+1] -= amount
		if q.Data[2*index+1] == 0 { //rebuild slice
			if 2*index != length-2 {
				copy(q.Data[2*index:2*length], q.Data[2*index+2:2*length+2])
			}
			q.Data = q.Data[0 : 2*length-2]
		}
		return true, nil
	}
	return false, nil
}

func (q *QuotationSlice) Serialize(zero *common.ZeroCopySink) {
	zero = common.NewZeroCopySink(nil, 2*len(q.Data))
	for _, v := range q.Data {
		zero.WriteUint64(v)
	}
}

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
