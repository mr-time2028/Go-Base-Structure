package validators

type Errors struct {
	MessageMap map[string][]string
	Code       int
}

func (e Errors) Add(field, message string) {
	e.MessageMap[field] = append(e.MessageMap[field], message)
}

func (e Errors) Get(field string) string {
	es := e.MessageMap[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
