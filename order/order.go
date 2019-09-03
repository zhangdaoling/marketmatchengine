package order

import "fmt"

type Order struct {
	ID                int32
	UserID            int32
	OrderTime         int64
	InitialPrice      int64
	InitialAmount     int64
	RemainAmount      int64
	LastMatchedPrice  int64
	LastMatchedAmount int64
	IsMarket          bool //market order or limit order
	IsBuy             bool //buy order or sell order
	IsCancel          bool //cancel order
	Canceled          bool
	Symbol            string
}

//for print
func (o *Order) String() string {
	return fmt.Sprintf("\n"+
		"ID: %d\n"+
		"UserID: %d\n"+
		"OrderTime: %d\n"+
		"InitialPrice: %d\n"+
		"RemainAmount: %d\n"+
		"LastMatchedPrice: %d\n"+
		"LastMatchedAmount: %d\n"+
		"IsMarket: %t\n"+
		"IsBuy: %t\n"+
		"IsCancel: %t\n"+
		"Canceled: %t\n"+
		"Symbol: %s\n",
		o.ID, o.UserID, o.OrderTime, o.InitialPrice, o.RemainAmount, o.LastMatchedPrice, o.LastMatchedAmount, o.IsMarket, o.IsBuy, o.IsCancel, o.Canceled, o.Symbol)
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

func (o *Order) Key() int32 {
	return o.ID
}

type MatchResult struct {
	BuyID      int32
	SellID     int32
	BuyUserID  int32
	SellUserID int32
	MatchTime  int64
	Price      int64
	Amount     int64
	IsCancel   bool //match result or cancel result
	Symbol     string
}

//for print
func (m *MatchResult) String() string {
	return fmt.Sprintf("\n"+
		"BuyID: %d\n"+
		"SellID: %d\n"+
		"BuyUserID: %d\n"+
		"SellUserID: %d\n"+
		"MatchTime: %d\n"+
		"Price: %d\n"+
		"Amount: %d\n"+
		"IsMatchResult: %t\n"+
		"Symbol: %s\n",
		m.BuyID, m.SellID, m.BuyUserID, m.SellUserID, m.MatchTime, m.Price, m.Amount, m.IsCancel, m.Symbol)
}

//to do: use which price when buy.InitialPrice >= sell.InitialPrice
func Match(lastPrice int64, buy *Order, sell *Order, time int64) (r *MatchResult) {
	if buy.Canceled || sell.Canceled {
		return nil
	}
	if buy.RemainAmount == 0 || sell.RemainAmount == 0 {
		return nil
	}
	var matchPrice, amount int64
	r = &MatchResult{
		BuyID:      buy.ID,
		SellID:     sell.ID,
		BuyUserID:  buy.UserID,
		SellUserID: sell.UserID,
		Symbol:     buy.Symbol,
		MatchTime:  time,
		IsCancel:   false,
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
	buy.LastMatchedAmount = amount
	buy.LastMatchedPrice = matchPrice
	sell.RemainAmount -= amount
	sell.LastMatchedAmount = amount
	sell.LastMatchedPrice = matchPrice
	return r
}

//ID is the time, samller time is bigger
func compareID(a int32, b int32) int {
	if a > b {
		return -1
	} else if a == b {
		return 0
	}
	return 1
}

func min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}
