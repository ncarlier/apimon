package rule

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type codeValidator struct {
	name     string
	def      string
	fromCode int
	toCode   int
	codes    []int
}

func newCodeValidator(param string) (Validator, error) {
	validator := &codeValidator{
		name: "Code",
		def:  "Code:" + param,
	}
	if strings.ContainsAny(param, "-") {
		codes := strings.SplitN(param, "-", 2)
		validator.fromCode, _ = strconv.Atoi(codes[0])
		validator.toCode, _ = strconv.Atoi(codes[1])
	} else {
		codes := strings.Split(param, ",")
		for _, c := range codes {
			code, err := strconv.Atoi(c)
			if err != nil {
				// Ignore error. Should not happen.
				continue
			}
			validator.codes = append(validator.codes, code)
		}
	}
	return validator, nil
}

func (v *codeValidator) Name() string {
	return v.name
}

func (v *codeValidator) Def() string {
	return v.def
}

func (v *codeValidator) Validate(status int, headers http.Header, body string) error {
	if len(v.codes) > 0 {
		for _, c := range v.codes {
			if status == c {
				return nil
			}
		}
		return fmt.Errorf("Unexpected status code: %d", status)
	} else if status >= v.fromCode && status <= v.toCode {
		return nil
	}
	return fmt.Errorf("Status out of range: %d", status)
}
