package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var Tests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"gq", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"ms", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"sa", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"mr", "/make-reservations", "GET", []postData{}, http.StatusOK},
	{"post-search-avail", "/search-availability", "POST", []postData{
		{key: "start", value: "2023-01-01"},
		{key: "end", value: "2023-01-02"},
	}, http.StatusOK},
	{"post-search-avail-json", "/search-availability-json", "POST", []postData{
		{key: "start", value: "2023-01-01"},
		{key: "end", value: "2023-01-02"},
	}, http.StatusOK},
	{"make reservation", "/make-reservations", "POST", []postData{
		{key: "first_name", value: "German"},
		{key: "last_name", value: "Rodriguez"},
		{key: "email", value: "me@here.com"},
		{key: "phone", value: "3015398906"},
	}, http.StatusOK},
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
				t.Errorf("for %s expected %d but got %d", te.name, te.expectedStatusCode, resp.StatusCode)
			}

		} else {

			values := url.Values{}

			for _, p := range te.params {
				// fmt.Printf("key %v, value %v", p.key, p.value)
				values.Add(p.key, p.value)
			}

			resp, err := ts.Client().PostForm(ts.URL+te.url, values)

			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != te.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", te.name, te.expectedStatusCode, resp.StatusCode)
			}

		}
	}

}
