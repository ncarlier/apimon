package rule_test

import (
	"net/http"
	"testing"

	"github.com/ncarlier/apimon/pkg/assert"
	"github.com/ncarlier/apimon/pkg/config"
	"github.com/ncarlier/apimon/pkg/rule"
)

func TestCertValidator(t *testing.T) {
	resp, err := http.Get("https://www.nunux.org/")
	assert.Nil(t, err, "Unable to get test URL")

	rules := []config.Rule{config.Rule{Name: "cert", Spec: "1"}}
	pipeline, err := rule.CreateValidatorPipeline(rules)
	assert.Nil(t, err, "Pipeline creation should not fail")
	assert.NotNil(t, pipeline, "Pipeline should be created")
	assert.Equal(t, 1, len(pipeline), "Invalid validator pipeline")
	validator := pipeline[0]
	assert.Equal(t, "cert", validator.Name(), "Invalid validator name")
	assert.Nil(t, validator.Validate("", resp), "Certificate validation should not fail")
}

func TestBadCertValidator(t *testing.T) {
	resp, err := http.Get("https://www.nunux.org/")
	assert.Nil(t, err, "Unable to get test URL")

	rules := []config.Rule{config.Rule{Name: "cert", Spec: "300"}}
	pipeline, err := rule.CreateValidatorPipeline(rules)
	validator := pipeline[0]
	err = validator.Validate("", resp)
	assert.NotNil(t, err, "Certificate validation should fail")
	assert.Equal(t, "certificate is about to expire", err.Error(), "Bad certificate validation error")
}
