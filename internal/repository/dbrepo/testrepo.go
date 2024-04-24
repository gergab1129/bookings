package dbrepo

import (
	"errors"

	"github.com/gergab1129/bookings/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts reservations into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {

	if res.RoomId > 2 {
		return 0, errors.New("error inserting reservation into database")
	}

	return 1, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *testDBRepo) InsertRoomRestriction(res models.RoomRestriction) error {

	if res.RoomId == 1 {
		return errors.New("error inserting room restriction into database")
	}

	return nil

}

// Return true if availability exists for room_id and false if no avaialbility
func (m *testDBRepo) SearchRoomAvailability(res models.RoomRestriction) (bool, error) {
	if res.RoomId == 2 {
		return false, errors.New("cannot search room availability")
	}
	return false, nil
}

func (m *testDBRepo) SearchAvailabilityByDates(res models.RoomRestriction) ([]models.Room, error) {
	var rooms []models.Room

	if res.StartDate.Format("2006-01-02") == "2050-01-01" {
		return nil, errors.New("error")
	}

	if res.StartDate.Format("2006-01-02") == "2050-01-02" {
		return rooms, nil
	}

	rooms = append(rooms, models.Room{
		Id:       1,
		RoomName: "General's Quarter",
	})

	return rooms, nil
}

func (m *testDBRepo) SearchRoomById(roomId *int) (string, error) {
	if *roomId > 2 {
		return "", errors.New("room not found")
	}
	return "", nil
}

func (m *testDBRepo) GetUserById(id int) (models.User, error) {
	return models.User{}, nil
}


func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

// Authenticate authenticates a user
func (m *testDBRepo) Authenticate(email, testPasswor string) (int, string, error) {
	return 0, "", nil
}