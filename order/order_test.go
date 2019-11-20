package order

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchMarket(t *testing.T) {
	var buyOrder, sellOrder *Order
	var expectTransaction, transaction *Transaction
	var symbol string = "A-B"
	var lastPrice uint64

	//market to market
	lastPrice = 1000
	buyOrder = &Order{
		Index:         1,
		OrderID:       3,
		OrderTime:     4,
		UserID:        5,
		InitialAmount: 10,
		Symbol:        symbol,
		IsMarket:      true,
		IsBuy:         true,
		InitialPrice:  10,
		RemainAmount:  10,
	}
	sellOrder = &Order{
		Index:         11,
		OrderID:       13,
		OrderTime:     14,
		UserID:        15,
		InitialAmount: 10,
		Symbol:        symbol,
		IsMarket:      true,
		IsBuy:         false,
		InitialPrice:  10,
		RemainAmount:  10,
	}
	expectTransaction = &Transaction{
		BuyIndex:     1,
		SellIndex:    11,
		MatchOrderID: buyOrder.OrderID,
		MatchTime:    buyOrder.OrderTime,
		BuyOrderID:   3,
		SellOrderID:  13,
		BuyUserID:    5,
		SellUserID:   15,
		Symbol:       symbol,
		IsBuy:        buyOrder.IsBuy,
		Price:        lastPrice,
		Amount:       10,
	}
	transaction = Match(lastPrice, buyOrder, buyOrder, sellOrder)
	equalTransaction(t, expectTransaction, transaction)

	//market to market
	lastPrice = 1000
	sellOrder = &Order{
		Index:         11,
		OrderID:       13,
		OrderTime:     14,
		UserID:        15,
		InitialAmount: 10,
		Symbol:        symbol,
		IsMarket:      false,
		IsBuy:         false,
		InitialPrice:  10,
		RemainAmount:  20,
	}
	expectTransaction = &Transaction{
		BuyIndex:     1,
		SellIndex:    11,
		MatchOrderID: buyOrder.OrderID,
		MatchTime:    buyOrder.OrderTime,
		BuyOrderID:   3,
		SellOrderID:  13,
		BuyUserID:    5,
		SellUserID:   15,
		Symbol:       symbol,
		IsBuy:        buyOrder.IsBuy,
		Price:        sellOrder.InitialPrice,
		Amount:       buyOrder.InitialAmount,
	}
	transaction = Match(lastPrice, buyOrder, buyOrder, sellOrder)
	equalTransaction(t, expectTransaction, transaction)
}

func TestMatchLimit(t *testing.T) {
	var buyOrder, sellOrder *Order
	var expectTransaction, transaction *Transaction
	var symbol string = "A-B"
	var lastPrice uint64

	//limit to market
	lastPrice = 1000
	buyOrder = &Order{
		Index:         1,
		OrderID:       3,
		OrderTime:     4,
		UserID:        5,
		InitialAmount: 10,
		Symbol:        symbol,
		IsMarket:      false,
		IsBuy:         true,
		InitialPrice:  10,
		RemainAmount:  10,
	}
	sellOrder = &Order{
		Index:         11,
		OrderID:       13,
		OrderTime:     14,
		UserID:        15,
		InitialAmount: 10,
		Symbol:        symbol,
		IsMarket:      true,
		IsBuy:         false,
		InitialPrice:  33,
		RemainAmount:  100,
	}
	expectTransaction = &Transaction{
		BuyIndex:     1,
		SellIndex:    11,
		MatchOrderID: buyOrder.OrderID,
		MatchTime:    buyOrder.OrderTime,
		BuyOrderID:   3,
		SellOrderID:  13,
		BuyUserID:    5,
		SellUserID:   15,
		Symbol:       symbol,
		IsBuy:        buyOrder.IsBuy,
		Price:        buyOrder.InitialPrice,
		Amount:       10,
	}
	transaction = Match(lastPrice, buyOrder, buyOrder, sellOrder)
	equalTransaction(t, expectTransaction, transaction)

	//limit to limit; buyPrice < sellPrice
	lastPrice = 1000
	sellOrder = &Order{
		Index:         11,
		OrderID:       13,
		OrderTime:     14,
		UserID:        15,
		InitialAmount: 10,
		Symbol:        symbol,
		IsMarket:      false,
		IsBuy:         false,
		InitialPrice:  33,
		RemainAmount:  5,
	}
	expectTransaction = nil
	transaction = Match(lastPrice, buyOrder, buyOrder, sellOrder)
	assert.Nil(t, transaction)

	//limit to limit; buyPrice = SellPrice
	lastPrice = 1000
	sellOrder = &Order{
		Index:         11,
		OrderID:       13,
		OrderTime:     14,
		UserID:        15,
		InitialAmount: 10,
		Symbol:        symbol,
		IsMarket:      false,
		IsBuy:         false,
		InitialPrice:  10,
		RemainAmount:  5,
	}
	expectTransaction = &Transaction{
		BuyIndex:     1,
		SellIndex:    11,
		MatchOrderID: buyOrder.OrderID,
		MatchTime:    buyOrder.OrderTime,
		BuyOrderID:   3,
		SellOrderID:  13,
		BuyUserID:    5,
		SellUserID:   15,
		Symbol:       symbol,
		IsBuy:        buyOrder.IsBuy,
		Price:        buyOrder.InitialPrice,
		Amount:       5,
	}
	transaction = Match(lastPrice, buyOrder, buyOrder, sellOrder)
	equalTransaction(t, expectTransaction, transaction)

	//limit to limit; buyPrice > SellPrice; direction = buy
	lastPrice = 1000
	sellOrder = &Order{
		Index:         11,
		OrderID:       13,
		OrderTime:     14,
		UserID:        15,
		InitialAmount: 10,
		Symbol:        symbol,
		IsMarket:      false,
		IsBuy:         false,
		InitialPrice:  5,
		RemainAmount:  5,
	}
	expectTransaction = &Transaction{
		BuyIndex:     1,
		SellIndex:    11,
		MatchOrderID: buyOrder.OrderID,
		MatchTime:    buyOrder.OrderTime,
		BuyOrderID:   3,
		SellOrderID:  13,
		BuyUserID:    5,
		SellUserID:   15,
		Symbol:       symbol,
		IsBuy:        buyOrder.IsBuy,
		Price:        sellOrder.InitialPrice,
		Amount:       5,
	}
	transaction = Match(lastPrice, buyOrder, buyOrder, sellOrder)
	equalTransaction(t, expectTransaction, transaction)

	//limit to limit; buyPrice > SellPrice; direction = false
	lastPrice = 1000
	sellOrder = &Order{
		Index:         11,
		OrderID:       13,
		OrderTime:     14,
		UserID:        15,
		InitialAmount: 10,
		Symbol:        symbol,
		IsMarket:      false,
		IsBuy:         false,
		InitialPrice:  5,
		RemainAmount:  5,
	}
	expectTransaction = &Transaction{
		BuyIndex:     1,
		SellIndex:    11,
		MatchOrderID: buyOrder.OrderID,
		MatchTime:    buyOrder.OrderTime,
		BuyOrderID:   3,
		SellOrderID:  13,
		BuyUserID:    5,
		SellUserID:   15,
		Symbol:       symbol,
		IsBuy:        buyOrder.IsBuy,
		Price:        buyOrder.InitialPrice,
		Amount:       5,
	}
	transaction = Match(lastPrice, buyOrder, buyOrder, sellOrder)
	equalTransaction(t, expectTransaction, transaction)
}

func TestCompare1(t *testing.T) {
	var order1, order2 *Order
	var r int

	order1 = &Order{
		IsMarket:     true,
		InitialPrice: 2,
		Index:        2,
		OrderID:      2,
	}

	//market and limit, always 1
	order2 = &Order{
		IsMarket:     false,
		InitialPrice: 2,
		Index:        2,
		OrderID:      2,
	}
	r = order1.Compare(order2)
	assert.Equal(t, r, 1, order2)

	//market and limit, always 1
	order2 = &Order{
		IsMarket:     false,
		InitialPrice: 3,
		Index:        2,
		OrderID:      2,
	}
	r = order1.Compare(order2)
	assert.Equal(t, r, 1, order2)

	//market and market, check id
	order2 = &Order{
		IsMarket:     true,
		InitialPrice: 1,
		Index:        2,
		OrderID:      2,
	}
	r = order1.Compare(order2)
	assert.Equal(t, r, 0, order2)

	//market and market, check id
	order2 = &Order{
		IsMarket:     true,
		InitialPrice: 3,
		Index:        2,
		OrderID:      2,
	}
	r = order1.Compare(order2)
	assert.Equal(t, r, 0, order2)

	//market and market, check id
	order2 = &Order{
		IsMarket:     true,
		InitialPrice: 2,
		Index:        1,
		OrderID:      1,
	}
	r = order1.Compare(order2)
	assert.Equal(t, r, -1, order2)

	//market and market, check id
	order2 = &Order{
		IsMarket:     true,
		InitialPrice: 2,
		Index:        2,
		OrderID:      2,
	}
	r = order1.Compare(order2)
	assert.Equal(t, r, 0, order2)

	//market and market, check id
	order2 = &Order{
		IsMarket:     true,
		InitialPrice: 2,
		OrderID:      3,
		Index:        3,
	}
	r = order1.Compare(order2)
	assert.Equal(t, r, 1, order2)
}

func TestCompare2(t *testing.T) {
	var order1, order2 *Order
	var r int

	order1 = &Order{
		IsMarket:     false,
		InitialPrice: 2,
		OrderID:      2,
		Index:        2,
	}

	//market and limit, always -1
	order2 = &Order{
		IsMarket:     true,
		InitialPrice: 2,
		OrderID:      2,
		Index:        2,
	}
	r = order1.Compare(order2)
	assert.Equal(t, r, -1, order2)

	//market and limit, always -1
	order2 = &Order{
		IsMarket:     true,
		InitialPrice: 2,
		OrderID:      2,
		Index:        2,
	}
	r = order1.Compare(order2)
	assert.Equal(t, r, -1, order2)

	//limit and limit, check price
	order2 = &Order{
		IsMarket:     false,
		InitialPrice: 1,
		OrderID:      2,
		Index:        2,
	}
	r = order1.Compare(order2)
	assert.Equal(t, r, -1, order2)

	//limit and limit, check price
	order2 = &Order{
		IsMarket:     false,
		InitialPrice: 3,
		OrderID:      2,
		Index:        2,
	}
	r = order1.Compare(order2)
	assert.Equal(t, r, 1, order2)

	//limit and limit, price = price, check id
	order2 = &Order{
		IsMarket:     false,
		InitialPrice: 2,
		OrderID:      1,
		Index:        1,
	}
	r = order1.Compare(order2)
	assert.Equal(t, r, -1, order2)

	//limit and limit, price = price, check id
	order2 = &Order{
		IsMarket:     false,
		InitialPrice: 2,
		OrderID:      3,
		Index:        3,
	}
	r = order1.Compare(order2)
	assert.Equal(t, r, 1, order2)

	//limit and limit, price = price, check id
	order2 = &Order{
		IsMarket:     false,
		InitialPrice: 2,
		OrderID:      2,
		Index:        2,
	}
	r = order1.Compare(order2)
	assert.Equal(t, r, 0, order2)
}

func TestSerialize(t *testing.T) {
	o1 := &Order{
		RemainAmount:  434234434567,
		Index:         2423534,
		OrderID:       243423534,
		CancelOrderID: 43434,
		OrderTime:     2083410,
		InitialPrice:  343207710,
		InitialAmount: 3430374723,
		IsMarket:      true,
		IsBuy:         false,
		Symbol:        "A-B",
	}
	o2 := &Order{}
	z := o1.Serialize()
	err := UnSerialize(z.Bytes(), o2)
	assert.Nil(t, err)
	assert.Equal(t, o1.RemainAmount, o2.RemainAmount)
	assert.Equal(t, o1.Index, o2.Index)
	assert.Equal(t, o1.OrderID, o2.OrderID)
	assert.Equal(t, o1.CancelOrderID, o2.CancelOrderID)
	assert.Equal(t, o1.OrderTime, o2.OrderTime)
	assert.Equal(t, o1.InitialPrice, o2.InitialPrice)
	assert.Equal(t, o1.InitialAmount, o2.InitialAmount)
	assert.Equal(t, o1.IsMarket, o2.IsMarket)
	assert.Equal(t, o1.IsBuy, o2.IsBuy)
	assert.Equal(t, o1.Symbol, o2.Symbol)
}

func equalTransaction(t *testing.T, t1 *Transaction, t2 *Transaction) {
	assert.Equal(t, t1.BuyIndex, t2.BuyIndex)
	assert.Equal(t, t1.SellIndex, t2.SellIndex)
	assert.Equal(t, t1.MatchOrderID, t2.MatchOrderID)
	assert.Equal(t, t1.MatchTime, t2.MatchTime)
	assert.Equal(t, t1.BuyOrderID, t2.BuyOrderID)
	assert.Equal(t, t1.SellOrderID, t2.SellOrderID)
	assert.Equal(t, t1.BuyUserID, t2.BuyUserID)
	assert.Equal(t, t1.SellUserID, t2.SellUserID)
	assert.Equal(t, t1.Price, t2.Price)
	assert.Equal(t, t1.Amount, t2.Amount)
	assert.Equal(t, t1.IsBuy, t2.IsBuy)
	assert.Equal(t, t1.Symbol, t2.Symbol)
}
