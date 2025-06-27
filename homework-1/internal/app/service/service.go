package service

import (
	"homework-1/internal/app/product"
	"homework-1/internal/pkg/model"
	"homework-1/internal/repository"
)

type ServiceMethods interface {
	Add(userID int64, sku uint32, count uint16) (total uint16, existed bool)
	Remove(userID int64, sku uint32) error
	Clear(userID int64) error
	Get(userID int64) (model.CartResponse, error)
}

type Service struct {
	repo           repository.CartRepository
	productService product.ProductValidator
}

func New(repo repository.CartRepository, productService product.ProductValidator) *Service {
	return &Service{
		repo:           repo,
		productService: productService,
	}
}

func (s *Service) Add(userID int64, sku uint32, count uint16) (total uint16, existed bool) {
	_, err := s.productService.ValidateProduct(sku)
	if err != nil {
		return 0, false
	}
	return s.repo.AddCart(userID, sku, count)
}

func (s *Service) Remove(userID int64, sku uint32) error {
	err := s.repo.RemoveCart(userID, sku)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Clear(userID int64) error {
	err := s.repo.ClearCart(userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Get(userID int64) (model.CartResponse, error) {
	var totalPrice uint32
	var cartItems []model.CartItem

	userCart, err := s.repo.GetCart(userID)
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
