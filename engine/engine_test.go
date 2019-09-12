package engine

import (
	"fmt"
	"github.com/zhangdaoling/marketmatchengine/order"
	"github.com/zhangdaoling/marketmatchengine/queue"
	"testing"

	"github.com/stretchr/testify/assert"
)

//var orders []*order.Order
var orderIDIndex uint32
var buyOrders, sellOrders []*order.Order
var result []*order.Transaction

func TestMatch(t *testing.T) {
	orderIDIndex = 1
	initBuy(2)
	initSell(2)
	e, err := NewEngine("usdt/btc", 100, 100, 100)
	assert.Nil(t, err)
	for _, o := range buyOrders {
		r := e.Match(o)
		t.Logf("%v\n", r)
	}
	for _, o := range sellOrders {
		r := e.Match(o)
		t.Logf("%v\n", r)
	}
	fmt.Println("xxxx")
}

//to do
func TestCancel(t *testing.T) {
	return
}

func TestPersit(t *testing.T) {
	orderIDIndex = 1
	initBuy(2)
	initSell(2)
	e1, err := NewEngine("usdt/btc", 100, 100, 100)
	assert.Nil(t, err)
	e1.BuyOrders = queue.NewPriorityList()
	e1.SellOrders = queue.NewPriorityList()
	for _, o := range buyOrders {
		e1.BuyOrders.Insert(o)
	}
	for _, o := range sellOrders {
		e1.SellOrders.Insert(o)
	}
	fileName, size, err := e1.Persist("/Users/zhangdaoling/work/go-work/marketmatchengine/")
	t.Logf("file: %s, size: %d \n", fileName, size)
	assert.Nil(t, err)
	e2, err := NewEngineFromFile("/Users/zhangdaoling/work/go-work/marketmatchengine/", fileName)
	assert.Nil(t, err)
	euqalEngine(t, e1, e2)
}

//to do
func TestQuotation(t *testing.T) {
	return
}

func TestSerialize(t *testing.T) {
	orderIDIndex = 1
	initBuy(2)
	initSell(2)
	e1, err := NewEngine("usdt/btc", 100, 100, 100)
	assert.Nil(t, err)
	e1.BuyOrders = queue.NewPriorityList()
	e1.SellOrders = queue.NewPriorityList()
	for _, o := range buyOrders {
		e1.BuyOrders.Insert(o)
	}
	for _, o := range sellOrders {
		e1.SellOrders.Insert(o)
	}

	zero := e1.serialize()
	t.Logf("size: %d \n", len(zero.Bytes()))
	e2, err := NewEngine("usdt/btc", 100, 100, 100)
	err = UnSerialize(zero.Bytes(), e2)
	assert.Nil(t, err)
	euqalEngine(t, e1, e2)
}

func TestSerializeList(t *testing.T) {
	orderIDIndex = 1
	initBuy(2)
	l1 := queue.NewPriorityList()
	q1 := order.NewQuotation(1000)
	for _, o := range buyOrders {
		l1.Insert(o)
		q1.Insert(o)
	}
	data := l1.Serialize()
	l2 := queue.NewPriorityList()
	q2 := order.NewQuotation(1000)
	err := unSerializeList(data.Bytes(), l2, q2)
	assert.Nil(t, err)
	equalQueue(t, l1, l2)
}

func euqalEngine(t *testing.T, e1 *Engine, e2 *Engine) {
	assert.Equal(t, e1.LastOrderID, e2.LastOrderID)
	assert.Equal(t, e1.LastMatchPrice, e2.LastMatchPrice)
	assert.Equal(t, e1.LastOrderTime, e2.LastOrderTime)
	assert.Equal(t, e1.Symbol, e2.Symbol)
	equalQueue(t, e1.BuyOrders.(*queue.PriorityList), e2.BuyOrders.(*queue.PriorityList))
	equalQueue(t, e1.SellOrders.(*queue.PriorityList), e2.SellOrders.(*queue.PriorityList))
	equalBytes(t, e1.BuyQuotations, e2.BuyQuotations)
	equalBytes(t, e1.SellQuotations, e2.SellQuotations)
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

func equalBytes(t *testing.T, b1 []byte, b2[]byte){
	assert.Equal(t, len(b1), len(b2), "b1 len: %d, b2 len: %d", len(b1), len(b2))
	for i := 0; i < len(b1); i++ {
		assert.Equal(t, b1[i], b2[i])
	}
}

func equalOrder(t *testing.T, o1 *order.Order, o2 *order.Order) {
	assert.Equal(t, o1.RemainAmount, o2.RemainAmount)
	assert.Equal(t, o1.ID, o2.ID)
	assert.Equal(t, o1.CancelOrderID, o2.CancelOrderID)
	assert.Equal(t, o1.OrderTime, o2.OrderTime)
	assert.Equal(t, o1.InitialPrice, o2.InitialPrice)
	assert.Equal(t, o1.InitialAmount, o2.InitialAmount)
	assert.Equal(t, o1.IsMarket, o2.IsMarket)
	assert.Equal(t, o1.IsBuy, o2.IsBuy)
	assert.Equal(t, o1.Symbol, o2.Symbol)
}

func initBuy(length int) {
	symbol := "usdt/btc"
	buyOrders = make([]*order.Order, 0, length)
	for i := 0; i < length/2; i++ {
		o := &order.Order{
			ID:            orderIDIndex,
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
