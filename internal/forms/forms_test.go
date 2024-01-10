package forms

import (
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	data := url.Values{}

	data.Add("a", "foo")
	data.Add("b", "bar")
	data.Add("c", "baz")

	form := New(data)

	test_form := Form{
		url.Values{},
		errors{},
	}
	
	compare := &test_form

	if reflect.TypeOf(form) != reflect.TypeOf(compare) {
		t.Error("form is not a pointer to *Form")
	}
}

func TestForm_MinLength(t *testing.T) {
	data := url.Values{}
	
	data.Add("a", "foo")
	data.Add("b", "bar")

	form := New(data)
	form.MinLegth("a", 3)

	if len(form.Errors) > 0 {
		t.Error("error checking for field length")
	}
}

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/", nil)

	testForm := New(r.PostForm)
	valid := testForm.IsValid()
	if !valid {
		t.Error("error added when not needed")
	}
}


func TestForm_Required(t *testing.T) {
	data := url.Values{}
	
	data.Add("a", "")
	data.Add("b", "bar")

	formRequired := New(data)
	formRequired.Required("a")
	if !(len(formRequired.Errors) > 0) {
		t.Error("form marked as valid when field required")
	}

	formNotRequired := New(data)

	formNotRequired.Required("b")
	// fmt.Println(formNotRequired.Errors)
	if (len(formNotRequired.Errors)>0) {
		t.Error("form marked as invalid when field is not required")
	}

}

func TestForm_Has(t *testing.T) {
	data := url.Values{}
	data.Add("a", "foo")

	form := New(data)

	if !form.Has("a") {
		t.Error("find field that does not exists")
	} 
}
