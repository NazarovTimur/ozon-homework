package main

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"homework-1/internal/app/loms"
	mw "homework-1/internal/http/loms"
	"homework-1/internal/pkg/constants"
	pb "homework-1/internal/proto/loms"
	"homework-1/internal/repository/order"
	"homework-1/internal/repository/stock"
	"log"
	"net"
	"net/http"
)

func main() {
	stockRepo, err := stock.NewStockFromJSON()
	if err != nil {
		log.Fatalf("init stock repository: %v", err)
	}
	orderRepo := order.New()
	lomsService := loms.New(orderRepo, stockRepo)

	lis, err := net.Listen("tcp", constants.GrpcPort)
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(mw.Logger))
	pb.RegisterLomsServiceServer(grpcServer, lomsService)
	fmt.Println("Server started at ", constants.GrpcPort)

	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	conn, err := grpc.NewClient(constants.GrpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	gwmux := runtime.NewServeMux()

	if err = pb.RegisterLomsServiceHandler(context.Background(), gwmux, conn); err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}
	gwServer := &http.Server{
		Addr:    constants.HttpPort,
		Handler: gwmux,
	}
	log.Printf("Serving gRPC-Gateway on:%s\n", gwServer.Addr)
	log.Fatalln(gwServer.ListenAndServe())
}
