package favicon

import (
	"github.com/donetkit/contrib-log/glog"
)

type option struct {
	logger      glog.ILoggerEntry
	faviconPath string
	routePaths  []string
}

type Option func(*option)

func WithLogger(logger glog.ILogger) Option {
	return func(o *option) {
		o.logger = logger.WithField("Favicon", "Favicon")
	}
}

// WithFaviconPath  faviconPath default ./favicon.ico
func WithFaviconPath(faviconPath string) Option {
	return func(o *option) {
		o.faviconPath = faviconPath
	}
}

// WithRoutePaths  routePaths default /favicon.ico
func WithRoutePaths(routePaths ...string) Option {
	return func(o *option) {
		for _, path := range routePaths {
			if path != "" {
				o.routePaths = append(o.routePaths, path)
			}
		}

	}
}
