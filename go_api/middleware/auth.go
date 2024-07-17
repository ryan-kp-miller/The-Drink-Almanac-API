package middleware

import (
	"net/http"
	"the-drink-almanac-api/service"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authService service.AuthService
}

// AuthUser extracts a userId from the Token header and adds the userId to the request context
func (m AuthMiddleware) AuthUser(c *gin.Context) {
	tokens := c.Request.Header["Token"]
	if len(tokens) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "the 'Token' header was not included in the request"})
		c.Abort()
		return
	}
	userId, err := m.authService.ValidateToken(tokens[0])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "the 'Token' header was invalid"})
		c.Abort()
		return
	}
	c.Set("userId", userId)
	c.Next()
}

func NewAuthMiddleware(authService service.AuthService) AuthMiddleware {
	return AuthMiddleware{
		authService: authService,
	}
}
