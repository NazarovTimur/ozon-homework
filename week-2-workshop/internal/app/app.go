package app

import (
	"context"
	"gitlab.ozon.dev/14/week-2-workshop/internal/clients/product"
	"gitlab.ozon.dev/14/week-2-workshop/internal/domain"
	memory_cart_repo "gitlab.ozon.dev/14/week-2-workshop/internal/repository/memory_cart_repository"
	"log"
	"net/http"
	"time"

	appHttp "gitlab.ozon.dev/14/week-2-workshop/internal/app/http"
	itemAdd "gitlab.ozon.dev/14/week-2-workshop/internal/service/cart/item/add"
)

type (
	mux interface {
		Handle(pattern string, handler http.Handler)
	}
	server interface {
		ListenAndServe() error
		Close() error
	}
	cartStorage interface {
		AddItem(_ context.Context, userID int64, item domain.Item) error
	}
	productProvider interface {
		GetProductInfo(ctx context.Context, sku uint32) (*domain.Product, error)
	}

	App struct {
		config   config
		mux      mux
		server   server
		storage  cartStorage
		products productProvider
	}
)

func NewApp(config config) (*App, error) {
	var (
		mux           = http.NewServeMux()
		products, err = product.New(config.productAddr, config.productToken)
	)

	if err != nil {
		return nil, err
	}

	return &App{
		config:   config,
		mux:      mux,
		server:   &http.Server{Addr: config.addr, Handler: wrapLogger(mux)},
		storage:  memory_cart_repo.NewMemoryStorage(),
		products: products,
	}, nil

}

func (a *App) ListenAndServe() error {
	a.mux.Handle(a.config.path.index, appHttp.NewIndexHandler())
	a.mux.Handle(a.config.path.cartItemAdd, appHttp.NewAddHandler(itemAdd.New(a.storage, a.products), a.config.path.cartItemAdd))

	return a.server.ListenAndServe()

}

func (a *App) Close() error {
	return a.server.Close()
}

func wrapLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, req)
		log.Printf("%s %s %s", req.Method, req.RequestURI, time.Since(start))
	})
}
