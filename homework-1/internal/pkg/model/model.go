package model

type CartItem struct {
	SkuID int64  `json:"sku_id"`
	Name  string `json:"name"`
	Count uint16 `json:"count"`
	Price uint32 `json:"price"`
}
type CartResponse struct {
	Items      []CartItem `json:"items"`
	TotalPrice uint32     `json:"total_price"`
}
