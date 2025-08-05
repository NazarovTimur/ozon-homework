package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"homework-1/api/proto/loms"
	"homework-1/cart/internal/app/cart"
	"homework-1/cart/internal/app/product"
	"homework-1/cart/internal/http/handler"
	"homework-1/cart/internal/http/middleware"
	"homework-1/cart/internal/pkg/config"
	"homework-1/cart/internal/pkg/retry"
	repoCart "homework-1/cart/internal/repository"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := godotenv.Load(".env.cart"); err != nil {
		log.Println("No .env.cart. file found or failed to load")
	}
	cfg := config.NewConfig()

	conn, err := grpc.NewClient(cfg.Grpc.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

	mux := http.NewServeMux()
	mux.HandleFunc("POST /user/{user_id}/cart/{sku_id}", h.AddItemToCart)
	mux.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", h.DeleteItemFromCart)
	mux.HandleFunc("DELETE /user/{user_id}/cart", h.ClearCart)
	mux.HandleFunc("GET /user/{user_id}/cart", h.GetCart)
	mux.HandleFunc("POST /user/{user_id}/checkout", h.CheckoutAll)

	loggedMux := middleware.LoggingMiddleware(mux)

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: loggedMux,
	}

	go func() {
		fmt.Printf("Сервер запущен на http://localhost:%s\n", cfg.Server.Port)
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
