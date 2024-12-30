package model

import "github.com/google/uuid"

type Sku = int64

type Review struct {
	SKU     Sku
	Comment string
	UserID  uuid.UUID
}
