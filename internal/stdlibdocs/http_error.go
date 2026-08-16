package stdlibdocs

import "fmt"

type httpStatusError struct {
	Endpoint   string
	StatusCode int
	Status     string
}

func (e *httpStatusError) Error() string {
	return fmt.Sprintf("GET %s returned %s", e.Endpoint, e.Status)
}
