package main

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"homework-1/api/proto/loms"
	lomsApp "homework-1/loms/internal/app/loms"
	stockServ "homework-1/loms/internal/app/stock"
	mw "homework-1/loms/internal/logger"
	"homework-1/loms/internal/pkg/config"
	"homework-1/loms/internal/repository/order"
	"homework-1/loms/internal/repository/stock"
	"homework-1/loms/internal/serverGRPC"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	if err := godotenv.Load(".env.loms"); err != nil {
		log.Println("No .env.loms file found or failed to load")
	}
	cfg := config.NewConfig()
	fmt.Println("Master DSN:", cfg.Database.MasterDSN)
	dbConnMaster, err := pgx.Connect(context.Background(), cfg.Database.MasterDSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	dbConnReplica, err := pgx.Connect(context.Background(), cfg.Database.ReplicaDSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	stockRep := stock.New(dbConnMaster, dbConnReplica)
	stockService := stockServ.New(stockRep)
	orderRepo := order.New(dbConnMaster, dbConnReplica)
	lomsService := lomsApp.New(orderRepo, stockService)
	lomsServer := serverGRPC.New(lomsService)

	lis, err := net.Listen("tcp", cfg.Grpc.Port)
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(mw.Logger))
	pb.RegisterLomsServiceServer(grpcServer, lomsServer)

	go func() {
		fmt.Println("Server started at ", cfg.Grpc.Port)
		if err = grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	conn, err := grpc.NewClient(cfg.Grpc.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	gwmux := runtime.NewServeMux()
	if err = pb.RegisterLomsServiceHandler(context.Background(), gwmux, conn); err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	gwServer := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: gwmux,
	}

	go func() {
		log.Printf("Serving gRPC-Gateway on:%s\n", gwServer.Addr)
		if err = gwServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("gateway error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Stopping gRPC server...")
	grpcServer.GracefulStop()

	log.Println("Shutting down HTTP gateway...")
	if err = gwServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Closing DB connections...")
	defer dbConnMaster.Close(ctx)
	defer dbConnReplica.Close(ctx)
	log.Println("shutdown complete")
}
