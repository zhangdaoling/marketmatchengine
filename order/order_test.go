package order

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		IndexTime:     249999999999534,
		OrderID:       243423534,
		CancelOrderID: 43434,
		OrderTime:     2083410,
		InitialPrice:  343207710,
		InitialAmount: 3430374723,
		IsMarket:      true,
		IsBuy:         false,
		Symbol:        "usdt-btc",
	}
	o2 := &Order{}
	z := o1.Serialize()
	err := UnSerialize(z.Bytes(), o2)
	assert.Nil(t, err)
	assert.Equal(t, o1.RemainAmount, o2.RemainAmount)
	assert.Equal(t, o1.Index, o2.Index)
	assert.Equal(t, o1.IndexTime, o2.IndexTime)
	assert.Equal(t, o1.OrderID, o2.OrderID)
	assert.Equal(t, o1.CancelOrderID, o2.CancelOrderID)
	assert.Equal(t, o1.OrderTime, o2.OrderTime)
	assert.Equal(t, o1.InitialPrice, o2.InitialPrice)
	assert.Equal(t, o1.InitialAmount, o2.InitialAmount)
	assert.Equal(t, o1.IsMarket, o2.IsMarket)
	assert.Equal(t, o1.IsBuy, o2.IsBuy)
	assert.Equal(t, o1.Symbol, o2.Symbol)
}

//to do
func TestMatch(t *testing.T) {
	return
}
