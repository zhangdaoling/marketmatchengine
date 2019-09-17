package engine

import (
	"crypto/md5"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/zhangdaoling/marketmatchengine/common"
	"github.com/zhangdaoling/marketmatchengine/order"
	"github.com/zhangdaoling/marketmatchengine/queue"
)

type Engine struct {
	LastIndex      uint64
	LastIndexTime  uint64
	LastOrderID    uint64
	LastOrderTime  uint64
	LastMatchPrice uint64
	Symbol         string
	BuyOrders      queue.PriorityQueue
	SellOrders     queue.PriorityQueue
	BuyQuotations  *order.QuotationSlice
	SellQuotations *order.QuotationSlice
	CheckSum       []byte
	lock           sync.Mutex
}

func NewEngineFromFile(fileName string, path string) (e *Engine, err error) {
	data, err := ioutil.ReadFile(fileName + path)
	if err != nil {
		log.Fatalf("error: read engine err: %s", err)
		return
	}
	e = &Engine{}
	err = UnSerialize(data, e)
	log.Fatalf("error: UnSerialize engine err: %s", err)
	return
}

func NewEngine(symbol string, lastIndex uint64, lastIndexTime uint64, lastPrice uint64) (engine *Engine, err error) {
	engine = &Engine{
		LastIndex:      lastIndex,
		LastIndexTime:  lastIndexTime,
		LastMatchPrice: lastPrice,
		Symbol:         symbol,
		BuyOrders:      queue.NewPriorityList(),
		SellOrders:     queue.NewPriorityList(),
		BuyQuotations:  order.NewQuotation(1024),
		SellQuotations: order.NewQuotation(1024),
	}
	return engine, nil
}

func (e *Engine) GetIndex() (index uint64) {
	e.lock.Lock()
	defer e.lock.Unlock()
	index = e.LastIndex
	return
}

func (e *Engine) Match(o *order.Order) (result []*order.Transaction, success bool, err error) {
	if o == nil {
		log.Printf("error: nil order\n")
		return nil, false, nil
	}
	if o.Symbol != e.Symbol {
		log.Printf("index: %d, id: %d, error: symbol: %s != %s\n", o.Index, o.OrderID, o.Symbol, e.Symbol)
		return nil, false, common.ErrSymbol
	}
	if o.Index <= e.LastIndex && e.LastIndex != 0 {
		log.Printf("index: %d, id: %d, warn: skip older order index: %d, current orderIndex: %d\n", o.Index, o.OrderID, o.Index, e.LastIndex)
		return nil, false, nil
	}

	e.lock.Lock()
	defer e.lock.Unlock()

	log.Printf("match order: %v\n", o)
	o.RemainAmount = o.InitialAmount
	if o.IsBuy {
		if e.BuyOrders.Search(o.OrderID) {
			log.Printf("warn: skip same index: %d, id: %d\n", o.Index, o.OrderID)
			return nil, false, nil
		}
		e.BuyOrders.Insert(o)
		if !o.IsMarket {
			e.BuyQuotations.Insert(o.IsBuy, o.InitialPrice, o.InitialAmount)
		}
	} else {
		if e.SellOrders.Search(o.OrderID) {
			log.Printf("warn: skip same index: %d, id: %d\n", o.Index, o.OrderID)
			return nil, false, nil
		}
		e.SellOrders.Insert(o)
		if !o.IsMarket {
			e.SellQuotations.Insert(o.IsBuy, o.InitialPrice, o.InitialAmount)
		}
	}
	result, err = e.match(o)
	if err != nil {
		log.Fatalf("error: engine match err: %s\n", err)
		return nil, false, err
	}
	e.LastIndex = o.Index
	e.LastIndexTime = o.OrderTime
	e.LastOrderID = o.Index
	e.LastOrderTime = o.OrderTime
	return result, true, nil
}

//to do 如果取消的订单找不到怎么办，需同样需要返回数据，要不然取消订单会一直等
func (e *Engine) Cancel(cancelOrder *order.Order) (result *order.CancelOrder, currentOrderID uint64, err error) {
	if cancelOrder == nil {
		log.Printf("error: nil order\n")
		return nil, e.LastIndex, nil
	}
	if cancelOrder.Symbol != e.Symbol {
		log.Printf("index: %d, id: %d, error: symbol: %s != %s\n", cancelOrder.Index, cancelOrder.OrderID, cancelOrder.Symbol, e.Symbol)
		return nil, e.LastIndex, common.ErrSymbol
	}
	if cancelOrder.Index <= e.LastIndex && e.LastIndex != 0 {
		log.Printf("index: %d, id: %d, warn: skip older order index: %d, current orderIndex: %d\n", cancelOrder.Index, cancelOrder.OrderID, cancelOrder.Index, e.LastIndex)
		return nil, e.LastIndex, nil
	}

	e.lock.Lock()
	defer e.lock.Unlock()

	log.Printf("cancel order: %v\n", cancelOrder)
	item := e.BuyOrders.Cancel(cancelOrder.CancelOrderID)
	if item == nil {
		item = e.SellOrders.Cancel(cancelOrder.CancelOrderID)
	}
	if item == nil {
		e.LastIndex = cancelOrder.Index
		e.LastIndexTime = cancelOrder.IndexTime
		e.LastOrderID = cancelOrder.OrderID
		e.LastOrderTime = cancelOrder.OrderTime
		return nil, e.LastIndex, nil
	}
	o := item.(*order.Order)
	result = &order.CancelOrder{
		OrderID:       cancelOrder.OrderID,
		CancelOrderID: o.OrderID,
		MatchTime:     cancelOrder.OrderTime,
		Price:         o.InitialPrice,
		Amount:        o.RemainAmount,
		IsBuy:         o.IsBuy,
		Symbol:        o.Symbol,
	}

	if result.IsBuy {
		isExist, err := e.BuyQuotations.SubAmount(result.Price, result.Amount, true)
		if err != nil {
			log.Fatalf("index: %d, id: %d, error: subAmount: %s\n", cancelOrder.Index, cancelOrder.OrderID, err)
			return nil, e.LastIndex, err
		}
		if !isExist {
			log.Fatalf("index: %d, id: %d, error: subAmount: %s\n", cancelOrder.Index, cancelOrder.OrderID, err)
			return nil, e.LastIndex, common.ErrNotExist
		}
	} else {
		isExist, err := e.SellQuotations.SubAmount(result.Price, result.Amount, false)
		if err != nil {
			log.Fatalf("index: %d, id: %d, error: subAmount: %s\n", cancelOrder.Index, cancelOrder.OrderID, err)
			return nil, e.LastIndex, err
		}
		if !isExist {
			log.Fatalf("index: %d, id: %d, error: subAmount: %s\n", cancelOrder.Index, cancelOrder.OrderID, err)
			return nil, e.LastIndex, common.ErrNotExist
		}
	}
	e.LastIndex = cancelOrder.Index
	e.LastIndexTime = cancelOrder.IndexTime
	e.LastOrderID = cancelOrder.OrderID
	e.LastOrderTime = cancelOrder.OrderTime
	return nil, e.LastIndex, nil
}

func (e *Engine) Persist(path string) (fileName string, size int, err error) {
	e.lock.Lock()
	defer e.lock.Unlock()
	//start := time.Now()
	//defer common.TimeConsume(start)
	fileName = "engine_" + time.Now().Format(time.RFC3339) + ".binary"
	f, err := os.OpenFile(path+fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Printf("error: create file: %s, err: %s\n", path+fileName, err)
		return
	}
	defer f.Close()

	data := e.serialize()
	size, err = f.Write(data.Bytes())
	if err != nil {
		log.Printf("error: write file: %s, err: %s\n", path+fileName, err)
		return
	}
	if size != len(data.Bytes()) {
	}
	return fileName, size, err
}

func (e *Engine) Quotation() (q *order.Quotation) {
	e.lock.Lock()
	defer e.lock.Unlock()
	buy := make([]uint64, len(e.BuyQuotations.Data))
	copy(buy, e.BuyQuotations.Data)
	sell := make([]uint64, len(e.SellQuotations.Data))
	copy(sell, e.BuyQuotations.Data)
	q = &order.Quotation{
		Time:               e.LastIndexTime,
		BuyQuotationSlice:  buy,
		SellQuotationSlice: sell,
	}
	return
}

func (e *Engine) match(o *order.Order) (result []*order.Transaction, err error) {
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
				isExist, err := e.BuyQuotations.SubAmount(buy.InitialPrice, matchResult.Amount, true)
				if err != nil {
					return nil, err
				}
				if !isExist {
					return nil, common.ErrNotExist
				}
			}
			if !sell.IsMarket {
				isExist, err := e.SellQuotations.SubAmount(sell.InitialPrice, matchResult.Amount, false)
				if err != nil {
					return nil, err
				}
				if !isExist {
					return nil, common.ErrNotExist
				}
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
	zero.WriteUint64(e.LastIndex)
	zero.WriteUint64(e.LastIndexTime)
	zero.WriteUint64(e.LastOrderID)
	zero.WriteUint64(e.LastOrderTime)
	zero.WriteUint64(e.LastMatchPrice)
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
	e.LastIndex, eof = zero.NextUint64()
	if eof {
		return common.ErrTooLarge
	}
	e.LastIndexTime, eof = zero.NextUint64()
	if eof {
		return common.ErrTooLarge
	}
	e.LastOrderID, eof = zero.NextUint64()
	if eof {
		return common.ErrTooLarge
	}
	e.LastOrderTime, eof = zero.NextUint64()
	if eof {
		return common.ErrTooLarge
	}
	e.LastMatchPrice, eof = zero.NextUint64()
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

func unSerializeList(data []byte, orders *queue.PriorityList, quotations *order.QuotationSlice) (err error) {
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
		if !o.IsMarket {
			quotations.Insert(o.IsBuy, o.InitialPrice, o.InitialAmount)
		}
	}
	return
}
