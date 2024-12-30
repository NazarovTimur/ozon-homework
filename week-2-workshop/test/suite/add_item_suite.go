package suite

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.ozon.dev/14/week-2-workshop/internal/clients/product"
	"gitlab.ozon.dev/14/week-2-workshop/internal/domain"
	memory_cart_repo "gitlab.ozon.dev/14/week-2-workshop/internal/repository/memory_cart_repository"
	"gitlab.ozon.dev/14/week-2-workshop/internal/service/cart/item/add"
	"gitlab.ozon.dev/14/week-2-workshop/internal/service/cart/list"
)

type ItemS struct {
	suite.Suite
	addHandler  *add.Handler
	listHandler *list.Handler
}

func (s *ItemS) SetupSuite() {
	storage := memory_cart_repo.NewMemoryStorage()
	productClient, err := product.New("http://route256.pavl.uk:8080", "testtoken")
	require.NoError(s.T(), err)

	s.addHandler = add.New(storage, productClient)
	s.listHandler = list.New(storage)
}

func (s *ItemS) TestAddItem() {
	ctx := context.Background()

	userID := int64(123)
	item1 := domain.Item{
		SKU:   773297411,
		Count: 2,
	}

	err := s.addHandler.AddItem(ctx, userID, item1)

	require.NoError(s.T(), err)

	itemList, _ := s.listHandler.ListItem(ctx, userID)

	require.Equal(s.T(), len(itemList), 1)
	require.Equal(s.T(), itemList[0].SKU, item1.SKU)
	require.Equal(s.T(), itemList[0].Count, item1.Count)

	item2 := domain.Item{
		SKU:   1148162,
		Count: 1,
	}

	err = s.addHandler.AddItem(ctx, userID, item2)

	require.NoError(s.T(), err)

	itemList, _ = s.listHandler.ListItem(ctx, userID)

	require.Equal(s.T(), len(itemList), 2)
	require.Equal(s.T(), itemList[0].SKU, item1.SKU)
	require.Equal(s.T(), itemList[0].Count, item1.Count)
	require.Equal(s.T(), itemList[1].SKU, item2.SKU)
	require.Equal(s.T(), itemList[1].Count, item2.Count)
}
