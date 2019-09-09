package engine

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/zhangdaoling/marketmatchengine/order"
	"github.com/zhangdaoling/marketmatchengine/queue"
	"testing"
	"time"
)

//var orders []*order.Order
var orderIDIndex uint32
var buyOrders, sellOrders []*order.Order
var result []*order.MatchResult

func TestEngine(t *testing.T) {
	orderIDIndex = 1
	initBuy(2)
	initSell(2)
	orderChan := make(chan *order.Order, 100)
	resultChan := make(chan *order.MatchResult, 100)
	shutChan := make(chan struct{})
	e, err := NewEngine(orderChan, resultChan, "usdt/btc", 100, 100, "/Users/zhangdaoling/work/go-work/marketmatchengine/")
	assert.Nil(t, err)
	//consume matchResult
	go getResult(resultChan)
	//consume order
	go getOrder(orderChan)
	//engine work
	go e.Loop(shutChan)

	var xxx int
	xxx = len(buyOrders) + len(sellOrders)/1000
	if xxx < 2 {
		xxx = 2
	}
	time.Sleep(time.Duration(xxx) * time.Second)
	fmt.Println("xxxx")
}

func TestSerialize(t *testing.T) {
	orderIDIndex = 1
	initBuy(2)
	initSell(2)
	orderChan := make(chan *order.Order, 100)
	resultChan := make(chan *order.MatchResult, 100)
	e1, err := NewEngine(orderChan, resultChan, "usdt/btc", 100, 100, "/Users/zhangdaoling/work/go-work/marketmatchengine")
	assert.Nil(t, err)
	e1.BuyQueue = queue.NewPriorityList()
	e1.SellQueue = queue.NewPriorityList()
	for _, o := range buyOrders {
		e1.BuyQueue.Insert(o)
	}
	for _, o := range sellOrders {
		e1.SellQueue.Insert(o)
	}

	zero := e1.Serialize()
	t.Logf("size: %d \n", len(zero.Bytes()))
	e2, _ := NewEngine(orderChan, resultChan, "usdt/btc", 100, 100, "/Users/zhangdaoling/work/go-work/marketmatchengine")
	err = UnSerialize(zero.Bytes(), e2)
	assert.Nil(t, err)
	euqalEngine(t, e1, e2)
}

func TestPersit(t *testing.T) {
	orderIDIndex = 1
	initBuy(2)
	initSell(2)
	orderChan := make(chan *order.Order, 100)
	resultChan := make(chan *order.MatchResult, 100)
	e1, err := NewEngine(orderChan, resultChan, "usdt/btc", 100, 10, "/Users/zhangdaoling/work/go-work/marketmatchengine/")
	assert.Nil(t, err)
	e1.BuyQueue = queue.NewPriorityList()
	e1.SellQueue = queue.NewPriorityList()
	for _, o := range buyOrders {
		e1.BuyQueue.Insert(o)
	}
	for _, o := range sellOrders {
		e1.SellQueue.Insert(o)
	}
	fileName, size, err := e1.Persist()
	t.Logf("file: %s, size: %d \n", fileName, size)
	assert.Nil(t, err)
	e2, err := NewEngineFromFile(10, fileName, "/Users/zhangdaoling/work/go-work/marketmatchengine/")
	assert.Nil(t, err)
	euqalEngine(t, e1, e2)
}

func TestSerializeList(t *testing.T) {
	orderIDIndex = 1
	initBuy(2)
	l1 := queue.NewPriorityList()
	for _, o := range buyOrders {
		l1.Insert(o)
	}
	data := l1.Serialize()
	l2 := queue.NewPriorityList()
	err := unSerializeList(data.Bytes(), l2)
	assert.Nil(t, err)
	equalQueue(t, l1, l2)
}

func euqalEngine(t *testing.T, e1 *Engine, e2 *Engine) {
	assert.Equal(t, e1.LastOrderID, e2.LastOrderID)
	assert.Equal(t, e1.LastMatchPrice, e2.LastMatchPrice)
	assert.Equal(t, e1.LastOrderTime, e2.LastOrderTime)
	assert.Equal(t, e1.Symbol, e2.Symbol)
	equalQueue(t, e1.BuyQueue.(*queue.PriorityList), e2.BuyQueue.(*queue.PriorityList))
	equalQueue(t, e1.SellQueue.(*queue.PriorityList), e2.SellQueue.(*queue.PriorityList))
}

func equalQueue(t *testing.T, l1 *queue.PriorityList, l2 *queue.PriorityList) {
	assert.Equal(t, l1.Len(), l2.Len())
	length := l2.Len()
	for i := 0; uint32(i) < length; i++ {
		e1 := l1.Pop()
		e2 := l2.Pop()
		assert.NotNil(t, e1)
		assert.NotNil(t, e2)
		o1 := e1.(*order.Order)
		o2 := e2.(*order.Order)
		equalOrder(t, o1, o2)
	}
}

func equalOrder(t *testing.T, o1 *order.Order, o2 *order.Order) {
	assert.Equal(t, o1.RemainAmount, o2.RemainAmount)
	assert.Equal(t, o1.ID, o2.ID)
	assert.Equal(t, o1.CancelID, o2.CancelID)
	assert.Equal(t, o1.UserID, o2.UserID)
	assert.Equal(t, o1.OrderTime, o2.OrderTime)
	assert.Equal(t, o1.InitialPrice, o2.InitialPrice)
	assert.Equal(t, o1.InitialAmount, o2.InitialAmount)
	assert.Equal(t, o1.IsMarket, o2.IsMarket)
	assert.Equal(t, o1.IsBuy, o2.IsBuy)
	assert.Equal(t, o1.Canceled, o2.Canceled)
	assert.Equal(t, o1.Symbol, o2.Symbol)
}

//to do
func TestCancel(t *testing.T) {
	return
}

func initBuy(length int) {
	symbol := "usdt/btc"
	buyOrders = make([]*order.Order, 0, length)
	for i := 0; i < length/2; i++ {
		o := &order.Order{
			ID:            orderIDIndex,
			UserID:        100000 + orderIDIndex,
			OrderTime:     2000000 + uint64(orderIDIndex),
			InitialPrice:  1 + 2*uint64(i),
			InitialAmount: 10,
			RemainAmount:  10,
			IsBuy:         true,
			Symbol:        symbol,
		}
		buyOrders = append(buyOrders, o)
		orderIDIndex++
	}
	for i := 0; i < length/2; i++ {
		o := &order.Order{
			ID:            orderIDIndex,
			UserID:        100000 + orderIDIndex,
			OrderTime:     2000000 + uint64(orderIDIndex),
			InitialPrice:  2 + 2*uint64(i),
			InitialAmount: 10,
			RemainAmount:  10,
			IsBuy:         true,
			Symbol:        symbol,
		}
		buyOrders = append(buyOrders, o)
		orderIDIndex++
	}
	return
}

func initSell(length int) (orders []*order.Order) {
	symbol := "usdt/btc"
	sellOrders = make([]*order.Order, 0, length)
	for i := 0; i < length/2; i++ {
		o := &order.Order{
			ID:            orderIDIndex,
			UserID:        100000 + orderIDIndex,
			OrderTime:     2000000 + uint64(orderIDIndex),
			InitialPrice:  1 + 2*uint64(i+length),
			InitialAmount: 10,
			RemainAmount:  10,
			IsBuy:         false,
			Symbol:        symbol,
		}
		sellOrders = append(sellOrders, o)
		orderIDIndex++
	}
	for i := 0; i < length/2; i++ {
		o := &order.Order{
			ID:            orderIDIndex,
			UserID:        100000 + orderIDIndex,
			OrderTime:     2000000 + uint64(orderIDIndex),
			InitialPrice:  2 + 2*uint64(i),
			InitialAmount: 10,
			RemainAmount:  10,
			IsBuy:         false,
			Symbol:        symbol,
		}
		sellOrders = append(sellOrders, o)
		orderIDIndex++
	}
	return
}

func getOrder(orderChan chan *order.Order) {
	for _, o := range buyOrders {
		fmt.Printf("get order: %v\n", o)
		orderChan <- o
	}
	for _, o := range sellOrders {
		fmt.Printf("get order: %v\n", o)
		orderChan <- o
	}
}

func getResult(resultChan chan *order.MatchResult) {
	var i int
	start := time.Now()
	for {
		select {
		case result := <-resultChan:
			i++
			fmt.Printf("count:%d, get result: %v\n", i, result)
			if i%(len(buyOrders)+len(sellOrders)) == 0 {
				fmt.Printf("buy: %d, sell: %d, match: %d, cost:%s second\n", len(buyOrders), len(sellOrders), i, time.Since(start).String())
			}
		}
	}
}
