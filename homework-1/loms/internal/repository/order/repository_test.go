package order

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"homework-1/loms/internal/pkg/model"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

type OrderTestSuite struct {
	suite.Suite
	ctx   context.Context
	conn  *pgx.Conn
	order *OrderRepository
}

func TestLOMSServiceSuite(t *testing.T) {
	suite.Run(t, new(OrderTestSuite))
}

func (s *OrderTestSuite) SetupSuite() {
	s.ctx = context.Background()
	var err error

	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(filename), "../../../../")
	envPath := filepath.Join(root, ".env.db.test")

	err = godotenv.Load(envPath)
	s.Require().NoError(err, "Не удалось загрузить .env.db.test")

	dsn := os.Getenv("TEST_DB_DSN")
	if !strings.Contains(dsn, "loms_test") {
		s.T().Fatal("Tests can only be run on loms_test!")
	}

	fmt.Println("Connecting to test DB:", dsn)
	s.conn, err = pgx.Connect(s.ctx, dsn)
	s.Require().NoError(err)

	s.order = New(s.conn, s.conn)
}

func (s *OrderTestSuite) SetupTest() {
	_, err := s.conn.Exec(s.ctx, `TRUNCATE orders_items, orders, items RESTART IDENTITY CASCADE`)
	s.Require().NoError(err)
}

func (s *OrderTestSuite) TearDownSuite() {
	_, err := s.conn.Exec(s.ctx, "TRUNCATE orders_items, orders, items RESTART IDENTITY CASCADE;")
	s.Require().NoError(err)
	_ = s.conn.Close(s.ctx)
}

var (
	testUserID int64 = 52
	testItems        = []model.OrderItem{{SKU: 1002, Count: 2}}
	testStatus       = "awaiting payment"
)

func (s *OrderTestSuite) TestCreate() {
	orderID, err := s.order.Create(s.ctx, testUserID, testItems)
	s.Require().NoError(err)
	s.Assert().Equal(orderID, int64(1))
}

func (s *OrderTestSuite) TestSetStatus() {
	orderID, err := s.order.Create(s.ctx, testUserID, testItems)
	s.Require().NoError(err)

	err = s.order.SetStatus(s.ctx, orderID, testStatus)
	s.Require().NoError(err)
}

func (s *OrderTestSuite) TestGetByID() {
	orderID, err := s.order.Create(s.ctx, testUserID, testItems)
	s.Require().NoError(err)

	status, UserID, items, err := s.order.GetByID(s.ctx, orderID)
	s.Require().NoError(err)
	s.Assert().Equal(status, "new")
	s.Assert().Equal(UserID, testUserID)
	s.Assert().Equal(items, testItems)
}
