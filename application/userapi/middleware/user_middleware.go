package middleware

import (
	"awesomeProject/application/userapi/utils"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//验证TOKEN中间件
func ValidTokenMiddleware(c *gin.Context) {
	tokenStr := c.GetHeader("Auth")
	err := validToken(tokenStr)
	if err != nil {
		c.JSON(429, &gin.H{"message": err})
		utils.ResponseError(c,1009)
		c.Abort()
	}
	c.Next()
}

func validToken(tokenStr string) error {
	_, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte("jey.sign"), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				fmt.Println("Timing is everything")
				return err
			} else {
				fmt.Println("token invalid:", err)
				return err
			}
		}
	}
	return nil
}