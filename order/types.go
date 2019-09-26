package order

import "fmt"

type Transaction struct {
	BuyIndex    uint64 `json:"buy_index"`
	SellIndex   uint64 `json:"sell_indxe"`
	BuyOrderID  uint64 `json:"buy_order_id"`
	SellOrderID uint64 `json:"sell_order_id"`
	BuyUserID   uint64 `json:"buy_user_id"`
	SellUserID  uint64 `json:"sell_user_id"`
	MatchTime   uint64 `json:"match_time"`
	Price       uint64 `json:"price"`
	Amount      uint64 `json:"amount"`
	IsBuy       bool   `json:"is_buy"` //direction
	Symbol      string `json:"symbol"`
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
