package cart

import (
	"context"
	"errors"
	"fmt"
	pb2 "homework-1/api/proto/loms"
	"homework-1/cart/internal/app/product"
	"homework-1/cart/internal/pkg/errorx"
	"homework-1/cart/internal/pkg/model"
	"homework-1/cart/internal/repository"
)

type ServiceMethods interface {
	Add(ctx context.Context, userID int64, sku uint32, count uint16) (total uint16, existed bool, error error)
	Remove(ctx context.Context, userID int64, sku uint32) error
	Clear(ctx context.Context, userID int64) error
	Get(ctx context.Context, userID int64) (model.CartResponse, error)
	Checkout(ctx context.Context, userID int64) (int64, error)
}

type Service struct {
	repo           repository.CartRepository
	productService product.ProductValidator
	lomsClient     pb2.LomsServiceClient
}

func New(repo repository.CartRepository, productService product.ProductValidator, lomsClient pb2.LomsServiceClient) *Service {
	return &Service{
		repo:           repo,
		productService: productService,
		lomsClient:     lomsClient,
	}
}

func (s *Service) Add(ctx context.Context, userID int64, sku uint32, count uint16) (total uint16, existed bool, error error) {
	_, err := s.productService.ValidateProduct(sku)
	if err != nil {
		return 0, false, err
	}
	resp, err := s.lomsClient.StocksInfo(ctx, &pb2.StocksInfoRequest{Sku: sku})
	if err != nil || resp == nil {
		return 0, false, err
	}
	if resp.Count < uint64(count) {
		return 0, false, errorx.ErrInsufficientStock
	}

	return s.repo.AddCart(ctx, userID, sku, count)
}

func (s *Service) Remove(ctx context.Context, userID int64, sku uint32) error {
	err := s.repo.RemoveCart(ctx, userID, sku)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Clear(ctx context.Context, userID int64) error {
	err := s.repo.ClearCart(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Get(ctx context.Context, userID int64) (model.CartResponse, error) {
	userCart, err := s.repo.GetCart(ctx, userID)
	if err != nil {
		return model.CartResponse{}, err
	}

	SKUs := make([]uint32, len(userCart))
	for i, cart := range userCart {
		SKUs[i] = cart.SkuID
	}

	productData, err := s.productService.ValidateProductParallel(SKUs)
	if err != nil {
		return model.CartResponse{}, err
	}

	var totalPrice uint32
	cartItems := make([]model.CartItem, len(userCart))

	for i, cart := range userCart {
		product, ok := productData[cart.SkuID]
		if !ok {
			return model.CartResponse{}, fmt.Errorf("product data not found for sku %d", cart.SkuID)
		}

		cartItems[i] = model.CartItem{
			SkuID: int64(cart.SkuID),
			Count: cart.Count,
			Name:  product.Name,
			Price: product.Price,
		}
		totalPrice += uint32(cart.Count) * product.Price
	}

	return model.CartResponse{Items: cartItems, TotalPrice: totalPrice}, nil
}

func (s *Service) Checkout(ctx context.Context, userID int64) (int64, error) {
	userCart, err := s.repo.GetCart(ctx, userID)
	if err != nil {
		return 0, err
	}
	items, err := convertItems(userCart)
	if err != nil {
		return 0, err
	}
	response, err := s.lomsClient.OrderCreate(ctx, &pb2.OrderCreateRequest{UserID: userID, Items: items})
	if err != nil || response == nil {
		return 0, err
	}
	err = s.repo.ClearCart(ctx, userID)
	if err != nil {
		return 0, err
	}
	return response.OrderID, nil
}

func convertItems(itemCart []model.ItemCart) ([]*pb2.OrderItem, error) {
	if len(itemCart) == 0 {
		return nil, errors.New("empty order items")
	}

	items := make([]*pb2.OrderItem, len(itemCart))
	for i, item := range itemCart {
		items[i] = &pb2.OrderItem{
			Sku:   item.SkuID,
			Count: uint32(item.Count),
		}
	}
	return items, nil
}
