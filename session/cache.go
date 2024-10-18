package session

import (
	"github.com/donetkit/contrib/utils/cache"
	"github.com/gorilla/sessions"
)

type ICacheStore interface {
	SessionsStore
}

// NewStore size: maximum number of idle connections.
// network: tcp or udp
// address: host:port
// password: redis-password
// Keys are defined in pairs to allow key rotation, but the common case is to set a single
// authentication key and optionally an encryption key.
//
// The first key in a pair is used for authentication and the second for encryption. The
// encryption key can be set to nil or omitted in the last pair, but the authentication key
// is required in all pairs.
//
// It is recommended to use an authentication key with 32 or 64 bytes. The encryption key,
// if set, must be either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256 modes.
func NewStore(cache cache.ICache, keyPairs ...[]byte) (ICacheStore, error) {
	store := NewCacheStore(cache, keyPairs...)
	return &cacheStore{store}, nil
}

type cacheStore struct {
	*CacheStore
}

func (c *cacheStore) Options(options Options) {
	c.CacheStore.Options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}
