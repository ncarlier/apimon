package monitoring_test

import (
	"crypto/tls"
	"testing"

	"github.com/ncarlier/apimon/pkg/assert"
	"github.com/ncarlier/apimon/pkg/config"
	"github.com/ncarlier/apimon/pkg/monitoring"
)

func TestTLSConfiguration(t *testing.T) {
	conf := config.Monitor{
		URL: "https://www.google.com",
		TLS: config.TLS{
			Unsafe:         true,
			ClientCertFile: "./cert.pem",
			ClientKeyFile:  "./key.pem",
			CACertFile:     "./ca.pem",
		},
	}

	tlsConfig, err := monitoring.NewTLSConfig(conf)
	assert.Nil(t, err, "should not fail")
	assert.NotNil(t, tlsConfig, "Configuration should be created")
	assert.Equal(t, true, tlsConfig.InsecureSkipVerify, "")
	assert.True(t, tlsConfig.ClientAuth == tls.NoClientCert, "")
}
