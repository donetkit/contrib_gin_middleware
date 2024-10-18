package logger

import (
	"bytes"
	"github.com/gin-gonic/gin"
)

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r bodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
