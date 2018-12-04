package middleware

import (
	"code.aliyun.com/zmdev/wechat_rank/errors"
	"code.aliyun.com/zmdev/wechat_rank/service"
	"github.com/gin-gonic/gin"
)

var (
	isLoginKey = "is_login"
	userIdKey  = "user_id"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isLogin := check(c)
		if !isLogin {
			_ = c.Error(errors.Unauthorized())
			c.Abort()
			return
		}
		c.Next()
	}
}

func check(c *gin.Context) bool {
	var (
		isLogin bool
	)
	if ticketId, err := c.Cookie("ticket_id"); err == nil {
		isValid, userId, err := service.TicketIsValid(c, ticketId)
		if err != nil {
			isLogin = false
		} else {
			isLogin = isValid
			c.Set(isLoginKey, isLogin)
			c.Set(userIdKey, userId)
		}
	} else {
		// cookie不存在
		isLogin = false
	}
	return isLogin
}

func CheckLogin(c *gin.Context) bool {
	isLogin, ok := c.Get(isLoginKey)
	if !ok {
		isLogin = check(c)
	}
	return isLogin.(bool)

}

func UserId(c *gin.Context) int64 {
	userId, ok := c.Get(userIdKey)
	if !ok {
		check(c)
		userId = c.GetInt64(userIdKey)
	}
	return userId.(int64)
}
