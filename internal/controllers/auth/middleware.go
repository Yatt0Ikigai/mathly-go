package auth

import (
	"mathly/internal/log"
	"mathly/internal/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtS service.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwtS.ValidateToken(tokenStr, service.Access)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.SetCookie(
			"access_token", // cookie name
			tokenStr,       // cookie value
			900,            // maxAge in seconds (15 * 60)
			"/",            // path
			"",             // domain (empty = current domain)
			true,           // secure (true in prod with HTTPS, false if localhost HTTP)
			true,           // httpOnly
		)
		log.Log.Infof("%s claims UserID", claims.UserID)
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
