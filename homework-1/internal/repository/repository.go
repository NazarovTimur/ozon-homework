package repository

import (
	"homework-1/internal/pkg/errorx"
	"homework-1/internal/pkg/model"
	"sync"
)

type CartRepository interface {
	AddCart(userID int64, sku uint32, count uint16) (total uint16, existed bool)
	RemoveCart(userID int64, sku uint32) error
	ClearCart(userID int64) error
	GetCart(userID int64) ([]model.ItemCart, error)
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

func (c *Cart) AddCart(userID int64, sku uint32, count uint16) (total uint16, existed bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.carts[userID]; !ok {
		c.carts[userID] = make(map[uint32]uint16)
	}

	_, existed = c.carts[userID][sku]
	c.carts[userID][sku] += count
	total = c.carts[userID][sku]

	return total, existed
}

func (c *Cart) RemoveCart(userID int64, sku uint32) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.carts[userID]; !ok {
		return errorx.ErrUserNotFound
	}
	delete(c.carts[userID], sku)
	return nil
}

func (c *Cart) ClearCart(userID int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.carts[userID]; !ok {
		return errorx.ErrUserNotFound
	}
	delete(c.carts, userID)
	return nil
}

func (c *Cart) GetCart(userID int64) ([]model.ItemCart, error) {
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
