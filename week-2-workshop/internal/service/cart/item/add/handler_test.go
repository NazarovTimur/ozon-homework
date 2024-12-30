package add

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/14/week-2-workshop/internal/domain"
	"gitlab.ozon.dev/14/week-2-workshop/internal/service/cart/item/add/mock"
	"testing"
)

func TestHandler_AddItem(t *testing.T) {
	ctx := context.Background()
	ctrl := minimock.NewController(t)
	productMock := mock.NewProductServiceMock(ctrl)
	repoMock := mock.NewRepositoryMock(ctrl)

	addHandler := New(repoMock, productMock)

	userID := int64(100)
	item := domain.Item{
		SKU:   1000,
		Count: 10,
	}
	product := domain.Product{
		Name:  "Книга",
		Price: 100,
	}

	productMock.GetProductInfoMock.Expect(ctx, 1000).Return(&product, nil)
	repoMock.AddItemMock.Expect(ctx, 100, item).Return(nil)
	err := addHandler.AddItem(ctx, userID, item)

	require.NoError(t, err)

	productInfoCount := productMock.GetProductInfoMock.Expect(ctx, 1000).Return(&product, nil).GetProductInfoAfterCounter()

	err = addHandler.AddItem(ctx, userID, item)
	err = addHandler.AddItem(ctx, userID, item)

	require.EqualValues(t, productInfoCount, 1)
}

func TestHandler_AddItem_Table(t *testing.T) {
	ctx := context.Background()

	type fields struct {
		productMock *mock.ProductServiceMock
		repoMock    *mock.RepositoryMock
	}

	type data struct {
		name    string
		userID  int64
		item    domain.Item
		product *domain.Product
		wantErr error
	}

	testData := []data{
		{
			name:   "valid add item",
			userID: 123,
			item: domain.Item{
				SKU:   100,
				Count: 2,
			},
			product: &domain.Product{
				Name:  "Книга",
				Price: 400,
			},
			wantErr: nil,
		},
		{
			name:   "product not found",
			userID: 123,
			item: domain.Item{
				SKU:   100,
				Count: 2,
			},
			product: nil,
			wantErr: ErrInvalidSKU,
		},
	}

	ctrl := minimock.NewController(t)
	productMock := mock.NewProductServiceMock(ctrl)
	repoMock := mock.NewRepositoryMock(ctrl)

	addHandler := New(repoMock, productMock)

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			productMock.GetProductInfoMock.Expect(ctx, tt.item.SKU).Return(tt.product, nil)
			repoMock.AddItemMock.Expect(ctx, tt.userID, tt.item).Return(nil)

			err := addHandler.AddItem(ctx, tt.userID, tt.item)

			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
