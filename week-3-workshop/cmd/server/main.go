package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"google.golang.org/grpc/credentials/insecure"

	"gitlab.ozon.dev/go/classroom-14/students/week-3-workshop/internal/mw"

	"google.golang.org/grpc/reflection"

	"gitlab.ozon.dev/go/classroom-14/students/week-3-workshop/internal/app/notes"

	desc "gitlab.ozon.dev/go/classroom-14/students/week-3-workshop/pkg/api/notes/v1"
	"google.golang.org/grpc"
)

const (
	grpcPort = 50051
	httpPort = 8081
)

func headerMatcher(key string) (string, bool) {
	switch strings.ToLower(key) {
	case "x-auth":
		return key, true
	default:
		return key, false
	}
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			mw.Panic,
			mw.Logger,
			mw.Panic,
			mw.Auth,
			mw.Validate,
		),
	)

	reflection.Register(grpcServer)

	controller := notes.NewService()

	desc.RegisterNotesServer(grpcServer, controller)

	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	conn, err := grpc.NewClient(fmt.Sprintf(":%d", grpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	gwmux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(headerMatcher))
	if err = desc.RegisterNotesHandler(context.Background(), gwmux, conn); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: gwmux,
	}
	log.Printf("Serving gRPC-Gateway on %s\n", gwServer.Addr)
	log.Fatalln(gwServer.ListenAndServe())
}
