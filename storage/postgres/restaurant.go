package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	pb "github.com/Mubinabd/reservation_service/genproto"
)

type RestaurantStorage struct {
	db *sql.DB
}

func NewRestaurantStorage(db *sql.DB) *RestaurantStorage {
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

func (r *RestaurantStorage) UpdateRestaurant(req *pb.CreateRestaurantReq) (*pb.Void, error) {
	baseQuery := "UPDATE restaurants SET"
	var conditions []string
	var args []interface{}
	paramIndex := 1

	if req.Name != "" && req.Name !="string"{
		conditions = append(conditions, fmt.Sprintf("name = $%d", paramIndex))
		args = append(args, req.Name)
		paramIndex++
	}
	if req.Address != "" && req.Name !="string" {
		conditions = append(conditions, fmt.Sprintf("address = $%d", paramIndex))
		args = append(args, req.Address)
		paramIndex++
	}
	if req.PhoneNumber != "" && req.Name !="string" {
		conditions = append(conditions, fmt.Sprintf("phone_number = $%d", paramIndex))
		args = append(args, req.PhoneNumber)
		paramIndex++
	}
	if req.Description != "" && req.Name !="string" {
		conditions = append(conditions, fmt.Sprintf("description = $%d", paramIndex))
		args = append(args, req.Description)
		paramIndex++
	}

	if len(conditions) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query := baseQuery + " " + strings.Join(conditions, ", ") + fmt.Sprintf(" WHERE id = $%d", paramIndex)
	args = append(args, req.Id)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return &pb.Void{}, nil
}

func (r *RestaurantStorage) DeleteRestaurant(id *pb.ById) (*pb.Void, error) {
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
	`
	var conditions []string
	var args []interface{}
	if req.Address != "" {
		conditions = append(conditions, fmt.Sprintf("address ILIKE $%d", len(args)+1))
		args = append(args, req.Address)
	}
	if len(conditions) > 0 {
		query += " where " + strings.Join(conditions, " AND ")
	}

	rows, err := r.db.Query(query, args...)
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

func (r *RestaurantStorage) GetRestaurant(req *pb.ById) (*pb.Restaurant, error) {
	query := `
		SELECT name, address, phone_number, description
		FROM restaurants
		WHERE id = $1
	`
	// log.Println(req, req.Id)
	row := r.db.QueryRow(query, req.Id)
	restaurant := &pb.Restaurant{}
	if err := row.Scan(&restaurant.Name, &restaurant.Address, &restaurant.PhoneNumber, &restaurant.Description); err != nil {
		return nil, err
	}
	return restaurant, nil
}
