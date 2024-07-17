package server

import "github.com/gin-gonic/gin"

// setUserIdInContext mocks the output of the auth middleware
// by setting the user id as a context variable for the handlers to use
func setUserIdInContext(testUserId string) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Set("userId", testUserId)
		c.Next()
	}
}
