package service

import (
	"context"

	pb "github.com/Mubinabd/reservation_service/genproto"
	"github.com/Mubinabd/reservation_service/storage"
	"github.com/google/uuid"
)

type RestaurantService struct {
	storage storage.StorageI
	pb.UnimplementedRestaurantServiceServer
}

func NewRestaurantService(storage storage.StorageI) *RestaurantService {
	return &RestaurantService{
		storage: storage,
	}
}

func (s *RestaurantService) CreateRestaurant(c context.Context, req *pb.CreateRestaurantReq) (*pb.Void, error) {
	id := uuid.NewString()
	req.Id = id
	return s.storage.Restaurant().CreateRestaurant(req)
}
func (s *RestaurantService) UpdateRestaurant(c context.Context, req *pb.CreateRestaurantReq) (*pb.Void, error) {
	return s.storage.Restaurant().UpdateRestaurant(req)
}
func (s *RestaurantService) DeleteRestaurant(c context.Context,id *pb.ById) (*pb.Void, error) {
	return s.storage.Restaurant().DeleteRestaurant(id)
}

func (s *RestaurantService) GetRestaurant(c context.Context,id *pb.ById) (*pb.Restaurant, error) {
	return s.storage.Restaurant().GetRestaurant(id)
}

func (s *RestaurantService) GetAllRestaurants(c context.Context, flt *pb.AddressFilter) (*pb.Restaurants, error) {
	return s.storage.Restaurant().GetAllRestaurants(flt)
}
