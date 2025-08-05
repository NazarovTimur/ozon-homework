package stock

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

type StockTestSuite struct {
	suite.Suite
	ctx   context.Context
	conn  *pgx.Conn
	stock *StockRepository
}

func TestLOMSServiceSuite(t *testing.T) {
	suite.Run(t, new(StockTestSuite))
}

func (s *StockTestSuite) SetupSuite() {
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

	s.stock = New(s.conn, s.conn)
}

func (s *StockTestSuite) SetupTest() {
	_, err := s.conn.Exec(s.ctx, `TRUNCATE stocks, reserved RESTART IDENTITY CASCADE`)
	s.Require().NoError(err)
}

func (s *StockTestSuite) TearDownSuite() {
	_, err := s.conn.Exec(s.ctx, "TRUNCATE stocks, reserved RESTART IDENTITY CASCADE;")
	s.Require().NoError(err)
	_ = s.conn.Close(s.ctx)
}

var (
	testSKU   uint32 = 773297411
	testCount uint16 = 30
)

func (s *StockTestSuite) TestReserve() {
	_, err := s.conn.Exec(s.ctx, `INSERT INTO stocks (sku, count) VALUES ($1, $2)`, testSKU, testCount)
	s.Require().NoError(err)

	check, err := s.stock.Reserve(s.ctx, testSKU, testCount)
	s.Require().NoError(err)
	s.Require().Equal(true, check)
}

func (s *StockTestSuite) TestReserveRemove() {
	_, err := s.conn.Exec(s.ctx, `INSERT INTO stocks (sku, count) VALUES ($1, $2)`, testSKU, testCount)
	s.Require().NoError(err)

	ok, err := s.stock.Reserve(s.ctx, testSKU, testCount)
	s.Require().NoError(err)
	s.Require().True(ok)

	err = s.stock.ReserveRemove(s.ctx, testSKU, testCount)
	s.Require().NoError(err)

	var reserved int
	err = s.conn.QueryRow(s.ctx, `SELECT count FROM reserved WHERE sku = $1`, testSKU).Scan(&reserved)
	s.Require().NoError(err)
}

func (s *StockTestSuite) TestReserveCancel() {
	_, err := s.conn.Exec(s.ctx, `INSERT INTO stocks (sku, count) VALUES ($1, $2)`, testSKU, testCount)
	s.Require().NoError(err)

	ok, err := s.stock.Reserve(s.ctx, testSKU, testCount)
	s.Require().NoError(err)
	s.Require().True(ok)

	err = s.stock.ReserveCancel(s.ctx, testSKU, testCount)
	s.Require().NoError(err)
}

func (s *StockTestSuite) TestGetBySKU() {
	_, err := s.conn.Exec(s.ctx, `INSERT INTO stocks (sku, count) VALUES ($1, $2)`, testSKU, testCount)
	s.Require().NoError(err)

	count, err := s.stock.GetBySKU(s.ctx, testSKU)
	s.Require().NoError(err)
	s.Require().Equal(testCount, count)
}
