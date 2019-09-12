package order

import (
	"fmt"
)

type Quotation []byte

func NewQuotation(cap uint32) (q Quotation) {
	q = make([]byte, 0, 8*cap)
	return
}

func (q Quotation) String() string {
	return fmt.Sprint("Quotation \n")
}

func (q *Quotation) Insert(o *Order) {
}

func (q *Quotation) SubAmount(price uint64, amount uint64) {
}
