package rule

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/antonmedv/expr"
)

type jsonExprValidator struct {
	name string
	spec string
}

func newJSONExprValidator(spec string) Validator {
	validator := &jsonExprValidator{
		name: "expr",
		spec: spec,
	}
	return validator
}

func (v *jsonExprValidator) Name() string {
	return v.name
}

func (v *jsonExprValidator) Spec() string {
	return v.spec
}

func (v *jsonExprValidator) Validate(status int, headers http.Header, body string) error {
	var m map[string]interface{}
	err := json.Unmarshal([]byte(body), &m)
	if err != nil {
		return fmt.Errorf("Unable to read body as JSON document: %s", err)
	}
	result, err := expr.Eval(v.spec, m)
	if err != nil {
		return fmt.Errorf("Unable to apply expression on JSON body: %s", err)
	}
	if toBool(result) {
		return nil
	}
	return fmt.Errorf("Bad expression result: %v", result)
}

func toBool(i1 interface{}) bool {
	if i1 == nil {
		return false
	}
	switch i2 := i1.(type) {
	default:
		return false
	case bool:
		return i2
	case string:
		return i2 == "true"
	case int:
		return i2 != 0
	case *bool:
		if i2 == nil {
			return false
		}
		return *i2
	case *string:
		if i2 == nil {
			return false
		}
		return *i2 == "true"
	case *int:
		if i2 == nil {
			return false
		}
		return *i2 != 0
	}
	return false
}
