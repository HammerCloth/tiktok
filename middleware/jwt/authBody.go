package jwt

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// AuthBody 鉴权中间件
// 若用户携带的token正确,解析token,将userId放入上下文context中并放行;否则,返回错误信息
func AuthBody() gin.HandlerFunc {
	return func(context *gin.Context) {
		auth := context.Request.PostFormValue("token")
		fmt.Printf("%v \n", auth)

		if len(auth) == 0 {
			context.Abort()
			context.JSON(http.StatusUnauthorized, Response{
				StatusCode: -1,
				StatusMsg:  "Unauthorized",
			})
		}
		auth = strings.Fields(auth)[1]
		token, err := parseToken(auth)
		if err != nil {
			context.Abort()
			context.JSON(http.StatusUnauthorized, Response{
				StatusCode: -1,
				StatusMsg:  "Token Error",
			})
		} else {
			println("token 正确")
		}
		context.Set("userId", token.Id)
		context.Next()
	}
}
