package test

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
		{"name": "cnt", "status": "DOWN", "messageCount": 101}
	]
}
`

var jsonPathValidationTests = []struct {
	spec     string
	body     string
	expected error
}{
	{"$.services[0].status", json, nil},
	{"$.services[0].missing", json, fmt.Errorf("unknown key missing")},
	{"$.services[?(@.status == \"UP\")]", json, nil},
	{"$.services[?(@.status == \"UP\")].name", json, nil},
	{"$.services[?(@.messageCount > 100)]", json, nil},
	{"$.services[?(@.status == \"ERROR\")].name", json, fmt.Errorf("body does not match JSON path")},
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

func TestUglyJSONPathValidator(t *testing.T) {
	var spec = "$..[?(@.b && @.messageCount<100)]"
	var uglyJSON = `
{
	"foo-a": {
		"paused": false,
		"messageCount": 0
	},
	"foo-b": {
		"b": true,
		"messageCount": 99
	}
}
`
	rules := []config.Rule{config.Rule{Name: "json-path", Spec: spec}}
	pipeline, err := rule.CreateValidatorPipeline(rules)
	assert.Nil(t, err, "Pipeline creation should not fail")
	assert.NotNil(t, pipeline, "Pipeline should be created")
	assert.Equal(t, 1, len(pipeline), "Invalid validator pipeline")
	validator := pipeline[0]
	assert.Equal(t, "json-path", validator.Name(), "Invalid validator name")
	err = validator.Validate(uglyJSON, nil)
	assert.Nil(t, err, "Validation should pass")
}
