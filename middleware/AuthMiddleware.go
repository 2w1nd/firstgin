package middleware

import (
	"com.w1nd/firstgin/common"
	"com.w1nd/firstgin/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
	//	获取authorization header
		tokenString := c.GetHeader("Authorization")

	//	validate token formate  不是bearer开头则不是一个正确的token
	//  oauth2.0规定的,Authorization的字符串开头必须要有Bearer
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer "){
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg": "权限不足",
			})
			c.Abort()
		}
		tokenString = tokenString[7:]

		token, claim, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg": "权限不足",
			})
			c.Abort()
			return
		}

	//	验证通过后获取claim中的userid
		userId := claim.UserId
		DB := common.GetDB()
		var user model.User
		DB.First(&user, userId)

	//	用户
		if user.ID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg": "权限不足",
			})
			c.Abort()
			return
		}

	//	用户存在
		c.Set("user",user)
		c.Next()
	}
}
