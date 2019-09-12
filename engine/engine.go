package engine

import (
	"crypto/md5"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/zhangdaoling/marketmatchengine/common"
	"github.com/zhangdaoling/marketmatchengine/order"
	"github.com/zhangdaoling/marketmatchengine/queue"
)

type Engine struct {
	LastOrderID    uint32
	LastMatchPrice uint64
	LastOrderTime  uint64
	Symbol         string
	BuyOrders      queue.PriorityQueue
	SellOrders     queue.PriorityQueue
	BuyQuotations  order.Quotation
	SellQuotations order.Quotation
	CheckSum       []byte
	lock           sync.Mutex
}

func NewEngineFromFile(fileName string, path string) (e *Engine, err error) {
	data, err := ioutil.ReadFile(fileName + path)
	if err != nil {
		return
	}
	e = &Engine{}
	err = UnSerialize(data, e)
	return
}

func NewEngine(symbol string, lastOrderID uint32, lastPrice uint64, lastOrderTime uint64) (engine *Engine, err error) {
	engine = &Engine{
		LastMatchPrice: lastPrice,
		Symbol:         symbol,
		BuyOrders:      queue.NewPriorityList(),
		SellOrders:     queue.NewPriorityList(),
		BuyQuotations:  order.NewQuotation(1000),
		SellQuotations: order.NewQuotation(1000),
	}
	return engine, nil
}

func (e *Engine) Match(o *order.Order) (reuslt []*order.Transaction) {
	e.lock.Lock()
	defer e.lock.Unlock()
	if o == nil {
		return
	}
	if o.ID <= e.LastOrderID {
		return
	}
	if o.Symbol != e.Symbol {
		return
	}

	e.LastOrderID = o.ID
	e.LastOrderTime = o.OrderTime

	if o.IsBuy {
		e.BuyOrders.Insert(o)
		e.BuyQuotations.Insert(o)
	} else {
		e.SellOrders.Insert(o)
		e.SellQuotations.Insert(o)
	}
	return e.match(o)
}

func (e *Engine) Cancel(cancelOrder *order.Order) (result *order.CancelOrder) {
	e.lock.Lock()
	defer e.lock.Unlock()
	if cancelOrder == nil {
		return
	}
	if cancelOrder.ID <= e.LastOrderID {
		return
	}
	if cancelOrder.Symbol != e.Symbol {
		return
	}

	e.LastOrderID = cancelOrder.ID
	e.LastOrderTime = cancelOrder.OrderTime

	item := e.BuyOrders.Cancel(cancelOrder.CancelOrderID)
	if item == nil {
		item = e.SellOrders.Cancel(cancelOrder.CancelOrderID)
	}
	if item == nil {
		return
	}
	o := item.(*order.Order)
	result = &order.CancelOrder{
		ID:            cancelOrder.ID,
		CancelOrderID: o.ID,
		MatchTime:     cancelOrder.OrderTime,
		Price:         o.InitialPrice,
		Amount:        o.RemainAmount,
		IsBuy:         o.IsBuy,
		Symbol:        o.Symbol,
	}

	if result.IsBuy {
		e.BuyQuotations.SubAmount(result.Price, result.Amount)
	} else {
		e.SellQuotations.SubAmount(result.Price, result.Amount)
	}

	return
}

func (e *Engine) Persist(path string) (fileName string, size int, err error) {
	e.lock.Lock()
	defer e.lock.Unlock()
	//start := time.Now()
	//defer common.TimeConsume(start)
	fileName = "engine_" + time.Now().Format(time.RFC3339) + ".binary"
	f, err := os.OpenFile(path+fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return
	}
	defer f.Close()

	data := e.serialize()
	size, err = f.Write(data.Bytes())
	if err != nil {
		return
	}
	if size != len(data.Bytes()) {
	}
	err = f.Sync()
	if err != nil {
		return
	}
	return fileName, size, err
}

func (e *Engine) Quotation()(data []byte){
	e.lock.Lock()
	defer e.lock.Unlock()
	zero := common.NewZeroCopySink(nil, len(e.BuyQuotations)+len(e.SellQuotations))
	zero.WriteUint32(uint32(len(e.BuyQuotations)))
	zero.WriteBytes(e.BuyQuotations)
	zero.WriteUint32(uint32(len(e.SellQuotations)))
	zero.WriteBytes(e.SellQuotations)
	return zero.Bytes()
}

func (e *Engine) match(o *order.Order) (result []*order.Transaction) {
	result = make([]*order.Transaction, 0, 2)
	for {
		buyItem := e.BuyOrders.First()
		sellItem := e.SellOrders.First()
		if buyItem == nil || sellItem == nil {
			return
		}
		buy := buyItem.(*order.Order)
		sell := sellItem.(*order.Order)

		matchResult := order.Match(e.LastMatchPrice, e.LastOrderTime, buy, sell, o.IsBuy)
		if matchResult != nil {
			result = append(result, matchResult)
			if !buy.IsMarket {
				e.BuyQuotations.SubAmount(matchResult.Price, matchResult.Amount)
			}
			if !sell.IsMarket {
				e.SellQuotations.SubAmount(matchResult.Price, matchResult.Amount)
			}
		}

		if buy.RemainAmount == 0 {
			e.BuyOrders.Pop()
		}
		if sell.RemainAmount == 0 {
			e.SellOrders.Pop()
		}
		if matchResult == nil {
			break
		}
	}
	return
}

func (e *Engine) serialize() (zero *common.ZeroCopySink) {
	zero = common.NewZeroCopySink(nil, 64*int(e.BuyOrders.Len()+e.SellOrders.Len()))
	zero.WriteUint32(e.LastOrderID)
	zero.WriteUint64(e.LastMatchPrice)
	zero.WriteUint64(e.LastOrderTime)
	zero.WriteString(e.Symbol)
	buyData := e.BuyOrders.Serialize()
	zero.WriteVarBytes(buyData.Bytes())
	sellData := e.SellOrders.Serialize()
	zero.WriteVarBytes(sellData.Bytes())
	//zero.WriteVarBytes(e.BuyQuotations)
	//zero.WriteVarBytes(e.SellQuotations)
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
	e.BuyOrders = queue.NewPriorityList()
	e.BuyQuotations = order.NewQuotation(1000)
	err = unSerializeList(buyBytes, e.BuyOrders.(*queue.PriorityList), e.BuyQuotations)
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
	e.SellOrders = queue.NewPriorityList()
	e.SellQuotations = order.NewQuotation(1000)
	err = unSerializeList(sellBytes, e.SellOrders.(*queue.PriorityList), e.SellQuotations)
	if err != nil {
		return
	}

	//calculate check sum
	dataSize := zero.Pos()
	zero.BackUp(dataSize)
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
	if common.IsByteSame(e.CheckSum, checkSum[:]) {
		return common.ErrEngineCheckSum
	}
	return
}

func unSerializeList(data []byte, orders *queue.PriorityList, quotations order.Quotation) (err error) {
	var eof, irregular bool
	var listType string
	var count uint32
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
	count, eof = zero.NextUint32()
	if eof {
		return common.ErrUnexpectedEOF
	}
	var orderBytes []byte
	for i := 0; uint32(i) < count; i++ {
		orderBytes, _, irregular, eof = zero.NextVarBytes()
		var o = &order.Order{}
		err = order.UnSerialize(orderBytes, o)
		if err != nil {
			return
		}
		orders.Insert(o)
		quotations.Insert(o)
	}
	return
}
