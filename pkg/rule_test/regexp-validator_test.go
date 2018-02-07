package rule_test

import (
	"fmt"
	"testing"

	"github.com/ncarlier/apimon/pkg/assert"
	"github.com/ncarlier/apimon/pkg/rule"
)

var notmatch = fmt.Errorf("body does not match the RegExp")

var regExpValidationTests = []struct {
	rule     string
	body     string
	expected error
}{
	{"RegExp:hello", "hello world", nil},
	{"RegExp:(foo){2}", "foofoobar", nil},
	{"RegExp:(foo){2}", "foobarbar", notmatch},
	{"RegExp:foo", "bar", notmatch},
}

func TestRegExpValidator(t *testing.T) {
	for idx, tt := range regExpValidationTests {
		pipeline, err := rule.CreateValidatorPipeline(tt.rule)
		assert.Nil(t, err, "Pipeline creation should not fail")
		assert.NotNil(t, pipeline, "Pipeline should be created")
		assert.Equal(t, 1, len(pipeline), "Invalid validator pipeline")
		validator := pipeline[0]
		assert.Equal(t, "RegExp", validator.Name(), "Invalid validator name")
		actual := validator.Validate(200, nil, tt.body)
		if (tt.expected == nil && actual != nil) ||
			(actual == nil && tt.expected != nil) ||
			(actual != nil && tt.expected.Error() != actual.Error()) {
			t.Errorf("Dataset(%d): expected %v, actual %v", idx, tt.expected, actual)
		}
	}
}

func TestBadRegExpValidator(t *testing.T) {
	_, err := rule.CreateValidatorPipeline("RegExp:(?!re)")
	assert.NotNil(t, err, "Pipeline creation should fail")
}
