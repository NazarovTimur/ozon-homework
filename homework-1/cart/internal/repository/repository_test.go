package repository

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"homework-1/cart/internal/pkg/model"
	"sync"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestCart_AddCart(t *testing.T) {
	testCart := make(map[int64]map[uint32]uint16)
	testCart[123] = make(map[uint32]uint16)
	testCart[123][52] = 6
	ctx := context.Background()

	cart := New()
	_, _, err := cart.AddCart(ctx, 123, 52, 6)
	if err != nil {
		t.Errorf("Error adding cart: %v", err)
	}
	require.Equal(t, cart.carts, testCart)
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
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &Cart{
				mu:    sync.RWMutex{},
				carts: tt.cart,
			}
			_, _, _ = repo.AddCart(ctx, tt.UserID, tt.SKU, tt.count)

			require.Equal(t, tt.expectedResult, repo.carts)
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
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &Cart{
				mu:    sync.RWMutex{},
				carts: tt.cart,
			}
			_ = repo.RemoveCart(ctx, tt.UserID, tt.SKU)

			require.Equal(t, tt.expectedResult, repo.carts)
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
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &Cart{
				mu:    sync.RWMutex{},
				carts: tt.cart,
			}
			_ = repo.ClearCart(ctx, tt.UserID)

			require.Equal(t, tt.expectedResult, repo.carts)
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
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &Cart{
				mu:    sync.RWMutex{},
				carts: tt.cart,
			}
			res, _ := repo.GetCart(ctx, tt.UserID)
			require.Equal(t, tt.expectedResult, res)
		})
	}
}

func TestConcurrentAccess(t *testing.T) {
	t.Parallel()
	repo := New()
	test := struct {
		ctx   context.Context
		SKU   uint32
		count uint16
	}{
		ctx:   context.Background(),
		SKU:   14,
		count: 6,
	}
	wg := sync.WaitGroup{}
	n := 1000

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int64) {
			defer wg.Done()
			repo.AddCart(test.ctx, i, test.SKU, test.count)
		}(int64(i))
	}

	wg.Wait()
	for i := 0; i < n; i++ {
		userID := int64(i)
		res, err := repo.GetCart(test.ctx, userID)
		if err != nil {
			t.Errorf("unexpected error for user %d: %v", userID, err)
			continue
		}
		if res[0].SkuID != test.SKU || res[0].Count != test.count {
			t.Errorf("unexpected cart item for user %d: got %+v, want SKU %d count %d",
				userID, res[0], test.SKU, test.count)
		}
	}
}

func TestCart_AddAndClear_Concurrent(t *testing.T) {
	t.Parallel()
	repo := New()
	userID := int64(3)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			_, _, err := repo.AddCart(context.Background(), int64(i)+userID, uint32(i), 1)
			if err != nil {
				t.Errorf("unexpected error for user %d: %v", userID, err)
			}
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			time.Sleep(25 * time.Microsecond)
			_ = repo.ClearCart(context.Background(), int64(i)*userID)
		}
	}()

	wg.Wait()

	cart, _ := repo.GetCart(context.Background(), userID)
	require.LessOrEqual(t, len(cart), 1000, "cart size should be limited due to concurrent clears")
}
