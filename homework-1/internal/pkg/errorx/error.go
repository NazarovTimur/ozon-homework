package errorx

import "errors"

var (
	ErrUserNotFound             = errors.New("user not found")
	ErrInsufficientStock        = errors.New("insufficient stock")
	ErrContextCanceled          = errors.New("operation cancelled by context")
	ErrRepositoryNotInitialized = errors.New("order repository not initialized")
	ErrOrderNotFound            = errors.New("order not found")
	ErrInvalidOrderStatus       = errors.New("invalid status")
	ErrOrderIDNotFound          = errors.New("ID not found")
	ErrStockSKUNotFound         = errors.New("sku dont find")
	ErrStockCancel              = errors.New("cancel dont finish")
)
