package dbrepo

import (
	"context"
	"time"

	"github.com/gergab1129/bookings/internal/models"
)

func (m *postgresDBrepo) AllUsers() bool {
	return true
}

// InsertReservation inserts reservations into the database
func (m *postgresDBrepo) InsertReservation(res models.Reservation) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var newID int

	stm := ` insert into reservations (first_name, last_name, email, phone, start_date
	, end_date, room_id, created_at, updated_at)
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	err := m.DB.QueryRowContext(ctx, stm, res.FirstName, res.LastName, res.Email, res.Phone, res.StartDate, res.EndDate,
		res.RoomId, time.Now(), time.Now()).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

//InsertRoomRestriction inserts a room restriction into the database
func (m *postgresDBrepo) InsertRoomRestriction(res models.RoomRestriction) error {
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stm := `
		INSERT INTO room_restrictions (
			start_date,
			end_date, 
			room_id, 
			reservation_id, 
			created_at,
			updated_at,
			restriction_id
		) values (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7);
	`

	_, err := m.DB.ExecContext(ctx, stm, res.StartDate, res.EndDate, res.RoomId, 
	res.ReservationId, time.Now(), time.Now(), res.RestrictionId)

	if err != nil {
		return err
	}

	return nil

}