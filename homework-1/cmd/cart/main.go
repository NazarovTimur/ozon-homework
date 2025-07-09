package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"homework-1/internal/app/cart"
	"homework-1/internal/app/product"
	"homework-1/internal/http/cart/handler"
	"homework-1/internal/http/cart/middleware"
	"homework-1/internal/pkg/config"
	"homework-1/internal/pkg/constants"
	"homework-1/internal/pkg/retry"
	pb "homework-1/internal/proto/loms"
	repoCart "homework-1/internal/repository/cart"
	"log"
	"net/http"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or failed to load")
	}
	cfg := config.NewConfig()

	conn, err := grpc.NewClient(constants.GrpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	lomsClient := pb.NewLomsServiceClient(conn)

	retryClient := retry.New(cfg.Retry.Count, cfg.Retry.Delay)
	productService := product.NewProductService(retryClient, cfg.Product.Url, cfg.Product.Token)
	rep := repoCart.New()
	cartService := cart.New(rep, productService, lomsClient)
	h := handler.New(cartService, productService)

	http.HandleFunc("POST /user/{user_id}/cart/{sku_id}", h.AddItemToCart)
	http.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", h.DeleteItemFromCart)
	http.HandleFunc("DELETE /user/{user_id}/cart", h.ClearCart)
	http.HandleFunc("GET /user/{user_id}/cart", h.GetCart)
	http.HandleFunc("POST /user/{user_id}/checkout", h.CheckoutAll)

	loggedMux := middleware.LoggingMiddleware(http.DefaultServeMux)
	fmt.Printf("Сервер запущен на http://localhost:%s\n", cfg.Server.Port)
	http.ListenAndServe(":"+cfg.Server.Port, loggedMux)
}
