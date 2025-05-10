package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func Authenticated(jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Authorization header is required",
			})
			return
		}

		bearerToken := token[len("Bearer "):]
		if bearerToken == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Bearer token is required",
			})
			return
		}

		jwtToken, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
			}
			return jwtSecret, nil
		})
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Invalid token",
			})
			return
		}
		if !jwtToken.Valid {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Token is not valid",
			})
			return
		}

		username := jwtToken.Claims.(jwt.MapClaims)["username"]
		if username == nil {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Username is required",
			})
			return
		}

		c.Set("username", jwtToken.Claims.(jwt.MapClaims)["username"])

		c.Next()
	}
}

func Authorize(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// implement proper authorization logic here
		// for the MVP just compare the username
		username := c.MustGet("username").(string)
		if username == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Username is required",
			})
			return
		}

		for _, allowedRole := range allowedRoles {
			if username == allowedRole {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(403, gin.H{
			"error": "You are not authorized to access this resource",
		})
		return
	}
}
