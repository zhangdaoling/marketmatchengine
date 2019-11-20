package order

import "fmt"

type Transaction struct {
	BuyIndex     uint64 `json:"buy_index"`
	SellIndex    uint64 `json:"sell_indxe"`
	MatchOrderID uint64 `json:"match_order_id"`
	MatchTime    uint64 `json:"match_time"`
	BuyOrderID   uint64 `json:"buy_order_id"`
	SellOrderID  uint64 `json:"sell_order_id"`
	BuyUserID    uint64 `json:"buy_user_id"`
	SellUserID   uint64 `json:"sell_user_id"`
	Price        uint64 `json:"price"`
	Amount       uint64 `json:"amount"`
	IsBuy        bool   `json:"is_buy"` //direction
	Symbol       string `json:"symbol"`
}

//for print
func (t Transaction) String() string {
	return fmt.Sprintf("Transaction \n"+
		"BuyIndex: %d\n"+
		"SellIndex: %d\n"+
		"MatchOrderID: %d\n"+
		"MatchTime: %d\n"+
		"BuyOrderID: %d\n"+
		"SellOrderID: %d\n"+
		"BuyUserID: %d\n"+
		"SellUserID: %d\n"+
		"Price: %d\n"+
		"Amount: %d\n"+
		"IsBuy: %t\n"+
		"Symbol: %s\n",
		t.BuyIndex, t.SellIndex, t.MatchOrderID, t.MatchTime, t.BuyOrderID, t.SellOrderID, t.BuyUserID, t.SellUserID, t.Price, t.Amount, t.IsBuy, t.Symbol)
}
