package gcors

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type gCors struct {
	allowAllOrigins            bool
	allowCredentials           bool
	allowOriginFunc            func(string) bool
	allowOriginWithContextFunc func(*gin.Context, string) bool
	allowOrigins               []string
	normalHeaders              http.Header
	preflightHeaders           http.Header
	wildcardOrigins            [][]string
	optionsResponseStatusCode  int
}

var (
	DefaultSchemas = []string{
		"http://",
		"https://",
	}
	ExtensionSchemas = []string{
		"chrome-extension://",
		"safari-extension://",
		"moz-extension://",
		"ms-browser-extension://",
	}
	FileSchemas = []string{
		"file://",
	}
	WebSocketSchemas = []string{
		"ws://",
		"wss://",
	}
)

func newCors(config Config) *gCors {
	if err := config.Validate(); err != nil {
		panic(err.Error())
	}

	for _, origin := range config.AllowOrigins {
		if origin == "*" {
			config.AllowAllOrigins = true
		}
	}

	if config.OptionsResponseStatusCode == 0 {
		config.OptionsResponseStatusCode = http.StatusNoContent
	}

	return &gCors{
		allowOriginFunc:            config.AllowOriginFunc,
		allowOriginWithContextFunc: config.AllowOriginWithContextFunc,
		allowAllOrigins:            config.AllowAllOrigins,
		allowCredentials:           config.AllowCredentials,
		allowOrigins:               normalize(config.AllowOrigins),
		normalHeaders:              generateNormalHeaders(config),
		preflightHeaders:           generatePreflightHeaders(config),
		wildcardOrigins:            config.parseWildcardRules(),
		optionsResponseStatusCode:  config.OptionsResponseStatusCode,
	}
}

func (gCors *gCors) applyCors(c *gin.Context) {
	origin := c.Request.Header.Get("Origin")
	if len(origin) == 0 {
		// request is not a CORS request
		return
	}
	host := c.Request.Host

	if origin == "http://"+host || origin == "https://"+host {
		// request is not a CORS request but have origin header.
		// for example, use fetch api
		return
	}

	if !gCors.isOriginValid(c, origin) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if c.Request.Method == "OPTIONS" {
		gCors.handlePreflight(c)
		defer c.AbortWithStatus(gCors.optionsResponseStatusCode)
	} else {
		gCors.handleNormal(c)
	}

	if !gCors.allowAllOrigins {
		c.Header("Access-Control-Allow-Origin", origin)
	}
}

func (gCors *gCors) validateWildcardOrigin(origin string) bool {
	for _, w := range gCors.wildcardOrigins {
		if w[0] == "*" && strings.HasSuffix(origin, w[1]) {
			return true
		}
		if w[1] == "*" && strings.HasPrefix(origin, w[0]) {
			return true
		}
		if strings.HasPrefix(origin, w[0]) && strings.HasSuffix(origin, w[1]) {
			return true
		}
	}

	return false
}

func (gCors *gCors) isOriginValid(c *gin.Context, origin string) bool {
	valid := gCors.validateOrigin(origin)
	if !valid && gCors.allowOriginWithContextFunc != nil {
		valid = gCors.allowOriginWithContextFunc(c, origin)
	}
	return valid
}

func (gCors *gCors) validateOrigin(origin string) bool {
	if gCors.allowAllOrigins {
		return true
	}
	for _, value := range gCors.allowOrigins {
		if value == origin {
			return true
		}
	}
	if len(gCors.wildcardOrigins) > 0 && gCors.validateWildcardOrigin(origin) {
		return true
	}
	if gCors.allowOriginFunc != nil {
		return gCors.allowOriginFunc(origin)
	}
	return false
}

func (gCors *gCors) handlePreflight(c *gin.Context) {
	header := c.Writer.Header()
	for key, value := range gCors.preflightHeaders {
		header[key] = value
	}
}

func (gCors *gCors) handleNormal(c *gin.Context) {
	header := c.Writer.Header()
	for key, value := range gCors.normalHeaders {
		header[key] = value
	}
}
