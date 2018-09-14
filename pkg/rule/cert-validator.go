package rule

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type certValidator struct {
	name string
	spec string
	nbf  time.Duration
}

func newCertValidator(spec string) (Validator, error) {
	validator := &certValidator{
		name: "cert",
		spec: spec,
	}
	i, err := strconv.Atoi(spec)
	if err != nil {
		return nil, err
	}
	validator.nbf = time.Duration(24*i) * time.Hour
	return validator, nil
}

func (v *certValidator) Name() string {
	return v.name
}

func (v *certValidator) Spec() string {
	return v.spec
}

func (v *certValidator) Validate(body string, resp *http.Response) error {
	if resp.TLS != nil {
		cert := resp.TLS.PeerCertificates[0]
		// fmt.Println("Certificate subject:", cert.Subject.String())
		// fmt.Println("Certificate expiration date:", cert.NotAfter)
		if cert.NotAfter.Before(time.Now().Add(v.nbf)) {
			return fmt.Errorf("certificate is about to expire")
		}
	}
	return nil
}
