package errors

import (
	"fmt"
	"net/http"
)

type ErrorInfo struct {
	Code    uint32 `json:"err_code"`
	Message string `json:"err_message"`
}

func NewErrorInfo(major, minor int, message string) ErrorInfo {
	return ErrorInfo{Code: uint32(major)*10000 + uint32(minor), Message: message}
}

func (err ErrorInfo) Error() string {
	return fmt.Sprintf("code: %d, message: %s", err.Code, err.Message)
}

func (err ErrorInfo) StatusCode() int {
	return int(err.Code / 10000)
}

var (
	ErrInvalidParameter = func(msg string) ErrorInfo { return NewErrorInfo(http.StatusBadRequest, 100, msg) }
	ErrInternal         = func(msg string) ErrorInfo {
		return NewErrorInfo(http.StatusInternalServerError, 100, "sorry, we made a mistake")
	}
	ErrUnimplemented = func() ErrorInfo {
		return NewErrorInfo(http.StatusNotImplemented, 100, "unimplemented")
	}
	ErrUnauthorized = func(msg string) ErrorInfo { return NewErrorInfo(http.StatusUnauthorized, 100, msg) }
	ErrTokenExpired = func() ErrorInfo { return NewErrorInfo(http.StatusUnauthorized, 101, "token失效") }
	ErrForbidden    = func(msg string) ErrorInfo { return NewErrorInfo(http.StatusForbidden, 100, msg) }
	// biz error
	ErrEmpNotFound     = func() ErrorInfo { return NewErrorInfo(http.StatusOK, 101, "工号不存在") }
	ErrPasswordInvalid = func() ErrorInfo { return NewErrorInfo(http.StatusOK, 102, "密码不正确") }
	ErrNotBind         = func() ErrorInfo { return NewErrorInfo(http.StatusOK, 103, "还未绑定工号") }
	ErrAlreadyBind     = func() ErrorInfo { return NewErrorInfo(http.StatusOK, 104, "已经被绑定工号") }
)
