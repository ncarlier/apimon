package rule

import (
	"fmt"
	"net/http"
	"regexp"
)

type regexpValidator struct {
	name string
	def  string
	re   *regexp.Regexp
}

func newRegexpValidator(param string) (Validator, error) {
	validator := &regexpValidator{
		name: "RegExp",
		def:  "RegExp:" + param,
	}
	var err error
	validator.re, err = regexp.Compile(param)
	if err != nil {
		return nil, err
	}
	return validator, nil
}

func (v *regexpValidator) Name() string {
	return v.name
}

func (v *regexpValidator) Def() string {
	return v.def
}

func (v *regexpValidator) Validate(status int, headers http.Header, body string) error {
	if v.re.MatchString(body) {
		return nil
	}
	return fmt.Errorf("body does not match the RegExp")
}
