package common

import (
	"errors"

	"Go-AIServiceSupport/common/e"
)

type AppError struct {
	Code	   int
	Message    string
	HTTPStatus int
	Cause      error
}

func (err *AppError) Error() string {
	if err.Cause != nil {
		return  err.Cause.Error()
	}

	return err.Message
}

// 根据错误码新建一个AppError
func NewAppError(code int) *AppError {
	return &AppError{
		Code: 		code,
		Message:    e.Message(code),
		HTTPStatus: e.HTTPStatus(code),
	}
}

// 根据错误码新建一个AppError，但是可自定义错误Message
func NewAppErrorWithMessage(code int, message string) *AppError {
	return &AppError{
		Code: 		code,
		Message:    e.Message(code),
		HTTPStatus: e.HTTPStatus(code),
	}
}

// 取出底层错误原因Cause
func (err *AppError) Unwrap() error {
	return err.Cause
}

// 根据底层错误原因新建一个AppError
func WrapAppError(code int, cause error) *AppError {
	return &AppError{
		Code: 		code,
		Message:    e.Message(code),
		HTTPStatus: e.HTTPStatus(code),
	}
}

// 从error中取出AppError（控制层使用）
func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}