package logger

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"math"
	"regexp"
	"runtime/debug"
	"time"
)

var cfg *config

type consoleColorModeValue int

type RequestLabelMappingFn func(c *gin.Context) string

// LogFormatter gives the signature of the formatter function passed to LoggerWithFormatter
type LogFormatter func(params LogFormatterParams) string

// LogFormatterParams is the structure any formatter will be handed when time to log comes
type LogFormatterParams struct {
	// TimeStamp shows the time after the webServe returns a response.
	TimeStamp time.Time
	// StatusCode is HTTP response code.
	StatusCode int
	// Latency is how much time the webServe cost to process a certain request.
	Latency time.Duration
	// ClientIP equals Context's ClientIP method.
	ClientIP string
	// Method is the HTTP method given to the request.
	Method string
	// Path is a path the client requests.
	Path string
	// ErrorMessage is set if error has occurred in processing the request.
	ErrorMessage string
	// isTerm shows whether does gin's output descriptor refers to a terminal.
	isTerm bool
	// BodySize is the size of the Response Body
	BodySize int
	// Keys are the keys set on the request's context.
	Keys map[string]interface{}

	RequestData      string
	RequestUserAgent string
	RequestReferer   string
	RequestProto     string

	RequestId string
	TraceId   string
	SpanId    string

	ResponseData string
}

// defaultLogFormatter is the default log format function Logger middleware uses.
var defaultLogFormatter = func(param LogFormatterParams) string {
	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}
	return fmt.Sprintf("%3d | %8v | %15s | %-7s %#v %s",
		param.StatusCode,
		param.Latency,
		param.ClientIP,
		param.Method,
		param.Path,
		param.ErrorMessage,
	)
}

// NewErrorLogger returns a handler func for any error type.
func NewErrorLogger(opts ...Option) gin.HandlerFunc {
	if cfg == nil {
		cfg = &config{
			endpointLabelMappingFn: func(c *gin.Context) string {
				return c.Request.URL.Path
			}}
	}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.formatter == nil {
		cfg.formatter = defaultLogFormatter
	}

	return ErrorLoggerT(gin.ErrorTypeAny)
}

// ErrorLoggerT returns a handler func for a given error type.
func ErrorLoggerT(typ gin.ErrorType) gin.HandlerFunc {
	isTerm := true
	return func(c *gin.Context) {
		defer func() {
			if errRecover := recover(); errRecover != nil {
				if cfg.logger == nil {
					return
				}
				var recoverErr = fmt.Sprintf("%s", errRecover)
				cfg.logger.Error(string(debug.Stack()))
				start := time.Now() // Start timer
				method := c.Request.Method
				endpoint := cfg.endpointLabelMappingFn(c)
				isOk := cfg.checkLabel(fmt.Sprintf("%d", c.Writer.Status()), cfg.excludeRegexStatus) && cfg.checkLabel(endpoint, cfg.excludeRegexEndpoint) && cfg.checkLabel(method, cfg.excludeRegexMethod)
				if !isOk {
					return
				}
				rawData, err := c.GetRawData()
				if err == nil {
					c.Request.Body = io.NopCloser(bytes.NewBuffer(rawData))
				}
				raw := c.Request.URL.RawQuery
				param := LogFormatterParams{
					isTerm: isTerm,
					Keys:   c.Keys,
				}
				// Stop timer
				param.ClientIP = c.ClientIP()
				param.Method = method
				param.StatusCode = c.Writer.Status()
				param.BodySize = c.Writer.Size()
				if raw != "" {
					endpoint = endpoint + "?" + raw
				}
				param.Path = endpoint
				param.TimeStamp = time.Now()
				param.Latency = param.TimeStamp.Sub(start)
				param.ErrorMessage = recoverErr
				param.RequestProto = c.Request.Proto
				param.RequestUserAgent = c.Request.UserAgent()
				param.RequestReferer = c.Request.Referer()
				param.RequestId = c.Request.Header.Get("X-Request-Id")

				writer := &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
				c.Writer = writer

				if len(rawData) <= cfg.bodyLength {
					param.RequestData = string(rawData)
				} else {
					param.ResponseData = fmt.Sprintf("request data is too large, limit size: %d \n%s", cfg.bodyLength, string(rawData[0:cfg.bodyLength]))
				}

				if writer.body.Len() <= cfg.rawDataLength {
					param.ResponseData = writer.body.String()
				} else {
					param.ResponseData = fmt.Sprintf("response data is too large, limit size: %d \n%s", cfg.rawDataLength, string(writer.body.Bytes()[0:cfg.rawDataLength]))
				}

				cfg.logger.Debugf("%s", cfg.formatter(param))

				if cfg.writerErrorFn != nil {
					code, msg := cfg.writerErrorFn(c, &param)
					c.JSON(code, msg)
					c.Abort()
					return
				}
				c.JSON(-1, param.ErrorMessage)
				c.Abort()
			}
		}()
		c.Next()

	}
}

// New instances a Logger middleware that will write the logs to gin.DefaultWriter. By default gin.DefaultWriter = os.Stdout.
func New(opts ...Option) gin.HandlerFunc {
	if cfg == nil {
		cfg = &config{
			rawDataLength: math.MaxInt,
			bodyLength:    math.MaxInt,
			endpointLabelMappingFn: func(c *gin.Context) string {
				return c.Request.URL.Path
			}}
	}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.formatter == nil {
		cfg.formatter = defaultLogFormatter
	}

	isTerm := true
	//gin.DefaultWriter = &writeLogger{pool: buffer.Pool{}}
	return func(c *gin.Context) {
		if cfg.logger == nil {
			return
		}
		start := time.Now() // Start timer
		method := c.Request.Method
		endpoint := cfg.endpointLabelMappingFn(c)
		isOk := cfg.checkLabel(fmt.Sprintf("%d", c.Writer.Status()), cfg.excludeRegexStatus) && cfg.checkLabel(endpoint, cfg.excludeRegexEndpoint) && cfg.checkLabel(method, cfg.excludeRegexMethod)
		if !isOk {
			return
		}
		rawData, err := c.GetRawData()
		if err == nil {
			c.Request.Body = io.NopCloser(bytes.NewBuffer(rawData))
		}
		writer := &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = writer
		// Process request
		c.Next()
		raw := c.Request.URL.RawQuery
		param := LogFormatterParams{
			isTerm: isTerm,
			Keys:   c.Keys,
		}
		// Stop timer
		param.ClientIP = c.ClientIP()
		param.Method = method
		param.StatusCode = c.Writer.Status()
		param.BodySize = c.Writer.Size()
		if raw != "" {
			endpoint = endpoint + "?" + raw
		}
		param.Path = endpoint
		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

		if len(rawData) <= cfg.bodyLength {
			param.RequestData = string(rawData)
		} else {
			param.ResponseData = fmt.Sprintf("request data is too large, limit size: %d \n%s", cfg.bodyLength, string(rawData[0:cfg.bodyLength]))
		}

		if writer.body.Len() <= cfg.rawDataLength {
			param.ResponseData = writer.body.String()
		} else {
			param.ResponseData = fmt.Sprintf("response data is too large, limit size: %d \n%s", cfg.rawDataLength, string(writer.body.Bytes()[0:cfg.rawDataLength]))
		}

		cfg.logger.Debugf("%s", cfg.formatter(param))

		if cfg.writerLogFn != nil {
			param.RequestProto = c.Request.Proto
			param.RequestUserAgent = c.Request.UserAgent()
			param.RequestReferer = c.Request.Referer()
			param.RequestId = c.Request.Header.Get("X-Request-Id")
			cfg.writerLogFn(c, &param)
		}

	}
}

// checkLabel returns the match result of labels.
// Return true if regex-pattern compiles failed.
func (c *config) checkLabel(label string, patterns []string) bool {
	if len(patterns) <= 0 {
		return true
	}
	for _, pattern := range patterns {
		if pattern == "" {
			return true
		}
		matched, err := regexp.MatchString(pattern, label)
		if err != nil {
			return true
		}
		if matched {
			return false
		}
	}
	return true
}
