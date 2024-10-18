package recovery

import (
	"fmt"
	"github.com/donetkit/contrib-log/glog"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// New returns a gin.HandlerFunc (middleware)
func New(logger glog.ILogger, stack ...bool) gin.HandlerFunc {
	logs := logger.WithField("Gin-Recover", "Gin-Recover")
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logs.Error(fmt.Sprintf("path: %s error: %s request: %s", c.Request.URL.Path, err, string(httpRequest)))
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint:
					c.Abort()
					return
				}
				if len(stack) > 0 && stack[0] {
					logs.Error("[Recovery from panic] %s error: %s request: %s stack: %s", time.Now().Format(time.RFC3339), err, string(httpRequest), string(debug.Stack()))

				} else {
					logs.Error("[Recovery from panic] %s error: %s request: %s", time.Now().Format(time.RFC3339), err, string(httpRequest))
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

func NewRecovery(handlers ...func(any any)) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				for _, f := range handlers {
					f(r)
				}
				c.Abort()
				return
			}
		}()
		c.Next()
	}
}
