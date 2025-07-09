package stock

import (
	"encoding/json"
	"fmt"
	"homework-1/internal/data"
	"homework-1/internal/pkg/model"
)

func NewStockFromJSON() (*StockRepository, error) {
	var stockItems []model.StockItem
	if err := json.Unmarshal(data.StockData, &stockItems); err != nil {
		return nil, fmt.Errorf("failed to parse stock data: %w", err)
	}

	stockData := make(map[uint32]uint16)
	reserveData := make(map[uint32]uint16)

	for _, item := range stockItems {
		stockData[uint32(item.SKU)] = item.TotalCount - item.Reserved
		reserveData[uint32(item.SKU)] = item.Reserved
	}

	return &StockRepository{
		stocks:   stockData,
		reserved: reserveData,
	}, nil
}
