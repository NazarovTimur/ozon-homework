package stock

import (
	"context"
	"homework-1/internal/pkg/errorx"
	"sync"
)

type StockInterface interface {
	Reserve(ctx context.Context, sku uint32, count uint16) bool
	ReserveRemove(ctx context.Context, sku uint32, count uint16) error
	ReserveCancel(ctx context.Context, sku uint32, count uint16) error
	GetBySKU(ctx context.Context, sku uint32) (uint16, error)
}
type StockRepository struct {
	mu       sync.RWMutex
	stocks   map[uint32]uint16
	reserved map[uint32]uint16
}

func (s *StockRepository) Reserve(ctx context.Context, sku uint32, count uint16) bool {
	if err := ctx.Err(); err != nil {
		return false
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if cnt, ok := s.stocks[sku]; ok {
		if count <= cnt {
			s.stocks[sku] = cnt - count
			s.reserved[sku] += count
			return true
		}
	}
	return false
}

func (s *StockRepository) ReserveRemove(ctx context.Context, sku uint32, count uint16) error {
	if err := ctx.Err(); err != nil {
		return errorx.ErrContextCanceled
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.reserved[sku]; ok {
		s.reserved[sku] = -count
		return nil
	}
	return errorx.ErrStockSKUNotFound
}

func (s *StockRepository) ReserveCancel(ctx context.Context, sku uint32, count uint16) error {
	if err := ctx.Err(); err != nil {
		return errorx.ErrContextCanceled
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.reserved[sku]; ok {
		s.reserved[sku] = -count
		s.stocks[sku] += count
		return nil
	}
	return errorx.ErrStockCancel
}

func (s *StockRepository) GetBySKU(ctx context.Context, sku uint32) (uint16, error) {
	if err := ctx.Err(); err != nil {
		return 0, errorx.ErrContextCanceled
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	if cnt, ok := s.stocks[sku]; ok {
		return cnt, nil
	}
	return 0, errorx.ErrStockSKUNotFound
}
