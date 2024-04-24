package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gergab1129/bookings/internal/config"
	"github.com/gergab1129/bookings/internal/driver"
	"github.com/gergab1129/bookings/internal/forms"
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

func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB: dbrepo.NewTestingRepo(
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

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := m.DB.SearchRoomById(&res.RoomId)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't find rooom")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room

	m.App.InfoLog.Println(res.Room.RoomName)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	}, r)
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't parse form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	// Go reference time 01/02 03:04:05PM '06 -0700

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't parse start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, ed)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't parse end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	roomId, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Invalid data")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	form := forms.New(r.PostForm)

	reservation := models.Reservation{
		FirstName: form.Values.Get("first_name"),
		LastName:  form.Values.Get("last_name"),
		Email:     form.Values.Get("email"),
		Phone:     form.Values.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomId:    roomId,
		Room:      models.Room{RoomName: form.Values.Get("room_name")},
	}

	form.Required("first_name", "last_name", "email", "phone")
	form.MinLegth("first_name", 3)
	form.Email("email")

	if !form.IsValid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		stringMap := make(map[string]string)
		stringMap["start_date"] = sd
		stringMap["end_date"] = ed
		http.Error(w, "form is not valid", http.StatusSeeOther)
		render.Template(w, "make-reservation.page.tmpl", &models.TemplateData{
			Form:      form,
			Data:      data,
			StringMap: stringMap,
		}, r)

		return
	}

	reservationId, err := m.DB.InsertReservation(reservation)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't insert reservation into database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomId:        roomId,
		ReservationId: reservationId,
		RestrictionId: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't insert room restriction into database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// send notifications

	htmlMessage := fmt.Sprintf(`
	<strong> Reservation Confirmation </strong><br>
	Dear %s: <br>
	This is to confirm your reservation from  %s to %s.
	`, reservation.FirstName, reservation.EndDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))

	msg := models.MailData{
		To:      reservation.Email,
		From:    "me@here.com",
		Subject: "Reservation Confirmation",
		Content: htmlMessage,
	}

	m.App.MailChan <- msg

	// send notification to property owner

	htmlMessage = fmt.Sprintf(`
	
	<strong> Your Room Is Booked </strong><br>

	Dear property owner:<br>

	Your room %s has been booked by %s %s. From %s to %s
	
	<br>
	You can contact them using the following phone number %s or through email: %s

	`, reservation.Room.RoomName, reservation.FirstName, reservation.LastName, reservation.StartDate.Format("2006-01-02"),
		reservation.EndDate.Format("2006-01-02"), reservation.Phone, reservation.Email)

	msg = models.MailData{
		To:       "property.owner@bookings.com",
		From:     "me@here.com",
		Subject:  "Room booked",
		Content:  htmlMessage,
		Template: "basic.html",
	}

	m.App.MailChan <- msg

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

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't parse form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, start)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't parse start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, end)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't parse end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	restriction := models.RoomRestriction{
		StartDate: startDate,
		EndDate:   endDate,
	}

	availableRooms, err := m.DB.SearchAvailabilityByDates(restriction)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't search rooms in database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if len(availableRooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No Availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
	}

	data := make(map[string]interface{})
	data["rooms"] = availableRooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	}, r)

	// data := make(map[string]interface{})
	// data["availableRooms"] = availableRooms

	// w.Write([]byte(fmt.Sprintf("Posted %s, %s", data["availableRooms"])))
}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomId    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "failed to parse form",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")

		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "cannot parse startdate",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")

		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: fmt.Sprint(err),
		}

		out, _ := json.MarshalIndent(resp, "", "     ")

		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	roomId, err := strconv.Atoi(r.Form.Get("room_id"))

	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "failed to parse form",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")

		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	res := models.RoomRestriction{
		StartDate: startDate,
		EndDate:   endDate,
		RoomId:    roomId,
	}

	available, err := m.DB.SearchRoomAvailability(res)

	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "can't connect to database ",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")

		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	resp := jsonResponse{OK: available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomId:    r.Form.Get("room_id")}

	out, _ := json.MarshalIndent(resp, "", "     ")

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

	stringMap := make(map[string]string)

	layout := "2006-01-02"
	sd := reservation.StartDate.Format(layout)
	ed := reservation.EndDate.Format(layout)

	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	// m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, "reservation-summary.page.tmpl",
		&models.TemplateData{
			Data:      data,
			StringMap: stringMap,
		}, r)
}

func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	if len(exploded) < 2 {
		m.App.Session.Put(r.Context(), "error", "Failed to get room id")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	roomId, err := strconv.Atoi(exploded[2])

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't parse room id")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		m.App.ErrorLog.Println("Can't get error from session")
		log.Println("Cannot get item from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	res.RoomId = roomId
	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservations", http.StatusSeeOther)
}

func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil {
		fmt.Println(err)
		m.App.Session.Put(r.Context(), "error", "Failed to get room id")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// sd := r.URL.Query().Get("s")
	// if sd == "" {
	// 	m.App.Session.Put(r.Context(), "error", "Failed to get reservation start")
	// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	// 	return
	// }

	// ed := r.URL.Query().Get("e")
	// if ed == "" {
	// 	m.App.Session.Put(r.Context(), "error", "Failed to get reservation end")
	// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	// 	return
	// }

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, r.URL.Query().Get("s"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Failed to get reservation start")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, r.URL.Query().Get("e"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Failed to get reservation start")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
		RoomId:    id,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservations", http.StatusSeeOther)

}

func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {

	render.Template(w, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	}, r)

}

// PostShowLogin handles the validation for user login
func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
    _ = m.App.Session.RenewToken(r.Context())

    err := r.ParseForm()
    if err != nil {
        log.Println(err)
    }

    form := forms.New(r.PostForm)
    form.Required("email", "password")
    form.Email("email")
    if !form.IsValid() {
    
        render.Template(w, "login.page.tmpl", &models.TemplateData{
            Form: form,
        }, r)
    }

    id, _, err := m.DB.Authenticate(form.Get("email"), form.Get("password"))
    if err != nil {
        log.Println(err)
        m.App.Session.Put(r.Context(), "error", "invalid login credentials")
        http.Redirect(w, r, "/user/login", http.StatusSeeOther)
    }

   m.App.Session.Put(r.Context(), "user_id", id)
   m.App.Session.Put(r.Context(), "flash", "logged in successfully")

   http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
   _ =  m.App.Session.Destroy(r.Context())
   _ = m.App.Session.RenewToken(r.Context())

   http.Redirect(w, r, "/", http.StatusSeeOther)
}


func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {

    render.Template(w, "admin-dashboard.page.tmpl", &models.TemplateData{}, r)

}
