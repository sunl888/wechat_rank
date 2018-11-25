package errors

import "github.com/zm-dev/gerrors"

const (
	ServiceName = "default"
)

// 参数绑定出错
func BindError(err error) error {
	return gerrors.BadRequest(10000, ServiceName, err.Error(), err)
}

func InternalServerError(msg string, err error) error {
	return gerrors.InternalServerError(10001, ServiceName, msg, err)
}

func BadRequest(msg string, err error) error {
	return gerrors.BadRequest(10002, ServiceName, msg, err)
}

func Unauthorized() error {
	return gerrors.Unauthorized(10003, ServiceName, "请先登录", nil)
}
