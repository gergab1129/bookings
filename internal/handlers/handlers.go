package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gergab1129/bookings/internal/config"
	"github.com/gergab1129/bookings/internal/driver"
	"github.com/gergab1129/bookings/internal/forms"
	"github.com/gergab1129/bookings/internal/helpers"
	"github.com/gergab1129/bookings/internal/models"
	"github.com/gergab1129/bookings/internal/render"
	"github.com/gergab1129/bookings/internal/repository"
	"github.com/gergab1129/bookings/internal/repository/dbrepo"
)

// Repo the repository used by handlers
var Repo *Repository

// Repository is the respository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB: dbrepo.NewPostgresRepo(db.SQL,
			a),
	}
}

// NewHandlers create a new repository
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	m.DB.AllUsers()
	render.Template(w, "home.page.tmpl", &models.TemplateData{}, r)
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	render.Template(w, "about.page.tmpl", &models.TemplateData{}, r)
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "contact.page.tmpl", &models.TemplateData{}, r)
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	emptyReservation := models.Reservation{}

	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.Template(w, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	}, r)
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	// Go reference time 01/02 03:04:05PM '06 -0700

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	endDate, err := time.Parse(layout, ed)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	roomId, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}


	reservation := models.Reservation{
		FirstName: form.Values.Get("first_name"),
		LastName:  form.Values.Get("last_name"),
		Email:     form.Values.Get("email"),
		Phone:     form.Values.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomId: roomId,
	}

	form.Required("first_name", "last_name", "email", "phone")
	form.MinLegth("first_name", 3)
	form.Email("email")

	if !form.IsValid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		}, r)

		return
	}

	reservationId, err := m.DB.InsertReservation(reservation)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	
	restriction := models.RoomRestriction{
		StartDate: startDate,
		EndDate: endDate,
		RoomId: roomId,
		ReservationId: reservationId,
		RestrictionId: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "majors.page.tmpl", &models.TemplateData{}, r)
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "generals.page.tmpl", &models.TemplateData{}, r)
}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "search-availability.page.tmpl",
		&models.TemplateData{},
		r,
	)

}

func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("Posted %s, %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json: "ok"`
	Message string `json: "message"`
}

func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	resp := jsonResponse{OK: true, Message: "Available"}

	out, err := json.MarshalIndent(resp, "", "     ")

	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {

	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		m.App.ErrorLog.Println("Can't get error from session")
		log.Println("Cannot get item from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, "reservation-summary.page.tmpl",
		&models.TemplateData{
			Data: data,
		}, r)
}
