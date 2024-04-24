package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gergab1129/bookings/internal/models"
)

// type postData struct {
// 	key   string
// 	value string
// }

var Tests = []struct {
	name   string
	url    string
	method string
	// params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"gq", "/generals-quarters", "GET", http.StatusOK},
	{"ms", "/majors-suite", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},

	// {"post-search-avail", "/search-availability", "POST", []postData{
	// 	{key: "start", value: "2023-01-01"},
	// 	{key: "end", value: "2023-01-02"},
	// }, http.StatusOK},
	// {"post-search-avail-json", "/search-availability-json", "POST", []postData{
	// 	{key: "start", value: "2023-01-01"},
	// 	{key: "end", value: "2023-01-02"},
	// }, http.StatusOK},
	// {"make reservation", "/make-reservations", "POST", []postData{
	// 	{key: "first_name", value: "German"},
	// 	{key: "last_name", value: "Rodriguez"},
	// 	{key: "email", value: "me@here.com"},
	// 	{key: "phone", value: "3015398906"},
	// }, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, te := range Tests {
		if te.method == "GET" {

			resp, err := ts.Client().Get(ts.URL + te.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != te.expectedStatusCode {
				t.Errorf(
					"for %s expected %d but got %d",
					te.name,
					te.expectedStatusCode,
					resp.StatusCode,
				)
			}
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, "2024-02-01")
	endDate, _ := time.Parse(layout, "2024-02-02")

	reservation := models.Reservation{
		RoomId:    1,
		StartDate: startDate,
		EndDate:   endDate,
		Room: models.Room{
			Id:       1,
			RoomName: "General's Quarter",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf(
			"Reservation handler returned wrong response code: got %d wanted %d",
			rr.Code,
			http.StatusOK,
		)
	}

	// Test case when reservation is not in session

	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf(
			"Reservation hanlder returned wrong response code: got %d wanted %d",
			rr.Code,
			http.StatusTemporaryRedirect,
		)
	}

	reservation = models.Reservation{
		RoomId:    3,
		StartDate: startDate,
		EndDate:   endDate,
		Room: models.Room{
			Id:       3,
			RoomName: "General's Quarter",
		},
	}

	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf(
			"Reservation hanlder returned wrong response code: got %d wanted %d",
			rr.Code,
			http.StatusTemporaryRedirect,
		)
	}
}

func TestRepository_PostReservation(t *testing.T) {
	rBody := "start_date=2050-01-01"
	rBody = fmt.Sprintf("%s&%s", rBody, "end_date=2050-01-02")
	rBody = fmt.Sprintf("%s&%s", rBody, "first_name=German")
	rBody = fmt.Sprintf("%s&%s", rBody, "last_name=Rodriguez")
	rBody = fmt.Sprintf("%s&%s", rBody, "email=me@me.com")
	rBody = fmt.Sprintf("%s&%s", rBody, "phone=1233338047")
	rBody = fmt.Sprintf("%s&%s", rBody, "room_id=2")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(rBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf(
			"Reservation hanlder returned wrong response code for post reservation: got %d wanted %d",
			rr.Code,
			http.StatusSeeOther,
		)
	}

	// test for missing post body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf(
			"Reservation hanlder returned wrong response code for missing body: got %d wanted %d",
			rr.Code,
			http.StatusTemporaryRedirect,
		)
	}

	// test for invalid start date

	rBody = "start_date=2050-0-01"
	rBody = fmt.Sprintf("%s&%s", rBody, "end_date=2050-01-02")
	rBody = fmt.Sprintf("%s&%s", rBody, "first_name=German")
	rBody = fmt.Sprintf("%s&%s", rBody, "last_name=Rodriguez")
	rBody = fmt.Sprintf("%s&%s", rBody, "email=me@me.com")
	rBody = fmt.Sprintf("%s&%s", rBody, "phone=1233338047")
	rBody = fmt.Sprintf("%s&%s", rBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(rBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf(
			"Reservation hanlder returned wrong response code for invalid start date: got %d wanted %d",
			rr.Code,
			http.StatusTemporaryRedirect,
		)
	}

	// test for invalid end-date

	rBody = "start_date=2050-01-01"
	rBody = fmt.Sprintf("%s&%s", rBody, "end_date=2050-14-02")
	rBody = fmt.Sprintf("%s&%s", rBody, "first_name=German")
	rBody = fmt.Sprintf("%s&%s", rBody, "last_name=Rodriguez")
	rBody = fmt.Sprintf("%s&%s", rBody, "email=me@me.com")
	rBody = fmt.Sprintf("%s&%s", rBody, "phone=1233338047")
	rBody = fmt.Sprintf("%s&%s", rBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(rBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf(
			"Reservation hanlder returned wrong response code for invalid end date: got %d wanted %d",
			rr.Code,
			http.StatusTemporaryRedirect,
		)
	}

	// test for invalid room_id

	rBody = "start_date=2050-01-01"
	rBody = fmt.Sprintf("%s&%s", rBody, "end_date=2050-01-02")
	rBody = fmt.Sprintf("%s&%s", rBody, "first_name=German")
	rBody = fmt.Sprintf("%s&%s", rBody, "last_name=Rodriguez")
	rBody = fmt.Sprintf("%s&%s", rBody, "email=me@me.com")
	rBody = fmt.Sprintf("%s&%s", rBody, "phone=1233338047")
	rBody = fmt.Sprintf("%s&%s", rBody, "room_id=a")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(rBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf(
			"Reservation hanlder returned wrong response code for invalid room id: got %d wanted %d",
			rr.Code,
			http.StatusTemporaryRedirect,
		)
	}

	// test for form is_invalid

	rBody = "start_date=2050-01-01"
	rBody = fmt.Sprintf("%s&%s", rBody, "end_date=2050-01-02")
	rBody = fmt.Sprintf("%s&%s", rBody, "first_name=Ge")
	rBody = fmt.Sprintf("%s&%s", rBody, "last_name=Rodriguez")
	rBody = fmt.Sprintf("%s&%s", rBody, "email=me@me.com")
	rBody = fmt.Sprintf("%s&%s", rBody, "phone=1233338047")
	rBody = fmt.Sprintf("%s&%s", rBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(rBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf(
			"Reservation hanlder returned wrong response code for invalid form: got %d wanted %d",
			rr.Code,
			http.StatusSeeOther,
		)
	}

	// test for room reservation insertion error

	rBody = "start_date=2050-01-01"
	rBody = fmt.Sprintf("%s&%s", rBody, "end_date=2050-01-02")
	rBody = fmt.Sprintf("%s&%s", rBody, "first_name=German")
	rBody = fmt.Sprintf("%s&%s", rBody, "last_name=Rodriguez")
	rBody = fmt.Sprintf("%s&%s", rBody, "email=me@me.com")
	rBody = fmt.Sprintf("%s&%s", rBody, "phone=1233338047")
	rBody = fmt.Sprintf("%s&%s", rBody, "room_id=3")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(rBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf(
			"Reservation hanlder returned wrong response code inserting room reservation: got %d wanted %d",
			rr.Code,
			http.StatusTemporaryRedirect,
		)
	}

	// test for room restriction insertion error

	rBody = "start_date=2050-01-01"
	rBody = fmt.Sprintf("%s&%s", rBody, "end_date=2050-01-02")
	rBody = fmt.Sprintf("%s&%s", rBody, "first_name=German")
	rBody = fmt.Sprintf("%s&%s", rBody, "last_name=Rodriguez")
	rBody = fmt.Sprintf("%s&%s", rBody, "email=me@me.com")
	rBody = fmt.Sprintf("%s&%s", rBody, "phone=1233338047")
	rBody = fmt.Sprintf("%s&%s", rBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(rBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf(
			"Reservation hanlder returned wrong response code inserting room restriction: got %d wanted %d",
			rr.Code,
			http.StatusTemporaryRedirect,
		)
	}
}

func TestRepoitory_PostAvailability(t *testing.T) {
	// test posting availability

	rBody := "start=2050-01-03"
	rBody = fmt.Sprintf("%s&%s", rBody, "end=2050-01-02")

	req, _ := http.NewRequest("POST", "/search-availability", strings.NewReader(rBody))

	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostAvailability)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("failed posting availability: got %d wanted %d", rr.Code, http.StatusOK)
	}

	// test for missing body

	// test for missing post body
	req, _ = http.NewRequest("POST", "/search-availability", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf(
			"Reservation hanlder returned wrong response code for missing body: got %d wanted %d",
			rr.Code,
			http.StatusTemporaryRedirect,
		)
	}

	// test for invalid start-date

	rBody = "start=2050-14-01"
	rBody = fmt.Sprintf("%s&%s", rBody, "end=2050-01-02")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(rBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf(
			"Reservation hanlder returned wrong response code for invalid end date: got %d wanted %d",
			rr.Code,
			http.StatusTemporaryRedirect,
		)
	}

	// test for invalid end-date

	rBody = "start=2050-01-01"
	rBody = fmt.Sprintf("%s&%s", rBody, "end=2050-14-02")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(rBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf(
			"Reservation hanlder returned wrong response code for invalid end date: got %d wanted %d",
			rr.Code,
			http.StatusTemporaryRedirect,
		)
	}

	// test error returning query

	rBody = "start=2050-01-01"
	rBody = fmt.Sprintf("%s&%s", rBody, "end=2050-01-02")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(rBody))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("failed retunring query: got %d wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for length 0 available rooms slice

	rBody = "start=2050-01-02"
	rBody = fmt.Sprintf("%s&%s", rBody, "end=2050-01-02")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(rBody))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf(
			"failed posting returning 0 legth slice: got %d wanted %d",
			rr.Code,
			http.StatusSeeOther,
		)
	}
}

func TestRerpository_AvailabilityJSON(t *testing.T) {
	// test post availabilty JSON
	rBody := "start=2050-01-01"
	rBody = fmt.Sprintf("%s&%s", rBody, "end=2050-01-02")
	rBody = fmt.Sprintf("%s&%s", rBody, "room_id=1")

	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(rBody))

	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.AvailabilityJSON)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("failed posting availability: got %d wanted %d", rr.Code, http.StatusOK)
	}

	// test post availabilty JSON no body

	req, _ = http.NewRequest("POST", "/search-availability-json", nil)

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	handler.ServeHTTP(rr, req)

	var j jsonResponse
	_ = json.Unmarshal(rr.Body.Bytes(), &j)

	if j.OK {
		t.Error("expected OK to be false, got true")
	}

	// test post availabilty JSON
	rBody = "start=2050-14-01"
	rBody = fmt.Sprintf("%s&%s", rBody, "end=2050-01-02")
	rBody = fmt.Sprintf("%s&%s", rBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(rBody))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	_ = json.Unmarshal(rr.Body.Bytes(), &j)

	if j.OK {
		t.Error("PostAvailabilityJSON returned true, expected false. Parsing start date")
	}

	// test post availabilty JSON parse end date

	rBody = "start=2050-01-01"
	rBody = fmt.Sprintf("%s&%s", rBody, "end=2050-24-02")
	rBody = fmt.Sprintf("%s&%s", rBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(rBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	_ = json.Unmarshal(rr.Body.Bytes(), &j)

	if j.OK {
		t.Error("PostAvailabilityJSON returned true, expected false. Parsing start date")
	}

	// test room id parse

	rBody = "start=2050-01-01"
	rBody = fmt.Sprintf("%s&%s", rBody, "end=2050-01-02")
	rBody = fmt.Sprintf("%s&%s", rBody, "room_id=a")

	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(rBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	_ = json.Unmarshal(rr.Body.Bytes(), &j)

	if j.OK {
		t.Error("PostAvailabilityJSON returned true, expected false. Parsing room_id")
	}

	// test connection to database

	rBody = "start=2050-01-01"
	rBody = fmt.Sprintf("%s&%s", rBody, "end=2050-01-02")
	rBody = fmt.Sprintf("%s&%s", rBody, "room_id=2")

	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(rBody))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	_ = json.Unmarshal(rr.Body.Bytes(), &j)

	if j.OK {
		t.Error("PostAvailabilityJSON returned true, expected false. Parsing room_id")
	}
}

func TestRepository_ReservationSummary(t *testing.T) {
	// test post reservation summary with context
	layout := "2006-01-02"

	startDate, _ := time.Parse(layout, "2050-01-01")
	endDate, _ := time.Parse(layout, "2050-01-02")

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	req, _ := http.NewRequest("GET", "/reservation-summary", nil)

	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	session.Put(ctx, "reservation", res)

	handler := http.HandlerFunc(Repo.ReservationSummary)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf(
			"Reservation summary returned wrong response code: got %d expected %d",
			rr.Code,
			http.StatusOK,
		)
	}

	// test reservation summary without context

	req, _ = http.NewRequest("GET", "/reservation-summary", nil)

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf(
			"Reservation summary returned wrong response code: got %d expected %d",
			rr.Code,
			http.StatusTemporaryRedirect,
		)
	}
}

func TestRespository_ChooseRoom(t *testing.T) {
	// test choose room
	req, _ := http.NewRequest("GET", "/choose-room", nil)
	req.RequestURI = "/choose-room/1"

	res := models.Reservation{}

	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.ChooseRoom)

	session.Put(ctx, "reservation", res)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf(
			"Choose room returned wrong response code: got %d expected %d",
			rr.Code,
			http.StatusSeeOther,
		)
	}

	// test choose room without url id

	req, _ = http.NewRequest("GET", "/choose-room", nil)

	res = models.Reservation{}

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	session.Put(ctx, "reservation", res)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf(
			"Choose room returned wrong response code: got %d expected %d",
			rr.Code,
			http.StatusTemporaryRedirect,
		)
	}

	// test choose room without reservation id

	req, _ = http.NewRequest("GET", "/choose-room", nil)
	req.RequestURI = "/choose-room/1"

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf(
			"Choose room returned wrong response code: got %d expected %d",
			rr.Code,
			http.StatusTemporaryRedirect,
		)
	}

	// test choose room with reservation bad

	req, _ = http.NewRequest("GET", "/choose-room", nil)
	req.RequestURI = "/choose-room/a"

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf(
			"Choose room returned wrong response code: got %d expected %d",
			rr.Code,
			http.StatusTemporaryRedirect,
		)
	}
}

func TestRepository_BookRoom(t *testing.T) {
	// test book room
	req, _ := http.NewRequest("GET", "/book-room?s=2050-01-01&e=2050-01-02&id=1", nil)

	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.BookRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf(
			"Book room returned wrong response code: got %d expected %d",
			rr.Code,
			http.StatusSeeOther,
		)
	}

	// test book room bad id

	req, _ = http.NewRequest("GET", "/book-room?s=2050-01-01&e=2050-01-02&id=a", nil)

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf(
			"Book room returned wrong response code: got %d expected %d",
			rr.Code,
			http.StatusTemporaryRedirect,
		)
	}

	// test book room bad start date

	req, _ = http.NewRequest("GET", "/book-room?s=&e=2050-01-02&id=1", nil)

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf(
			"Book room returned wrong response code: got %d expected %d",
			rr.Code,
			http.StatusTemporaryRedirect,
		)
	}

	// test book room bad end date

	req, _ = http.NewRequest("GET", "/book-room?s=2050-01-01&e=&id=1", nil)

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf(
			"Book room bad end date returned wrong response code: got %d expected %d",
			rr.Code,
			http.StatusTemporaryRedirect,
		)
	}
}

func getCtx(r *http.Request) context.Context {
	ctx, err := session.Load(r.Context(), r.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
