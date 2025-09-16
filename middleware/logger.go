package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger logs method, path, status, latency, user_id and request_id
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		rid, _ := c.Get("request_id")
		userID := c.GetUint("user_id")
		log.Printf("rid=%v method=%s path=%s status=%d latency=%s user_id=%d", rid, c.Request.Method, c.FullPath(), status, latency, userID)
	}
}


