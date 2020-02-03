package rule

import (
	"fmt"

	"github.com/ncarlier/apimon/pkg/config"
)

// CreateValidatorPipeline create a pipeline of rule validator regarding the definition
func CreateValidatorPipeline(rules []config.Rule) ([]Validator, error) {
	if len(rules) == 0 {
		rules = []config.Rule{config.Rule{Name: "code", Spec: "200"}}
	}
	pipeline := make([]Validator, len(rules))
	for i, rule := range rules {
		validator, err := createValidator(rule)
		if err != nil {
			return pipeline, err
		}
		pipeline[i] = validator
	}
	return pipeline, nil
}

func createValidator(rule config.Rule) (Validator, error) {
	var result Validator
	var err error
	switch rule.Name {
	case "code":
		result, err = newCodeValidator(rule.Spec)
	case "regexp":
		result, err = newRegexpValidator(rule.Spec)
	case "json-path":
		result, err = newJSONPathValidator(rule.Spec)
	case "json-expr":
		result = newJSONExprValidator(rule.Spec)
	case "cert":
		result, err = newCertValidator(rule.Spec)
	default:
		err = fmt.Errorf("unknown rule name: %s", rule.Name)
	}
	return result, err
}
