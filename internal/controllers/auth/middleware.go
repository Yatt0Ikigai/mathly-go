package auth

import (
	"mathly/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtS service.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		jwt_token, err := c.Cookie("jwt")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		claims, err := jwtS.ValidateToken(jwt_token, service.Access)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
