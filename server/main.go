package main

import (
	"log"
	"net"

	pb "github.com/Mubinabd/reservation_service/genproto"
	"github.com/Mubinabd/reservation_service/service"
	"github.com/Mubinabd/reservation_service/storage/postgres"
	"google.golang.org/grpc"
)

func main() {
	db, err := postgres.ConnectDB()
	if err != nil {
		panic(err)
	}
	liss, err := net.Listen("tcp", ":8088")
	if err != nil {
		panic(err)
	}
	
	s := grpc.NewServer()
	pb.RegisterReservationServiceServer(s, service.NewReservationService(db))
	pb.RegisterRestaurantServiceServer(s, service.NewRestaurantService(db))
	pb.RegisterMenuServiceServer(s, service.NewMenuService(db))
	pb.RegisterOrderServiceServer(s, service.NewOrderService(db))

	log.Printf("server listening at %v", liss.Addr())
	if err := s.Serve(liss); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
