package rule

import (
	"fmt"
	"strings"
)

// CreateValidatorPipeline create a pipeline of rule validator regarding the definition
func CreateValidatorPipeline(def string) ([]Validator, error) {
	defs := strings.Split(def, ";")
	pipeline := make([]Validator, len(defs))
	for i, d := range defs {
		validator, err := createValidator(d)
		if err != nil {
			return pipeline, err
		}
		pipeline[i] = validator
	}
	return pipeline, nil
}

func createValidator(def string) (Validator, error) {
	defParts := strings.SplitN(def, ":", 2)
	if len(defParts) != 2 {
		return nil, fmt.Errorf("rule is empty")
	}
	name := defParts[0]
	param := defParts[1]
	var rule Validator
	var err error
	switch name {
	case "Code":
		rule, err = newCodeValidator(param)
	case "RegExp":
		rule, err = newRegexpValidator(param)
	case "JSONPath":
		rule = newJSONPathValidator(param)
	default:
		err = fmt.Errorf("unknown rule name: %s", name)
	}
	return rule, err
}
