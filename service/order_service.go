package service

import (
	"context"

	pb "github.com/Mubinabd/reservation_service/genproto"
	"github.com/Mubinabd/reservation_service/storage"
	"github.com/google/uuid"
)

type OrderService struct {
	storage storage.StorageI
	pb.UnimplementedOrderServiceServer
}

func NewOrderService(storage storage.StorageI) *OrderService {
	return &OrderService{
		storage: storage,
	}
}

func (s *OrderService)CreateOrder(c context.Context, req *pb.Order) (*pb.Void, error) {
	id := uuid.NewString()
	req.Id = id
	return s.storage.Orders().CreateOrder(req)
}
func (s *OrderService)UpdateOrder(c context.Context, req *pb.Order) (*pb.Void, error) {
	return s.storage.Orders().UpdateOrder(req)
}
func (s *OrderService)DeleteOrder(c context.Context,id *pb.ById) (*pb.Void, error) {
	return s.storage.Orders().DeleteOrder(id)
}

func (s *OrderService)GetOrder(c context.Context,id *pb.ById) (*pb.Order, error) {
	return s.storage.Orders().GetOrder(id)
}

func (s *OrderService)GetAllOrders(c context.Context, flt *pb.Void) (*pb.Orders, error) {
	return s.storage.Orders().GetAllOrders(flt)
}