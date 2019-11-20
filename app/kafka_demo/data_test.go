package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/zhangdaoling/marketmatchengine/order"
)

var orderIDIndex uint64

//init balance sql
func TestInitBalance(*testing.T) {
	orderIDIndex = 1
	symbol := []string{"A", "B"}
	count := 5000000
	limit := 1000
	headerString := "INSERT user_balance (user_id, symbol, amount) VALUES"

	f, err := os.OpenFile("./mysql/balance.sql", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for _, s := range symbol {
		for i := 0; i < count; i++ {
			if ((i+1)%limit == 0) || i >= count {
				headerString = fmt.Sprintf(`%s (%d, "%s", %d);`, headerString, i+1, s, 100000000)
				_, err = f.WriteString(headerString + "\n")
				if err != nil {
					log.Println(err)
				}
				headerString = "INSERT user_balance (user_id, symbol, amount) VALUES"
			} else {
				headerString = fmt.Sprintf(`%s (%d, "%s", %d),`, headerString, i+1, s, 100000000)
			}
		}
	}
	f.Sync()
}

//init order sql
func TestInitOrder(*testing.T) {
	orderIDIndex = 1
	orders := initBuy(50000)
	o2 := initSell(50000)
	orders = append(orders, o2...)

	limit := 1000
	headerString := "INSERT orders (user_id, symbol, amount, price, is_market, is_buy) VALUES"

	f, err := os.OpenFile("./mysql/order.sql", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for i, o := range orders {
		if ((i+1)%limit == 0) || i >= len(orders) {
			headerString = fmt.Sprintf(`%s (%d, "%s", %d, %d, %t, %t);`, headerString, o.UserID, o.Symbol, o.InitialAmount, o.InitialPrice, o.IsMarket, o.IsBuy)
			_, err = f.WriteString(headerString + "\n")
			if err != nil {
				log.Println(err)
			}
			headerString = "INSERT orders (user_id, symbol, amount, price, is_market, is_buy) VALUES"
		} else {
			headerString = fmt.Sprintf(`%s (%d, "%s", %d, %d, %t, %t),`, headerString, o.UserID, o.Symbol, o.InitialAmount, o.InitialPrice, o.IsMarket, o.IsBuy)
		}
	}
	f.Sync()
}

func initBuy(length int) (orders []*order.Order) {
	symbol := "A-B"
	orders = make([]*order.Order, 0, length)
	for i := 0; i < length/2; i++ {
		o := &order.Order{
			Index:         orderIDIndex,
			OrderID:       orderIDIndex,
			OrderTime:     2000000 + uint64(orderIDIndex),
			UserID:        orderIDIndex,
			InitialPrice:  1 + 2*uint64(i),
			InitialAmount: 10,
			RemainAmount:  10,
			IsBuy:         true,
			Symbol:        symbol,
		}
		orders = append(orders, o)
		orderIDIndex++
	}
	for i := 0; i < length/2; i++ {
		o := &order.Order{
			Index:         orderIDIndex,
			OrderID:       orderIDIndex,
			OrderTime:     2000000 + uint64(orderIDIndex),
			UserID:        orderIDIndex,
			InitialPrice:  uint64(length - i),
			InitialAmount: 10,
			RemainAmount:  10,
			IsBuy:         true,
			Symbol:        symbol,
		}
		orders = append(orders, o)
		orderIDIndex++
	}
	return
}

func initSell(length int) (orders []*order.Order) {
	symbol := "A-B"
	orders = make([]*order.Order, 0, length)
	for i := 0; i < length/2; i++ {
		o := &order.Order{
			Index:         orderIDIndex,
			OrderID:       orderIDIndex,
			OrderTime:     2000000 + uint64(orderIDIndex),
			UserID:        orderIDIndex,
			InitialPrice:  1 + 2*uint64(i+length),
			InitialAmount: 10,
			RemainAmount:  10,
			IsBuy:         false,
			Symbol:        symbol,
		}
		orders = append(orders, o)
		orderIDIndex++
	}
	for i := 0; i < length/2; i++ {
		o := &order.Order{
			Index:         orderIDIndex,
			OrderID:       orderIDIndex,
			OrderTime:     2000000 + uint64(orderIDIndex),
			UserID:        orderIDIndex,
			InitialPrice:  2 + 2*uint64(length/2-i),
			InitialAmount: 5,
			RemainAmount:  5,
			IsBuy:         false,
			Symbol:        symbol,
		}
		orders = append(orders, o)
		orderIDIndex++
	}
	return
}
