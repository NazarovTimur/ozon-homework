package cart

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"homework-1/cart/internal/app/product/mock"
	model2 "homework-1/cart/internal/pkg/model"
	mock2 "homework-1/cart/internal/repository/mock"
	mockLoms "homework-1/loms/internal/app/loms/mock"
	"homework-1/loms/internal/proto/loms"
	"reflect"
	"testing"
)

func TestServer_Add(t *testing.T) {
	ctrl := minimock.NewController(t)
	mockRepo := mock2.NewCartRepositoryMock(ctrl)
	mockService := mock.NewProductValidatorMock(ctrl)
	mockLomsService := mockLoms.NewLomsServiceClientMock(ctrl)

	serverCart := New(mockRepo, mockService, mockLomsService)

	ctx := context.Background()
	userID := int64(25)
	skuID := uint32(10)
	count := uint16(10)
	responseProduct := model2.ProductResponse{
		Name:  "Timi",
		Price: uint32(10),
	}
	requestLoms := pb.StocksInfoRequest{Sku: skuID}
	responseLoms := pb.StocksInfoResponse{Count: uint64(count)}

	mockService.ValidateProductMock.Expect(skuID).Return(&responseProduct, nil)
	mockLomsService.StocksInfoMock.Expect(ctx, &requestLoms).Return(&responseLoms, nil)
	mockRepo.AddCartMock.Expect(ctx, userID, skuID, count).Return(10, false, nil)

	total, existed, err := serverCart.Add(ctx, userID, skuID, count)
	if total != 10 || existed || err != nil {
		t.Errorf("Total %d, existed %v", total, existed)
	}
}

func TestServer_Remove(t *testing.T) {
	ctrl := minimock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mock2.NewCartRepositoryMock(ctrl)
	mockService := mock.NewProductValidatorMock(ctrl)
	mockLomsService := mockLoms.NewLomsServiceClientMock(ctrl)

	serverCart := New(mockRepo, mockService, mockLomsService)
	userID := int64(25)
	skuID := uint32(10)
	ctx := context.Background()

	mockRepo.RemoveCartMock.Expect(ctx, userID, skuID).Return(nil)

	serverCart.Remove(ctx, userID, skuID)
}

func TestServer_Clear(t *testing.T) {
	ctrl := minimock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mock2.NewCartRepositoryMock(ctrl)
	mockService := mock.NewProductValidatorMock(ctrl)
	mockLomsService := mockLoms.NewLomsServiceClientMock(ctrl)
	ctx := context.Background()

	serverCart := New(mockRepo, mockService, mockLomsService)
	userID := int64(25)

	mockRepo.ClearCartMock.Expect(ctx, userID).Return(nil)
	serverCart.Clear(ctx, userID)
}

func TestServer_Get(t *testing.T) {
	ctrl := minimock.NewController(t)
	mockRepo := mock2.NewCartRepositoryMock(ctrl)
	mockService := mock.NewProductValidatorMock(ctrl)
	mockLomsService := mockLoms.NewLomsServiceClientMock(ctrl)
	ctx := context.Background()

	serverCart := New(mockRepo, mockService, mockLomsService)
	userID := int64(25)
	responseProduct := model2.ProductResponse{
		Name:  "Timi",
		Price: uint32(10),
	}

	var responseCart = model2.CartResponse{
		Items: []model2.CartItem{
			{SkuID: int64(10),
				Name:  "Timi",
				Count: uint16(1),
				Price: uint32(10)},
		},
		TotalPrice: 10,
	}
	mockService.ValidateProductMock.Expect(uint32(10)).Return(&responseProduct, nil)
	mockRepo.GetCartMock.Expect(ctx, userID).Return([]model2.ItemCart{{SkuID: 10, Count: 1}}, nil)

	answer, err := serverCart.Get(ctx, userID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(answer, responseCart) {
		t.Errorf("unexpected response:\n got: %+v\nwant: %+v", answer, responseCart)
	}
}
