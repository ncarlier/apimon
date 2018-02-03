package rule_test

import (
	"fmt"
	"testing"

	"github.com/ncarlier/apimon/pkg/assert"
	"github.com/ncarlier/apimon/pkg/rule"
)

var validationTests = []struct {
	rule     string
	code     int
	expected error
}{
	{"Code:200", 200, nil},
	{"Code:200,201", 201, nil},
	{"Code:200,201,204", 204, nil},
	{"Code:200-204", 200, nil},
	{"Code:200-204", 202, nil},
	{"Code:200,204", 300, fmt.Errorf("Unexpected status code: %d", 300)},
	{"Code:200-204", 300, fmt.Errorf("Status out of range: %d", 300)},
}

func TestCodeValidator(t *testing.T) {
	for idx, tt := range validationTests {
		pipeline, err := rule.CreateValidatorPipeline(tt.rule)
		assert.Nil(t, err, "Pipeline creation should not fail")
		assert.NotNil(t, pipeline, "Pipeline should be created")
		assert.Equal(t, 1, len(pipeline), "Invalid validator pipeline")
		validator := pipeline[0]
		assert.Equal(t, "Code", validator.Name(), "Invalid validator name")
		actual := validator.Validate(tt.code, nil, "")
		if (tt.expected == nil && actual != nil) ||
			(actual == nil && tt.expected != nil) ||
			(actual != nil && tt.expected.Error() != actual.Error()) {
			t.Errorf("Dataset(%d): expected %v, actual %v", idx, tt.expected, actual)
		}
	}
}
