package engine

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/zhangdaoling/marketmatchengine/order"
	"github.com/zhangdaoling/marketmatchengine/queue"
	"testing"
	"time"
)

var result []*order.MatchResult
var orders []*order.Order
var testLength int

func TestEngine(t *testing.T) {
	initBuyAndSell(2)
	orderChan := make(chan *order.Order, 100)
	resultChan := make(chan *order.MatchResult, 100)
	shutChan := make(chan struct{})
	e, _ := NewEngine(orderChan, resultChan, "usdt/btc", 100, 100, "/Users/zhangdaoling/work/go-work/marketmatchengine")
	go getResult(resultChan)
	go getOrder(orderChan)
	go e.Loop(shutChan)

	var xxx int
	if testLength/1000 < 2 {
		xxx = 2
	}
	time.Sleep(time.Duration(xxx) * time.Second)
	fmt.Println("dfdsfdsfd")
}

func TestSerializeEngine(t *testing.T) {
	initBuy(2)
	orderChan := make(chan *order.Order, 100)
	resultChan := make(chan *order.MatchResult, 100)
	shutChan := make(chan struct{})
	e1, _ := NewEngine(orderChan, resultChan, "usdt/btc", 100, 100, "/Users/zhangdaoling/work/go-work/marketmatchengine")
	go getResult(resultChan)
	go getOrder(orderChan)
	go e1.Loop(shutChan)
	var xxx int
	if testLength/1000 < 2 {
		xxx = 2
	}
	time.Sleep(time.Duration(xxx) * time.Second)
	zero := e1.Serialize()
	e2, _ := NewEngine(orderChan, resultChan, "usdt/btc", 100, 100, "/Users/zhangdaoling/work/go-work/marketmatchengine")
	err := UnSerialize(zero.Bytes(), e2)
	assert.Nil(t, err)
	assert.Equal(t, e1.LastOrderID, e2.LastOrderID)
	assert.Equal(t, e1.LastMatchPrice, e2.LastMatchPrice)
	assert.Equal(t, e1.LastOrderTime, e2.LastOrderTime)
	assert.Equal(t, e1.Symbol, e2.Symbol)
	fmt.Println("1111111111")
	checkQueue(t, e1.BuyQueue.(*queue.PriorityList), e2.BuyQueue.(*queue.PriorityList))
	fmt.Println("2222222222")
	checkQueue(t, e1.SellQueue.(*queue.PriorityList), e2.SellQueue.(*queue.PriorityList))
}

func TestSerializeList(t *testing.T) {
	initBuy(10)
	l1 := queue.NewPriorityList()
	l2 := queue.NewPriorityList()
	for _, o := range orders {
		l1.Insert(o)
	}
	data := l1.Serialize()
	err := unSerializeList(data.Bytes(), l2)
	assert.Nil(t, err)
	checkQueue(t, l1, l2)
}

func checkQueue(t *testing.T, l1 *queue.PriorityList, l2 *queue.PriorityList) {
	assert.Equal(t, l1.Len(), l2.Len())
	for i := 0; uint32(i) < l2.Len(); i++ {
		e1 := l1.Pop()
		e2 := l2.Pop()
		assert.NotNil(t, e1)
		assert.NotNil(t, e2)
		o1 := e1.(*order.Order)
		o2 := e2.(*order.Order)
		fmt.Println("o1: ", o1)
		fmt.Println("o2: ", o2)
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
}

//to do
func TestCancel(t *testing.T) {
	return
}

func initBuy(length int) (orders []*order.Order) {
	symbol := "usdt/btc"
	orders = make([]*order.Order, 0, 2*length)
	for i := 0; i < testLength/2; i++ {
		o := &order.Order{
			ID:            uint32(i + length),
			UserID:        10000 + uint32(i+length),
			OrderTime:     200000 + uint64(i+length),
			InitialPrice:  1 + 2*uint64(i),
			InitialAmount: 10,
			RemainAmount:  10,
			IsBuy:         true,
			Symbol:        symbol,
		}
		orders = append(orders, o)
	}
	return
}

func initBuyAndSell(l int) {
	testLength = l
	symbol := "usdt/btc"
	orders = make([]*order.Order, 0, 2*testLength)
	length := len(orders)
	for i := 0; i < testLength/2; i++ {
		o := &order.Order{
			ID:            1 + uint32(i+length),
			UserID:        10000 + uint32(i+length),
			OrderTime:     200000 + uint64(i+length),
			InitialPrice:  1 + 2*uint64(i),
			InitialAmount: 10,
			RemainAmount:  10,
			IsBuy:         true,
			Symbol:        symbol,
		}
		orders = append(orders, o)
	}
	length = len(orders)
	for i := 0; i < testLength/2; i++ {
		o := &order.Order{
			ID:            1 + uint32(i+length),
			UserID:        10000 + uint32(i+length),
			OrderTime:     200000 + uint64(i+length),
			InitialPrice:  2 + 2*uint64(i),
			InitialAmount: 10,
			RemainAmount:  10,
			IsBuy:         true,
			Symbol:        symbol,
		}
		orders = append(orders, o)
	}

	length = len(orders)
	for i := 0; i < testLength/2; i++ {
		o := &order.Order{
			ID:            1 + uint32(i+length),
			UserID:        10000 + uint32(i+length),
			OrderTime:     200000 + uint64(i+length),
			InitialPrice:  1 + 2*uint64(i+1000),
			InitialAmount: 10,
			RemainAmount:  10,
			Symbol:        symbol,
		}
		orders = append(orders, o)
	}
	length = len(orders)
	for i := 0; i < testLength/2; i++ {
		o := &order.Order{
			ID:            1 + uint32(i+length),
			UserID:        10000 + uint32(i+length),
			OrderTime:     200000 + uint64(i+length),
			InitialPrice:  uint64(testLength - i),
			InitialAmount: 10,
			RemainAmount:  10,
			Symbol:        symbol,
		}
		orders = append(orders, o)
	}
}

func getOrder(orderChan chan *order.Order) {
	for _, o := range orders {
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
			if i%(testLength/2) == 0 {
				fmt.Printf("buy: %d, sell: %d, match: %d, cost:%s second\n", testLength, testLength, testLength, time.Since(start).String())
			}
		}
	}
}
