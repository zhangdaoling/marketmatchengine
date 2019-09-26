package engine

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/zhangdaoling/marketmatchengine/order"
	"github.com/zhangdaoling/marketmatchengine/queue"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//var orders []*order.Order
var orderIDIndex uint64
var buyOrders, sellOrders []*order.Order
var result []*order.Transaction

//to do, not complete
func TestMatch(t *testing.T) {
	orderIDIndex = 1
	initBuy(10)
	initSell(10)
	e, err := NewEngine("A-B", 0, 0, 100)
	assert.Nil(t, err)
	t.Logf("%v\n", time.Now())
	for _, o := range buyOrders {
		_, _, err := e.Match(o)
		assert.Nil(t, err)
		//t.Logf("%v\n", r)
	}
	for _, o := range sellOrders {
		_, _, err := e.Match(o)
		assert.Nil(t, err)
		//t.Logf("%v\n", r)
	}
	t.Logf("%v\n", time.Now())
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
	e1, err := NewEngine("A-B", 100, 100, 100)
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
	e1, err := NewEngine("A-B", 100, 100, 100)
	assert.Nil(t, err)
	e1.BuyOrders = queue.NewPriorityList()
	e1.SellOrders = queue.NewPriorityList()
	for _, o := range buyOrders {
		e1.BuyOrders.Insert(o)
		e1.BuyQuotations.Insert(o.IsBuy, o.InitialPrice, o.InitialAmount)
	}
	for _, o := range sellOrders {
		e1.SellOrders.Insert(o)
		e1.SellQuotations.Insert(o.IsBuy, o.InitialPrice, o.InitialAmount)
	}

	zero := e1.serialize()
	t.Logf("size: %d \n", len(zero.Bytes()))
	e2, err := NewEngine("A-B", 100, 100, 100)
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
		q1.Insert(o.IsBuy, o.InitialPrice, o.InitialAmount)
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
	equalArray(t, e1.BuyQuotations.Data, e2.BuyQuotations.Data)
	equalArray(t, e1.SellQuotations.Data, e2.SellQuotations.Data)
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

func equalArray(t *testing.T, b1 []uint64, b2 []uint64) {
	assert.Equal(t, len(b1), len(b2), "b1 len: %d, b2 len: %d", len(b1), len(b2))
	for i := 0; i < len(b1); i++ {
		assert.Equal(t, b1[i], b2[i])
	}
}

func equalOrder(t *testing.T, o1 *order.Order, o2 *order.Order) {
	assert.Equal(t, o1.RemainAmount, o2.RemainAmount)
	assert.Equal(t, o1.Index, o2.Index)
	assert.Equal(t, o1.IndexTime, o2.IndexTime)
	assert.Equal(t, o1.OrderID, o2.OrderID)
	assert.Equal(t, o1.OrderTime, o2.OrderTime)
	assert.Equal(t, o1.UserID, o2.UserID)
	assert.Equal(t, o1.InitialPrice, o2.InitialPrice)
	assert.Equal(t, o1.InitialAmount, o2.InitialAmount)
	assert.Equal(t, o1.CancelOrderID, o2.CancelOrderID)
	assert.Equal(t, o1.IsMarket, o2.IsMarket)
	assert.Equal(t, o1.IsBuy, o2.IsBuy)
	assert.Equal(t, o1.Symbol, o2.Symbol)
}

func initBuy(length int) {
	symbol := "A-B"
	buyOrders = make([]*order.Order, 0, length)
	for i := 0; i < length/2; i++ {
		o := &order.Order{
			Index:         orderIDIndex,
			OrderID:       orderIDIndex,
			OrderTime:     2000000 + uint64(orderIDIndex),
			UserID:        orderIDIndex,
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
			Index:         orderIDIndex,
			OrderID:       orderIDIndex,
			OrderTime:     2000000 + uint64(orderIDIndex),
			UserID:        orderIDIndex,
			InitialPrice:  uint64(length - i),
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
	symbol := "A-B"
	sellOrders = make([]*order.Order, 0, length)
	for i := 0; i < length/2; i++ {
		o := &order.Order{
			Index:         orderIDIndex,
			OrderID:       orderIDIndex,
			OrderTime:     2000000 + uint64(orderIDIndex),
			UserID:        orderIDIndex,
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
			Index:         orderIDIndex,
			OrderID:       orderIDIndex,
			OrderTime:     2000000 + uint64(orderIDIndex),
			UserID:        orderIDIndex,
			InitialPrice:  2 + 2*uint64(length/2-i),
			InitialAmount: 5,
			RemainAmount:  5,
			IsBuy:         false,
			Symbol:        symbol,
		}
		sellOrders = append(sellOrders, o)
		orderIDIndex++
	}
	return
}

//push order to kafka
func TestData(t *testing.T) {
	orderIDIndex = 1
	initBuy(50000)
	initSell(50000)
	buyOrders = append(buyOrders, sellOrders...)

	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	var b []byte
	for _, o := range buyOrders {
		b, err = json.Marshal(o)
		if err != nil {
			log.Printf("json marshal error: %v\n", err)
			return
		}
		msg := &sarama.ProducerMessage{
			Topic:     "order_A-B",
			Value:     sarama.ByteEncoder(b),
			Timestamp: time.Now(),
		}
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.Printf("FAILED to send message: %s\n", err)
		} else {
			log.Printf("> message sent to partition %d at offset %d\n", partition, offset)
		}
	}
}
