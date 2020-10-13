package errors

import (
	"fmt"
	"net/http"
)

type ErrorInfo struct {
	Code    uint32 `json:"code"` // 0 成功，非0 失败
	Message string `json:"message"`
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
	ErrFaceNotFound     = func() ErrorInfo { return NewErrorInfo(http.StatusOK, 101, "没有识别到人脸") }
	ErrBodyNotFound     = func() ErrorInfo { return NewErrorInfo(http.StatusOK, 102, "没有识别到人像") }

)
