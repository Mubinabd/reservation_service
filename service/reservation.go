package service

import (
	"context"

	pb "github.com/Mubinabd/reservation_service/genproto"
	"github.com/Mubinabd/reservation_service/storage"
	"github.com/google/uuid"
)

type ReservationService struct {
	stg storage.StorageI
	pb.UnimplementedReservationServiceServer
}

func NewReservationService(stg storage.StorageI) *ReservationService {
	return &ReservationService{
		stg: stg,
	}
}


func (s *ReservationService) CreateReservation(ctx context.Context, req *pb.ReservationCreate) (*pb.Void, error) {
	id := uuid.NewString()
	req.Id = id
	return s.stg.Reservation().CreateReservation(req)
}

func (s *ReservationService) GetReservation(ctx context.Context, req *pb.ById) (*pb.Reservation, error) {
	return s.stg.Reservation().GetReservation(req)
}

func (s *ReservationService) UpdateReservation(ctx context.Context, req *pb.ReservationCreate) (*pb.Void, error) {
	return s.stg.Reservation().UpdateReservation(req)
}

func (s *ReservationService) DeleteReservation(ctx context.Context, req *pb.ById) (*pb.Void, error) {
	return s.stg.Reservation().DeleteReservation(req)
}

func (s *ReservationService) GetAllReservation(ctx context.Context, req *pb.FilterByTime) (*pb.Reservations, error) {
	return s.stg.Reservation().GetReservationByFilter(req)
}

func (s *ReservationService) GetTotalSum(ctx context.Context, req *pb.ById) (*pb.Total, error) {
    return s.stg.Reservation().GetTotalSum(req)
}

func (s *ReservationService) CheckReservation(ctx context.Context, req *pb.ResrvationTime) (*pb.Void, error) {
    return s.stg.Reservation().CheckReservation(req)
}