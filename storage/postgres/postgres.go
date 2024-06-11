package postgres

import (
	"database/sql"
	"fmt"

	"github.com/Mubinabd/reservation_service/config"
	"github.com/Mubinabd/reservation_service/storage"
	_ "github.com/lib/pq"

)

type Storage struct {
	db           *sql.DB
	ReservationS storage.ReservationI
	RestaurantS storage.RestaurantI
	MenuS storage.MenuI
}

func ConnectDB() (*Storage, error) {
	cfg := config.Load()
	dbConn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDatabase)
	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	resS := NewReservationStorage(db)
	restS := NewRestaurantStorage(db)
	menuS := NewMenuRepo(db)
	return &Storage{
		db:    db,
		ReservationS: resS,
		RestaurantS: restS,
		MenuS: menuS,
	}, nil
}
func (s *Storage) Reservation() storage.ReservationI {
	if s.ReservationS == nil {
		s.ReservationS = NewReservationStorage(s.db)
	}
	return s.ReservationS
}

func (s *Storage) Restaurant() storage.RestaurantI {
	if s.RestaurantS == nil {
        s.RestaurantS = NewRestaurantStorage(s.db)
    }
    return s.RestaurantS
}

func (s *Storage) Menu() storage.MenuI {
	if s.MenuS == nil {
        s.MenuS = NewMenuRepo(s.db)
    }
    return s.MenuS
}