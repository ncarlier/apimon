package rule_test

import (
	"fmt"
	"testing"

	"github.com/ncarlier/apimon/pkg/assert"
	"github.com/ncarlier/apimon/pkg/config"
	"github.com/ncarlier/apimon/pkg/rule"
)

var exprJSON = `
{
	"status": "DOWN",
	"serviceA": {
		"status": "UP"
	}
}
`

var jsonExprValidationTests = []struct {
	spec     string
	body     string
	expected error
}{
	{"status == \"DOWN\"", exprJSON, nil},
	{"serviceA.status == \"UP\"", exprJSON, nil},
	{"status == \"UP\"", exprJSON, fmt.Errorf("Bad expression result: false")},
}

func TestJSONExprValidator(t *testing.T) {
	for idx, tt := range jsonExprValidationTests {
		rules := []config.Rule{config.Rule{Name: "json-expr", Spec: tt.spec}}
		pipeline, err := rule.CreateValidatorPipeline(rules)
		assert.Nil(t, err, "Pipeline creation should not fail")
		assert.NotNil(t, pipeline, "Pipeline should be created")
		assert.Equal(t, 1, len(pipeline), "Invalid validator pipeline")
		validator := pipeline[0]
		assert.Equal(t, "expr", validator.Name(), "Invalid validator name")
		actual := validator.Validate(200, nil, tt.body)
		if (tt.expected == nil && actual != nil) ||
			(actual == nil && tt.expected != nil) ||
			(actual != nil && tt.expected.Error() != actual.Error()) {
			t.Errorf("Dataset(%d): expected %v, actual %v", idx, tt.expected, actual)
		}
	}
}
