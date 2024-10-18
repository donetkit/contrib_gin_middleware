package httpsign

import (
	"github.com/donetkit/contrib_gin_middleware/httpsign/crypto"
)

// KeyID define type
type KeyID string

// Secret define secret key and algorithm that key use
type Secret struct {
	Key       string
	Algorithm crypto.Crypto
}

// Secrets map with keyID and secret
type Secrets map[KeyID]*Secret
