package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gergab1129/bookings/internal/config"
)

var app *config.AppConfig

// NewHelpers sets app config for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Authenticated validates if user is authenticated
func Authenticated(r *http.Request) bool {
    exits := app.Session.Exists(r.Context(), "user_id")

    return exits
}
