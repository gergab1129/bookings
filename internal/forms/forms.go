package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Forms create a custom From struct
type Form struct {
	url.Values
	Errors errors
}

// type FormConfig map[string]interface{}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors{},
	}
}

// Has checks if a field is filled in form
func (f *Form) Has(field string) bool {

	if f.Values.Get(field) == "" {
		return false
	} else {
		return true
	}

}

// IsValid returns true if form is valid otherwiser returns false
func (f *Form) IsValid() bool {
	return len(f.Errors) == 0
}

// Required ranges trough the slice of field names and checks if blank
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)

		if strings.TrimSpace(value) == "" {
			f.Errors.add(field, "This field cannot be blank")
		}
	}
}

// MinLegth checks for string minimum length
func (f *Form) MinLegth(field string, length int) bool {

	fieldValue := f.Values.Get(field)
	if len([]rune(fieldValue)) < length {
		f.Errors.add(field, fmt.Sprintf("This field should have at least %v characters", length))
		return false
	} else {
		return true
	}
	// return len([]rune(fieldValue)) > length
}

// Email checks for valid e-mail address
func (f *Form) Email(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.add(field, "Invalid email address")
	}
}
