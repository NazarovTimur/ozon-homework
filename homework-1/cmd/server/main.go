package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"homework-1/internal/app/product"
	"homework-1/internal/app/service"
	"homework-1/internal/http/handler"
	"homework-1/internal/http/middleware"
	"homework-1/internal/pkg/config"
	"homework-1/internal/pkg/retry"
	"homework-1/internal/repository"
	"log"
	"net/http"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or failed to load")
	}
	cfg := config.NewConfig()

	retryClient := retry.New(cfg.Retry.Count, cfg.Retry.Delay)
	productService := product.NewProductService(retryClient, cfg.Product.Url, cfg.Product.Token)
	rep := repository.New()
	cartService := service.New(rep, productService)
	h := handler.New(cartService, productService)

	http.HandleFunc("POST /user/{user_id}/cart/{sku_id}", h.AddItemToCart)
	http.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", h.DeleteItemFromCart)
	http.HandleFunc("DELETE /user/{user_id}/cart", h.ClearCart)
	http.HandleFunc("GET /user/{user_id}/cart", h.GetCart)

	loggedMux := middleware.LoggingMiddleware(http.DefaultServeMux)
	fmt.Printf("Сервер запущен на http://localhost:%s\n", cfg.Server.Port)
	http.ListenAndServe(":"+cfg.Server.Port, loggedMux)
}
