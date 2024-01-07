package common

import (
	"time"

	"github.com/pkg/errors"
)

const (
	SuccessDuration            = 3 * time.Hour
	SecretName                 = "keptn-certs"
	CertificatesSecretEmptyErr = "certificates secret is empty"
)

var ErrCouldNotUpdateCRD = errors.New("could not update CRD config")
