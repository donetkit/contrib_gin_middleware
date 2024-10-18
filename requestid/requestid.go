package requestid

import (
	"github.com/donetkit/contrib/utils/uuid"
	"github.com/gin-gonic/gin"
)

var headerXRequestID string

var headerRequestIdKey = "X-Request-Id"

// Config defines the config for RequestID middleware
type config struct {
	// Generator defines a function to generate an ID.
	// Optional. Default: func() string {
	//   return uuid.New().String()
	// }
	generator Generator
	headerKey string
}

// New initializes the RequestID middleware.
func New(opts ...Option) gin.HandlerFunc {
	cfg := &config{
		generator: func() string {
			return uuid.NewUUID()
		},
		headerKey: headerRequestIdKey,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	headerXRequestID = cfg.headerKey
	return func(c *gin.Context) {
		// Get id from request
		rid := c.GetHeader(cfg.headerKey)
		if rid == "" {
			rid = cfg.generator()
		}
		// Set the id to ensure that the requestid is in the response
		c.Header(cfg.headerKey, rid)
		c.Next()
	}
}

// Get returns the request identifier
func Get(c *gin.Context) string {
	return c.Writer.Header().Get(headerXRequestID)
}
