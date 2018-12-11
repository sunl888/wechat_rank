package handler

import (
	"code.aliyun.com/zmdev/wechat_rank/errors"
	"code.aliyun.com/zmdev/wechat_rank/model"
	"code.aliyun.com/zmdev/wechat_rank/service"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

type Auth struct {
}

func (a *Auth) Login(c *gin.Context) {
	type Req struct {
		Account  string `form:"account"`
		Password string `form:"password"`
	}
	req := &Req{}
	if err := c.ShouldBind(req); err != nil {
		_ = c.Error(errors.BindError(err))
		return
	}
	resp, err := service.UserLogin(c, strings.TrimSpace(req.Account), strings.TrimSpace(req.Password))
	if err != nil {
		_ = c.Error(err)
		return
	}
	setAuthCookie(c, resp.Id, resp.UserId, int(resp.ExpiredAt.Sub(time.Now()).Seconds()))
	c.JSON(204, resp)
}

func setAuthCookie(c *gin.Context, ticketId string, userId int64, maxAge int) {
	c.SetCookie("ticket_id", ticketId, maxAge, "", "", false, true)
	c.SetCookie("user_id", strconv.FormatInt(userId, 10), maxAge, "", "", false, false)
}

func removeAuthCookie(c *gin.Context) {
	c.SetCookie("ticket_id", "", -1, "", "", false, true)
	c.SetCookie("user_id", "", -1, "", "", false, false)
}

func (a *Auth) Logout(c *gin.Context) {
	ticketId, err := c.Cookie("ticket_id")
	if err != nil {
		c.JSON(204, nil)
		return
	}
	removeAuthCookie(c)
	err = service.TicketDestroy(c, ticketId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(204, nil)
}

func (a *Auth) Register(c *gin.Context) {
	type Req struct {
		Account  string `form:"account"`
		Password string `form:"password"`
	}
	req := &Req{}

	if err := c.ShouldBind(req); err != nil {
		_ = c.Error(err)
		return
	}
	regResp, err := service.UserRegister(c, strings.TrimSpace(req.Account), model.CertificateType(0), req.Password)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(201, regResp)
}

func NewAuth() *Auth {
	return &Auth{}
}
