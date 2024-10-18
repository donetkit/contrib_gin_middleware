package ip_white

import (
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"strings"
)

func New(opts ...Option) gin.HandlerFunc {
	cfg := &option{}
	for _, opt := range opts {
		opt(cfg)
	}
	return func(c *gin.Context) {
		if !isIPWhite(c.ClientIP(), cfg.WhiteList) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
	}
}

func isIPWhite(ip string, whitelist []string) bool {
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return false
	}

	for _, allowedIP := range whitelist {
		if strings.Contains(allowedIP, "/") {
			_, ipNet, err := net.ParseCIDR(allowedIP)
			if err != nil {
				continue
			}
			if ipNet.Contains(ipAddr) {
				return true
			}
		} else {
			allowedIP = strings.TrimSpace(allowedIP)
			if allowedIP == ip {
				return true
			}
		}
	}

	return false
}
