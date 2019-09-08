package engine

import (
	"fmt"
	"github.com/zhangdaoling/marketmatchengine/order"
	"testing"
	"time"
)

var orders []*order.Order
var result []*order.MatchResult
var testLength = 10000

//to do
func TestCancel(t *testing.T) {
	return
}

func TestEngine(t *testing.T) {
	InitOrders()
	orderChan := make(chan *order.Order, 100)
	resultChan := make(chan *order.MatchResult, 100)
	shutChan := make(chan struct{})
	e, _ := NewEngine(orderChan, resultChan, "usdt/btc", 100, 100, "/Users/zhangdaoling/work/go-work/marketmatchengine")
	go getResult(resultChan)
	go getOrder(orderChan)
	go e.Loop(shutChan)

	time.Sleep(10 * time.Second)
}

func InitOrders() {
	symbol := "usdt/btc"
	orders = make([]*order.Order, 0, 2*testLength)
	length := len(orders)
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
	length = len(orders)
	for i := 0; i < testLength/2; i++ {
		o := &order.Order{
			ID:            uint32(i + length),
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
			ID:            uint32(i + length),
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
			ID:            uint32(i + length),
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
