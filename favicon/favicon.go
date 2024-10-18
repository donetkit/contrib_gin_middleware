package favicon

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func New(opts ...Option) gin.HandlerFunc {
	cfg := &option{
		faviconPath: "./favicon.ico",
		routePaths:  []string{"/favicon.ico"},
	}
	for _, opt := range opts {
		opt(cfg)
	}
	path := filepath.FromSlash(cfg.faviconPath)
	if len(path) > 0 && !os.IsPathSeparator(path[0]) {
		wd, err := os.Getwd()
		if err == nil {
			path = filepath.Join(wd, path)
		} else {
			if cfg.logger != nil {
				cfg.logger.Error(err)
			}
		}
	}
	info, err := os.Stat(path)
	if err != nil {
		if cfg.logger != nil {
			cfg.logger.Error(err)
		}
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		if cfg.logger != nil {
			cfg.logger.Error(err)
		}
	}
	return func(c *gin.Context) {
		var notExist = true
		for _, url := range cfg.routePaths {
			if c.Request.RequestURI == url {
				notExist = false
				break
			}
		}
		if notExist {
			return
		}
		if c.Request.Method != "GET" && c.Request.Method != "HEAD" {
			status := http.StatusOK
			if c.Request.Method != "OPTIONS" {
				status = http.StatusMethodNotAllowed
			}
			c.Header("Allow", "GET,HEAD,OPTIONS")
			c.AbortWithStatus(status)
			return
		}
		if info == nil || info.IsDir() || file == nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.Header("Content-Type", "image/x-icon")
		http.ServeContent(c.Writer, c.Request, "favicon.ico", info.ModTime(), bytes.NewReader(file))
	}
}
