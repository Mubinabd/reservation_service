package storage

import pb "github.com/Mubinabd/reservation_service/genproto"


type StorageI interface {
	Reservation() ReservationI
	Restaurant() RestaurantI
	Menu() MenuI
}

type ReservationI interface {
	CreateReservation(res *pb.ReservationCreate) (*pb.Void, error)
	UpdateReservation(res *pb.ReservationCreate) (*pb.Void, error)
	GetReservation(id *pb.ById) (*pb.Reservation, error)
	DeleteReservation(id *pb.ById) (*pb.Void, error)
	GetReservationByFilter(filter *pb.FilterByTime) (*pb.Reservations, error)
}
type RestaurantI interface {
	CreateRestaurant(req *pb.CreateRestaurantReq) (*pb.Void, error)
	UpdateRestaurant(req *pb.CreateRestaurantReq) (*pb.Void, error)
	DeleteRestaurant(id *pb.ById) (*pb.Void, error)
	GetRestaurant(req *pb.ById) (*pb.Restaurant, error)
	GetAllRestaurants(req *pb.AddressFilter) (*pb.Restaurants, error)
}
type MenuI interface{
	Create(menu *pb.Menu)(*pb.Void,error)
	Update(menu *pb.Menu)(*pb.Void,error)
	Delete(menu *pb.ById)(*pb.Void,error)
	GetById(menu *pb.ById)(*pb.Menu,error)
	GetAll(menu *pb.MenuFilter)(*pb.Menus,error)
}