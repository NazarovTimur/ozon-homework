package model

type StockItem struct {
	SKU        int32  `json:"sku"`
	TotalCount uint16 `json:"total_count"`
	Reserved   uint16 `json:"reserved"`
}
