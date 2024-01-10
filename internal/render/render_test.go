package render

import (
	"net/http"
	"testing"

	"github.com/gergab1129/bookings/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	
	var td models.TemplateData
	
	r, err := getSession()

	if err != nil {
		t.Error(err)
	}
	
	session.Put(r.Context(), "flash", "123")
	result := AddDefaultData(&td, r)

	if result.Flash != "123" {
		t.Error("flash value of 123 not found in session")
	}

}

func TestRenderTemplate(t *testing.T)  {

	pathToTemplates =  "./../../templates"

	tc, err := CreateTemplateCache()

	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	r, err := getSession()

	if err != nil {
		t.Error(err)
	}
	
	var ww myWriter
	err = RenderTemplate(ww, "home.page.tmpl", &models.TemplateData{}, r)

	if err != nil { 
		t.Error("error writing template to browser")
	}

	err = RenderTemplate(ww, "err-home.page.tmpl", &models.TemplateData{}, r)

	if err == nil { 
		t.Error("render template that don't exist")
	}
}


func getSession() (*http.Request, error) {

	r, err  :=  http.NewRequest("GET", "/", nil)
	
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil

}

func TestNewTemplates(t *testing.T){
	NewTemplates(app)
}

func TestCreateTempleateCache(t *testing.T) {
	pathToTemplates = "./../../templates"

	_, err := CreateTemplateCache()

	if err != nil {
		t.Error(err)
	}
}