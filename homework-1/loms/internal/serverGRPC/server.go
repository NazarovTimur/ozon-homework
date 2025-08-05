package serverGRPC

import (
	"context"
	"homework-1/api/proto/loms"
	"homework-1/loms/internal/app/loms"
)

type LomsServerInterface interface {
	OrderCreate(ctx context.Context, request *pb.OrderCreateRequest) (*pb.OrderCreateResponse, error)
	OrderInfo(ctx context.Context, request *pb.OrderInfoRequest) (*pb.OrderInfoResponse, error)
	OrderPay(ctx context.Context, request *pb.OrderPayRequest) (*pb.OrderPayResponse, error)
	OrderCancel(ctx context.Context, request *pb.OrderCancelRequest) (*pb.OrderCancelResponse, error)
	StocksInfo(ctx context.Context, request *pb.StocksInfoRequest) (*pb.StocksInfoResponse, error)
}

type Server struct {
	pb.UnimplementedLomsServiceServer
	lomsService loms.LomsInterface
}

func New(lomsService loms.LomsInterface) *Server {
	return &Server{
		lomsService: lomsService,
	}
}

func (s *Server) OrderCreate(ctx context.Context, request *pb.OrderCreateRequest) (*pb.OrderCreateResponse, error) {
	return s.lomsService.OrderCreate(ctx, request)
}

func (s *Server) OrderInfo(ctx context.Context, request *pb.OrderInfoRequest) (*pb.OrderInfoResponse, error) {
	return s.lomsService.OrderInfo(ctx, request)
}

func (s *Server) OrderPay(ctx context.Context, request *pb.OrderPayRequest) (*pb.OrderPayResponse, error) {
	return s.lomsService.OrderPay(ctx, request)
}

func (s *Server) OrderCancel(ctx context.Context, request *pb.OrderCancelRequest) (*pb.OrderCancelResponse, error) {
	return s.lomsService.OrderCancel(ctx, request)
}

func (s *Server) StocksInfo(ctx context.Context, request *pb.StocksInfoRequest) (*pb.StocksInfoResponse, error) {
	return s.lomsService.StocksInfo(ctx, request)
}
