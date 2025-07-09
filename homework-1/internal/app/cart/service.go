package cart

import (
	"context"
	"errors"
	"homework-1/internal/app/product"
	"homework-1/internal/pkg/errorx"
	"homework-1/internal/pkg/model"
	pb "homework-1/internal/proto/loms"
	"homework-1/internal/repository/cart"
)

type ServiceMethods interface {
	Add(ctx context.Context, userID int64, sku uint32, count uint16) (total uint16, existed bool, error error)
	Remove(ctx context.Context, userID int64, sku uint32) error
	Clear(ctx context.Context, userID int64) error
	Get(ctx context.Context, userID int64) (model.CartResponse, error)
	Checkout(ctx context.Context, userID int64) (int64, error)
}

type Service struct {
	repo           cart.CartRepository
	productService product.ProductValidator
	lomsClient     pb.LomsServiceClient
}

func New(repo cart.CartRepository, productService product.ProductValidator, lomsClient pb.LomsServiceClient) *Service {
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
	resp, err := s.lomsClient.StocksInfo(ctx, &pb.StocksInfoRequest{Sku: sku})
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
	var totalPrice uint32
	var cartItems []model.CartItem

	userCart, err := s.repo.GetCart(ctx, userID)
	if err != nil {
		return model.CartResponse{}, err
	}
	for _, cart := range userCart {
		productData, err := s.productService.ValidateProduct(cart.SkuID)
		if err != nil {
			return model.CartResponse{}, err
		}
		totalPrice += uint32(cart.Count) * productData.Price
		cartItems = append(cartItems, model.CartItem{Name: productData.Name, SkuID: int64(cart.SkuID), Count: cart.Count, Price: productData.Price})
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
	response, err := s.lomsClient.OrderCreate(ctx, &pb.OrderCreateRequest{UserID: userID, Items: items})
	if err != nil || response == nil {
		return 0, err
	}
	err = s.repo.ClearCart(ctx, userID)
	if err != nil {
		return 0, err
	}
	return response.OrderID, nil
}

func convertItems(itemCart []model.ItemCart) ([]*pb.OrderItem, error) {
	if len(itemCart) == 0 {
		return nil, errors.New("empty order items")
	}

	items := make([]*pb.OrderItem, len(itemCart))
	for i, item := range itemCart {
		items[i] = &pb.OrderItem{
			Sku:   item.SkuID,
			Count: uint32(item.Count),
		}
	}
	return items, nil
}
