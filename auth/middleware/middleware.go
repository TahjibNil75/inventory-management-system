package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/inventory-management-system/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the JWT from cookie
		tokenstring, err := c.Cookie("AdminJWT")
		if tokenstring == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "request does not contain an access token",
			})
			c.Abort()
		}
		err = utils.ValidateToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
