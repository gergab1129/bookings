package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gergab1129/bookings/internal/config"
	"github.com/gergab1129/bookings/internal/forms"
	"github.com/gergab1129/bookings/internal/models"
	"github.com/gergab1129/bookings/internal/render"
)

// Repo the repository used by handlers
var Repo *Repository

// Repository is the respository type 
type Repository struct {
	App *config.AppConfig
}

//NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository {
		App:  a,
	}
}

// NewHandlers create a new repository
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{}, r)
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	stringMap := make(map[string]string)
	stringMap["test"] =  "Hello, again."

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")

	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap, 
	}, r)
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "contact.page.tmpl", &models.TemplateData{}, r)
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	
	emptyReservation  := models.Reservation{}
	
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.RenderTemplate(w, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	}, r)
}

//PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		fmt.Println(err)
		return
	}

	form := forms.New(r.PostForm)

	reservation := models.Reservation{
		FirstName: form.Values.Get("first_name"),
		LastName: form.Values.Get("last_name"),
		Email: form.Values.Get("email"),
		Phone: form.Values.Get("phone"),
	}

	form.Required("first_name", "last_name", "email", "phone")
	form.MinLegth("first_name", 3)
	form.Email("email")
	
	if !form.IsValid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTemplate(w, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		}, r)

		return
	}

}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "majors.page.tmpl", &models.TemplateData{}, r)
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "generals.page.tmpl", &models.TemplateData{}, r)
}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "search-availability.page.tmpl",
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
	OK bool `json: "ok"`;
	Message string `json: "message"`
}

func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{OK: true, Message: "Available"}

	out, err := json.MarshalIndent(resp, "", "     ")

	if err != nil {
		log.Fatal("Error parsing json")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}