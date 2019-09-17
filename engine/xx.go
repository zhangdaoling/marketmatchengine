package engine

import (
	"time"

	"github.com/zhangdaoling/marketmatchengine/order"
)

type Enginexx struct {
	OrderChan       chan *order.Order
	TransactionChan chan *order.Transaction
	CancelOrderChan chan *order.CancelOrder
	PersistTime     int
	PersistPath     string
	Engine          Engine
}

func NewEnginexxx(orderChan chan *order.Order, transactionChan chan *order.Transaction, cancelOrderChan chan *order.CancelOrder, persistTime int, persistPath string) (engine *Enginexx, err error) {
	engine = &Enginexx{
		OrderChan:       orderChan,
		TransactionChan: transactionChan,
		CancelOrderChan: cancelOrderChan,
		PersistTime:     persistTime,
		PersistPath:     persistPath,
	}
	return engine, nil
}

func (e *Enginexx) Loop(shutdown chan struct{}) {
	timer := time.NewTimer(time.Duration(e.PersistTime) * time.Minute)
	for {
		select {
		case <-shutdown:
			return
		case o := <-e.OrderChan:
			e.ProcessOrder(o)
		case <-timer.C:
		}
	}
}

func (e *Enginexx) ProcessOrder(o *order.Order) (err error) {
	return
}
