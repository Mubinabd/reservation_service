package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	pb "github.com/Mubinabd/reservation_service/genproto"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	menuRepo := NewMenuRepo(db)

	menu := &pb.Menu{
		RestaurantId: "restaurant_123",
		Name:         "Test Menu",
		Description:  "Test Description",
		Price:        10.0,
	}

	mock.ExpectExec(`INSERT INTO menu`).
		WithArgs(sqlmock.AnyArg(), menu.RestaurantId, menu.Name, menu.Description, menu.Price).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = menuRepo.Create(menu)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	menuRepo := NewMenuRepo(db)

	menu := &pb.Menu{
		Id:           "menu_123",
		RestaurantId: "restaurant_123",
		Name:         "Updated Menu",
		Description:  "Updated Description",
		Price:        20.0,
	}

	mock.ExpectExec(`update menu set`).
		WithArgs(menu.RestaurantId, menu.Name, menu.Description, menu.Price, menu.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = menuRepo.Update(menu)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	menuRepo := NewMenuRepo(db)

	id := &pb.ById{Id: "menu_123"}

	mock.ExpectExec(`update menu set deleted_at`).
		WithArgs(id.GetId()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = menuRepo.Delete(id)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	menuRepo := NewMenuRepo(db)

	id := &pb.ById{Id: "menu_123"}

	rows := sqlmock.NewRows([]string{"id", "restaurants_id", "name", "description", "price"}).
		AddRow("menu_123", "restaurant_123", "Test Menu", "Test Description", 10.0)

	mock.ExpectQuery(`select id, restaurants_id, name, description, price from menu where id = \$1`).
		WithArgs(id.GetId()).
		WillReturnRows(rows)

	menu, err := menuRepo.GetById(id)
	assert.NoError(t, err)
	assert.NotNil(t, menu)
	assert.Equal(t, "menu_123", menu.Id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	menuRepo := NewMenuRepo(db)

	flt := &pb.MenuFilter{
		RestaurantId: "restaurant_123",
		PriceFrom:    5,
		PriceTo:      15,
	}

	rows := sqlmock.NewRows([]string{"id", "restaurant_id", "name", "description", "price"}).
		AddRow("menu_123", "restaurant_123", "Test Menu 1", "Test Description 1", 10.0).
		AddRow("menu_124", "restaurant_123", "Test Menu 2", "Test Description 2", 12.0)

	mock.ExpectQuery(`SELECT id, restaurant_id, name, description, price FROM menu WHERE deleted_at = 0 AND restaurant_id = \$1 AND price > \$2 AND price < \$3`).
		WithArgs(flt.RestaurantId, flt.PriceFrom, flt.PriceTo).
		WillReturnRows(rows)

	menus, err := menuRepo.GetAll(flt)
	assert.NoError(t, err)
	assert.NotNil(t, menus)
	assert.Equal(t, 2, len(menus.Menus))
	assert.NoError(t, mock.ExpectationsWereMet())
}