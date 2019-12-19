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
	LastOrderIndex uint64
	LastOrderID    uint64
	LastOrderTime  uint64
	LastMatchPrice uint64
	Symbol         string
	BuyOrders      queue.PriorityQueue
	SellOrders     queue.PriorityQueue
	BuyQuotations  *order.BookSlice
	SellQuotations *order.BookSlice
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

func NewEngine(symbol string, lastIndex uint64, lastPrice uint64) (engine *Engine, err error) {
	engine = &Engine{
		LastOrderIndex: lastIndex,
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
	index = e.LastOrderIndex
	return
}

func (e *Engine) update(o *order.Order) {
	e.LastOrderIndex = o.OrderIndex
	e.LastOrderID = o.OrderID
	e.LastOrderTime = o.OrderTime
}

func (e *Engine) Match(o *order.Order) (result []*order.Transaction, next bool, err error) {
	if o == nil {
		log.Printf("error: nil order\n")
		return nil, true, nil
	}
	if o.Symbol != e.Symbol {
		log.Printf("index: %d, id: %d, error: symbol: %s != %s\n", o.OrderIndex, o.OrderID, o.Symbol, e.Symbol)
		return nil, true, common.ErrSymbol
	}
	if o.OrderIndex <= e.LastOrderIndex && e.LastOrderIndex != 0 {
		log.Printf("index: %d, id: %d, warn: skip older order index: %d, current orderIndex: %d\n", o.OrderIndex, o.OrderID, o.OrderIndex, e.LastOrderIndex)
		return nil, true, nil
	}
	if o.InitialAmount < 0 {
		log.Printf("index: %d, id: %d, warn: amount<0 , current orderIndex: %d, engineIndex: %d\n", o.OrderIndex, o.OrderID, o.OrderIndex, e.LastOrderIndex)
		return nil, true, nil
	}
	if o.InitialAmount == 0 {
		log.Printf("index: %d, id: %d, warn: amount=0 , current orderIndex: %d, engineIndex: %d\n", o.OrderIndex, o.OrderID, o.OrderIndex, e.LastOrderIndex)
		return nil, true, nil
	}

	e.lock.Lock()
	defer e.lock.Unlock()

	o.RemainAmount = o.InitialAmount
	if o.IsBuy {
		if e.BuyOrders.Search(o.OrderID) {
			log.Printf("warn: skip same index: %d, id: %d\n", o.OrderIndex, o.OrderID)
			return nil, true, nil
		}
		e.BuyOrders.Insert(o)
		if !o.IsMarket {
			e.BuyQuotations.Insert(o.IsBuy, o.InitialPrice, o.InitialAmount)
		}
	} else {
		if e.SellOrders.Search(o.OrderID) {
			log.Printf("warn: skip same index: %d, id: %d\n", o.OrderIndex, o.OrderID)
			return nil, true, nil
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
	e.update(o)
	return result, true, nil
}

//to do 如果取消的订单找不到怎么办，需同样需要返回数据，要不然取消订单会一直等
func (e *Engine) Cancel(cancelOrder *order.Order) (result *order.CancelTransaction, currentOrderID uint64, err error) {
	/*
		if cancelOrder == nil {
			log.Printf("error: nil order\n")
			return nil, e.LastOrderIndex, nil
		}
		if cancelOrder.Symbol != e.Symbol {
			log.Printf("index: %d, id: %d, error: symbol: %s != %s\n", cancelOrder.OrderIndex, cancelOrder.OrderID, cancelOrder.Symbol, e.Symbol)
			return nil, e.LastOrderIndex, common.ErrSymbol
		}
		if cancelOrder.OrderIndex <= e.LastOrderIndex && e.LastOrderIndex != 0 {
			log.Printf("index: %d, id: %d, warn: skip older order index: %d, current orderIndex: %d\n", cancelOrder.OrderIndex, cancelOrder.OrderID, cancelOrder.OrderIndex, e.LastOrderIndex)
			return nil, e.LastOrderIndex, nil
		}

		e.lock.Lock()
		defer e.lock.Unlock()

		log.Printf("cancel order: %v\n", cancelOrder)
		item := e.BuyOrders.Cancel(cancelOrder.CancelOrderID)
		if item == nil {
			item = e.SellOrders.Cancel(cancelOrder.CancelOrderID)
		}
		if item == nil {
			e.update(cancelOrder)
			return nil, e.LastOrderIndex, nil
		}
		o := item.(*order.Order)
		result = &order.CancelTransaction{
			OrderID:       cancelOrder.OrderID,
			CancelOrderID: o.OrderID,
			MatchTime:     cancelOrder.OrderTime,
			Price:         o.InitialPrice,
			Amount:        o.RemainAmount,
			IsBuy:         o.IsBuy,
			Symbol:        o.Symbol,
		}

		if result.IsBuy {
			isExist, err := e.BuyQuotations.SubAmount(true, result.Price, result.Amount)
			if err != nil {
				log.Fatalf("index: %d, id: %d, error: subAmount: %s\n", cancelOrder.OrderIndex, cancelOrder.OrderID, err)
				return nil, e.LastOrderIndex, err
			}
			if !isExist {
				log.Fatalf("index: %d, id: %d, error: subAmount: %s\n", cancelOrder.OrderIndex, cancelOrder.OrderID, err)
				return nil, e.LastOrderIndex, common.ErrNotExist
			}
		} else {
			isExist, err := e.SellQuotations.SubAmount(true, result.Price, result.Amount)
			if err != nil {
				log.Fatalf("index: %d, id: %d, error: subAmount: %s\n", cancelOrder.OrderIndex, cancelOrder.OrderID, err)
				return nil, e.LastOrderIndex, err
			}
			if !isExist {
				log.Fatalf("index: %d, id: %d, error: subAmount: %s\n", cancelOrder.OrderIndex, cancelOrder.OrderID, err)
				return nil, e.LastOrderIndex, common.ErrNotExist
			}
		}
		e.update(cancelOrder)
		return nil, e.LastOrderIndex, nil
	*/
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

func (e *Engine) Quotation() (q *order.OrderBook) {
	e.lock.Lock()
	defer e.lock.Unlock()
	buy := make([]uint64, len(e.BuyQuotations.Data))
	copy(buy, e.BuyQuotations.Data)
	sell := make([]uint64, len(e.SellQuotations.Data))
	copy(sell, e.SellQuotations.Data)
	q = &order.OrderBook{
		Index:              e.LastOrderIndex,
		MatchOrderID:       e.LastOrderID,
		MatchTime:          e.LastOrderTime,
		MatchPrice:         e.LastMatchPrice,
		BuyQuotationSlice:  buy,
		SellQuotationSlice: sell,
	}
	return
}

//to do： how to update orderIndex，if order process err
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

		matchResult := order.Match(e.LastMatchPrice, o, buy, sell)
		if matchResult != nil {
			if buy.RemainAmount < matchResult.Amount || sell.RemainAmount < matchResult.Amount {
				return nil, common.ErrAmount
			}
			buy.RemainAmount -= matchResult.Amount
			sell.RemainAmount -= matchResult.Amount
			result = append(result, matchResult)
			e.LastMatchPrice = matchResult.Price
			if !buy.IsMarket {
				isExist, index, amount := e.BuyQuotations.BinarySearch(true, buy.InitialPrice)
				if !isExist {
					return nil, common.ErrNotExist
				}
				if amount < matchResult.Amount {
					return nil, common.ErrAmount
				}
				e.BuyQuotations.SubAmount(index, matchResult.Amount)
			}
			if !sell.IsMarket {
				isExist, index, amount := e.SellQuotations.BinarySearch(false, sell.InitialPrice)
				if !isExist {
					return nil, common.ErrNotExist
				}
				if amount < matchResult.Amount {
					return nil, common.ErrAmount
				}
				e.SellQuotations.SubAmount(index, matchResult.Amount)
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
	zero.WriteUint64(e.LastOrderIndex)
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
	e.LastOrderIndex, eof = zero.NextUint64()
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

func unSerializeList(data []byte, orders *queue.PriorityList, quotations *order.BookSlice) (err error) {
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
