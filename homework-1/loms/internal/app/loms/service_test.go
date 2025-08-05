package loms

import (
	"context"
	"errors"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/suite"
	"homework-1/api/proto/loms"
	"homework-1/loms/internal/pkg/model"
	mockOrd "homework-1/loms/internal/repository/order/mock"
	mockSt "homework-1/loms/internal/repository/stock/mock"
	"testing"
)

type LomsServiceTestSuite struct {
	suite.Suite
	ctrl        *minimock.Controller
	mockOrder   *mockOrd.OrderInterfaceMock
	mockStock   *mockSt.StockInterfaceMock
	serviceLoms *Service
}

func TestLOMSServiceSuite(t *testing.T) {
	suite.Run(t, new(LomsServiceTestSuite))
}

func (s *LomsServiceTestSuite) SetupTest() {
	s.ctrl = minimock.NewController(s.T())
	s.mockOrder = mockOrd.NewOrderInterfaceMock(s.ctrl)
	s.mockStock = mockSt.NewStockInterfaceMock(s.ctrl)
	stockAdapter := &stockAdapterForTest{mock: s.mockStock}
	s.serviceLoms = New(s.mockOrder, stockAdapter)
}

func (s *LomsServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

const (
	testUserID  int64  = 14
	testSKU     uint32 = 52
	testCount   uint32 = 44
	testOrderID int64  = 1
	testStatus         = "awaiting payment"
)

func (s *LomsServiceTestSuite) TestOrderCreate() {
	ctx := context.Background()
	items := []*pb.OrderItem{{Sku: testSKU, Count: testCount}}
	checkItems := []model.OrderItem{{SKU: testSKU, Count: uint16(testCount)}}
	requestLoms := &pb.OrderCreateRequest{UserID: testUserID, Items: items}
	responseLoms := &pb.OrderCreateResponse{OrderID: 1}

	s.mockOrder.CreateMock.Expect(ctx, testUserID, checkItems).Return(1, nil)
	s.mockOrder.SetStatusMock.Expect(ctx, testOrderID, testStatus).Return(nil)
	s.mockStock.ReserveMock.Expect(ctx, testSKU, uint16(testCount)).Return(true, nil)

	resp, err := s.serviceLoms.OrderCreate(ctx, requestLoms)
	s.NoError(err)
	s.Equal(responseLoms, resp)
}

func (s *LomsServiceTestSuite) TestOrderCreate_Failed() {
	ctx := context.Background()
	requestLoms := &pb.OrderCreateRequest{UserID: 0, Items: nil}

	resp, err := s.serviceLoms.OrderCreate(ctx, requestLoms)

	s.Nil(resp)
	s.Error(err)
}

func (s *LomsServiceTestSuite) TestOrderInfo() {
	ctx := context.Background()
	items := []*pb.OrderItem{{Sku: testSKU, Count: testCount}}
	checkItems := []model.OrderItem{{SKU: testSKU, Count: uint16(testCount)}}
	requestLoms := &pb.OrderInfoRequest{OrderID: 1}
	responseLoms := &pb.OrderInfoResponse{
		Status: testStatus,
		UserID: testUserID,
		Items:  items,
	}

	s.mockOrder.GetByIDMock.Expect(ctx, testOrderID).Return(testStatus, testUserID, checkItems, nil)

	resp, err := s.serviceLoms.OrderInfo(ctx, requestLoms)
	s.NoError(err)
	s.Equal(resp, responseLoms)
}

func (s *LomsServiceTestSuite) TestOrderInfo_Failed() {
	ctx := context.Background()
	requestLoms := &pb.OrderInfoRequest{OrderID: testOrderID}
	expectedErr := errors.New("repository failure")

	s.mockOrder.GetByIDMock.Expect(ctx, testOrderID).Return("", 0, nil, expectedErr)

	resp, err := s.serviceLoms.OrderInfo(ctx, requestLoms)
	s.Equal(&pb.OrderInfoResponse{}, resp)
	s.EqualError(err, expectedErr.Error())
}

func (s *LomsServiceTestSuite) TestOrderPay() {
	ctx := context.Background()
	checkItems := []model.OrderItem{{SKU: testSKU, Count: uint16(testCount)}}
	requestLoms := &pb.OrderPayRequest{OrderID: 1}
	responseLoms := &pb.OrderPayResponse{}

	s.mockOrder.GetByIDMock.Expect(ctx, testOrderID).Return(testStatus, testUserID, checkItems, nil)
	s.mockOrder.SetStatusMock.Expect(ctx, testOrderID, "payed").Return(nil)
	s.mockStock.ReserveRemoveMock.Expect(ctx, testSKU, uint16(testCount)).Return(nil)

	resp, err := s.serviceLoms.OrderPay(ctx, requestLoms)
	s.NoError(err)
	s.Equal(resp, responseLoms)
}

func (s *LomsServiceTestSuite) TestOrderCancel() {
	ctx := context.Background()
	checkItems := []model.OrderItem{{SKU: testSKU, Count: uint16(testCount)}}
	requestLoms := &pb.OrderCancelRequest{OrderID: 1}
	responseLoms := &pb.OrderCancelResponse{}

	s.mockOrder.GetByIDMock.Expect(ctx, testOrderID).Return(testStatus, testUserID, checkItems, nil)
	s.mockOrder.SetStatusMock.Expect(ctx, testOrderID, "cancelled").Return(nil)
	s.mockStock.ReserveCancelMock.Expect(ctx, testSKU, uint16(testCount)).Return(nil)

	resp, err := s.serviceLoms.OrderCancel(ctx, requestLoms)
	s.NoError(err)
	s.Equal(resp, responseLoms)
}

func (s *LomsServiceTestSuite) TestStocksInfo() {
	ctx := context.Background()
	requestLoms := &pb.StocksInfoRequest{Sku: testSKU}
	responseLoms := &pb.StocksInfoResponse{Count: uint64(testCount)}

	s.mockStock.GetBySKUMock.Expect(ctx, testSKU).Return(uint16(testCount), nil)

	resp, err := s.serviceLoms.StocksInfo(ctx, requestLoms)
	s.NoError(err)
	s.Equal(resp, responseLoms)
}

// stockAdapterForTest адаптирует mockStock к интерфейсу, который ожидает слайс OrderItem
type stockAdapterForTest struct {
	mock *mockSt.StockInterfaceMock
}

func (a *stockAdapterForTest) Reserve(ctx context.Context, items []model.OrderItem) (bool, error) {
	for _, item := range items {
		ok, err := a.mock.Reserve(ctx, item.SKU, item.Count)
		if err != nil || !ok {
			return false, err
		}
	}
	return true, nil
}

func (a *stockAdapterForTest) ReserveRemove(ctx context.Context, sku uint32, count uint16) error {
	return a.mock.ReserveRemove(ctx, sku, count)
}

func (a *stockAdapterForTest) ReserveCancel(ctx context.Context, sku uint32, count uint16) error {
	return a.mock.ReserveCancel(ctx, sku, count)
}

func (a *stockAdapterForTest) GetBySKU(ctx context.Context, sku uint32) (uint16, error) {
	return a.mock.GetBySKU(ctx, sku)
}
