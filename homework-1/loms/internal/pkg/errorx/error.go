package errorx

import "errors"

var (
	ErrContextCanceled    = errors.New("operation cancelled by context")
	ErrOrderNotFound      = errors.New("order not found")
	ErrOrderCreating      = errors.New("error creating order")
	ErrOrderCreatingItem  = errors.New("error creating order")
	ErrInvalidOrderStatus = errors.New("invalid status")
	ErrTransactionBegin   = errors.New("error beginning transaction")
	ErrTransactionCommit  = errors.New("error committing transaction")
	ErrSettingStatus      = errors.New("error setting status")

	ErrStockNotFound     = errors.New("stock not found")
	ErrInsufficientStock = errors.New("not enough stock available")
	ErrReserveNotFound   = errors.New("no reserve found to update")
)
