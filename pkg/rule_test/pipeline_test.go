package rule_test

import (
	"testing"

	"github.com/ncarlier/apimon/pkg/assert"
	"github.com/ncarlier/apimon/pkg/rule"
)

func TestBadPipeline(t *testing.T) {
	def := "Code:200;Foo:Bar"
	expected := "unknown rule name: Foo"
	_, err := rule.CreateValidatorPipeline(def)
	assert.NotNil(t, err, "Pipeline creation should fail")
	assert.Equal(t, expected, err.Error(), "Unexpected error")
}

func TestEmptyPipeline(t *testing.T) {
	def := ""
	expected := "rule is empty"
	_, err := rule.CreateValidatorPipeline(def)
	assert.NotNil(t, err, "Pipeline creation should fail")
	assert.Equal(t, expected, err.Error(), "Unexpected error")
}

func TestSimplePipeline(t *testing.T) {
	def := "Code:200"
	pipeline, err := rule.CreateValidatorPipeline(def)
	assert.Nil(t, err, "Pipeline creation should not fail")
	assert.NotNil(t, pipeline, "Pipeline should  be created")
	assert.Equal(t, 1, len(pipeline), "Invalid validator pipeline")
	validator := pipeline[0]
	assert.Equal(t, def, validator.Def(), "Invalid validator definition")
}
