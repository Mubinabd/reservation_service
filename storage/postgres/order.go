package postgres

import (
	"database/sql"

	pb "github.com/Mubinabd/reservation_service/genproto"
	"github.com/google/uuid"
)

type OrderStorage struct {
	db *sql.DB
}

func NewOrderStorage(db *sql.DB) *OrderStorage {
	return &OrderStorage{db: db}
}

func (s *OrderStorage)CreateOrder(req *pb.Order)(*pb.Void,error){
	id := uuid.NewString()
	query := `INSERT INTO orders (id, reservation_id, menu_item_id, quantity) VALUES ($1, $2, $3, $4)`

	_, err := s.db.Exec(query, id, req.ReservationId, req.MenuItemId, req.Quantity)
	if err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}

func(s *OrderStorage)GetOrder(id *pb.ById)(*pb.Order,error){
	query := `SELECT id, reservation_id, menu_item_id, quantity FROM orders WHERE id = $1`	
	row := s.db.QueryRow(query, id.Id)
	order := &pb.Order{}
	err := row.Scan(&order.Id, &order.ReservationId, &order.MenuItemId, &order.Quantity)
	if err != nil {
		return nil, err
	}
	return order, nil

}

func (s *OrderStorage)UpdateOrder(req *pb.Order)(*pb.Void,error){
	query := `UPDATE orders SET reservation_id = $1, menu_item_id = $2, quantity = $3 WHERE id = $4`
	_, err := s.db.Exec(query, req.ReservationId, req.MenuItemId, req.Quantity, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}

func(s *OrderStorage)DeleteOrder(id *pb.ById)(*pb.Void,error){
	query := `DELETE FROM orders WHERE id = $1`
	_, err := s.db.Exec(query, id.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}

func(s *OrderStorage)GetAllOrders(*pb.Void)(*pb.Orders,error){
	query := `SELECT id, reservation_id, menu_item_id, quantity FROM orders`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	orders := &pb.Orders{}
	for rows.Next() {
		order := &pb.Order{}
		err := rows.Scan(&order.Id, &order.ReservationId, &order.MenuItemId, &order.Quantity)
		if err != nil {
			return nil, err
		}
		orders.Orders = append(orders.Orders, order)
	}
	return orders, nil
}