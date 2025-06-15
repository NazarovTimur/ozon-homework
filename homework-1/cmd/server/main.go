package main

import (
	"fmt"
	"homework-1/internal/app/cart"
	"homework-1/internal/app/product"
	"homework-1/internal/http/handler"
	"homework-1/internal/http/middleware"
	"homework-1/internal/pkg/retry"
	"net/http"
	"time"
)

func main() {
	retryClient := retry.New(3, 200*time.Millisecond)
	productService := product.NewProductService(retryClient, "http://route256.pavl.uk:8080/get_product", "testtoken")
	cartService := cart.New()
	handler := handler.New(cartService, productService)

	http.HandleFunc("POST /user/{user_id}/cart/{sku_id}", handler.AddItemToCart)
	http.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", handler.DeleteItemFromCart)
	http.HandleFunc("DELETE /user/{user_id}/cart", handler.ClearCart)
	http.HandleFunc("GET /user/{user_id}/cart", handler.GetCart)

	loggedMux := middleware.LoggingMiddleware(http.DefaultServeMux)
	fmt.Println("Сервер запущен на http://localhost:8082")
	http.ListenAndServe(":8082", loggedMux)

}
