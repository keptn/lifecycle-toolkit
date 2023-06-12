package certificates

import (
	"crypto/x509"
	"encoding/pem"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/certificatehandler_mock.go . ICertificateHandler
type ICertificateHandler interface {
	Decode(data []byte) (p *pem.Block, rest []byte)
	Parse(der []byte) (*x509.Certificate, error)
}

type defaultCertificateHandler struct {
}

func (c defaultCertificateHandler) Decode(data []byte) (p *pem.Block, rest []byte) {
	return pem.Decode(data)
}
func (c defaultCertificateHandler) Parse(der []byte) (*x509.Certificate, error) {
	return x509.ParseCertificate(der)
}
