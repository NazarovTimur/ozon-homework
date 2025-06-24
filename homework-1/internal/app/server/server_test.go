package server

import (
	"github.com/gojuno/minimock/v3"
	"homework-1/internal/app/product/mock"
	"homework-1/internal/pkg/model"
	"homework-1/internal/repository"
	"reflect"
	"testing"
)

func TestServer_Add(t *testing.T) {
	ctrl := minimock.NewController(t)
	mockRepo := mock.NewCartRepositoryMock(ctrl)
	mockService := mock.NewProductValidatorMock(ctrl)

	serverCart := New(mockRepo, mockService)

	userID := int64(25)
	skuID := uint32(10)
	count := uint16(10)
	responseProduct := repository.ProductResponse{
		Name:  "Timi",
		Price: uint32(10),
	}

	mockService.ValidateProductMock.Expect(skuID).Return(&responseProduct, nil)
	mockRepo.AddCartMock.Expect(userID, skuID, count).Return(10, false)

	total, existed := serverCart.Add(userID, skuID, count)
	if total != 10 || existed {
		t.Errorf("Total %d, existed %v", total, existed)
	}
}

func TestServer_Remove(t *testing.T) {
	ctrl := minimock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mock.NewCartRepositoryMock(ctrl)
	mockService := mock.NewProductValidatorMock(ctrl)

	serverCart := New(mockRepo, mockService)
	userID := int64(25)
	skuID := uint32(10)

	mockRepo.RemoveCartMock.Expect(userID, skuID).Return()

	serverCart.Remove(userID, skuID)
}

func TestServer_Clear(t *testing.T) {
	ctrl := minimock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mock.NewCartRepositoryMock(ctrl)
	mockService := mock.NewProductValidatorMock(ctrl)

	serverCart := New(mockRepo, mockService)
	userID := int64(25)

	mockRepo.ClearCartMock.Expect(userID).Return()
	serverCart.Clear(userID)
}

func TestServer_Get(t *testing.T) {
	ctrl := minimock.NewController(t)
	mockRepo := mock.NewCartRepositoryMock(ctrl)
	mockService := mock.NewProductValidatorMock(ctrl)

	serverCart := New(mockRepo, mockService)
	userID := int64(25)
	responseProduct := repository.ProductResponse{
		Name:  "Timi",
		Price: uint32(10),
	}

	var responseCart = model.CartResponse{
		Items: []model.CartItem{
			{SkuID: int64(10),
				Name:  "Timi",
				Count: uint16(1),
				Price: uint32(10)},
		},
		TotalPrice: 10,
	}
	mockService.ValidateProductMock.Expect(uint32(10)).Return(&responseProduct, nil)
	mockRepo.GetCartMock.Expect(userID).Return(map[uint32]uint16{10: 1}, nil)

	answer, err := serverCart.Get(userID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(answer, responseCart) {
		t.Errorf("unexpected response:\n got: %+v\nwant: %+v", answer, responseCart)
	}
}
