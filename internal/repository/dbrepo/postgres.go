package dbrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

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
		res.RoomId, time.Now(), time.Now()).
		Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

// InsertRoomRestriction inserts a room restriction into the database
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

// Return true if availability exists for room_id and false if no avaialbility
func (m *postgresDBrepo) SearchRoomAvailability(res models.RoomRestriction) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stm := `SELECT count(id)
	from room_restrictions 
	where start_date <= $1 and end_date >= $2
	and room_id = $3 
	group by room_id`

	row := m.DB.QueryRowContext(ctx, stm, res.EndDate, res.StartDate, res.RoomId)

	var numOfRows int

	err := row.Scan(&numOfRows)

	switch {
	case err == sql.ErrNoRows:
		return true, nil

	case err != nil:
		return false, err

	default:
		return false, nil
	}
}

func (m *postgresDBrepo) SearchAvailabilityByDates(
	res models.RoomRestriction,
) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `

	Select id, room_name, price
	from public.rooms 
	where id not in (select room_id 
		from public.room_restrictions 
		where start_date <= $2 and end_date >= $1);
	`

	var rooms []models.Room
	fmt.Printf("Now: %s", time.Now())
	fmt.Printf("startDate: %s \n endDate: %s", res.StartDate, res.EndDate)
	err := m.DB.SelectContext(ctx, &rooms, query, res.StartDate, res.EndDate)
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (m *postgresDBrepo) SearchRoomById(roomId *int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT room_name FROM rooms where id = $1
	`
	rows, err := m.DB.QueryxContext(ctx, query, roomId)
	if err != nil {
		return "", err
	}

	var roomName string
	for rows.Next() {
		err = rows.Scan(&roomName)
		m.App.InfoLog.Println(roomName)
	}

	if err != nil {
		return "", err
	}
	return roomName, nil
}

func (m *postgresDBrepo) GetUserById(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, first_name, last_name, email, password, access_level,
			created_at, updated_at
		FROM public.users
		where id = $1
	`

	var u models.User

	err := m.DB.QueryRowxContext(ctx, query, id).StructScan(&u)
	if err != nil {
		return u, err
	}
	return u, nil
}

func (m *postgresDBrepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		UPDATE public.useres SET first_name = $1, last_name = $2, email = $3, acces_level = $4, updated_at = $5
	`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

// Authenticate authenticates a user
func (m *postgresDBrepo) Authenticate(email, testPasswor string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashPassword string

	row := m.DB.QueryRowxContext(ctx, "Select id, password from users where email = $1", email)
	err := row.Scan(&id, &hashPassword)
	if err != nil {
		return 0, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(testPasswor))

	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashPassword, nil
}
