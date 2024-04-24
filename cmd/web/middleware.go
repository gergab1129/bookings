package main

import (
	"net/http"

	"github.com/gergab1129/bookings/internal/helpers"
	"github.com/justinas/nosurf"
)

// func WriteToConsole(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Println("Hit the page")
// 		next.ServeHTTP(w, r)
// 	})
// }

// NoSurf adds CSRF protection to all POST request
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduciton,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad loads and saves the session in every request
func SessionLoad(next http.Handler) http.Handler {

	return session.LoadAndSave(next)

}


func Auth(next http.Handler) http.Handler {

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    
        if !helpers.Authenticated(r) {
            session.Put(r.Context(), "error", "Log in first!")        
            http.Redirect(w, r, "/user/login", http.StatusSeeOther)
            return
        }

        next.ServeHTTP(w, r)
    })
}
