package requests

import (
	"errors"
	"net/http"
	"testing"
)

func TestErrorTypes(t *testing.T) {
	t.Run("ErrorType_String", func(t *testing.T) {
		testCases := []struct {
			errorType ErrorType
			expected  string
		}{
			{ErrorTypeNetwork, "NetworkError"},
			{ErrorTypeTimeout, "TimeoutError"},
			{ErrorTypeAuth, "AuthError"},
			{ErrorTypeRateLimit, "RateLimitError"},
			{ErrorTypeServerError, "ServerError"},
			{ErrorTypeClientError, "ClientError"},
			{ErrorTypeValidation, "ValidationError"},
			{ErrorTypeSerialization, "SerializationError"},
			{ErrorTypeRedirect, "RedirectError"},
			{ErrorTypeInternal, "InternalError"},
			{ErrorType(999), "UnknownError"},
		}

		for _, tc := range testCases {
			t.Run(tc.expected, func(t *testing.T) {
				if result := tc.errorType.String(); result != tc.expected {
					t.Errorf("Expected %s, got %s", tc.expected, result)
				}
			})
		}
	})
}

func TestRequestError(t *testing.T) {
	t.Run("NewRequestError", func(t *testing.T) {
		cause := errors.New("underlying error")
		err := NewRequestError(ErrorTypeNetwork, "network failed", cause)

		if err.Type != ErrorTypeNetwork {
			t.Errorf("Expected ErrorTypeNetwork, got %v", err.Type)
		}
		if err.Message != "network failed" {
			t.Errorf("Expected 'network failed', got %s", err.Message)
		}
		if err.Cause != cause {
			t.Errorf("Expected cause to be set correctly")
		}
	})

	t.Run("Error_WithCause", func(t *testing.T) {
		cause := errors.New("connection refused")
		err := NewNetworkError("failed to connect", cause)

		expected := "NetworkError: failed to connect (caused by: connection refused)"
		if err.Error() != expected {
			t.Errorf("Expected %s, got %s", expected, err.Error())
		}
	})

	t.Run("Error_WithoutCause", func(t *testing.T) {
		err := NewValidationError("invalid URL format", nil)

		expected := "ValidationError: invalid URL format"
		if err.Error() != expected {
			t.Errorf("Expected %s, got %s", expected, err.Error())
		}
	})

	t.Run("Unwrap", func(t *testing.T) {
		cause := errors.New("original error")
		err := NewTimeoutError("request timeout", cause)

		if unwrapped := err.Unwrap(); unwrapped != cause {
			t.Errorf("Expected unwrapped error to be %v, got %v", cause, unwrapped)
		}
	})

	t.Run("Is", func(t *testing.T) {
		err1 := NewNetworkError("network error", nil)
		err2 := NewNetworkError("another network error", nil)
		err3 := NewTimeoutError("timeout error", nil)

		// 相同类型的错误应该匹配
		if !err1.Is(err2) {
			t.Error("Expected errors of same type to match")
		}

		// 不同类型的错误不应该匹配
		if err1.Is(err3) {
			t.Error("Expected errors of different types not to match")
		}

		// nil 不应该匹配
		if err1.Is(nil) {
			t.Error("Expected error not to match nil")
		}
	})
}

func TestErrorCreators(t *testing.T) {
	t.Run("NewNetworkError", func(t *testing.T) {
		cause := errors.New("connection failed")
		err := NewNetworkError("network issue", cause)

		if err.Type != ErrorTypeNetwork {
			t.Errorf("Expected ErrorTypeNetwork, got %v", err.Type)
		}
		if err.Message != "network issue" {
			t.Errorf("Expected 'network issue', got %s", err.Message)
		}
		if err.Cause != cause {
			t.Errorf("Expected cause to be set")
		}
	})

	t.Run("NewServerError", func(t *testing.T) {
		err := NewServerError(500, "internal server error")

		if err.Type != ErrorTypeServerError {
			t.Errorf("Expected ErrorTypeServerError, got %v", err.Type)
		}
		if err.StatusCode != 500 {
			t.Errorf("Expected status code 500, got %d", err.StatusCode)
		}
		if err.Message != "internal server error" {
			t.Errorf("Expected 'internal server error', got %s", err.Message)
		}
	})

	t.Run("NewClientError", func(t *testing.T) {
		err := NewClientError(404, "not found")

		if err.Type != ErrorTypeClientError {
			t.Errorf("Expected ErrorTypeClientError, got %v", err.Type)
		}
		if err.StatusCode != 404 {
			t.Errorf("Expected status code 404, got %d", err.StatusCode)
		}
	})
}

func TestErrorWrappers(t *testing.T) {
	t.Run("WrapWithRequest", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "https://example.com", nil)
		err := NewNetworkError("network failed", nil)

		wrapped := err.WrapWithRequest(req)

		if wrapped.Request != req {
			t.Error("Expected request to be set")
		}
		if wrapped.URL.String() != "https://example.com" {
			t.Errorf("Expected URL to be set to https://example.com, got %s", wrapped.URL.String())
		}
	})

	t.Run("WrapWithResponse", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: 500,
		}
		err := NewServerError(500, "server error")

		wrapped := err.WrapWithResponse(resp)

		if wrapped.Response != resp {
			t.Error("Expected response to be set")
		}
		if wrapped.StatusCode != 500 {
			t.Errorf("Expected status code 500, got %d", wrapped.StatusCode)
		}
	})
}

func TestErrorTypeCheckers(t *testing.T) {
	t.Run("IsNetworkError", func(t *testing.T) {
		networkErr := NewNetworkError("network failed", nil)
		timeoutErr := NewTimeoutError("timeout", nil)
		regularErr := errors.New("regular error")

		if !IsNetworkError(networkErr) {
			t.Error("Expected IsNetworkError to return true for network error")
		}
		if IsNetworkError(timeoutErr) {
			t.Error("Expected IsNetworkError to return false for timeout error")
		}
		if IsNetworkError(regularErr) {
			t.Error("Expected IsNetworkError to return false for regular error")
		}
	})

	t.Run("IsTimeoutError", func(t *testing.T) {
		timeoutErr := NewTimeoutError("timeout", nil)
		networkErr := NewNetworkError("network failed", nil)

		if !IsTimeoutError(timeoutErr) {
			t.Error("Expected IsTimeoutError to return true for timeout error")
		}
		if IsTimeoutError(networkErr) {
			t.Error("Expected IsTimeoutError to return false for network error")
		}
	})

	t.Run("IsServerError", func(t *testing.T) {
		serverErr := NewServerError(500, "server error")
		clientErr := NewClientError(404, "not found")

		// 测试直接的服务器错误
		if !IsServerError(serverErr) {
			t.Error("Expected IsServerError to return true for server error")
		}
		if IsServerError(clientErr) {
			t.Error("Expected IsServerError to return false for client error")
		}

		// 测试基于状态码的判断
		statusErr := &RequestError{
			Type:       ErrorTypeInternal,
			StatusCode: 502,
		}
		if !IsServerError(statusErr) {
			t.Error("Expected IsServerError to return true for 5xx status code")
		}
	})

	t.Run("IsClientError", func(t *testing.T) {
		clientErr := NewClientError(404, "not found")
		serverErr := NewServerError(500, "server error")

		if !IsClientError(clientErr) {
			t.Error("Expected IsClientError to return true for client error")
		}
		if IsClientError(serverErr) {
			t.Error("Expected IsClientError to return false for server error")
		}

		// 测试基于状态码的判断
		statusErr := &RequestError{
			Type:       ErrorTypeInternal,
			StatusCode: 401,
		}
		if !IsClientError(statusErr) {
			t.Error("Expected IsClientError to return true for 4xx status code")
		}
	})

	t.Run("IsValidationError", func(t *testing.T) {
		validationErr := NewValidationError("invalid input", nil)
		networkErr := NewNetworkError("network failed", nil)

		if !IsValidationError(validationErr) {
			t.Error("Expected IsValidationError to return true for validation error")
		}
		if IsValidationError(networkErr) {
			t.Error("Expected IsValidationError to return false for network error")
		}
	})
}

func TestGetStatusCode(t *testing.T) {
	t.Run("WithStatusCode", func(t *testing.T) {
		err := NewServerError(500, "server error")

		if code := GetStatusCode(err); code != 500 {
			t.Errorf("Expected status code 500, got %d", code)
		}
	})

	t.Run("WithoutStatusCode", func(t *testing.T) {
		err := NewNetworkError("network failed", nil)

		if code := GetStatusCode(err); code != 0 {
			t.Errorf("Expected status code 0, got %d", code)
		}
	})

	t.Run("RegularError", func(t *testing.T) {
		err := errors.New("regular error")

		if code := GetStatusCode(err); code != 0 {
			t.Errorf("Expected status code 0 for regular error, got %d", code)
		}
	})
}

func TestErrorIntegration(t *testing.T) {
	t.Run("CompleteErrorFlow", func(t *testing.T) {
		// 模拟一个完整的错误流程
		cause := errors.New("connection refused")
		err := NewNetworkError("failed to connect to server", cause)

		// 添加请求信息
		req, _ := http.NewRequest("GET", "https://api.example.com/users", nil)
		err = err.WrapWithRequest(req)

		// 验证错误信息完整性
		if err.Type != ErrorTypeNetwork {
			t.Error("Error type not preserved")
		}
		if err.Request != req {
			t.Error("Request not attached")
		}
		if err.URL.String() != "https://api.example.com/users" {
			t.Error("URL not extracted from request")
		}

		// 验证错误可以被正确识别
		if !IsNetworkError(err) {
			t.Error("Error not correctly identified as network error")
		}

		// 验证原始错误可以被解包
		if unwrapped := errors.Unwrap(err); unwrapped != cause {
			t.Error("Original error not preserved")
		}
	})
}
