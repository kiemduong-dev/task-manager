package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

const headerRequestID = "X-Request-ID"

func generateRequestID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

// RequestID sets a unique request id into context and response header
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader(headerRequestID)
		if rid == "" {
			rid = generateRequestID()
		}
		c.Set("request_id", rid)
		c.Writer.Header().Set(headerRequestID, rid)
		c.Next()
	}
}
