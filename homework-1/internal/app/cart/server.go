package cart

import (
	"sync"
)

type CartService struct {
	mu    sync.RWMutex
	carts map[int64]map[uint32]uint16
}

func New() *CartService {
	return &CartService{
		carts: make(map[int64]map[uint32]uint16),
	}
}

func (cs *CartService) Add(userID int64, sku uint32, count uint16) (total uint16, existed bool) {
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

func (cs *CartService) Remove(userID int64, sku uint32) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if _, ok := cs.carts[userID]; ok {
		delete(cs.carts[userID], sku)
	}
}

func (cs *CartService) Clear(userID int64) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	delete(cs.carts, userID)
}

func (cs *CartService) Get(userID int64) (map[uint32]uint16, bool) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	cart, ok := cs.carts[userID]
	return cart, ok
}
