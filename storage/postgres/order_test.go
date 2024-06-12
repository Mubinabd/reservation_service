package postgres_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Mubinabd/reservation_service/genproto"
	"github.com/Mubinabd/reservation_service/storage/postgres"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := postgres.NewOrderStorage(db)

	order := &genproto.Order{
		ReservationId: "res123",
		MenuItemId:    "menu123",
		Quantity:      2,
	}

	mock.ExpectExec(`INSERT INTO orders`).
		WithArgs(sqlmock.AnyArg(), order.ReservationId, order.MenuItemId, order.Quantity).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = storage.CreateOrder(order)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := postgres.NewOrderStorage(db)

	orderID := uuid.NewString()
	rows := sqlmock.NewRows([]string{"id", "reservation_id", "menu_item_id", "quantity"}).
		AddRow(orderID, "res123", "menu123", 2)

	mock.ExpectQuery(`SELECT id, reservation_id, menu_item_id, quantity FROM orders WHERE id = \$1`).
		WithArgs(orderID).
		WillReturnRows(rows)

	req := &genproto.ById{Id: orderID}
	order, err := storage.GetOrder(req)
	assert.NoError(t, err)
	assert.Equal(t, orderID, order.Id)
	assert.Equal(t, "res123", order.ReservationId)
	assert.Equal(t, "menu123", order.MenuItemId)
	assert.Equal(t, int32(2), order.Quantity)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := postgres.NewOrderStorage(db)

	order := &genproto.Order{
		Id:            "order123",
		ReservationId: "res123",
		MenuItemId:    "menu123",
		Quantity:      3,
	}

	mock.ExpectExec(`UPDATE orders SET reservation_id = \$1, menu_item_id = \$2, quantity = \$3 WHERE id = \$4`).
		WithArgs(order.ReservationId, order.MenuItemId, order.Quantity, order.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = storage.UpdateOrder(order)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := postgres.NewOrderStorage(db)

	orderID := "order123"
	mock.ExpectExec(`DELETE FROM orders WHERE id = \$1`).
		WithArgs(orderID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	req := &genproto.ById{Id: orderID}
	_, err = storage.DeleteOrder(req)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllOrders(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := postgres.NewOrderStorage(db)

	rows := sqlmock.NewRows([]string{"id", "reservation_id", "menu_item_id", "quantity"}).
		AddRow("order1", "res123", "menu123", 2).
		AddRow("order2", "res124", "menu124", 3)

	mock.ExpectQuery(`SELECT id, reservation_id, menu_item_id, quantity FROM orders`).
		WillReturnRows(rows)

	resp, err := storage.GetAllOrders(&genproto.Void{})
	assert.NoError(t, err)
	assert.Len(t, resp.Orders, 2)
	assert.Equal(t, "order1", resp.Orders[0].Id)
	assert.Equal(t, "order2", resp.Orders[1].Id)
	assert.NoError(t, mock.ExpectationsWereMet())
}