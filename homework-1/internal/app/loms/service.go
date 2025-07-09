package loms

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"homework-1/internal/pkg/constants"
	"homework-1/internal/pkg/model"
	pb "homework-1/internal/proto/loms"
	"homework-1/internal/repository/order"
	"homework-1/internal/repository/stock"
)

type LomsInterface interface {
	OrderCreate(ctx context.Context, request *pb.OrderCreateRequest) (*pb.OrderCreateResponse, error)
	OrderInfo(ctx context.Context, request *pb.OrderInfoRequest) (*pb.OrderInfoResponse, error)
	OrderPay(ctx context.Context, request *pb.OrderPayRequest) (*pb.OrderPayResponse, error)
	OrderCancel(ctx context.Context, request *pb.OrderCancelRequest) (*pb.OrderCancelResponse, error)
	StocksInfo(ctx context.Context, request *pb.StocksInfoRequest) (*pb.StocksInfoResponse, error)
}

type Service struct {
	pb.UnimplementedLomsServiceServer
	stockRepo stock.StockInterface
	orderRepo order.OrderInterface
}

func New(orderRepo order.OrderInterface, stockRepo stock.StockInterface) *Service {
	return &Service{
		orderRepo: orderRepo,
		stockRepo: stockRepo,
	}
}

func (s *Service) OrderCreate(ctx context.Context, request *pb.OrderCreateRequest) (*pb.OrderCreateResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}
	items, err := convertItems(request.Items)
	if err != nil {
		return nil, err
	}

	orderID, err := s.orderRepo.Create(ctx, request.UserID, items)
	if err != nil {
		return &pb.OrderCreateResponse{}, err
	}
	var reserved []model.OrderItem
	for _, item := range items {
		if ok := s.stockRepo.Reserve(ctx, item.SKU, item.Count); !ok {
			for _, r := range reserved {
				_ = s.stockRepo.ReserveCancel(ctx, r.SKU, r.Count)
			}
			_ = s.orderRepo.SetStatus(ctx, orderID, constants.StatusFailed)
		}
		reserved = append(reserved, item)
	}
	if err = s.orderRepo.SetStatus(ctx, orderID, constants.StatusAwaitingPayment); err != nil {
		return nil, err
	}
	return &pb.OrderCreateResponse{OrderID: orderID}, nil
}

func (s *Service) OrderInfo(ctx context.Context, request *pb.OrderInfoRequest) (*pb.OrderInfoResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	stat, user, itemsOrder, err := s.orderRepo.GetByID(ctx, request.OrderID)
	if err != nil {
		return &pb.OrderInfoResponse{}, err
	}
	items, err := backConvert(itemsOrder)
	if err != nil {
		return nil, err
	}
	return &pb.OrderInfoResponse{Status: stat, UserID: user, Items: items}, nil
}

func (s *Service) OrderPay(ctx context.Context, request *pb.OrderPayRequest) (*pb.OrderPayResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}
	_, _, items, err := s.orderRepo.GetByID(ctx, request.OrderID)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if err = s.stockRepo.ReserveRemove(ctx, item.SKU, item.Count); err != nil {
			return nil, err
		}
	}
	if err = s.orderRepo.SetStatus(ctx, request.OrderID, constants.StatusPayed); err != nil {
		return nil, err
	}
	return &pb.OrderPayResponse{}, nil
}

func (s *Service) OrderCancel(ctx context.Context, request *pb.OrderCancelRequest) (*pb.OrderCancelResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}
	_, _, items, err := s.orderRepo.GetByID(ctx, request.OrderID)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if err = s.stockRepo.ReserveCancel(ctx, item.SKU, item.Count); err != nil {
			return nil, err
		}
	}
	if err = s.orderRepo.SetStatus(ctx, request.OrderID, constants.StatusCancelled); err != nil {
		return nil, err
	}
	return &pb.OrderCancelResponse{}, nil
}

func (s *Service) StocksInfo(ctx context.Context, request *pb.StocksInfoRequest) (*pb.StocksInfoResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}
	count, err := s.stockRepo.GetBySKU(ctx, request.Sku)
	if err != nil {
		return nil, err
	}
	return &pb.StocksInfoResponse{Count: uint64(count)}, nil
}

func convertItems(protoItems []*pb.OrderItem) ([]model.OrderItem, error) {
	if len(protoItems) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty order items")
	}

	items := make([]model.OrderItem, len(protoItems))
	for i, item := range protoItems {
		items[i] = model.OrderItem{
			SKU:   item.Sku,
			Count: uint16(item.Count),
		}
	}
	return items, nil
}

func backConvert(orderItems []model.OrderItem) ([]*pb.OrderItem, error) {
	if len(orderItems) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty order items")
	}

	items := make([]*pb.OrderItem, len(orderItems))
	for i, item := range orderItems {
		items[i] = &pb.OrderItem{
			Sku:   item.SKU,
			Count: uint32(item.Count),
		}
	}
	return items, nil
}
