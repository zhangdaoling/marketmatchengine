package order

import "fmt"

type Transaction struct {
	BuyIndex    uint64
	SellIndex   uint64
	BuyOrderID  uint64
	SellOrderID uint64
	BuyUserID   uint64
	SellUserID  uint64
	MatchTime   uint64
	Price       uint64
	Amount      uint64
	IsBuy       bool //direction
	Symbol      string
}

//for print
func (t Transaction) String() string {
	return fmt.Sprintf("Transaction \n"+
		"BuyIndex: %d\n"+
		"SellIndex: %d\n"+
		"BuyOrderID: %d\n"+
		"SellOrderID: %d\n"+
		"BuyUserID: %d\n"+
		"SellUserID: %d\n"+
		"MatchTime: %d\n"+
		"Price: %d\n"+
		"Amount: %d\n"+
		"IsBuy: %t\n"+
		"Symbol: %s\n",
		t.BuyIndex, t.SellIndex, t.BuyOrderID, t.SellOrderID, t.BuyUserID, t.SellUserID, t.MatchTime, t.Price, t.Amount, t.IsBuy, t.Symbol)
}

type CancelOrder struct {
	Index         uint64
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
func (c CancelOrder) String() string {
	return fmt.Sprintf("CancelOrder \n"+
		"Index： %d\n"+
		"OrderID: %d\n"+
		"CancelOrderID: %d\n"+
		"UserID: %d\n"+
		"MatchTime: %d\n"+
		"Price: %d\n"+
		"Amount: %d\n"+
		"IsBuy: %t\n"+
		"Symbol: %s\n",
		c.Index, c.OrderID, c.CancelOrderID, c.UserID, c.MatchTime, c.Price, c.Amount, c.IsBuy, c.Symbol)
}
