package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gergab1129/bookings/internal/config"
	"github.com/gergab1129/bookings/internal/driver"
	"github.com/gergab1129/bookings/internal/handlers"
	"github.com/gergab1129/bookings/internal/helpers"
	"github.com/gergab1129/bookings/internal/models"
	"github.com/gergab1129/bookings/internal/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber string = ":8080"

var infoLog *log.Logger
var errorLog *log.Logger
var app config.AppConfig
var session *scs.SessionManager

func main() {

	db, err := run()

	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	fmt.Println("Port number: %", portNumber)
	// _ = http.ListenAndServe(portNumber, nil)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

}

func run() (*driver.DB, error) {

	// what I am going to put in the session

	gob.Register(models.Reservation{})
	gob.Register(models.Room{})
	gob.Register(models.RoomRestriction{})
	gob.Register(models.Reservation{})
	gob.Register(models.Restriction{})

	// set to true when in production
	app.InProduciton = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduciton

	app.Session = session

	// connect to dabatabse

	log.Println("Connecting to database...")

	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=german sslmode=disable")

	if err != nil {
		log.Fatal("Cannot connect to database")
	}

	log.Println("Connected to database!")

	tc, err := render.CreateTemplateCache()

	if err != nil {
		fmt.Println("Cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
