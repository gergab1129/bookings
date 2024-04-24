package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/gergab1129/bookings/internal/models"
)

// AppConfig holds the application config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduciton  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData
}
