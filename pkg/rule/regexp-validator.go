package rule

import (
	"fmt"
	"net/http"
	"regexp"
)

type regexpValidator struct {
	name string
	spec string
	re   *regexp.Regexp
}

func newRegexpValidator(spec string) (Validator, error) {
	validator := &regexpValidator{
		name: "regexp",
		spec: spec,
	}
	var err error
	validator.re, err = regexp.Compile(spec)
	if err != nil {
		return nil, err
	}
	return validator, nil
}

func (v *regexpValidator) Name() string {
	return v.name
}

func (v *regexpValidator) Spec() string {
	return v.spec
}

func (v *regexpValidator) Validate(status int, headers http.Header, body string) error {
	if v.re.MatchString(body) {
		return nil
	}
	return fmt.Errorf("body does not match the RegExp")
}
