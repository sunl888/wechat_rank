package service

import (
	"code.aliyun.com/zmdev/wechat_rank/errors"
	"code.aliyun.com/zmdev/wechat_rank/model"
	"code.aliyun.com/zmdev/wechat_rank/pkg/hasher"
	"github.com/gin-gonic/gin"
)

type userService struct {
	us   model.UserStore
	cs   model.CertificateStore
	tSvc model.TicketService
	h    hasher.Hasher
}

func (uSvc *userService) UserLogin(account, password string) (ticket *model.Ticket, err error) {
	c, err := uSvc.cs.CertificateLoadByAccount(account)
	if err != nil {
		if uSvc.cs.CertificateIsNotExistErr(err) { //账号不存在
			err = errors.ErrAccountNotFound()
		}
		return nil, err
	}
	user, err := uSvc.us.UserLoad(c.UserId)
	if err != nil {
		return nil, err
	}
	if uSvc.h.Check(password, user.Password) {
		// 登录成功
		return uSvc.tSvc.TicketGen(user.Id)
	}

	return nil, errors.ErrPassword()
}

func (uSvc *userService) UserRegister(account string, certificateType model.CertificateType, password string) (userId int64, err error) {
	if exist, err := uSvc.cs.CertificateExist(account); err != nil {
		return 0, err
	} else if exist {
		return 0, errors.ErrAccountAlreadyExisted()
	}
	user := &model.User{Password: uSvc.h.Make(password), PwPlain: password}
	if err := uSvc.us.UserCreate(user); err != nil {
		return 0, err
	}
	certificate := &model.Certificate{UserId: user.Id, Account: account, Type: certificateType}
	if err := uSvc.cs.CertificateCreate(certificate); err != nil {
		return 0, err
	}
	return user.Id, nil
}

func (uSvc *userService) UserUpdatePassword(userId int64, newPassword string) error {
	return uSvc.us.UserUpdate(&model.User{
		Id:       userId,
		Password: uSvc.h.Make(newPassword),
		PwPlain:  newPassword,
	})
}

func UserLogin(ctx *gin.Context, account, password string) (*model.Ticket, error) {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.UserLogin(account, password)
	}
	return nil, ServiceError
}

func UserRegister(ctx *gin.Context, account string, certificateType model.CertificateType, password string) (userId int64, err error) {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.UserRegister(account, certificateType, password)
	}
	return 0, ServiceError
}

func UserUpdatePassword(ctx *gin.Context, userId int64, newPassword string) error {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.UserUpdatePassword(userId, newPassword)
	}
	return nil
}

func NewUserService(us model.UserStore, cs model.CertificateStore, tSvc model.TicketService, h hasher.Hasher) model.UserService {
	return &userService{us: us, cs: cs, tSvc: tSvc, h: h}
}
