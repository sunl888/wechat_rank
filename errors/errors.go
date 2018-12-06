package errors

import "github.com/zm-dev/gerrors"

// 参数绑定出错
func BindError(err error) error {
	return gerrors.BadRequest(10000, err.Error(), err)
}

// 清博大数据Api error
func QingboError(name, msg string, code, status int) error {
	return gerrors.New(code, msg, name, status)
}

func InternalServerError(msg string, err error) error {
	return gerrors.InternalServerError(10001, msg, err)
}

func BadRequest(msg string, err error) error {
	return gerrors.BadRequest(10002, msg, err)
}

func Unauthorized() error {
	return gerrors.Unauthorized(10003, "请先登录", nil)
}

func ErrAccountAlreadyExisted() error {
	return gerrors.BadRequest(20001, "account already existed", nil)
}

func ErrPassword() error {
	return gerrors.BadRequest(20002, "error password", nil)
}

func ErrAccountNotFound() error {
	return gerrors.NotFound(20003, "account not found", nil)
}
