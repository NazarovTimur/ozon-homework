package order

import (
	"context"
	"homework-1/internal/pkg/constants"
	"homework-1/internal/pkg/errorx"
	"homework-1/internal/pkg/model"
	"sync"
)

type OrderInterface interface {
	Create(ctx context.Context, userID int64, items []model.OrderItem) (int64, error)
	SetStatus(ctx context.Context, orderID int64, status string) error
	GetByID(ctx context.Context, orderID int64) (string, int64, []model.OrderItem, error)
}

type OrderRepository struct {
	mu        sync.RWMutex
	orders    map[int64]Order
	currentID int64
}

func New() *OrderRepository {
	return &OrderRepository{
		orders: make(map[int64]Order),
	}
}

type Order struct {
	ID     int64
	UserID int64
	Items  []model.OrderItem
	Status string
}

func (o *OrderRepository) Create(ctx context.Context, userID int64, items []model.OrderItem) (int64, error) {
	if err := ctx.Err(); err != nil {
		return 0, errorx.ErrContextCanceled
	}

	o.mu.Lock()
	defer o.mu.Unlock()
	o.currentID++
	newOrder := Order{
		ID:     o.currentID,
		UserID: userID,
		Items:  items,
		Status: constants.StatusNew,
	}
	if o.orders == nil {
		return 0, errorx.ErrRepositoryNotInitialized
	}
	o.orders[newOrder.ID] = newOrder
	return newOrder.ID, nil
}

func (o *OrderRepository) SetStatus(ctx context.Context, orderID int64, status string) error {
	if err := ctx.Err(); err != nil {
		return errorx.ErrContextCanceled
	}

	o.mu.Lock()
	defer o.mu.Unlock()
	order, ok := o.orders[orderID]
	if ok != true {
		return errorx.ErrOrderNotFound
	}
	if status == constants.StatusAwaitingPayment || status == constants.StatusFailed || status == constants.StatusPayed || status == constants.StatusCancelled {
		order.Status = status
		o.orders[orderID] = order
		return nil
	}
	return errorx.ErrInvalidOrderStatus
}

func (o *OrderRepository) GetByID(ctx context.Context, orderID int64) (string, int64, []model.OrderItem, error) {
	if err := ctx.Err(); err != nil {
		return "", 0, []model.OrderItem{}, errorx.ErrContextCanceled
	}

	o.mu.RLock()
	defer o.mu.RUnlock()
	if order, ok := o.orders[orderID]; ok {
		return order.Status, order.UserID, order.Items, nil
	}
	return "", 0, []model.OrderItem{}, errorx.ErrOrderIDNotFound
}
