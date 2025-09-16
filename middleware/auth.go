package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "devsecret"
	}

	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization"})
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			// enforce HS256
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok || t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, jwt.ErrTokenUnverifiable
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}
		// validate exp if present
		if expVal, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(expVal) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
				return
			}
		}
		// user_id usually float64 in MapClaims
		uidFloat, _ := claims["user_id"].(float64)
		role, _ := claims["role"].(string)
		c.Set("user_id", uint(uidFloat))
		c.Set("role", role)
		c.Next()
	}
}

// Authorize enforces that the authenticated user has one of the allowed roles.
// Usage: router.Group("/admin", AuthMiddleware(), Authorize("admin"))
func Authorize(allowedRoles ...string) gin.HandlerFunc {
	roleAllowed := func(role string) bool {
		if len(allowedRoles) == 0 {
			return true
		}
		for _, r := range allowedRoles {
			if r == role {
				return true
			}
		}
		return false
	}
	return func(c *gin.Context) {
		role := c.GetString("role")
		if !roleAllowed(role) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}
