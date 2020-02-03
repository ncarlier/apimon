package rule

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PaesslerAG/gval"
	"github.com/PaesslerAG/jsonpath"
)

type jsonPathValidator struct {
	name string
	spec string
	eval gval.Evaluable
}

func newJSONPathValidator(spec string) (Validator, error) {
	builder := gval.Full(jsonpath.PlaceholderExtension())
	eval, err := builder.NewEvaluable(spec)
	if err != nil {
		return nil, err
	}
	validator := &jsonPathValidator{
		name: "json-path",
		spec: spec,
		eval: eval,
	}
	return validator, nil
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

	res, err := v.eval(context.Background(), jsonData)
	if err != nil {
		return err
	}
	// fmt.Println("RES=", res)
	if val, ok := res.([]interface{}); ok && len(val) == 0 {
		return fmt.Errorf("body does not match JSON path")
	}
	return nil
}
