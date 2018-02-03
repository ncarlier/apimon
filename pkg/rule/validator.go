package rule

import (
	"net/http"
)

// Validator interface of a rule validator
type Validator interface {
	Name() string
	Def() string
	Validate(status int, headers http.Header, body string) error
}
