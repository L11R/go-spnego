package spnego

import (
	"encoding/base64"
	"net/http"

	"github.com/alexbrainman/sspi/negotiate"
)

// SSPI implements spnego.Provider interface on Windows OS
type sspi struct{}

// New constructs OS specific implementation of spnego.Provider interface
func New() Provider {
	return &sspi{}
}

// SetSPNEGOHeader puts the SPNEGO authorization header on HTTP request object
func (s *sspi) SetSPNEGOHeader(req *http.Request) error {
	h, err := canonicalizeHostname(req.URL.Hostname())
	if err != nil {
		return err
	}

	header, err := s.GetSPNEGOHeader(h)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", header)
	return nil
}

// GetSPNEGOHeader returns the SPNEGO authorization header
func (s *sspi) GetSPNEGOHeader(hostname string) (string, error) {
	spn := "HTTP/" + hostname

	cred, err := negotiate.AcquireCurrentUserCredentials()
	if err != nil {
		return "", err
	}
	defer cred.Release()

	secctx, token, err := negotiate.NewClientContext(cred, spn)
	if err != nil {
		return "", err
	}
	defer secctx.Release()

	return "Negotiate " + base64.StdEncoding.EncodeToString(token), nil
}
