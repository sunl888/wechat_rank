package gerrors

import (
	"encoding/json"
)

type GlobalError struct {
	Code        int
	ServiceName string
	Message     string
	Detail      interface{}
	StatusCode  int
}

func (ge *GlobalError) Error() string {
	b, _ := json.Marshal(ge)
	return string(b)
}

func New(code int, serviceName, message string, detail interface{}, statusCode int) error {
	return &GlobalError{
		Code:        code,
		ServiceName: serviceName,
		Message:     message,
		Detail:      detail,
		StatusCode:  statusCode,
	}
}

// BadRequest generates a 400 error.
func BadRequest(code int, serviceName, message string, detail interface{}) error {
	return &GlobalError{
		Code:        code,
		ServiceName: serviceName,
		StatusCode:  400,
		Message:     message,
		Detail:      detail,
	}
}

// Unauthorized generates a 401 error.
func Unauthorized(code int, serviceName, message string, detail interface{}) error {
	return &GlobalError{
		Code:        code,
		ServiceName: serviceName,
		Message:     message,
		StatusCode:  401,
		Detail:      detail,
	}
}

// Forbidden generates a 403 error.
func Forbidden(code int, serviceName, message string, detail interface{}) error {
	return &GlobalError{
		Code:        code,
		ServiceName: serviceName,
		Message:     message,
		StatusCode:  403,
		Detail:      detail,
	}
}

// NotFound generates a 404 error.
func NotFound(code int, serviceName, message string, detail interface{}) error {
	return &GlobalError{
		Code:        code,
		ServiceName: serviceName,
		Message:     message,
		StatusCode:  404,
		Detail:      detail,
	}
}

// InternalServerError generates a 500 error.
func InternalServerError(code int, serviceName, message string, detail interface{}) error {
	return &GlobalError{
		Code:        code,
		ServiceName: serviceName,
		Message:     message,
		StatusCode:  500,
		Detail:      detail,
	}
}

// Conflict generates a 409 error.
func Conflict(code int, serviceName, message string, detail interface{}) error {
	return &GlobalError{
		Code:        code,
		ServiceName: serviceName,
		Message:     message,
		StatusCode:  409,
		Detail:      detail,
	}
}

type ValidateError map[string][]string

// UnprocessableEntity generates a 422 error.
func UnprocessableEntity(code int, serviceName string, ve ValidateError) error {
	return &GlobalError{
		Code:        code,
		ServiceName: serviceName,
		Message:     "The given data failed to pass validation.",
		StatusCode:  422,
		Detail:      ve,
	}
}
