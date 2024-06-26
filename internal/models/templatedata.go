package models

import "github.com/gergab1129/bookings/internal/forms"

// TemplateDate holds data set from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float64
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
	Form      *forms.Form
    IsAuthenticated int
}
