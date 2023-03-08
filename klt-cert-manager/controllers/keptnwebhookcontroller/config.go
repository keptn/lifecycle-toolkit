package keptnwebhookcontroller

import (
	"time"
)

const (
	successDuration            = 3 * time.Hour
	secretName                 = "klt-certs"
	certificatesSecretEmptyErr = "certificates secret is empty"
	couldNotUpdateCRDErr       = "could not update crd config"
)
