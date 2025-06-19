package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"homework-1/internal/app/product"
	"homework-1/internal/app/server"
	"homework-1/internal/config"
	"homework-1/internal/http/handler"
	"homework-1/internal/http/middleware"
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
	cartService := server.New(rep, productService)
	handler := handler.New(cartService, productService)

	http.HandleFunc("POST /user/{user_id}/cart/{sku_id}", handler.AddItemToCart)
	http.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", handler.DeleteItemFromCart)
	http.HandleFunc("DELETE /user/{user_id}/cart", handler.ClearCart)
	http.HandleFunc("GET /user/{user_id}/cart", handler.GetCart)

	loggedMux := middleware.LoggingMiddleware(http.DefaultServeMux)
	fmt.Printf("Сервер запущен на http://localhost:%s\n", cfg.Server.Port)
	http.ListenAndServe(":"+cfg.Server.Port, loggedMux)
}
