package engine

import (
	"fmt"
	"github.com/zhangdaoling/marketmatchengine/order"
	"github.com/zhangdaoling/marketmatchengine/queue"
	"time"
)

//to do
func NewEngineFromFile(engineFile string, orderChan chan *order.Order, matchResultChan chan *order.MatchResult) (engine *Engine, err error) {
	return
}

//to do
func(e *Engine) Serialize() (data []byte){
	return
}

//to do
func UnSerialize(data []byte, e *Engine) (err error){
	return
}

type Engine struct {
	OrderChan       chan *order.Order
	MatchResultChan chan *order.MatchResult
	LastMatchPrice  int64
	LastOrderTime   int64
	LastOrderID     int32
	Symbol          string
	BuyQueue        queue.PriorityQueue
	SellQueue       queue.PriorityQueue
}

func NewEngine(orderChan chan *order.Order, matchResultChan chan *order.MatchResult, symbol string, lastPrice int64) (engine *Engine, err error) {
	sellQueue := queue.NewPriorityList()
	buyQueue := queue.NewPriorityList()
	engine = &Engine{
		OrderChan:       orderChan,
		MatchResultChan: matchResultChan,
		Symbol:          symbol,
		BuyQueue:        buyQueue,
		SellQueue:       sellQueue,
		LastMatchPrice:  lastPrice,
	}
	return engine, nil
}

func (e *Engine) Loop(shutdown chan struct{}) {
	timer := time.NewTimer(100*time.Second)
	for {
		select {
		case <-shutdown:
			return
		case o := <-e.OrderChan:
			e.processOrder(o)
		case <-timer.C:
			e.Serialize()
		}
	}
}

func (e *Engine) processOrder(o *order.Order) (err error) {
	//start := time.Now()
	//defer TimeConsume(start)
	if o == nil {
		return
	}
	if o.ID <= e.LastOrderID {
		return fmt.Errorf("id is old")
	}
	if o.Symbol != e.Symbol {
		return fmt.Errorf("symbol is error")
	}

	e.LastOrderID = o.ID
	e.LastOrderTime = o.OrderTime
	if o.CancelID != 0 {
		return e.cancel(o)
	}
	if o.IsBuy {
		e.BuyQueue.Insert(o)
		return e.match()
	}

	e.SellQueue.Insert(o)
	return e.match()
}

func (e *Engine) cancel(cancelOrder *order.Order) (err error) {
	var item queue.Item
	if cancelOrder.IsBuy {
		item = e.BuyQueue.Cancel(cancelOrder.CancelID)
	} else {
		item = e.SellQueue.Cancel(cancelOrder.CancelID)
	}
	if item == nil {
		return
	}
	o := item.(*order.Order)
	result := &order.MatchResult{
		CancelID:  o.ID,
		Price:     o.InitialPrice,
		Amount:    o.RemainAmount,
		MatchTime: e.LastOrderTime,
		Symbol:    o.Symbol,
	}
	if cancelOrder.IsBuy {
		result.BuyID = o.ID
		result.BuyUserID = o.UserID
	} else {
		result.SellID = o.ID
		result.SellUserID = o.UserID
	}
	e.MatchResultChan <- result

	return
}

func (e *Engine) match() (err error) {
	for {
		buyItem := e.BuyQueue.First()
		sellItem := e.SellQueue.First()
		if buyItem == nil || sellItem == nil {
			return
		}
		buy := buyItem.(*order.Order)
		sell := sellItem.(*order.Order)

		matchResult := order.Match(e.LastMatchPrice, buy, sell, e.LastOrderTime)

		if buy.RemainAmount == 0 || buy.Canceled {
			e.BuyQueue.Pop()
		}
		if sell.RemainAmount == 0 || sell.Canceled {
			e.SellQueue.Pop()
		}
		if matchResult == nil {
			break
		}
		e.MatchResultChan <- matchResult
	}
	return
}

func TimeConsume(start time.Time) {
	fmt.Printf("cost %s\n", time.Since(start).String())
}
