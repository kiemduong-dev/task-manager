package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var visitors = make(map[string]*rate.Limiter)
var mu sync.Mutex

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	lim, exists := visitors[ip]
	if !exists {
		lim = rate.NewLimiter(1, 5) // 1 req/sec, burst 5
		visitors[ip] = lim
	}
	return lim
}

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		lim := getLimiter(ip)
		if !lim.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error":"rate limit"})
			return
		}
		c.Next()
	}
}
