package common

import (
	"github.com/pkg/errors"
	"time"
)

const (
	SuccessDuration            = 3 * time.Hour
	SecretName                 = "klt-certs"
	CertificatesSecretEmptyErr = "certificates secret is empty"
)

var ErrCouldNotUpdateCRD = errors.New("could not update CRD config")
