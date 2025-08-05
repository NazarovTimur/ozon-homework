package stock

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"homework-1/loms/internal/pkg/errorx"
)

type StockInterface interface {
	Reserve(ctx context.Context, sku uint32, count uint16) (bool, error)
	ReserveRemove(ctx context.Context, sku uint32, count uint16) error
	ReserveCancel(ctx context.Context, sku uint32, count uint16) error
	GetBySKU(ctx context.Context, sku uint32) (uint16, error)
}
type StockRepository struct {
	connMaster  *pgx.Conn
	connReplica *pgx.Conn
}

func New(connMaster *pgx.Conn, connReplica *pgx.Conn) *StockRepository {
	return &StockRepository{
		connMaster:  connMaster,
		connReplica: connReplica,
	}
}

const (
	selectCount    = `SELECT count from stocks WHERE sku = $1 FOR UPDATE;`
	updateCount    = `UPDATE stocks SET count = count - $2 WHERE sku = $1;`
	returnCount    = `UPDATE stocks SET count = count + $2 WHERE sku = $1;`
	insertReserved = `INSERT INTO reserved (sku, count) VALUES ($1, $2) ON CONFLICT (sku) DO UPDATE SET count = reserved.count + $2;`
	updateReserve  = `UPDATE reserved SET count = count - $2 WHERE sku = $1 AND count >= $2;`
	getBySKU       = `SELECT count FROM stocks WHERE sku = $1;`
)

func (s *StockRepository) Reserve(ctx context.Context, sku uint32, count uint16) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, errorx.ErrContextCanceled
	}

	tx, err := s.connMaster.Begin(ctx)
	if err != nil {
		return false, errorx.ErrTransactionBegin
	}
	defer tx.Rollback(ctx)

	var countDB int
	err = tx.QueryRow(ctx, selectCount, sku).Scan(&countDB)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("query stock count: %w", err)
	}

	if uint16(countDB) < count {
		return false, errorx.ErrInsufficientStock
	}

	_, err = tx.Exec(ctx, updateCount, sku, count)
	if err != nil {
		return false, fmt.Errorf("update stock count: %w", err)
	}

	_, err = tx.Exec(ctx, insertReserved, sku, count)
	if err != nil {
		return false, fmt.Errorf("insert into reserved: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return false, errorx.ErrTransactionCommit
	}
	return true, nil
}

func (s *StockRepository) ReserveRemove(ctx context.Context, sku uint32, count uint16) error {
	if err := ctx.Err(); err != nil {
		return errorx.ErrContextCanceled
	}

	cmdTag, err := s.connMaster.Exec(ctx, updateReserve, sku, count)
	if err != nil {
		return fmt.Errorf("update reserve: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return errorx.ErrReserveNotFound
	}

	return nil
}

func (s *StockRepository) ReserveCancel(ctx context.Context, sku uint32, count uint16) error {
	if err := ctx.Err(); err != nil {
		return errorx.ErrContextCanceled
	}

	cmdTag, err := s.connMaster.Exec(ctx, updateReserve, sku, count)
	if err != nil {
		return fmt.Errorf("update reserve: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return errorx.ErrReserveNotFound
	}

	_, err = s.connMaster.Exec(ctx, returnCount, sku, count)
	if err != nil {
		return fmt.Errorf("return stock: %w", err)
	}

	return nil
}

func (s *StockRepository) GetBySKU(ctx context.Context, sku uint32) (uint16, error) {
	if err := ctx.Err(); err != nil {
		return 0, errorx.ErrContextCanceled
	}

	var countDB int
	err := s.connReplica.QueryRow(ctx, getBySKU, sku).Scan(&countDB)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, errorx.ErrStockNotFound
	}
	if err != nil {
		return 0, fmt.Errorf("get stock by SKU: %w", err)
	}

	return uint16(countDB), nil
}
