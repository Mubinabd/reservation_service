package postgres

import (
	"database/sql"
	"log"

	pb "github.com/Mubinabd/reservation_service/genproto"
)

type RestaurantStorage struct {
	db *sql.DB
}

func NewRestaurantStorage (db *sql.DB) *RestaurantStorage {
	return &RestaurantStorage{db: db}
}

func (r *RestaurantStorage) CreateRestaurant(req *pb.CreateRestaurantReq) (*pb.Void, error) {
	query := `
		INSERT INTO restaurants
		(id, name, address, phone_number, description)
		VALUES
		($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(query, req.Id, req.Name, req.Address, req.PhoneNumber, req.Description)
	if err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}

func (r *RestaurantStorage)UpdateRestaurant(req *pb.CreateRestaurantReq) (*pb.Void, error) {
	query := `
		UPDATE restaurants
		SET name = $2, address = $3, phone_number = $4, description = $5
		WHERE id = $1	
	`
	_, err := r.db.Exec(query, req.Id, req.Name, req.Address, req.PhoneNumber, req.Description)
	if err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}

func (r *RestaurantStorage)DeleteRestaurant(id *pb.ById) (*pb.Void, error) {
	query := `
		DELETE FROM restaurants
		WHERE id = $1
	`
	_, err := r.db.Exec(query, id.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}

func (r *RestaurantStorage) GetAllRestaurants(req *pb.AddressFilter) (*pb.Restaurants, error) {
	query := `
		SELECT name, address, phone_number, description
		FROM restaurants
		WHERE address = $1
	`
	rows, err := r.db.Query(query, req.Address)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var restaurants []*pb.Restaurant
	for rows.Next() {
		restaurant := &pb.Restaurant{}
		if err := rows.Scan(&restaurant.Name, &restaurant.Address, &restaurant.PhoneNumber, &restaurant.Description); err != nil {
			return nil, err
		}
		restaurants = append(restaurants, restaurant)
	}
	return &pb.Restaurants{Restaurants: restaurants}, nil
}


func (r *RestaurantStorage)GetRestaurant(req *pb.ById) (*pb.Restaurant, error) {
	query := `
		SELECT id, name, phone_number, description
		FROM restaurants
		WHERE id = $1
	`
	log.Println(req,req.Id)
	row := r.db.QueryRow(query, req.Id)
	restaurant := &pb.Restaurant{}
	if err := row.Scan(&restaurant.Name, &restaurant.Address, &restaurant.PhoneNumber, &restaurant.Description); err != nil {
		return nil, err
	}
	return restaurant, nil
}