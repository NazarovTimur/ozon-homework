package server

import (
	"homework-1/internal/app/product"
	"homework-1/internal/pkg/model"
	"homework-1/internal/repository"
)

type Server struct {
	repo           repository.CartRepository
	productService product.ProductValidator
}

func New(repo repository.CartRepository, productService product.ProductValidator) *Server {
	return &Server{
		repo:           repo,
		productService: productService,
	}
}

func (s *Server) Add(userID int64, sku uint32, count uint16) (total uint16, existed bool) {
	_, err := s.productService.ValidateProduct(sku)
	if err != nil {
		return 0, false
	}
	return s.repo.AddCart(userID, sku, count)
}

func (s *Server) Remove(userID int64, sku uint32) {
	s.repo.RemoveCart(userID, sku)
}

func (s *Server) Clear(userID int64) {
	s.repo.ClearCart(userID)
}

func (s *Server) Get(userID int64) (model.CartResponse, error) {
	var totalPrice uint32
	var cartItems []model.CartItem

	userCart, err := s.repo.GetCart(userID)
	if err != nil {
		return model.CartResponse{}, err
	}
	for sku, count := range userCart {
		productData, err := s.productService.ValidateProduct(sku)
		if err != nil {
			return model.CartResponse{}, err
		}
		totalPrice += uint32(count) * productData.Price
		cartItems = append(cartItems, model.CartItem{Name: productData.Name, SkuID: int64(sku), Count: count, Price: productData.Price})
	}
	return model.CartResponse{Items: cartItems, TotalPrice: totalPrice}, nil
}
