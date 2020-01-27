package rule_test

import (
	"fmt"
	"testing"

	"github.com/ncarlier/apimon/pkg/assert"
	"github.com/ncarlier/apimon/pkg/config"
	"github.com/ncarlier/apimon/pkg/rule"
)

var json = `
{
	"services": [
		{"name": "foo", "status": "UP"},
		{"name": "bar", "status": "DOWN"},
		{"name": "cnt", "status": "DOWN", "messageCount": 100}
	]
}
`

var jsonPathValidationTests = []struct {
	spec     string
	body     string
	expected error
}{
	{"$.services[0].status", json, nil},
	{"$.services[0].missing", json, fmt.Errorf("key error: missing not found in object")},
	{"$.services[?(@.status == 'UP')]", json, nil},
	{"$.services[?(@.status == 'UP')].name", json, nil},
	{"$.services[?(@.messageCount < 1000)]", json, nil},
	{"$.services[?(@.status == 'ERROR')].name", json, fmt.Errorf("body does not match JSON path")},
	{"...", json, fmt.Errorf("should start with '$'")},
}

func TestJSONPathValidator(t *testing.T) {
	for idx, tt := range jsonPathValidationTests {
		rules := []config.Rule{config.Rule{Name: "json-path", Spec: tt.spec}}
		pipeline, err := rule.CreateValidatorPipeline(rules)
		assert.Nil(t, err, "Pipeline creation should not fail")
		assert.NotNil(t, pipeline, "Pipeline should be created")
		assert.Equal(t, 1, len(pipeline), "Invalid validator pipeline")
		validator := pipeline[0]
		assert.Equal(t, "json-path", validator.Name(), "Invalid validator name")
		actual := validator.Validate(tt.body, nil)
		if (tt.expected == nil && actual != nil) ||
			(actual == nil && tt.expected != nil) ||
			(actual != nil && tt.expected.Error() != actual.Error()) {
			t.Errorf("Dataset(%d): expected %v, actual %v", idx, tt.expected, actual)
		}
	}
}
