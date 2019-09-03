package order

import (
	"gotest.tools/assert"
	"testing"
)

func TestCompare1(t *testing.T) {
	var order1, order2 *Order
	var r int

	order1 = &Order{
		IsMarket:     true,
		InitialPrice: 2,
		ID:           2,
	}

	//market and limit, always 1
	order2 = &Order{
		IsMarket:     false,
		InitialPrice: 2,
		ID:           2,
	}
	r = order1.Compare(order2)
	assert.Equal(t, r, 1, order2)

	//market and limit, always 1
	order2 = &Order{
		IsMarket:     false,
		InitialPrice: 3,
		ID:           2,
	}
	r = order1.Compare(order2)
	assert.Equal(t, r, 1, order2)

	//market and market, check id
	order2 = &Order{
		IsMarket:     true,
		InitialPrice: 1,
		ID:           2,
	}
	r = order1.Compare(order2)
	assert.Equal(t, r, 0, order2)

	//market and market, check id
	order2 = &Order{
		IsMarket:     true,
		InitialPrice: 3,
		ID:           2,
	}
	r = order1.Compare(order2)
	assert.Equal(t, r, 0, order2)

	//market and market, check id
	order2 = &Order{
		IsMarket:     true,
		InitialPrice: 2,
		ID:           1,
	}
	r = order1.Compare(order2)
	assert.Equal(t, r, -1, order2)

	//market and market, check id
	order2 = &Order{
		IsMarket:     true,
		InitialPrice: 2,
		ID:           2,
	}
	r = order1.Compare(order2)
	assert.Equal(t, r, 0, order2)

	//market and market, check id
	order2 = &Order{
		IsMarket:     true,
		InitialPrice: 2,
		ID:           3,
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
		ID:           2,
	}

	//market and limit, always -1
	order2 = &Order{
		IsMarket:     true,
		InitialPrice: 2,
		ID:           2,
	}
	r = order1.Compare(order2)
	assert.Equal(t, r, -1, order2)

	//market and limit, always -1
	order2 = &Order{
		IsMarket:     true,
		InitialPrice: 2,
		ID:           2,
	}
	r = order1.Compare(order2)
	assert.Equal(t, r, -1, order2)

	//limit and limit, check price
	order2 = &Order{
		IsMarket:     false,
		InitialPrice: 1,
		ID:           2,
	}
	r = order1.Compare(order2)
	assert.Equal(t, r, -1, order2)

	//limit and limit, check price
	order2 = &Order{
		IsMarket:     false,
		InitialPrice: 3,
		ID:           2,
	}
	r = order1.Compare(order2)
	assert.Equal(t, r, 1, order2)

	//limit and limit, price = price, check id
	order2 = &Order{
		IsMarket:     false,
		InitialPrice: 2,
		ID:           1,
	}
	r = order1.Compare(order2)
	assert.Equal(t, r, -1, order2)

	//limit and limit, price = price, check id
	order2 = &Order{
		IsMarket:     false,
		InitialPrice: 2,
		ID:           3,
	}
	r = order1.Compare(order2)
	assert.Equal(t, r, 1, order2)

	//limit and limit, price = price, check id
	order2 = &Order{
		IsMarket:     false,
		InitialPrice: 2,
		ID:           2,
	}
	r = order1.Compare(order2)
	assert.Equal(t, r, 0, order2)
}
