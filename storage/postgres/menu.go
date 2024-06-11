package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	pb "github.com/Mubinabd/reservation_service/genproto"
	"github.com/google/uuid"
)

type MenuRepo struct {
	db *sql.DB
}

func NewMenuRepo(db *sql.DB) *MenuRepo {
	return &MenuRepo{db: db}
}

func (m *MenuRepo) Create(menu *pb.Menu) (*pb.Void, error) {
	menu.Id = uuid.NewString()
	query :=
		`
	INSERT INTO menu(
		id,
		restaurants_id,
		name,
		description,
		price) 
		VALUES($1,$2,$3,$4,$5)
	`

	_, err := m.db.Exec(query, menu.Id, menu.RestaurantId, menu.Name, menu.Description, menu.Price)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	return nil, nil
}

func (m *MenuRepo) Update(menu *pb.Menu) (*pb.Void, error) {
	query :=
		`
	update menu set 
		restaurants_id = $1,
		name = $2,
		description = $3,
		price = $4
	where 
		id = $5
	`
	_, err := m.db.Exec(query, menu.GetRestaurantId(), menu.GetName(), menu.GetDescription(), menu.GetPrice(), menu.GetId())
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return nil, nil
}

func (m *MenuRepo) Delete(id *pb.ById) (*pb.Void, error) {
	query :=
		`
		update menu set deleted_at = EXTRACT(EPOCH FROM NOW()) where id = $1
		`

	_, err := m.db.Exec(query, id.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}

func (m *MenuRepo) GetById(id *pb.ById) (*pb.Menu, error) {
	query :=
		`
		select 
			id,
			restaurants_id,
			name,
			description,
			price 
			from menu 
			where id = $1
	`
	menu := pb.Menu{}
	err := m.db.QueryRow(query, id.GetId()).Scan(&menu.Id, &menu.RestaurantId, &menu.Name, &menu.Description, &menu.Price)

	if err != nil {
		return nil, err
	}

	return &menu, nil
}

func (m *MenuRepo) GetAll(flt *pb.MenuFilter) (*pb.Menus, error) {

	baseQuery := `
		SELECT id, restaurant_id, name, description, price 
		FROM menu 
		WHERE deleted_at = 0`

	var conditions []string
	var args []interface{}
	paramIndex := 1

	if flt.RestaurantId != "" {
		conditions = append(conditions, fmt.Sprintf("restaurant_id = $%d", paramIndex))
		args = append(args, flt.RestaurantId)
		paramIndex++
	}
	if flt.PriceFrom > 0 {
		conditions = append(conditions, fmt.Sprintf("price > $%d", paramIndex))
		args = append(args, flt.PriceFrom)
		paramIndex++
	}
	if flt.PriceTo > 0 {
		conditions = append(conditions, fmt.Sprintf("price < $%d", paramIndex))
		args = append(args, flt.PriceTo)
		paramIndex++
	}

	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	rows, err := m.db.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	menus := pb.Menus{}
	for rows.Next() {
		menu := pb.Menu{}
		err := rows.Scan(&menu.Id, &menu.RestaurantId, &menu.Name, &menu.Description, &menu.Price)
		if err != nil {
			return nil, err
		}
		menus.Menus = append(menus.Menus, &menu)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &menus, nil
}
