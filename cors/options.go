package cors

import "time"

// Option for queue system
type Option func(*Config)

// WithAllowOrigins set origins function default *
func WithAllowOrigins(origins []string) Option {
	return func(cfg *Config) {
		cfg.AllowOrigins = origins
	}
}

// WithAllowMethods set methods function default "GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"
func WithAllowMethods(methods []string) Option {
	return func(cfg *Config) {
		cfg.AllowMethods = methods
	}
}

// WithAllowHeaders set headers function default "Origin", "Content-Type", "Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Cache-Control", "Content-Language"
func WithAllowHeaders(headers []string) Option {
	return func(cfg *Config) {
		cfg.AllowHeaders = headers
	}
}

// WithExposeHeaders set headers function
func WithExposeHeaders(headers []string) Option {
	return func(cfg *Config) {
		cfg.ExposeHeaders = headers
	}
}

// WithAllowCredentials set allowCredentials function default true
func WithAllowCredentials(allowCredentials bool) Option {
	return func(cfg *Config) {
		cfg.AllowCredentials = allowCredentials
	}
}

// WithAllowAllOrigins set allowAllOrigins function default false
func WithAllowAllOrigins(allowAllOrigins bool) Option {
	return func(cfg *Config) {
		cfg.AllowAllOrigins = allowAllOrigins
	}
}

// WithAllowWildcard set allowWildcard function default false
func WithAllowWildcard(allowWildcard bool) Option {
	return func(cfg *Config) {
		cfg.AllowWildcard = allowWildcard
	}
}

// WithAllowOriginFunc set origin function default true
func WithAllowOriginFunc(originFunc func(origin string) bool) Option {
	return func(cfg *Config) {
		if originFunc != nil {
			cfg.AllowOriginFunc = originFunc
		}
	}
}

// WithMaxAge set maxAge function default 12 * time.Hour
func WithMaxAge(maxAge time.Duration) Option {
	return func(cfg *Config) {
		cfg.MaxAge = maxAge
	}
}
