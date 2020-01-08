package order

import "fmt"

type Transaction struct {
	MatchOrderIndex uint64 `json:"match_order_indxe"`
	MatchOrderID    uint64 `json:"match_order_id"`
	MatchOrderTime  uint64 `json:"match_order_time"`
	BuyOrderID      uint64 `json:"buy_order_id"`
	SellOrderID     uint64 `json:"sell_order_id"`
	BuyUserID       uint64 `json:"buy_user_id"`
	SellUserID      uint64 `json:"sell_user_id"`
	Price           uint64 `json:"price"`
	Amount          uint64 `json:"amount"`
	IsBuy           bool   `json:"is_buy"` //direction
	Symbol          string `json:"symbol"`
}

//for print
func (t Transaction) String() string {
	return fmt.Sprintf("Transaction \n"+
		"MatchOrderIndex: %d\n"+
		"MatchOrderID: %d\n"+
		"MatchOrderTime: %d\n"+
		"BuyOrderID: %d\n"+
		"SellOrderID: %d\n"+
		"BuyUserID: %d\n"+
		"SellUserID: %d\n"+
		"Price: %d\n"+
		"Amount: %d\n"+
		"IsBuy: %t\n"+
		"Symbol: %s\n",
		t.MatchOrderIndex, t.MatchOrderID, t.MatchOrderTime, t.BuyOrderID, t.SellOrderID, t.BuyUserID, t.SellUserID, t.Price, t.Amount, t.IsBuy, t.Symbol)
}
