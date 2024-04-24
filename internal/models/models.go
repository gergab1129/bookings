package models

import (
	"time"
)

// Users is the user model
type User struct {
	Id          int       `db:"id"`
	FirstName   string    `db:"first_name"`
	LastName    string    `db:"last_name"`
	Email       string    `db:"email"`
	Password    string    `db:"password"`
	AccessLevel int       `db:"access_level"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// Rooms is the room model
type Room struct {
	Id        int
	RoomName  string    `db:"room_name"`
	Price     float32   `db:"price"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Restrictions is the restrictions model
type Restriction struct {
	Id              int
	RestrictionName string    `db:"restriction_name"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

// RoomRestrictions is the room restrictions model
type RoomRestriction struct {
	Id            int       `db:"id"`
	StartDate     time.Time `db:"start_date"`
	EndDate       time.Time `db:"end_date"`
	RoomId        int       `db:"room_id"`
	ReservationId int       `db:"reservation_id"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
	RestrictionId int       `db:"restriction_id"`
	Room          Room
	Reservation   Reservation
	Restriction   Restriction
}

// Reservations is the room restrictions model
type Reservation struct {
	Id        int       `db:"id"`
	StartDate time.Time `db:"start_date"`
	EndDate   time.Time `db:"end_date"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	Phone     string    `db:"phone"`
	RoomId    int       `db:"room_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Room      Room
}

// MailData holds and e-mail message
type MailData struct {
	To       string
	From     string
	Subject  string
	Content  string
	Template string
}
