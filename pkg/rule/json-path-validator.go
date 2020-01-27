package rule

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/oliveagle/jsonpath"
)

type jsonPathValidator struct {
	name string
	spec string
}

func newJSONPathValidator(spec string) Validator {
	validator := &jsonPathValidator{
		name: "json-path",
		spec: spec,
	}
	return validator
}

func (v *jsonPathValidator) Name() string {
	return v.name
}

func (v *jsonPathValidator) Spec() string {
	return v.spec
}

func (v *jsonPathValidator) Validate(body string, resp *http.Response) error {
	var jsonData interface{}
	err := json.Unmarshal([]byte(body), &jsonData)
	if err != nil {
		return fmt.Errorf("body is not valid JSON")
	}
	res, err := jsonpath.JsonPathLookup(jsonData, v.spec)
	if err != nil {
		return err
	}
	if val, ok := res.([]interface{}); ok == true {
		if len(val) == 0 {
			return fmt.Errorf("body does not match JSON path")
		}
	}
	// fmt.Println("RES=", res)
	return nil
}
