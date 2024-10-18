package requestid

// Option for queue system
type Option func(*config)

type Generator func() string

// WithGenerator set generator function
func WithGenerator(g Generator) Option {
	return func(cfg *config) {
		cfg.generator = g
	}
}

// WithCustomHeaderStrKey set custom header key for request id
func WithCustomHeaderStrKey(s string) Option {
	return func(cfg *config) {
		cfg.headerKey = s
	}
}
