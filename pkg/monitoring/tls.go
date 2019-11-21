package monitoring

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/ncarlier/apimon/pkg/config"
)

// NewTLSConfig create client TLS configuration
func NewTLSConfig(conf config.Monitor) (*tls.Config, error) {
	tlsConfig := tls.Config{
		InsecureSkipVerify: conf.TLS.Unsafe,
	}

	// Load CA cert
	if conf.TLS.CACertFile != "" {
		caCert, err := ioutil.ReadFile(conf.TLS.CACertFile)
		if err != nil {
			return &tlsConfig, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
	}

	// Load client cert
	if conf.TLS.ClientCertFile != "" {
		cert, err := tls.LoadX509KeyPair(conf.TLS.ClientCertFile, conf.TLS.ClientKeyFile)
		if err != nil {
			return &tlsConfig, err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
		tlsConfig.BuildNameToCertificate()
	}

	return &tlsConfig, nil
}
