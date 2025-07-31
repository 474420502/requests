package requests

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// ErrorType 定义错误类型枚举
type ErrorType int

const (
	// ErrorTypeNetwork 网络连接错误
	ErrorTypeNetwork ErrorType = iota
	// ErrorTypeTimeout 超时错误
	ErrorTypeTimeout
	// ErrorTypeAuth 认证错误
	ErrorTypeAuth
	// ErrorTypeRateLimit 限流错误
	ErrorTypeRateLimit
	// ErrorTypeServerError 服务器错误 (5xx)
	ErrorTypeServerError
	// ErrorTypeClientError 客户端错误 (4xx)
	ErrorTypeClientError
	// ErrorTypeValidation 参数验证错误
	ErrorTypeValidation
	// ErrorTypeSerialization 序列化/反序列化错误
	ErrorTypeSerialization
	// ErrorTypeRedirect 重定向错误
	ErrorTypeRedirect
	// ErrorTypeInternal 内部错误
	ErrorTypeInternal
)

// String 返回错误类型的字符串表示
func (et ErrorType) String() string {
	switch et {
	case ErrorTypeNetwork:
		return "NetworkError"
	case ErrorTypeTimeout:
		return "TimeoutError"
	case ErrorTypeAuth:
		return "AuthError"
	case ErrorTypeRateLimit:
		return "RateLimitError"
	case ErrorTypeServerError:
		return "ServerError"
	case ErrorTypeClientError:
		return "ClientError"
	case ErrorTypeValidation:
		return "ValidationError"
	case ErrorTypeSerialization:
		return "SerializationError"
	case ErrorTypeRedirect:
		return "RedirectError"
	case ErrorTypeInternal:
		return "InternalError"
	default:
		return "UnknownError"
	}
}

// RequestError 统一的请求错误类型
type RequestError struct {
	// Type 错误类型
	Type ErrorType
	// Message 错误消息
	Message string
	// Cause 原始错误
	Cause error
	// Request 失败的请求对象
	Request *http.Request
	// Response 响应对象（如果有的话）
	Response *http.Response
	// URL 请求的 URL
	URL *url.URL
	// StatusCode HTTP 状态码（如果有的话）
	StatusCode int
}

// Error 实现 error 接口
func (e *RequestError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type.String(), e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type.String(), e.Message)
}

// Unwrap 支持 Go 1.13+ 的错误解包
func (e *RequestError) Unwrap() error {
	return e.Cause
}

// Is 支持 errors.Is 检查
func (e *RequestError) Is(target error) bool {
	if target == nil {
		return false
	}

	if t, ok := target.(*RequestError); ok {
		return e.Type == t.Type
	}
	return false
}

// NewRequestError 创建一个新的请求错误
func NewRequestError(errorType ErrorType, message string, cause error) *RequestError {
	return &RequestError{
		Type:    errorType,
		Message: message,
		Cause:   cause,
	}
}

// NewNetworkError 创建网络错误
func NewNetworkError(message string, cause error) *RequestError {
	return NewRequestError(ErrorTypeNetwork, message, cause)
}

// NewTimeoutError 创建超时错误
func NewTimeoutError(message string, cause error) *RequestError {
	return NewRequestError(ErrorTypeTimeout, message, cause)
}

// NewAuthError 创建认证错误
func NewAuthError(message string, cause error) *RequestError {
	return NewRequestError(ErrorTypeAuth, message, cause)
}

// NewValidationError 创建验证错误
func NewValidationError(message string, cause error) *RequestError {
	return NewRequestError(ErrorTypeValidation, message, cause)
}

// NewSerializationError 创建序列化错误
func NewSerializationError(message string, cause error) *RequestError {
	return NewRequestError(ErrorTypeSerialization, message, cause)
}

// NewServerError 创建服务器错误
func NewServerError(statusCode int, message string) *RequestError {
	err := NewRequestError(ErrorTypeServerError, message, nil)
	err.StatusCode = statusCode
	return err
}

// NewClientError 创建客户端错误
func NewClientError(statusCode int, message string) *RequestError {
	err := NewRequestError(ErrorTypeClientError, message, nil)
	err.StatusCode = statusCode
	return err
}

// WrapWithRequest 为错误添加请求信息
func (e *RequestError) WrapWithRequest(req *http.Request) *RequestError {
	e.Request = req
	if req != nil {
		e.URL = req.URL
	}
	return e
}

// WrapWithResponse 为错误添加响应信息
func (e *RequestError) WrapWithResponse(resp *http.Response) *RequestError {
	e.Response = resp
	if resp != nil {
		e.StatusCode = resp.StatusCode
	}
	return e
}

// IsNetworkError 检查是否为网络错误
func IsNetworkError(err error) bool {
	var reqErr *RequestError
	if errors.As(err, &reqErr) {
		return reqErr.Type == ErrorTypeNetwork
	}
	return false
}

// IsTimeoutError 检查是否为超时错误
func IsTimeoutError(err error) bool {
	var reqErr *RequestError
	if errors.As(err, &reqErr) {
		return reqErr.Type == ErrorTypeTimeout
	}
	return false
}

// IsAuthError 检查是否为认证错误
func IsAuthError(err error) bool {
	var reqErr *RequestError
	if errors.As(err, &reqErr) {
		return reqErr.Type == ErrorTypeAuth
	}
	return false
}

// IsServerError 检查是否为服务器错误 (5xx)
func IsServerError(err error) bool {
	var reqErr *RequestError
	if errors.As(err, &reqErr) {
		return reqErr.Type == ErrorTypeServerError ||
			(reqErr.StatusCode >= 500 && reqErr.StatusCode < 600)
	}
	return false
}

// IsClientError 检查是否为客户端错误 (4xx)
func IsClientError(err error) bool {
	var reqErr *RequestError
	if errors.As(err, &reqErr) {
		return reqErr.Type == ErrorTypeClientError ||
			(reqErr.StatusCode >= 400 && reqErr.StatusCode < 500)
	}
	return false
}

// IsValidationError 检查是否为验证错误
func IsValidationError(err error) bool {
	var reqErr *RequestError
	if errors.As(err, &reqErr) {
		return reqErr.Type == ErrorTypeValidation
	}
	return false
}

// GetStatusCode 从错误中提取状态码，如果没有返回 0
func GetStatusCode(err error) int {
	var reqErr *RequestError
	if errors.As(err, &reqErr) {
		return reqErr.StatusCode
	}
	return 0
}
