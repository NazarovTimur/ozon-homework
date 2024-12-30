package real_db_repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"gitlab.ozon.dev/14/week-2-workshop/internal/domain"
	"log"
	"os"
)

type (
	RealDBStorage struct {
		conn *pgx.Conn
	}
)

func NewRealDBStorage(ctx context.Context, connStr string) *RealDBStorage {
	log.Println(connStr)
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil
	}

	if err := conn.Ping(ctx); err != nil {
		log.Fatalf("Error with connect %v", err.Error())
	}
	log.Println("Success connect")

	return &RealDBStorage{
		conn: conn,
	}
}

func (m *RealDBStorage) AddItem(_ context.Context, userID int64, item domain.Item) error {
	return nil
}

func (m *RealDBStorage) ListItem(ctx context.Context, userID int64) []domain.Item {
	rows, err := m.conn.Query(ctx, "SELECT item_id, item_count FROM user_items")
	if err != nil {
		log.Fatalf("error with select item %v", err.Error())
	}
	defer rows.Close()

	result := make([]domain.Item, 0)
	var (
		itemId    uint32
		itemCount uint16
	)

	for rows.Next() {
		err := rows.Scan(&itemId, &itemCount)
		if err != nil {
			panic(err)
		}

		result = append(result, domain.Item{
			SKU:   itemId,
			Count: itemCount,
		})
	}

	return result
}
