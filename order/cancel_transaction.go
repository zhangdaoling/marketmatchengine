package order

import "fmt"

type CancelTransaction struct {
	OrderIndex    uint64
	OrderID       uint64
	CancelOrderID uint64
	UserID        uint64
	MatchTime     uint64
	Price         uint64
	Amount        uint64
	IsBuy         bool
	Symbol        string
}

//for print
func (c CancelTransaction) String() string {
	return fmt.Sprintf("CancelTransaction \n"+
		"OrderIndex： %d\n"+
		"OrderID: %d\n"+
		"CancelOrderID: %d\n"+
		"UserID: %d\n"+
		"MatchTime: %d\n"+
		"Price: %d\n"+
		"Amount: %d\n"+
		"IsBuy: %t\n"+
		"Symbol: %s\n",
		c.OrderIndex, c.OrderID, c.CancelOrderID, c.UserID, c.MatchTime, c.Price, c.Amount, c.IsBuy, c.Symbol)
}
