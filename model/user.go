package model

import "errors"

type User struct {
	Id       int64
	Password string
	PwPlain  string // 密码明文
}

type UserStore interface {
	UserExist(int64) (bool, error)
	UserLoad(int64) (*User, error)
	UserIsNotExistErr(error) bool
	UserUpdate(*User) error
	UserCreate(*User) error
}

type UserService interface {
	UserLogin(account, password string) (*Ticket, error)
	UserRegister(account string, certificateType CertificateType, password string) (userID int64, err error)
	UserUpdatePassword(userId int64, newPassword string) error
}

var ErrUserNotExist = errors.New("user not exist")

func UserIsNotExistErr(err error) bool {
	return err == ErrUserNotExist
}
