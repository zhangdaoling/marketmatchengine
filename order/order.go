package order

import (
	"fmt"
	"github.com/zhangdaoling/marketmatchengine/common"
)

type Order struct {
	RemainAmount  uint64
	ID            uint32
	CancelID      uint32
	UserID        uint32
	OrderTime     uint64
	InitialPrice  uint64
	InitialAmount uint64
	IsMarket      bool //market order or limit order
	IsBuy         bool //buy order or sell order
	Canceled      bool //canceled
	Symbol        string
	Data          *common.ZeroCopySink
}

//market first, price second, id third
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

func (o *Order) Key() uint32 {
	return o.ID
}

//for print
func (o *Order) String() string {
	return fmt.Sprintf("\n"+
		"ID: %d\n"+
		"CancelID: %d\n"+
		"UserID: %d\n"+
		"OrderTime: %d\n"+
		"InitialPrice: %d\n"+
		"RemainAmount: %d\n"+
		"IsMarket: %t\n"+
		"IsBuy: %t\n"+
		"Canceled: %t\n"+
		"Symbol: %s\n",
		o.ID, o.CancelID, o.UserID, o.OrderTime, o.InitialPrice, o.RemainAmount, o.IsMarket, o.IsBuy, o.Canceled, o.Symbol)
}

func (o *Order) Serialize() (zero *common.ZeroCopySink) {
	o.Data = o.serialize()
	return o.Data
}

func UnSerialize(data []byte, o *Order) (err error) {
	zero := common.NewZeroCopySource(data)
	var eof, irregular bool
	o.RemainAmount, eof = zero.NextUint64()
	if eof {
		return common.ErrTooLarge
	}
	o.ID, eof = zero.NextUint32()
	if eof {
		return common.ErrTooLarge
	}
	o.CancelID, eof = zero.NextUint32()
	if eof {
		return common.ErrTooLarge
	}
	o.UserID, eof = zero.NextUint32()
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
	o.Canceled, irregular, eof = zero.NextBool()
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
	zero = common.NewZeroCopySink(nil)
	zero.WriteUint64(o.RemainAmount)
	zero.WriteUint32(o.ID)
	zero.WriteUint32(o.CancelID)
	zero.WriteUint32(o.UserID)
	zero.WriteUint64(o.OrderTime)
	zero.WriteUint64(o.InitialPrice)
	zero.WriteUint64(o.InitialAmount)
	zero.WriteBool(o.IsMarket)
	zero.WriteBool(o.IsBuy)
	zero.WriteBool(o.Canceled)
	zero.WriteString(o.Symbol)
	return
}

type MatchResult struct {
	CancelID   uint32
	BuyID      uint32
	SellID     uint32
	BuyUserID  uint32
	SellUserID uint32
	MatchTime  uint64
	Price      uint64
	Amount     uint64
	Symbol     string
}

//for print
func (m *MatchResult) String() string {
	return fmt.Sprintf("\n"+
		"CancelID: %d\n"+
		"BuyID: %d\n"+
		"SellID: %d\n"+
		"BuyUserID: %d\n"+
		"SellUserID: %d\n"+
		"MatchTime: %d\n"+
		"Price: %d\n"+
		"Amount: %d\n"+
		"Symbol: %s\n",
		m.CancelID, m.BuyID, m.SellID, m.BuyUserID, m.SellUserID, m.MatchTime, m.Price, m.Amount, m.Symbol)
}

//to do: use which price when buy.InitialPrice >= sell.InitialPrice
func Match(lastPrice uint64, buy *Order, sell *Order, time uint64) (r *MatchResult) {
	if buy.Canceled || sell.Canceled {
		return nil
	}
	if buy.RemainAmount == 0 || sell.RemainAmount == 0 {
		return nil
	}
	var matchPrice, amount uint64
	r = &MatchResult{
		BuyID:      buy.ID,
		SellID:     sell.ID,
		BuyUserID:  buy.UserID,
		SellUserID: sell.UserID,
		Symbol:     buy.Symbol,
		MatchTime:  time,
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
		matchPrice = sell.InitialPrice
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
