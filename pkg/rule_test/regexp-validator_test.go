package rule_test

import (
	"fmt"
	"testing"

	"github.com/ncarlier/apimon/pkg/assert"
	"github.com/ncarlier/apimon/pkg/config"
	"github.com/ncarlier/apimon/pkg/rule"
)

var notmatch = fmt.Errorf("body does not match the RegExp")

var regExpValidationTests = []struct {
	spec     string
	body     string
	expected error
}{
	{"hello", "hello world", nil},
	{"(foo){2}", "foofoobar", nil},
	{"(foo){2}", "foobarbar", notmatch},
	{"foo", "bar", notmatch},
}

func TestRegExpValidator(t *testing.T) {
	for idx, tt := range regExpValidationTests {
		rules := []config.Rule{config.Rule{Name: "regexp", Spec: tt.spec}}
		pipeline, err := rule.CreateValidatorPipeline(rules)
		assert.Nil(t, err, "Pipeline creation should not fail")
		assert.NotNil(t, pipeline, "Pipeline should be created")
		assert.Equal(t, 1, len(pipeline), "Invalid validator pipeline")
		validator := pipeline[0]
		assert.Equal(t, "regexp", validator.Name(), "Invalid validator name")
		actual := validator.Validate(tt.body, nil)
		if (tt.expected == nil && actual != nil) ||
			(actual == nil && tt.expected != nil) ||
			(actual != nil && tt.expected.Error() != actual.Error()) {
			t.Errorf("Dataset(%d): expected %v, actual %v", idx, tt.expected, actual)
		}
	}
}

func TestBadRegExpValidator(t *testing.T) {
	rules := []config.Rule{config.Rule{Name: "regexp", Spec: "(?!re)"}}
	_, err := rule.CreateValidatorPipeline(rules)
	assert.NotNil(t, err, "Pipeline creation should fail")
}
