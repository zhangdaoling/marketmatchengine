package order

import (
	"fmt"

	"github.com/zhangdaoling/marketmatchengine/common"
)

type Order struct {
	RemainAmount  uint64
	ID            uint32
	CancelOrderID uint32
	OrderTime     uint64
	InitialPrice  uint64
	InitialAmount uint64
	IsMarket      bool //market order or limit order
	IsBuy         bool //buy order or sell order
	Symbol        string
	Data          *common.ZeroCopySink
}

//for print
func (o Order) String() string {
	return fmt.Sprintf("order:\n"+
		"ID: %d\n"+
		"CancelOrderID: %d\n"+
		"OrderTime: %d\n"+
		"InitialPrice: %d\n"+
		"RemainAmount: %d\n"+
		"IsMarket: %t\n"+
		"IsBuy: %t\n"+
		"Symbol: %s\n",
		o.ID, o.CancelOrderID, o.OrderTime, o.InitialPrice, o.RemainAmount, o.IsMarket, o.IsBuy, o.Symbol)
}

//for queue.Item interface
//market first, price second, id third
//both must be buy or sell
func (o *Order) Compare(item interface{}) int {
	i := item.(*Order)
	if o.IsMarket && i.IsMarket {
		return compareID(o.ID, i.ID)
	} else if o.IsMarket {
		return 1
	} else if i.IsMarket {
		return -1
	}
	// !o.IsMarket && !i.IsMarket{
	if o.InitialPrice > i.InitialPrice {
		if o.IsBuy {
			return 1
		}
		return -1
	} else if o.InitialPrice < i.InitialPrice {
		if o.IsBuy {
			return -1
		}
		return 1
	}
	//o.InitialPrice  == o.InitialPrice
	return compareID(o.ID, i.ID)
}

//for queue.Item interface
func (o *Order) Key() uint32 {
	return o.ID
}

//for queue.Item interface
func (o *Order) Serialize() (zero *common.ZeroCopySink) {
	o.Data = o.serialize()
	return o.Data
}

func UnSerialize(data []byte, o *Order) (err error) {
	var eof, irregular bool
	zero := common.NewZeroCopySource(data)
	o.RemainAmount, eof = zero.NextUint64()
	if eof {
		return common.ErrTooLarge
	}
	o.ID, eof = zero.NextUint32()
	if eof {
		return common.ErrTooLarge
	}
	o.CancelOrderID, eof = zero.NextUint32()
	if eof {
		return common.ErrTooLarge
	}
	o.OrderTime, eof = zero.NextUint64()
	if eof {
		return common.ErrTooLarge
	}
	o.InitialPrice, eof = zero.NextUint64()
	if eof {
		return common.ErrTooLarge
	}
	o.InitialAmount, eof = zero.NextUint64()
	if eof {
		return common.ErrTooLarge
	}
	o.IsMarket, irregular, eof = zero.NextBool()
	if irregular {
		return common.ErrIrregularData
	}
	if eof {
		return common.ErrTooLarge
	}
	o.IsBuy, irregular, eof = zero.NextBool()
	if irregular {
		return common.ErrIrregularData
	}
	if eof {
		return common.ErrTooLarge
	}
	o.Symbol, _, irregular, eof = zero.NextString()
	if irregular {
		return common.ErrIrregularData
	}
	if eof {
		return common.ErrUnexpectedEOF
	}
	return nil
}

func (o *Order) serialize() (zero *common.ZeroCopySink) {
	zero = common.NewZeroCopySink(nil, 64)
	zero.WriteUint64(o.RemainAmount)
	zero.WriteUint32(o.ID)
	zero.WriteUint32(o.CancelOrderID)
	zero.WriteUint64(o.OrderTime)
	zero.WriteUint64(o.InitialPrice)
	zero.WriteUint64(o.InitialAmount)
	zero.WriteBool(o.IsMarket)
	zero.WriteBool(o.IsBuy)
	zero.WriteString(o.Symbol)
	return
}

//to do: use which price when buy.InitialPrice >= sell.InitialPrice
//Direction ture = buy
func Match(lastPrice uint64, time uint64, buy *Order, sell *Order, direction bool) (r *Transaction) {
	if buy.RemainAmount == 0 || sell.RemainAmount == 0 {
		return nil
	}
	var matchPrice, amount uint64
	r = &Transaction{
		BuyID:     buy.ID,
		SellID:    sell.ID,
		Symbol:    buy.Symbol,
		MatchTime: time,
		IsBuy:     direction,
	}
	if buy.IsMarket && sell.IsMarket {
		matchPrice = lastPrice
		amount = min(buy.RemainAmount, sell.RemainAmount)

	} else if buy.IsMarket {
		matchPrice = sell.InitialPrice
		amount = min(buy.RemainAmount, sell.RemainAmount)
	} else if sell.IsMarket {
		matchPrice = buy.InitialPrice
		amount = min(buy.RemainAmount, sell.RemainAmount)
	} else {
		if buy.InitialPrice < sell.InitialPrice {
			return nil
		}
		if direction {
			matchPrice = sell.InitialPrice
		} else {
			matchPrice = buy.InitialPrice
		}
		amount = min(buy.RemainAmount, sell.RemainAmount)
	}
	r.Price = matchPrice
	r.Amount = amount
	buy.RemainAmount -= amount
	sell.RemainAmount -= amount
	return r
}

//ID is the time, samller time is bigger
func compareID(a uint32, b uint32) int {
	if a > b {
		return -1
	} else if a == b {
		return 0
	}
	return 1
}

func min(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}
