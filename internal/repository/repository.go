package repository

import "github.com/gergab1129/bookings/internal/models"

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(res models.RoomRestriction) error
	SearchRoomAvailability(res models.RoomRestriction) (bool, error)
	SearchAvailabilityByDates(res models.RoomRestriction) ([]models.Room, error)
	SearchRoomById(roomId *int) (string, error)
}
