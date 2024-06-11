package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pb "github.com/Mubinabd/reservation_service/genproto"
)

func TestCreateReservation(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	reservationStorage := NewReservationStorage(db)

	res := &pb.ReservationCreate{
		Id:             "1",
		UserId:         "user1",
		RestaurantId:   "restaurant1",
		ReservationTime: "2024-06-11T12:00:00Z",
		Status:         "confirmed",
	}

	mock.ExpectExec("insert into reservations").
		WithArgs(res.Id, res.UserId, res.RestaurantId, res.ReservationTime, res.Status).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = reservationStorage.CreateReservation(res)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func TestUpdateReservation(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	reservationStorage := NewReservationStorage(db)

	res := &pb.ReservationCreate{
		Id:              "1",
		UserId:          "user1",
		RestaurantId:    "restaurant1",
		ReservationTime: "2024-06-11T12:00:00Z",
		Status:          "confirmed",
	}

	mock.ExpectExec(`update reservations set user_id = \$1, restaurant_id = \$2, reservation_time = \$3, status = \$4, updated_at = now\(\) where id = \$5 and deleted_at = 0`).
		WithArgs(res.UserId, res.RestaurantId, res.ReservationTime, res.Status, res.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = reservationStorage.UpdateReservation(res)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetReservation(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	reservationStorage := NewReservationStorage(db)

	id := &pb.ById{Id: "1"}

	mock.ExpectQuery(`select user_id, restaurant_id, reservation_time, status from reservations where id = \$1 and deleted_at = 0`).
		WithArgs(id.Id).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "restaurant_id", "reservation_time", "status"}).
			AddRow("user1", "restaurant1", "2024-06-11T12:00:00Z", "confirmed"))

	res, err := reservationStorage.GetReservation(id)
	require.NoError(t, err)
	assert.Equal(t, "user1", res.UserId)
	assert.Equal(t, "restaurant1", res.RestaurantId)
	assert.Equal(t, "2024-06-11T12:00:00Z", res.ReservationTime)
	assert.Equal(t, "confirmed", res.Status)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestDeleteReservation(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	reservationStorage := NewReservationStorage(db)

	id := &pb.ById{Id: "1"}

	mock.ExpectExec(`UPDATE reservations SET deleted_at = EXTRACT\(EPOCH FROM NOW\(\)\) WHERE id = \$1`).
		WithArgs(id.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = reservationStorage.DeleteReservation(id)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetReservationByFilter(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	reservationStorage := NewReservationStorage(db)

	filter := &pb.FilterByTime{
		RestaurantId:    "restaurant1",
		ReservationFrom: "2024-06-10T00:00:00Z",
		ReservationTo:   "2024-06-12T00:00:00Z",
	}

	query := `
        SELECT 
            user_id,
            restaurant_id,
            reservation_time,
            status 
        FROM reservations 
        WHERE deleted_at = 0 AND restaurant_id = \$1 AND reservation_time >= \$2 AND reservation_time <= \$3`
	mock.ExpectQuery(query).
		WithArgs(filter.RestaurantId, filter.ReservationFrom, filter.ReservationTo).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "restaurant_id", "reservation_time", "status"}).
			AddRow("user1", "restaurant1", "2024-06-11T12:00:00Z", "confirmed"))

	reservations, err := reservationStorage.GetReservationByFilter(filter)
	require.NoError(t, err)
	assert.Len(t, reservations.Reservations, 1)
	assert.Equal(t, "user1", reservations.Reservations[0].UserId)
	assert.Equal(t, "restaurant1", reservations.Reservations[0].RestaurantId)
	assert.Equal(t, "2024-06-11T12:00:00Z", reservations.Reservations[0].ReservationTime)
	assert.Equal(t, "confirmed", reservations.Reservations[0].Status)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
