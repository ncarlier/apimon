package rule_test

import (
	"fmt"
	"testing"

	"github.com/ncarlier/apimon/pkg/assert"
	"github.com/ncarlier/apimon/pkg/rule"
)

var json = `
{
	"services": [
		{"name": "foo", "status": "UP"},
		{"name": "bar", "status": "DOWN"}
	]
}
`

var jsonPathValidationTests = []struct {
	rule     string
	body     string
	expected error
}{
	{"JSONPath:$.services[0].status", json, nil},
	{"JSONPath:$.services[0].missing", json, fmt.Errorf("key error: missing not found in object")},
	{"JSONPath:$.services[?(@.status == 'UP')]", json, nil},
	{"JSONPath:$.services[?(@.status == 'UP')].name", json, nil},
	{"JSONPath:$.services[?(@.status == 'ERROR')].name", json, fmt.Errorf("body does not match JSON path")},
	{"JSONPath:...", json, fmt.Errorf("should start with '$'")},
}

func TestJSONPathValidator(t *testing.T) {
	for idx, tt := range jsonPathValidationTests {
		pipeline, err := rule.CreateValidatorPipeline(tt.rule)
		assert.Nil(t, err, "Pipeline creation should not fail")
		assert.NotNil(t, pipeline, "Pipeline should be created")
		assert.Equal(t, 1, len(pipeline), "Invalid validator pipeline")
		validator := pipeline[0]
		assert.Equal(t, "JSONPath", validator.Name(), "Invalid validator name")
		actual := validator.Validate(200, nil, tt.body)
		if (tt.expected == nil && actual != nil) ||
			(actual == nil && tt.expected != nil) ||
			(actual != nil && tt.expected.Error() != actual.Error()) {
			t.Errorf("Dataset(%d): expected %v, actual %v", idx, tt.expected, actual)
		}
	}
}
