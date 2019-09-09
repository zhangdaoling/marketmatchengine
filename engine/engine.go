package engine

import (
	"crypto/md5"
	"fmt"
	"github.com/zhangdaoling/marketmatchengine/common"
	"github.com/zhangdaoling/marketmatchengine/order"
	"github.com/zhangdaoling/marketmatchengine/queue"
	"time"
)

type Engine struct {
	OrderChan       chan *order.Order
	MatchResultChan chan *order.MatchResult
	LastOrderID     uint32
	LastMatchPrice  uint64
	LastOrderTime   uint64
	Symbol          string
	BuyQueue        queue.PriorityQueue
	SellQueue       queue.PriorityQueue
	CheckSum        []byte
	PersistTime     uint64
	PersistPath     string
}

//to do
func NewEngineFromFile(engineFile string, orderChan chan *order.Order, matchResultChan chan *order.MatchResult) (engine *Engine, err error) {
	return
}

func NewEngine(orderChan chan *order.Order, matchResultChan chan *order.MatchResult, symbol string, lastPrice uint64, persistTime uint64, persistPath string) (engine *Engine, err error) {
	sellQueue := queue.NewPriorityList()
	buyQueue := queue.NewPriorityList()
	engine = &Engine{
		OrderChan:       orderChan,
		MatchResultChan: matchResultChan,
		LastMatchPrice:  lastPrice,
		Symbol:          symbol,
		BuyQueue:        buyQueue,
		SellQueue:       sellQueue,
		PersistTime:     persistTime,
		PersistPath:     persistPath,
	}
	return engine, nil
}

func (e *Engine) Loop(shutdown chan struct{}) {
	timer := time.NewTimer(100 * time.Second)
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
	//fmt.Println("process order: ", o)
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

//cancel order
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

func (e *Engine) Serialize() (zero *common.ZeroCopySink) {
	zero = common.NewZeroCopySink(nil, 64*int(e.BuyQueue.Len()+e.SellQueue.Len()))
	zero.WriteUint32(e.LastOrderID)
	zero.WriteUint64(e.LastMatchPrice)
	zero.WriteUint64(e.LastOrderTime)
	zero.WriteString(e.Symbol)
	buyData := e.BuyQueue.Serialize()
	zero.WriteVarBytes(buyData.Bytes())
	sellData := e.SellQueue.Serialize()
	zero.WriteVarBytes(sellData.Bytes())
	sum := md5.Sum(zero.Bytes())
	zero.WriteVarBytes(sum[:])
	return
}

func UnSerialize(data []byte, e *Engine) (err error) {
	var irregular, eof bool
	zero := common.NewZeroCopySource(data)
	e.LastOrderID, eof = zero.NextUint32()
	if eof {
		return common.ErrTooLarge
	}
	e.LastMatchPrice, eof = zero.NextUint64()
	if eof {
		return common.ErrTooLarge
	}
	e.LastOrderTime, eof = zero.NextUint64()
	if eof {
		return common.ErrTooLarge
	}
	e.Symbol, _, irregular, eof = zero.NextString()
	if irregular {
		return common.ErrIrregularData
	}
	if eof {
		return common.ErrUnexpectedEOF
	}

	var buyBytes, sellBytes []byte
	buyBytes, _, irregular, eof = zero.NextVarBytes()
	if irregular {
		return common.ErrIrregularData
	}
	if eof {
		return common.ErrTooLarge
	}
	err = unSerializeList(buyBytes, e.BuyQueue.(*queue.PriorityList))
	if err != nil {
		return
	}
	sellBytes, _, irregular, eof = zero.NextVarBytes()
	if irregular {
		return common.ErrIrregularData
	}
	if eof {
		return common.ErrTooLarge
	}
	err = unSerializeList(sellBytes, e.SellQueue.(*queue.PriorityList))
	if err != nil {
		return
	}

	//calculate check sum
	dataSize := zero.Size() - zero.Len()
	zero.Skip(dataSize)
	dataByte, eof := zero.NextBytes(dataSize)
	if eof {
		return common.ErrTooLarge
	}
	checkSum := md5.Sum(dataByte)
	//get check sum
	e.CheckSum, _, irregular, eof = zero.NextVarBytes()
	if irregular {
		return common.ErrIrregularData
	}
	if eof {
		return common.ErrTooLarge
	}
	if isByteSame(e.CheckSum, checkSum[:]){
		return common.ErrEngineCheckSum
	}
	return
}

func unSerializeList(data []byte, p *queue.PriorityList) (err error) {
	var eof, irregular bool
	var listType string
	var count uint32
	var o = &order.Order{}
	zero := common.NewZeroCopySource(data)
	listType, _, irregular, eof = zero.NextString()
	if irregular {
		return common.ErrIrregularData
	}
	if eof {
		return common.ErrUnexpectedEOF
	}
	if listType != queue.List_Queue_Name {
		return common.ErrQueueType
	}
	fmt.Println(listType)
	count, eof = zero.NextUint32()
	if eof {
		return common.ErrUnexpectedEOF
	}
	var orderBytes []byte
	for i := 0; uint32(i) < count; i++ {
		orderBytes, _, irregular, eof = zero.NextVarBytes()
		err = order.UnSerialize(orderBytes, o)
		if err != nil {
			return
		}
		p.Insert(o)
	}
	return
}

func isByteSame(data1 []byte, data2 []byte) bool{
	if len(data1) != len(data2){
		return false
	}
	for i:=0;i<len(data1);i++{
		if data1[i] != data2[2]{
			return false
		}
	}
	return true
}



