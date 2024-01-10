package forms

import (
	"net/url"
	"testing"
)

func TestError_Get(t *testing.T){
	data := url.Values{}
	data.Add("a", "foo")

	form := New(data)
	form.MinLegth("a", 4)

	if form.Errors.Get("a") == "" {
		t.Error("error not returned properly")
	}
}