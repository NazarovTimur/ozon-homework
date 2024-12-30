package suite

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gitlab.ozon.dev/14/week-2-workshop/internal/repository/real_db_repository"
	"gitlab.ozon.dev/14/week-2-workshop/internal/service/cart/list"
	"log"
	"path/filepath"
	"time"
)

type TCSuite struct {
	suite.Suite
	listHandler *list.Handler
}

func (s *TCSuite) SetupTest() {
	ctx := context.Background()

	dbName := "cart"
	dbUser := "user"
	dbPassword := "pass"
	postgresContainter, err := postgres.Run(ctx,
		"docker.io/postgres:16-alpine",
		postgres.WithInitScripts(filepath.Join("..", "test_data", "init_db.sql")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithLogConsumers(new(StdoutLogConsumer)),
		testcontainers.WithWaitStrategy(
			wait.ForExposedPort(),
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(50*time.Second),
		),
	)

	if err != nil {
		log.Fatalf("failed to run testcontainer %v", err.Error())
	}

	connStr, err := postgresContainter.ConnectionString(ctx, "sslmode=disable")

	storage := real_db_repository.NewRealDBStorage(ctx, connStr)

	s.listHandler = list.New(storage)
}

func (s *TCSuite) TestListItem() {
	ctx := context.Background()
	userID := int64(123)
	itemList, _ := s.listHandler.ListItem(ctx, userID)

	require.Equal(s.T(), len(itemList), 1)
	require.EqualValues(s.T(), itemList[0].Count, 2)
	require.EqualValues(s.T(), itemList[0].SKU, 100000)
}

type StdoutLogConsumer struct {
}

func (lc *StdoutLogConsumer) Accept(l testcontainers.Log) {
	fmt.Print(string(l.Content))
}
