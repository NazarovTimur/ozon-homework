package model

type Order struct {
	ID     int64
	UserID int64
	Items  []OrderItem
	Status string
}

type OrderItem struct {
	SKU   uint32
	Count uint16
}
