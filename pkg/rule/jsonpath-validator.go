package rule

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/oliveagle/jsonpath"
)

type jsonpathValidator struct {
	name string
	spec string
}

func newJSONPathValidator(spec string) Validator {
	validator := &jsonpathValidator{
		name: "json-path",
		spec: spec,
	}
	return validator
}

func (v *jsonpathValidator) Name() string {
	return v.name
}

func (v *jsonpathValidator) Spec() string {
	return v.spec
}

func (v *jsonpathValidator) Validate(status int, headers http.Header, body string) error {
	var jsonData interface{}
	err := json.Unmarshal([]byte(body), &jsonData)
	if err != nil {
		return fmt.Errorf("body is not valid JSON")
	}
	res, err := jsonpath.JsonPathLookup(jsonData, v.spec)
	if err != nil {
		return err
	}
	if v, ok := res.([]interface{}); ok == true {
		if len(v) == 0 {
			return fmt.Errorf("body does not match JSON path")
		}
	}
	// fmt.Println("RES=", res)
	return nil
}
