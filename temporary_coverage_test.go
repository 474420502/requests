package requests

import (
	"strings"
	"testing"
)

// TestTemporary_ZeroCoverageMethods 测试Temporary中0%覆盖率的方法
func TestTemporary_ZeroCoverageMethods(t *testing.T) {
	session := NewSession()

	t.Run("SetBodyUrlencoded", func(t *testing.T) {
		temp := NewTemporary(session, "http://httpbin.org/post")
		temp.SetMethod("POST")

		data := map[string]string{
			"key1": "value1",
			"key2": "value2",
		}

		temp.SetBodyUrlencoded(data)

		// 验证内容类型是否正确设置
		if temp.req.header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Error("Expected Content-Type to be application/x-www-form-urlencoded")
		}
	})

	t.Run("SetBodyPlain", func(t *testing.T) {
		temp := NewTemporary(session, "http://httpbin.org/post")
		temp.SetMethod("POST")

		plainText := "This is plain text data"
		temp.SetBodyPlain(plainText)

		// 验证内容类型是否正确设置
		if temp.req.header.Get("Content-Type") != "text/plain" {
			t.Error("Expected Content-Type to be text/plain")
		}
	})

	t.Run("SetBodyStream", func(t *testing.T) {
		temp := NewTemporary(session, "http://httpbin.org/post")
		temp.SetMethod("POST")

		// SetBodyStream expects string, []byte, or []rune, not io.Reader
		streamData := "Stream data content"
		temp.SetBodyStream(streamData)

		// 验证内容类型是否正确设置
		if temp.req.header.Get("Content-Type") != "application/octet-stream" {
			t.Errorf("Expected Content-Type to be application/octet-stream, got %s", temp.req.header.Get("Content-Type"))
		}
	})

	t.Run("CreateBodyMultipart", func(t *testing.T) {
		temp := NewTemporary(session, "http://httpbin.org/post")
		temp.SetMethod("POST")

		multipart := temp.CreateBodyMultipart()
		if multipart == nil {
			t.Error("Expected CreateBodyMultipart to return non-nil MultipartFormData")
		}

		// CreateBodyMultipart只是创建MultipartFormData对象，不会自动设置Content-Type
		// Content-Type是在实际使用multipart数据时设置的
		// 所以这里我们只验证multipart对象是否被创建
	})

	t.Run("Error", func(t *testing.T) {
		// 测试正常情况
		temp := NewTemporary(session, "http://httpbin.org/get")
		err := temp.Error()
		if err != nil && temp.err == nil {
			t.Errorf("Expected no error for valid temporary, got: %v", err)
		}

		// 测试错误情况
		tempWithError := NewTemporary(session, "://invalid-url")
		err = tempWithError.Error()
		if err == nil {
			t.Error("Expected error for invalid URL")
		}
	})

	t.Run("Execute", func(t *testing.T) {
		temp := NewTemporary(session, "http://httpbin.org/get")
		temp.SetMethod("GET")

		// Execute方法应该委托给内部的Request
		resp, err := temp.Execute()
		// 由于httpbin.org可能不可达，我们主要测试方法是否存在且不panic
		_ = resp
		_ = err
		// 不检查具体结果，因为网络可能不可用
	})
}

// TestTemporary_BuildRequestAndTestExecute 测试构建请求和测试执行
func TestTemporary_BuildRequestAndTestExecute(t *testing.T) {
	session := NewSession()

	t.Run("BuildRequestValid", func(t *testing.T) {
		temp := NewTemporary(session, "http://httpbin.org/post")
		temp.SetMethod("POST")
		temp.SetBodyPlain("test data")

		httpReq, err := temp.BuildRequest()
		if err != nil {
			t.Errorf("Expected no error from BuildRequest, got: %v", err)
		}

		if httpReq == nil {
			t.Error("Expected non-nil HTTP request")
		}

		if httpReq.Method != "POST" {
			t.Errorf("Expected method POST, got %s", httpReq.Method)
		}
	})

	t.Run("BuildRequestWithError", func(t *testing.T) {
		temp := NewTemporary(session, "://invalid-url")
		temp.SetMethod("GET")

		httpReq, err := temp.BuildRequest()
		if err == nil {
			t.Error("Expected error from BuildRequest with invalid URL")
		}

		if httpReq != nil {
			t.Error("Expected nil HTTP request when error occurs")
		}
	})
}

// TestTemporary_SetMethodVariations 测试SetMethod的各种情况
func TestTemporary_SetMethodVariations(t *testing.T) {
	session := NewSession()

	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS", "TRACE"}

	for _, method := range methods {
		t.Run("Method_"+method, func(t *testing.T) {
			temp := NewTemporary(session, "http://httpbin.org/anything")
			temp.SetMethod(method)

			if temp.Method != method {
				t.Errorf("Expected method %s, got %s", method, temp.Method)
			}

			// 验证内部Request也设置了正确的方法
			httpReq, err := temp.BuildRequest()
			if err != nil {
				t.Errorf("BuildRequest failed: %v", err)
				return
			}

			if httpReq.Method != method {
				t.Errorf("Expected HTTP request method %s, got %s", method, httpReq.Method)
			}
		})
	}
}

// TestTemporary_ChainedCalls 测试链式调用
func TestTemporary_ChainedCalls(t *testing.T) {
	session := NewSession()

	t.Run("ChainedMethodCalls", func(t *testing.T) {
		temp := NewTemporary(session, "http://httpbin.org/post")

		// 测试方法链式调用
		result := temp.SetMethod("POST")
		if result != temp {
			t.Error("Expected SetMethod to return the same Temporary instance")
		}

		// 构建请求验证链式调用是否正确工作
		httpReq, err := temp.BuildRequest()
		if err != nil {
			t.Errorf("BuildRequest failed after chained calls: %v", err)
		}

		if httpReq == nil {
			t.Error("Expected non-nil HTTP request after chained calls")
		}
	})
}

// TestTemporary_BodyVariations 测试不同类型的body设置
func TestTemporary_BodyVariations(t *testing.T) {
	session := NewSession()

	t.Run("EmptyBodyMethods", func(t *testing.T) {
		temp := NewTemporary(session, "http://httpbin.org/post")
		temp.SetMethod("POST")

		// 测试空数据
		temp.SetBodyPlain("")
		temp.SetBodyUrlencoded(map[string]string{})

		// 这些调用不应该panic
		httpReq, err := temp.BuildRequest()
		if err != nil {
			t.Errorf("BuildRequest failed with empty body: %v", err)
		}

		if httpReq == nil {
			t.Error("Expected non-nil HTTP request with empty body")
		}
	})

	t.Run("LargeBodyData", func(t *testing.T) {
		temp := NewTemporary(session, "http://httpbin.org/post")
		temp.SetMethod("POST")

		// 测试大量数据
		largeData := strings.Repeat("Large data content ", 1000)
		temp.SetBodyPlain(largeData)

		// 应该能够处理大量数据
		httpReq, err := temp.BuildRequest()
		if err != nil {
			t.Errorf("BuildRequest failed with large body: %v", err)
		}

		if httpReq == nil {
			t.Error("Expected non-nil HTTP request with large body")
		}
	})
}

// TestTemporary_ErrorPropagation 测试错误传播
func TestTemporary_ErrorPropagation(t *testing.T) {
	session := NewSession()

	t.Run("ErrorInConstruction", func(t *testing.T) {
		// 使用无效URL创建Temporary
		temp := NewTemporary(session, "://malformed-url")

		// 错误应该在构造时被捕获
		if temp.err == nil {
			t.Error("Expected error in temporary construction")
		}

		// 后续操作应该保持错误状态
		temp.SetMethod("GET")
		temp.SetBodyPlain("test")

		// BuildRequest应该返回错误
		httpReq, err := temp.BuildRequest()
		if err == nil {
			t.Error("Expected BuildRequest to return error")
		}

		if httpReq != nil {
			t.Error("Expected nil HTTP request when error exists")
		}

		// Execute应该返回错误
		resp, err := temp.Execute()
		if err == nil {
			t.Error("Expected Execute to return error")
		}

		if resp != nil {
			t.Error("Expected nil response when error exists")
		}
	})
}
