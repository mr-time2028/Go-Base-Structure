package validators

import "net/http"

type Validation struct {
	Errors Errors
}

func New() Validation {
	return Validation{
		Errors{
			map[string][]string{},
			http.StatusBadRequest,
		},
	}
}

func (v *Validation) Valid() bool {
	return len(v.Errors.MessageMap) == 0
}
