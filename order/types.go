package order

import "fmt"

type Transaction struct {
	CancelID   uint32
	BuyID      uint32
	SellID     uint32
	BuyUserID  uint32
	SellUserID uint32
	MatchTime  uint64
	Price      uint64
	Amount     uint64
	IsBuy      bool //direction
	Symbol     string
}

//for print
func (t Transaction) String() string {
	return fmt.Sprintf("Transaction \n"+
		"CancelOrderID: %d\n"+
		"BuyID: %d\n"+
		"SellID: %d\n"+
		"BuyUserID: %d\n"+
		"SellUserID: %d\n"+
		"MatchTime: %d\n"+
		"Price: %d\n"+
		"Amount: %d\n"+
		"IsBuy: %t\n"+
		"Symbol: %s\n",
		t.CancelID, t.BuyID, t.SellID, t.BuyUserID, t.SellUserID, t.MatchTime, t.Price, t.Amount, t.IsBuy, t.Symbol)
}

type CancelOrder struct {
	ID            uint32
	CancelOrderID uint32
	UserID        uint32
	MatchTime     uint64
	Price         uint64
	Amount        uint64
	IsBuy         bool
	Symbol        string
}

//for print
func (c CancelOrder) String() string {
	return fmt.Sprintf("CancelOrder \n"+
		"ID: %d\n"+
		"CancelOrderID: %d\n"+
		"UserID: %d\n"+
		"MatchTime: %d\n"+
		"Price: %d\n"+
		"Amount: %d\n"+
		"IsBuy: %t\n"+
		"Symbol: %s\n",
		c.ID, c.CancelOrderID, c.UserID, c.MatchTime, c.Price, c.Amount, c.IsBuy, c.Symbol)
}
