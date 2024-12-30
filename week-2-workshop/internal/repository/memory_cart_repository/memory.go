package memory_cart_repo

import (
	"context"
	"gitlab.ozon.dev/14/week-2-workshop/internal/domain"
	"sync"
)

// itemsMap is index sku -> item.
type (
	itemsMap map[uint32]domain.Item

	MemoryStorage struct {
		items map[int64]itemsMap
		mtx   sync.RWMutex
	}
)

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		items: make(map[int64]itemsMap),
		mtx:   sync.RWMutex{},
	}
}

func (m *MemoryStorage) AddItem(_ context.Context, userID int64, item domain.Item) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if m.items[userID] == nil {
		m.items[userID] = itemsMap{}
	}

	m.items[userID][item.SKU] = domain.Item{
		SKU:   item.SKU,
		Count: m.items[userID][item.SKU].Count + item.Count,
	}
	//for key, val := range m.items[userID] {
	//	log.Printf("%v:%v\n", key, val)
	//}
	return nil
}

func (m *MemoryStorage) ListItem(_ context.Context, userID int64) []domain.Item {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	itemMap := m.items[userID]

	result := make([]domain.Item, 0, len(itemMap))

	for _, value := range itemMap {
		result = append(result, value)
	}

	return result
}
