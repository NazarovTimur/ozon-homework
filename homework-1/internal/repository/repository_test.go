package repository

import (
	"homework-1/internal/pkg/model"
	"reflect"
	"testing"
)

func TestCart_AddCart(t *testing.T) {
	testCart := make(map[int64]map[uint32]uint16)
	testCart[123] = make(map[uint32]uint16)
	testCart[123][52] = 6

	cart := New()
	cart.AddCart(123, 52, 6)
	if !reflect.DeepEqual(testCart, cart.carts) {
		t.Errorf("unexpected response:\n got: %v\nwant: %v", cart.carts, testCart)
	}
}

func TestCart_AddCart2(t *testing.T) {
	tests := []struct {
		name           string
		cart           map[int64]map[uint32]uint16
		UserID         int64
		SKU            uint32
		count          uint16
		expectedResult map[int64]map[uint32]uint16
	}{{
		name:           "add cart",
		cart:           map[int64]map[uint32]uint16{44: {21: 5}},
		UserID:         44,
		SKU:            21,
		count:          10,
		expectedResult: map[int64]map[uint32]uint16{44: {21: 15}},
	},
		{
			name:           "add new SKU",
			cart:           map[int64]map[uint32]uint16{10: {10: 10}},
			UserID:         10,
			SKU:            20,
			count:          20,
			expectedResult: map[int64]map[uint32]uint16{10: {10: 10, 20: 20}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &Cart{carts: tt.cart}
			repo.AddCart(tt.UserID, tt.SKU, tt.count)

			if !reflect.DeepEqual(repo.carts, tt.expectedResult) {
				t.Errorf("unexpected response:\n got: %v\nwant: %v", repo.carts, tt.expectedResult)
			}
		})
	}
}

func TestCart_RemoveCart(t *testing.T) {
	tests := []struct {
		name           string
		cart           map[int64]map[uint32]uint16
		UserID         int64
		SKU            uint32
		expectedResult map[int64]map[uint32]uint16
	}{
		{
			name:           "remove cart",
			cart:           map[int64]map[uint32]uint16{44: {21: 5}},
			UserID:         44,
			SKU:            21,
			expectedResult: map[int64]map[uint32]uint16{44: {}},
		},
		{
			name:           "remove 0",
			cart:           map[int64]map[uint32]uint16{44: {33: 4}},
			UserID:         44,
			SKU:            21,
			expectedResult: map[int64]map[uint32]uint16{44: {33: 4}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &Cart{carts: tt.cart}
			repo.RemoveCart(tt.UserID, tt.SKU)

			if !reflect.DeepEqual(repo.carts, tt.expectedResult) {
				t.Errorf("unexpected response:\n got: %v\nwant: %v", repo.carts, tt.expectedResult)
			}
		})
	}
}

func TestCart_Clear(t *testing.T) {
	tests := []struct {
		name           string
		cart           map[int64]map[uint32]uint16
		UserID         int64
		expectedResult map[int64]map[uint32]uint16
	}{
		{
			name:           "remove cart",
			cart:           map[int64]map[uint32]uint16{44: {21: 5}},
			UserID:         44,
			expectedResult: map[int64]map[uint32]uint16{},
		},
		{
			name:           "remove 0",
			cart:           map[int64]map[uint32]uint16{42: {12: 52}},
			UserID:         44,
			expectedResult: map[int64]map[uint32]uint16{42: {12: 52}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &Cart{carts: tt.cart}
			repo.ClearCart(tt.UserID)

			if !reflect.DeepEqual(repo.carts, tt.expectedResult) {
				t.Errorf("unexpected response:\n got: %v\nwant: %v", repo.carts, tt.expectedResult)
			}
		})
	}
}

func TestCart_GetCart(t *testing.T) {
	tests := []struct {
		name           string
		cart           map[int64]map[uint32]uint16
		UserID         int64
		expectedResult []model.ItemCart
	}{
		{
			name:   "get cart",
			cart:   map[int64]map[uint32]uint16{44: {21: 5}},
			UserID: 44,
			expectedResult: []model.ItemCart{{
				SkuID: 21,
				Count: 5,
			}},
		},
		{
			name:           "get 0",
			cart:           map[int64]map[uint32]uint16{42: {12: 52}},
			UserID:         44,
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &Cart{carts: tt.cart}
			res, _ := repo.GetCart(tt.UserID)
			if !reflect.DeepEqual(res, tt.expectedResult) {
				t.Errorf("unexpected response:\n got: %v\nwant: %v", res, tt.expectedResult)
			}
		})
	}
}
