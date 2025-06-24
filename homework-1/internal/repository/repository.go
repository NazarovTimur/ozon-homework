package repository

import (
	"errors"
	"sync"
)

type CartRepository interface {
	AddCart(userID int64, sku uint32, count uint16) (total uint16, existed bool)
	RemoveCart(userID int64, sku uint32)
	ClearCart(userID int64)
	GetCart(userID int64) (map[uint32]uint16, error)
}

type Cart struct {
	mu    sync.RWMutex
	carts map[int64]map[uint32]uint16
}

func New() *Cart {
	return &Cart{
		carts: make(map[int64]map[uint32]uint16),
	}
}

func (cs *Cart) AddCart(userID int64, sku uint32, count uint16) (total uint16, existed bool) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if _, ok := cs.carts[userID]; !ok {
		cs.carts[userID] = make(map[uint32]uint16)
	}

	_, existed = cs.carts[userID][sku]
	cs.carts[userID][sku] += count
	total = cs.carts[userID][sku]

	return total, existed
}

func (cs *Cart) RemoveCart(userID int64, sku uint32) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if _, ok := cs.carts[userID]; ok {
		delete(cs.carts[userID], sku)
	}
}

func (cs *Cart) ClearCart(userID int64) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	delete(cs.carts, userID)
}

func (cs *Cart) GetCart(userID int64) (map[uint32]uint16, error) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	if cart, ok := cs.carts[userID]; ok {
		result := make(map[uint32]uint16, len(cart))
		for k, v := range cart {
			result[k] = v
		}
		return result, nil
	}
	return map[uint32]uint16{}, errors.New("user not found")
}
