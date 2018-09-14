package rule

import (
	"net/http"
)

// Validator interface of a rule validator
type Validator interface {
	Name() string
	Spec() string
	Validate(body string, resp *http.Response) error
}
