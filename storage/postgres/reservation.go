package postgres

import (
	"database/sql"
	"fmt"
	"time"

	// "log"
	"strings"

	pb "github.com/Mubinabd/reservation_service/genproto"
)

type ReservationStorage struct {
	db *sql.DB
}

func NewReservationStorage(db *sql.DB) *ReservationStorage {
	return &ReservationStorage{
		db: db,
	}
}

func (r *ReservationStorage) CreateReservation(res *pb.ReservationCreate) (*pb.Void, error) {
	_,err := r.CheckReservation(&pb.ResrvationTime{ReservationTime: res.ReservationTime,RestaurantId: res.RestaurantId})
	if err!= nil {
        return nil, err
    }
	query := `insert into reservations(
		id,
		user_id,
		restaurant_id,
		reservation_time,
		status
	) values($1,$2,$3,$4,$5) `
	_, err = r.db.Exec(query, res.Id, res.UserId, res.RestaurantId, res.ReservationTime, res.Status)
	if err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}

func (r *ReservationStorage) UpdateReservation(res *pb.ReservationCreate) (*pb.Void, error) {
	_,err := r.CheckReservation(&pb.ResrvationTime{ReservationTime: res.ReservationTime,RestaurantId: res.RestaurantId})
	if err!= nil {
        return nil, err
    }

	query := `update reservations set
        user_id = $1,
        restaurant_id = $2,
        reservation_time = $3,
        status = $4,
		updated_at = now()
    where id = $5 and deleted_at = 0`
	_, err = r.db.Exec(query, res.UserId, res.RestaurantId, res.ReservationTime, res.Status, res.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}

func (r *ReservationStorage) GetReservation(id *pb.ById) (*pb.Reservation, error) {
	query := `select 
		user_id,
		restaurant_id,
		reservation_time,
		status 
		from reservations 
		where id = $1 and deleted_at = 0`
	row := r.db.QueryRow(query, id.Id)
	res := &pb.Reservation{}
	err := row.Scan(&res.UserId, &res.RestaurantId, &res.ReservationTime, &res.Status)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *ReservationStorage) DeleteReservation(id *pb.ById) (*pb.Void, error) {
	query := `UPDATE reservations SET deleted_at = EXTRACT(EPOCH FROM NOW()) WHERE id = $1`

	_, err := r.db.Exec(query, id.Id)
	if err != nil {
		return nil, err
	}

	return &pb.Void{}, nil
}

func (r *ReservationStorage) GetReservationByFilter(filter *pb.FilterByTime) (*pb.Reservations, error) {
	baseQuery := `
        SELECT 
            user_id,
            restaurant_id,
            reservation_time,
            status 
        FROM reservations 
        WHERE deleted_at = 0`

	var conditions []string
	var args []interface{}
	paramIndex := 1

	if filter.RestaurantId != "" {
		conditions = append(conditions, fmt.Sprintf("restaurant_id = $%d", paramIndex))
		args = append(args, filter.RestaurantId)
		paramIndex++
	}
	if filter.ReservationFrom != "" {
		conditions = append(conditions, fmt.Sprintf("reservation_time >= $%d", paramIndex))
		args = append(args, filter.ReservationFrom)
		paramIndex++
	}
	if filter.ReservationTo != "" {
		conditions = append(conditions, fmt.Sprintf("reservation_time <= $%d", paramIndex))
		args = append(args, filter.ReservationTo)
		paramIndex++
	}

	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations pb.Reservations
	for rows.Next() {
		var reservation pb.Reservation
		err := rows.Scan(&reservation.UserId, &reservation.RestaurantId, &reservation.ReservationTime, &reservation.Status)
		if err != nil {
			return nil, err
		}
		reservations.Reservations = append(reservations.Reservations, &reservation)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &reservations, nil
}

func (r *ReservationStorage) GetTotalSum(id *pb.ById) (*pb.Total, error) {
	query := `
	SELECT 
		SUM(m.price * ro.quantity) AS total_sum
	FROM 
		reservationsorders ro
	JOIN 
		menu m ON ro.menu_item_id = m.id
	WHERE 
		ro.reservation_id = $1
	AND 
		ro.deleted_at = 0
	AND 
		m.deleted_at = 0;`

	var total float32
	err := r.db.QueryRow(query, id.Id).Scan(&total)
	// log.Println(111111,total,1111111111,err)
	if err != nil {
		return nil, err
	}
	return &pb.Total{
		Total: total,
	}, nil
}

func parseTime(timeStr string) (time.Time, error) {
	formats := []string{
		time.RFC3339,           // "2006-01-02T15:04:05Z07:00"
		"2006-01-02 15:04:05",  // "2006-01-02 15:04:05"
	}
	var t time.Time
	var err error
	for _, format := range formats {
		t, err = time.Parse(format, timeStr)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid time format: %v", timeStr)
}

func (r *ReservationStorage) CheckReservation(timeReq *pb.ResrvationTime) (*pb.Void, error) {

	reservationTime, err := parseTime(timeReq.ReservationTime)
	if err != nil {
		return nil, fmt.Errorf("invalid reservation time format: %w", err)
	}

	endTime := reservationTime.Add(59 * time.Minute)

	query := `
		SELECT COUNT(*)
		FROM reservations
		WHERE restaurant_id = $1
		AND reservation_time >= $2
		AND reservation_time <= $3
		AND deleted_at = 0
	`

	var count int
	err = r.db.QueryRow(query, timeReq.RestaurantId, reservationTime, endTime).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("error checking reservation: %w", err)
	}

	if count > 0 {
		return nil, fmt.Errorf("the time slot from %s to %s is already reserved", reservationTime, endTime)
	}

	return &pb.Void{}, nil
}