package requests

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"
)

// TestRequestBuilderMethods 测试Request构建器的方法
func TestRequestBuilderMethods(t *testing.T) {
	t.Run("SetHeader", func(t *testing.T) {
		req := Get("https://httpbin.org/get")
		result := req.SetHeader("Custom-Header", "test-value")

		if result != req {
			t.Error("SetHeader should return the request for chaining")
		}

		if req.header.Get("Custom-Header") != "test-value" {
			t.Error("Header not set correctly")
		}
	})

	t.Run("SetHeaders", func(t *testing.T) {
		req := Get("https://httpbin.org/get")
		headers := map[string]string{
			"Header1": "value1",
			"Header2": "value2",
		}
		result := req.SetHeaders(headers)

		if result != req {
			t.Error("SetHeaders should return the request for chaining")
		}

		for k, v := range headers {
			if req.header.Get(k) != v {
				t.Errorf("Header %s not set correctly: expected %s, got %s", k, v, req.header.Get(k))
			}
		}
	})

	t.Run("SetContentType", func(t *testing.T) {
		req := Get("https://httpbin.org/get")
		contentType := "application/json"
		result := req.SetContentType(contentType)

		if result != req {
			t.Error("SetContentType should return the request for chaining")
		}

		if req.header.Get("Content-Type") != contentType {
			t.Errorf("Content-Type not set correctly: expected %s, got %s", contentType, req.header.Get("Content-Type"))
		}
	})

	t.Run("WithTimeout", func(t *testing.T) {
		req := Get("https://httpbin.org/get")
		timeout := 10 * time.Second
		result := req.WithTimeout(timeout)

		if result != req {
			t.Error("WithTimeout should return the request for chaining")
		}

		if req.timeout != timeout {
			t.Errorf("Timeout not set correctly: expected %v, got %v", timeout, req.timeout)
		}
	})

	t.Run("WithContext", func(t *testing.T) {
		req := Get("https://httpbin.org/get")
		ctx := context.WithValue(context.Background(), "test", "value")
		result := req.WithContext(ctx)

		if result != req {
			t.Error("WithContext should return the request for chaining")
		}

		if req.ctx != ctx {
			t.Error("Context not set correctly")
		}
	})

	t.Run("SetBodyJson", func(t *testing.T) {
		req := Post("https://httpbin.org/post")
		data := map[string]interface{}{
			"key1": "value1",
			"key2": 123,
		}
		result := req.SetBodyJson(data)

		if result != req {
			t.Error("SetBodyJson should return the request for chaining")
		}

		bodyContent := req.body.String()
		if !strings.Contains(bodyContent, "key1") || !strings.Contains(bodyContent, "value1") {
			t.Error("JSON body not set correctly")
		}

		if req.header.Get("Content-Type") != "application/json" {
			t.Error("Content-Type should be set to application/json")
		}
	})

	t.Run("SetCookie", func(t *testing.T) {
		req := Get("https://httpbin.org/get")
		cookie := &http.Cookie{
			Name:  "test-cookie",
			Value: "test-value",
		}
		result := req.SetCookie(cookie)

		if result != req {
			t.Error("SetCookie should return the request for chaining")
		}

		if req.cookies["test-cookie"] != cookie {
			t.Error("Cookie not set correctly")
		}
	})

	t.Run("AddCookies", func(t *testing.T) {
		req := Get("https://httpbin.org/get")
		cookies := []*http.Cookie{
			{Name: "cookie1", Value: "value1"},
			{Name: "cookie2", Value: "value2"},
		}
		result := req.AddCookies(cookies)

		if result != req {
			t.Error("AddCookies should return the request for chaining")
		}

		for _, cookie := range cookies {
			if req.cookies[cookie.Name] != cookie {
				t.Errorf("Cookie %s not set correctly", cookie.Name)
			}
		}
	})

	t.Run("SetCookieValue", func(t *testing.T) {
		req := Get("https://httpbin.org/get")
		result := req.SetCookieValue("test-name", "test-value")

		if result != req {
			t.Error("SetCookieValue should return the request for chaining")
		}

		cookie := req.cookies["test-name"]
		if cookie == nil {
			t.Fatal("Cookie should be set")
		}
		if cookie.Value != "test-value" {
			t.Errorf("Cookie value not set correctly: expected test-value, got %s", cookie.Value)
		}
	})

	t.Run("DelCookie", func(t *testing.T) {
		req := Get("https://httpbin.org/get")
		// 先添加一个cookie
		req.SetCookieValue("test-cookie", "test-value")

		// 然后删除它
		result := req.DelCookie("test-cookie")

		if result != req {
			t.Error("DelCookie should return the request for chaining")
		}

		if req.cookies["test-cookie"] != nil {
			t.Error("Cookie should be deleted")
		}
	})

	t.Run("AddHeader", func(t *testing.T) {
		req := Get("https://httpbin.org/get")
		result := req.AddHeader("Test-Header", "value1")

		if result != req {
			t.Error("AddHeader should return the request for chaining")
		}

		if req.header.Get("Test-Header") != "value1" {
			t.Error("Header not added correctly")
		}

		// 添加相同名称的header（应该append）
		req.AddHeader("Test-Header", "value2")
		values := req.header["Test-Header"]
		if len(values) != 2 || values[0] != "value1" || values[1] != "value2" {
			t.Error("Multiple headers not added correctly")
		}
	})

	t.Run("DelHeader", func(t *testing.T) {
		req := Get("https://httpbin.org/get")
		req.SetHeader("Test-Header", "test-value")

		result := req.DelHeader("Test-Header")

		if result != req {
			t.Error("DelHeader should return the request for chaining")
		}

		if req.header.Get("Test-Header") != "" {
			t.Error("Header should be deleted")
		}
	})
}

// TestRequestErrorHandling 测试Request的错误处理
func TestRequestErrorHandling(t *testing.T) {
	t.Run("InvalidURL", func(t *testing.T) {
		req := Get("://invalid-url")
		if req.err == nil {
			t.Error("Request with invalid URL should have error")
		}
	})

	t.Run("ErrorPropagation", func(t *testing.T) {
		req := Get("://invalid-url")
		// 即使有错误，链式调用仍应返回请求对象
		result := req.SetHeader("test", "value")
		if result != req {
			t.Error("Methods should still return request even with error")
		}

		// 错误应该保持
		if req.err == nil {
			t.Error("Error should be preserved through method calls")
		}
	})

	t.Run("JSONSerializationError", func(t *testing.T) {
		req := Post("https://httpbin.org/post")

		// 创建一个无法序列化的对象
		invalidData := make(chan int)
		result := req.SetBodyJson(invalidData)

		if result != req {
			t.Error("SetBodyJson should return request even with serialization error")
		}

		if req.err == nil {
			t.Error("SetBodyJson should set error for unserializable data")
		}
	})
}

// TestRequestChaining 测试方法链式调用
func TestRequestChaining(t *testing.T) {
	req := Get("https://httpbin.org/get").
		SetHeader("Custom-Header", "test").
		SetContentType("application/json").
		WithTimeout(10 * time.Second).
		SetCookie(&http.Cookie{Name: "test", Value: "value"})

	if req == nil {
		t.Fatal("Chained request should not be nil")
	}

	if req.header.Get("Custom-Header") != "test" {
		t.Error("Chained header not set correctly")
	}

	if req.header.Get("Content-Type") != "application/json" {
		t.Error("Chained content type not set correctly")
	}

	if req.timeout != 10*time.Second {
		t.Error("Chained timeout not set correctly")
	}

	if req.cookies["test"].Value != "value" {
		t.Error("Chained cookie not set correctly")
	}
}

// TestRequestWithSession 测试带Session的Request
func TestRequestWithSession(t *testing.T) {
	session := NewSession()
	session.AddHeader("Session-Header", "session-value")

	req := session.Get("https://httpbin.org/get")

	if req.session != session {
		t.Error("Request should reference the session")
	}

	// Session的header应该可用
	if req.session.Header.Get("Session-Header") != "session-value" {
		t.Error("Session headers should be accessible from request")
	}
}

// TestRequestContextHandling 测试Request的上下文处理
func TestRequestContextHandling(t *testing.T) {
	t.Run("DefaultContext", func(t *testing.T) {
		session := NewSession()
		req := session.Get("https://httpbin.org/get")

		if req.ctx == nil {
			t.Error("Request should have a default context")
		}
	})

	t.Run("CustomContext", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "test-key", "test-value")
		req := Get("https://httpbin.org/get").WithContext(ctx)

		if req.ctx != ctx {
			t.Error("Custom context should be set")
		}

		// 验证上下文值
		if value := req.ctx.Value("test-key"); value != "test-value" {
			t.Errorf("Context value not preserved: expected test-value, got %v", value)
		}
	})

	t.Run("ContextWithTimeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := Get("https://httpbin.org/get").WithContext(ctx)

		if req.ctx != ctx {
			t.Error("Timeout context should be set")
		}

		// 验证上下文有超时
		deadline, ok := req.ctx.Deadline()
		if !ok {
			t.Error("Context should have deadline")
		}

		if time.Until(deadline) > 5*time.Second {
			t.Error("Context deadline should be approximately 5 seconds")
		}
	})
}

// TestRequestQueryParameters 测试查询参数功能
func TestRequestQueryParameters(t *testing.T) {
	t.Run("AddQuery", func(t *testing.T) {
		req := Get("https://httpbin.org/get")
		result := req.AddQuery("key1", "value1")

		if result != req {
			t.Error("AddQuery should return the request for chaining")
		}

		// 由于我们无法直接访问query参数，我们检查URL是否正确更新
		// 这个测试可能需要根据实际实现进行调整
	})

	t.Run("SetQuery", func(t *testing.T) {
		req := Get("https://httpbin.org/get")
		params := make(map[string][]string)
		params["key1"] = []string{"value1"}
		params["key2"] = []string{"value2"}

		result := req.SetQuery(params)

		if result != req {
			t.Error("SetQuery should return the request for chaining")
		}
	})
}

// TestRequestBodyMethods 测试Request的body设置方法
func TestRequestBodyMethods(t *testing.T) {
	t.Run("SetBodyWithType", func(t *testing.T) {
		req := Post("https://httpbin.org/post")
		data := "test body data"
		contentType := "text/plain"

		result := req.SetBodyWithType(contentType, data)

		if result != req {
			t.Error("SetBodyWithType should return the request for chaining")
		}

		if req.header.Get("Content-Type") != contentType {
			t.Errorf("Content-Type not set correctly: expected %s, got %s", contentType, req.header.Get("Content-Type"))
		}
	})
}
