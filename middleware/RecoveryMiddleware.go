package middleware

import (
	"com.w1nd/firstgin/response"
	"fmt"
	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				response.Fail(c, fmt.Sprint(err), nil)
			}
		}()
		c.Next()
	}
}
