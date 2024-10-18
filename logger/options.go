package logger

import (
	"github.com/donetkit/contrib-log/glog"
	"github.com/gin-gonic/gin"
)

// Config defines the config for logger middleware
type config struct {
	// Optional. Default value is gin.defaultLogFormatter
	formatter              LogFormatter
	logger                 glog.ILoggerEntry
	excludeRegexStatus     []string
	excludeRegexEndpoint   []string
	excludeRegexMethod     []string
	endpointLabelMappingFn RequestLabelMappingFn
	writerLogFn            WriterLogFn
	writerErrorFn          WriterErrorFn
	bodyLength             int
	rawDataLength          int
}

// Option for queue system
type Option func(*config)

type WriterLogFn func(c *gin.Context, log *LogFormatterParams)

type WriterErrorFn func(c *gin.Context, log *LogFormatterParams) (int, interface{})

// WithLogger set logger function
func WithLogger(logger glog.ILogger) Option {
	return func(cfg *config) {
		cfg.logger = logger.WithField("Gin-Logger", "Gin-Logger")
	}
}

// WithExcludeRegexMethod set excludeRegexMethod function regexp
func WithExcludeRegexMethod(excludeRegexMethod []string) Option {
	return func(cfg *config) {
		cfg.excludeRegexMethod = excludeRegexMethod
	}
}

// WithExcludeRegexStatus set excludeRegexStatus function regexp
func WithExcludeRegexStatus(excludeRegexStatus []string) Option {
	return func(cfg *config) {
		cfg.excludeRegexStatus = excludeRegexStatus
	}
}

// WithExcludeRegexEndpoint set excludeRegexEndpoint function regexp
func WithExcludeRegexEndpoint(excludeRegexEndpoint []string) Option {
	return func(cfg *config) {
		cfg.excludeRegexEndpoint = excludeRegexEndpoint
	}
}

// WithEndpointLabelMappingFn set endpointLabelMappingFn function
func WithEndpointLabelMappingFn(endpointLabelMappingFn RequestLabelMappingFn) Option {
	return func(cfg *config) {
		cfg.endpointLabelMappingFn = endpointLabelMappingFn
	}
}

// WithFormatter set formatter function
func WithFormatter(formatter LogFormatter) Option {
	return func(cfg *config) {
		cfg.formatter = formatter
	}
}

// WithWriterLogFn set fn WriterLogFn
func WithWriterLogFn(fn WriterLogFn) Option {
	return func(cfg *config) {
		cfg.writerLogFn = fn
	}
}

// WithWriterErrorFn set fn WriterErrorFn
func WithWriterErrorFn(fn WriterErrorFn) Option {
	return func(cfg *config) {
		cfg.writerErrorFn = fn
	}
}

// WithBodyLength set fn bodyLength
func WithBodyLength(bodyLength int) Option {
	return func(cfg *config) {
		cfg.bodyLength = bodyLength
	}
}

// WithRawDataLength set fn rawDataLength
func WithRawDataLength(rawDataLength int) Option {
	return func(cfg *config) {
		cfg.rawDataLength = rawDataLength
	}
}
