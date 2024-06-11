package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	pb "github.com/Mubinabd/reservation_service/genproto"
	"github.com/stretchr/testify/assert"
)

func TestCreateRestaurant(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewRestaurantStorage(db)

	req := &pb.CreateRestaurantReq{
		Id:          "1",
		Name:        "Amarant",
		Address:     "Tashkent",
		PhoneNumber: "+9989925699823",
		Description: "Bonquet Hall",
	}

	mock.ExpectExec(`INSERT INTO restaurants`).
		WithArgs(req.Id, req.Name, req.Address, req.PhoneNumber, req.Description).
		WillReturnResult(sqlmock.NewResult(1, 1))

	resp, err := storage.CreateRestaurant(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestUpdateRestaurant(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewRestaurantStorage(db)

	req := &pb.CreateRestaurantReq{
		Id:          "1",
		Name:        "Updated Restaurant",
		Address:     "123 Updated Street",
		PhoneNumber: "123-456-7891",
		Description: "An updated restaurant",
	}

	mock.ExpectExec(`UPDATE restaurants`).
		WithArgs(req.Id, req.Name, req.Address, req.PhoneNumber, req.Description).
		WillReturnResult(sqlmock.NewResult(1, 1))

	resp, err := storage.UpdateRestaurant(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestDeleteRestaurant(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewRestaurantStorage(db)

	req := &pb.ById{
		Id: "1",
	}

	mock.ExpectExec(`DELETE FROM restaurants WHERE id = \$1`).
		WithArgs(req.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	resp, err := storage.DeleteRestaurant(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestGetRestaurant(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewRestaurantStorage(db)

	req := &pb.ById{
		Id: "a0c83280-d2f0-4168-9e61-47ab98345233",
	}

	mock.ExpectQuery(`SELECT id, name, address, phone_number, description FROM restaurants WHERE id = \$1`).
		WithArgs(req.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "address", "phone_number", "description"}).
			AddRow("a0c83280-d2f0-4168-9e61-47ab98345233", "Amarant", "Tashkent", "+9989925699823", "Bonquet Hall"))

	resp, err := storage.GetRestaurant(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "a0c83280-d2f0-4168-9e61-47ab98345233", req.Id)
	assert.Equal(t, "Amarant", resp.Name)
	assert.Equal(t, "Tashkent", resp.Address)
	assert.Equal(t, "+9989925699823", resp.PhoneNumber)
	assert.Equal(t, "Bonquet Hall", resp.Description)
}

