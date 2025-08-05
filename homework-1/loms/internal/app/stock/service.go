package stock

import (
	"context"
	"homework-1/loms/internal/pkg/model"
	"homework-1/loms/internal/repository/stock"
)

type StockServiceInterface interface {
	Reserve(ctx context.Context, items []model.OrderItem) (bool, error)
	ReserveRemove(ctx context.Context, sku uint32, count uint16) error
	ReserveCancel(ctx context.Context, sku uint32, count uint16) error
	GetBySKU(ctx context.Context, sku uint32) (uint16, error)
}

type Service struct {
	stockRepo stock.StockInterface
}

func New(stockRepo stock.StockInterface) *Service {
	return &Service{
		stockRepo: stockRepo,
	}
}

func (s *Service) Reserve(ctx context.Context, items []model.OrderItem) (bool, error) {
	reserved := make([]model.OrderItem, 0, len(items))
	for _, item := range items {
		ok, err := s.stockRepo.Reserve(ctx, item.SKU, item.Count)
		if err != nil {
			return false, err
		}
		if !ok {
			for _, r := range reserved {
				if err = s.stockRepo.ReserveCancel(ctx, r.SKU, r.Count); err != nil {
					return false, err
				}
			}
			return false, nil
		}
		reserved = append(reserved, item)
	}
	return true, nil
}

func (s *Service) ReserveRemove(ctx context.Context, sku uint32, count uint16) error {
	if err := s.stockRepo.ReserveRemove(ctx, sku, count); err != nil {
		return err
	}
	return nil
}

func (s *Service) ReserveCancel(ctx context.Context, sku uint32, count uint16) error {
	if err := s.stockRepo.ReserveCancel(ctx, sku, count); err != nil {
		return err
	}
	return nil
}

func (s *Service) GetBySKU(ctx context.Context, sku uint32) (uint16, error) {
	count, err := s.stockRepo.GetBySKU(ctx, sku)
	if err != nil {
		return 0, err
	}
	return count, nil
}
