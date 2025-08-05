package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"homework-1/loms/internal/pkg/constants"
	"homework-1/loms/internal/pkg/errorx"
	"homework-1/loms/internal/pkg/model"
)

type OrderInterface interface {
	Create(ctx context.Context, userID int64, items []model.OrderItem) (int64, error)
	SetStatus(ctx context.Context, orderID int64, status string) error
	GetByID(ctx context.Context, orderID int64) (string, int64, []model.OrderItem, error)
}

type OrderRepository struct {
	connMaster  *pgx.Conn
	connReplica *pgx.Conn
}

func New(connMaster *pgx.Conn, connReplica *pgx.Conn) *OrderRepository {
	return &OrderRepository{
		connMaster:  connMaster,
		connReplica: connReplica,
	}
}

const (
	insertOrder  = `INSERT INTO orders (user_id, status) VALUES ($1, $2) RETURNING id;`
	insertItems  = `INSERT INTO items (sku, count) VALUES ($1, $2) ON CONFLICT (sku) DO UPDATE SET count = items.count + excluded.count RETURNING id;`
	insertOrders = `INSERT INTO orders_items (order_id, items_id) VALUES ($1, $2);`
	updateStatus = `UPDATE orders SET status = $1 WHERE id = $2;`
	getOrderByID = `SELECT id, user_id, status FROM orders WHERE id = $1;`

	getOrderItems = `
		SELECT i.sku, i.count
		FROM orders_items oi
		JOIN items i ON oi.items_id = i.id
		WHERE oi.order_id = $1;
	`
)

func (o *OrderRepository) Create(ctx context.Context, userID int64, items []model.OrderItem) (int64, error) {
	if err := ctx.Err(); err != nil {
		return 0, errorx.ErrContextCanceled
	}

	tx, err := o.connMaster.Begin(ctx)
	if err != nil {
		return 0, errorx.ErrTransactionBegin
	}
	defer tx.Rollback(ctx)

	var orderID int64
	err = tx.QueryRow(ctx, insertOrder, userID, constants.StatusNew).Scan(&orderID)
	if err != nil {
		return 0, errorx.ErrOrderCreating
	}
	for _, item := range items {
		var itemID int64
		err = tx.QueryRow(ctx, insertItems, item.SKU, item.Count).Scan(&itemID)
		if err != nil {
			return 0, errorx.ErrOrderCreatingItem
		}

		_, err = tx.Exec(ctx, insertOrders, orderID, itemID)
		if err != nil {
			return 0, errorx.ErrOrderCreating
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, errorx.ErrTransactionCommit
	}

	return orderID, nil
}

func (o *OrderRepository) SetStatus(ctx context.Context, orderID int64, status string) error {
	if err := ctx.Err(); err != nil {
		return errorx.ErrContextCanceled
	}

	if !(status == constants.StatusAwaitingPayment || status == constants.StatusFailed || status == constants.StatusPayed || status == constants.StatusCancelled) {
		return errorx.ErrInvalidOrderStatus
	}

	cmdTag, err := o.connMaster.Exec(ctx, updateStatus, status, orderID)
	if err != nil {
		return errorx.ErrSettingStatus
	}

	if cmdTag.RowsAffected() == 0 {
		return errorx.ErrInvalidOrderStatus
	}

	return nil
}

func (o *OrderRepository) GetByID(ctx context.Context, orderID int64) (string, int64, []model.OrderItem, error) {
	if err := ctx.Err(); err != nil {
		return "", 0, []model.OrderItem{}, errorx.ErrContextCanceled
	}

	var order model.Order
	err := o.connReplica.QueryRow(ctx, getOrderByID, orderID).Scan(&order.ID, &order.UserID, &order.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", 0, []model.OrderItem{}, errorx.ErrOrderNotFound
		}
		return "", 0, []model.OrderItem{}, fmt.Errorf("error getting order: %w", err)
	}

	rows, err := o.connReplica.Query(ctx, getOrderItems, orderID)
	if err != nil {
		return "", 0, []model.OrderItem{}, fmt.Errorf("error getting order items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item model.OrderItem
		if err = rows.Scan(&item.SKU, &item.Count); err != nil {
			return "", 0, nil, fmt.Errorf("failed to scan item: %w", err)
		}
		order.Items = append(order.Items, item)
	}

	if err = rows.Err(); err != nil {
		return "", 0, []model.OrderItem{}, fmt.Errorf("rows iteration error: %w", err)
	}

	return order.Status, order.UserID, order.Items, nil
}
