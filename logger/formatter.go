package logger

import (
	"github.com/donetkit/contrib-log/glog"
	"github.com/donetkit/contrib/utils/buffer"
	"github.com/gin-gonic/gin"
	"strings"
)

func SetGinDefaultWriter(logger glog.ILogger) {
	gin.DefaultWriter = &writeLogger{pool: buffer.Pool{}, logger: logger.WithField("Gin-Logger", "Gin-Logger")}
}

type writeLogger struct {
	logger glog.ILoggerEntry
	pool   buffer.Pool
}

// Write implements io.Writer.
func (l *writeLogger) Write(p []byte) (n int, err error) {
	buf := l.pool.Get()
	defer l.pool.Put(buf)
	n, err = buf.Write(p)
	if l.logger != nil {
		msg := buf.String()
		if strings.HasSuffix(msg, "\n") {
			msg = msg[:len(msg)-1]
		}
		l.logger.Debug(msg)
	}
	return n, err
}
