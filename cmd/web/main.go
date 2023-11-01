package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gergab1129/bookings/pkg/config"
	"github.com/gergab1129/bookings/pkg/handlers"
	"github.com/gergab1129/bookings/pkg/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber string = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {


	// set to true when in production
	app.InProduciton = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduciton

	app.Session = session

	tc, err :=  render.CreateTemplateCache()

	if err != nil {
		fmt.Println("Cannot create template cache")
		os.Exit(1)
	}
	
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)

	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	fmt.Println("Port number: %", portNumber) 
	// _ = http.ListenAndServe(portNumber, nil)

	srv := &http.Server {
		Addr: portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

}
