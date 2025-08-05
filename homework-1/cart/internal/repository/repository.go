package repository

import (
	"context"
	"homework-1/cart/internal/pkg/errorx"
	"homework-1/cart/internal/pkg/model"
	"sync"
)

type CartRepository interface {
	AddCart(ctx context.Context, userID int64, sku uint32, count uint16) (total uint16, existed bool, err error)
	RemoveCart(ctx context.Context, userID int64, sku uint32) error
	ClearCart(ctx context.Context, userID int64) error
	GetCart(ctx context.Context, userID int64) ([]model.ItemCart, error)
}

type Cart struct {
	mu    sync.RWMutex
	carts map[int64]map[uint32]uint16 // userID → (sku → quantity)
}

func New() *Cart {
	return &Cart{
		carts: make(map[int64]map[uint32]uint16),
	}
}

func (c *Cart) AddCart(ctx context.Context, userID int64, sku uint32, count uint16) (total uint16, existed bool, err error) {
	if err = ctx.Err(); err != nil {
		return 0, false, errorx.ErrContextCanceled
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.carts[userID]; !ok {
		c.carts[userID] = make(map[uint32]uint16)
	}

	_, existed = c.carts[userID][sku]
	c.carts[userID][sku] += count
	total = c.carts[userID][sku]

	return total, existed, nil
}

func (c *Cart) RemoveCart(ctx context.Context, userID int64, sku uint32) error {
	if err := ctx.Err(); err != nil {
		return errorx.ErrContextCanceled
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.carts[userID]; !ok {
		return errorx.ErrUserNotFound
	}
	delete(c.carts[userID], sku)
	return nil
}

func (c *Cart) ClearCart(ctx context.Context, userID int64) error {
	if err := ctx.Err(); err != nil {
		return errorx.ErrContextCanceled
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.carts[userID]; !ok {
		return errorx.ErrUserNotFound
	}
	delete(c.carts, userID)
	return nil
}

func (c *Cart) GetCart(ctx context.Context, userID int64) ([]model.ItemCart, error) {
	if err := ctx.Err(); err != nil {
		return nil, errorx.ErrContextCanceled
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if cart, ok := c.carts[userID]; ok {
		items := make([]model.ItemCart, 0, len(cart))
		for sku, count := range cart {
			items = append(items, model.ItemCart{
				SkuID: sku,
				Count: count,
			})
		}
		return items, nil
	}
	return nil, errorx.ErrUserNotFound
}
