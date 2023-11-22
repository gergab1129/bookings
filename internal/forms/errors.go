package forms


type errors map[string][]string

// Add errors message for a given form field
func (e errors) add(field string, msg string) {

	e[field] = append(e[field], msg)
}

//Get returns the first error message
func (e errors) Get(field string) string {
	es := e[field]

	if len(es) == 0 {
		return ""
	} else {
		return es[0]	
	}
}