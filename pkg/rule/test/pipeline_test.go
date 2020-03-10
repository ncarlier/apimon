package test

import (
	"testing"

	"github.com/ncarlier/apimon/pkg/assert"
	"github.com/ncarlier/apimon/pkg/config"
	"github.com/ncarlier/apimon/pkg/rule"
)

func TestBadPipeline(t *testing.T) {
	rules := []config.Rule{
		config.Rule{Name: "code", Spec: "200"},
		config.Rule{Name: "foo", Spec: "bar"},
	}
	expected := "unknown rule name: foo"
	_, err := rule.CreateValidatorPipeline(rules)
	assert.NotNil(t, err, "Pipeline creation should fail")
	assert.Equal(t, expected, err.Error(), "Unexpected error")
}

func TestDefaultPipeline(t *testing.T) {
	rules := []config.Rule{}
	pipeline, err := rule.CreateValidatorPipeline(rules)
	assert.Nil(t, err, "Pipeline creation should not fail")
	assert.NotNil(t, pipeline, "Pipeline should  be created")
	assert.Equal(t, 1, len(pipeline), "Invalid validator pipeline")
	validator := pipeline[0]
	assert.Equal(t, "code", validator.Name(), "Invalid validator name")
	assert.Equal(t, "200", validator.Spec(), "Invalid validator spec")
}

func TestSimplePipeline(t *testing.T) {
	rules := []config.Rule{
		config.Rule{Name: "code", Spec: "202"},
	}
	pipeline, err := rule.CreateValidatorPipeline(rules)
	assert.Nil(t, err, "Pipeline creation should not fail")
	assert.NotNil(t, pipeline, "Pipeline should  be created")
	assert.Equal(t, 1, len(pipeline), "Invalid validator pipeline")
	validator := pipeline[0]
	assert.Equal(t, "code", validator.Name(), "Invalid validator name")
	assert.Equal(t, "202", validator.Spec(), "Invalid validator spec")
}
