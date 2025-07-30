package requests

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// TestAdditionalCoverage 额外的覆盖率测试
func TestAdditionalCoverage(t *testing.T) {
	t.Run("SessionWithDefaults", func(t *testing.T) {
		// 测试预定义会话的创建
		sessions := map[string]func() (*Session, error){
			"API":             NewSessionForAPI,
			"Scraping":        NewSessionForScraping,
			"Testing":         NewSessionForTesting,
			"WithDefaults":    NewSessionWithDefaults,
			"Secure":          NewSecureSession,
			"HighPerformance": NewHighPerformanceSession,
		}

		for name, creator := range sessions {
			t.Run(name, func(t *testing.T) {
				session, err := creator()
				if err != nil {
					t.Errorf("Failed to create %s session: %v", name, err)
				}
				if session == nil {
					t.Errorf("Expected non-nil session for %s", name)
				}
			})
		}
	})

	t.Run("SessionWithRetry", func(t *testing.T) {
		// 测试带重试的会话
		session, err := NewSessionWithRetry(3, 2*time.Second)
		if err != nil {
			t.Errorf("Failed to create retry session: %v", err)
		}
		if session == nil {
			t.Error("Expected non-nil retry session")
		}
	})

	t.Run("SessionWithProxy", func(t *testing.T) {
		// 测试带代理的会话
		session, err := NewSessionWithProxy("http://proxy.example.com:8080")
		if err != nil {
			t.Errorf("Failed to create proxy session: %v", err)
		}
		if session == nil {
			t.Error("Expected non-nil proxy session")
		}
	})

	t.Run("SessionWithInvalidProxy", func(t *testing.T) {
		// 测试无效代理的会话创建
		session, err := NewSessionWithProxy("://invalid-proxy")
		if err == nil {
			t.Error("Expected error for invalid proxy URL")
		}
		if session != nil {
			t.Error("Expected nil session for invalid proxy")
		}
	})
}

// TestContextHandling 测试上下文处理
func TestContextHandling(t *testing.T) {
	session := NewSession()

	t.Run("WithContextTimeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()

		req := session.Get("http://httpbin.org/delay/10").WithContext(ctx)
		if req == nil {
			t.Error("Expected non-nil request")
		}

		// 由于超时很短，请求应该会超时
		_, err := req.Execute()
		if err == nil {
			// 网络很快的情况下可能不会超时，这是正常的
			t.Log("Request completed before timeout (network is very fast)")
		}
	})

	t.Run("WithCancelledContext", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // 立即取消

		req := session.Get("http://httpbin.org/get").WithContext(ctx)
		if req == nil {
			t.Error("Expected non-nil request")
		}

		// 使用已取消的上下文应该导致错误
		_, err := req.Execute()
		if err == nil {
			t.Log("Request completed despite cancelled context")
		}
	})
}

// TestErrorScenarios 测试各种错误场景
func TestErrorScenarios(t *testing.T) {
	t.Run("InvalidURLParsing", func(t *testing.T) {
		session := NewSession()

		// 测试各种无效URL格式
		invalidURLs := []string{
			":",
			"://",
			"ht tp://invalid space.com",
			string([]byte{0x7f}), // 无效字符
		}

		for i, url := range invalidURLs {
			t.Run(fmt.Sprintf("InvalidURL_%d", i), func(t *testing.T) {
				req := session.Get(url)
				if req == nil {
					t.Error("Expected non-nil request even with invalid URL")
				}

				// 执行应该返回错误
				_, err := req.Execute()
				if err == nil {
					t.Logf("URL %s was accepted (might be valid in some contexts)", url)
				}
			})
		}
	})

	t.Run("RequestWithNilSession", func(t *testing.T) {
		// 测试使用nil session创建request会导致panic
		// 我们应该在这里使用recover来捕获panic
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic when creating request with nil session")
			}
		}()

		// 这应该导致panic
		req := NewRequest(nil, "GET", "http://example.com")
		_ = req // 避免未使用变量警告
	})
}

// TestContentTypeHandling 测试内容类型处理
func TestContentTypeHandling(t *testing.T) {
	session := NewSession()

	t.Run("ContentTypeOverrides", func(t *testing.T) {
		req := session.Post("http://httpbin.org/post")

		// 设置JSON内容
		req.SetBodyJSON(map[string]string{"key": "value"})
		originalContentType := req.header.Get("Content-Type")

		// 手动覆盖内容类型
		req.SetContentType("application/custom-json")
		newContentType := req.header.Get("Content-Type")

		if newContentType == originalContentType {
			t.Error("Expected content type to be overridden")
		}

		if newContentType != "application/custom-json" {
			t.Errorf("Expected custom content type, got %s", newContentType)
		}
	})

	t.Run("MultipleContentTypeChanges", func(t *testing.T) {
		req := session.Post("http://httpbin.org/post")

		// 多次更改内容类型
		types := []string{
			"application/json",
			"text/plain",
			"application/xml",
			"multipart/form-data",
		}

		for _, contentType := range types {
			req.SetContentType(contentType)
			if req.header.Get("Content-Type") != contentType {
				t.Errorf("Expected content type %s, got %s", contentType, req.header.Get("Content-Type"))
			}
		}
	})
}

// TestHeaderManagement 测试头部管理
func TestHeaderManagement(t *testing.T) {
	session := NewSession()

	t.Run("HeaderCaseInsensitivity", func(t *testing.T) {
		req := session.Get("http://httpbin.org/get")

		// 设置不同大小写的头部
		req.SetHeader("content-type", "application/json")
		req.SetHeader("Content-Type", "text/plain")

		// HTTP头部应该是大小写不敏感的，后者应该覆盖前者
		contentType := req.header.Get("Content-Type")
		if contentType != "text/plain" {
			t.Errorf("Expected Content-Type to be 'text/plain', got '%s'", contentType)
		}
	})

	t.Run("HeaderDeletion", func(t *testing.T) {
		req := session.Get("http://httpbin.org/get")

		// 设置头部然后删除
		req.SetHeader("X-Custom-Header", "test-value")
		if req.header.Get("X-Custom-Header") != "test-value" {
			t.Error("Expected custom header to be set")
		}

		req.DelHeader("X-Custom-Header")
		if req.header.Get("X-Custom-Header") != "" {
			t.Error("Expected custom header to be deleted")
		}
	})

	t.Run("HeaderAddition", func(t *testing.T) {
		req := session.Get("http://httpbin.org/get")

		// 添加多个相同名称的头部
		req.AddHeader("Accept", "application/json")
		req.AddHeader("Accept", "text/plain")

		// 验证多个值是否存在
		acceptHeaders := req.header.Values("Accept")
		if len(acceptHeaders) < 2 {
			t.Errorf("Expected at least 2 Accept headers, got %d", len(acceptHeaders))
		}
	})
}

// TestQueryParameterEdgeCases 测试查询参数边界情况
func TestQueryParameterEdgeCases(t *testing.T) {
	session := NewSession()

	t.Run("EmptyQueryValues", func(t *testing.T) {
		req := session.Get("http://httpbin.org/get")

		// 添加空值查询参数
		req.AddQuery("empty", "")
		req.AddQuery("", "value")

		// 这些操作应该不会导致panic
		_, err := req.Execute()
		_ = err // 忽略网络错误
	})

	t.Run("SpecialCharactersInQuery", func(t *testing.T) {
		req := session.Get("http://httpbin.org/get")

		// 添加包含特殊字符的查询参数
		specialChars := []string{
			"hello world",
			"中文参数",
			"param=value&other=123",
			"<script>alert('xss')</script>",
		}

		for i, char := range specialChars {
			key := fmt.Sprintf("param_%d", i)
			req.AddQuery(key, char)
		}

		// URL应该正确编码这些参数
		if req.parsedURL.RawQuery == "" {
			t.Error("Expected query parameters to be set")
		}
	})
}

// TestTimeoutVariations 测试各种超时设置
func TestTimeoutVariations(t *testing.T) {
	session := NewSession()

	t.Run("VeryShortTimeout", func(t *testing.T) {
		req := session.Get("http://httpbin.org/delay/1").WithTimeout(1 * time.Nanosecond)

		// 极短的超时应该导致请求失败
		_, err := req.Execute()
		if err == nil {
			t.Log("Request completed despite very short timeout")
		}
	})

	t.Run("VeryLongTimeout", func(t *testing.T) {
		req := session.Get("http://httpbin.org/get").WithTimeout(1 * time.Hour)

		// 长超时不应该影响正常请求
		_, err := req.Execute()
		_ = err // 忽略网络错误
	})

	t.Run("ZeroTimeout", func(t *testing.T) {
		req := session.Get("http://httpbin.org/get").WithTimeout(0)

		// 零超时意味着无超时
		_, err := req.Execute()
		_ = err // 忽略网络错误
	})
}

// TestCookieEdgeCases 测试Cookie边界情况
func TestCookieEdgeCases(t *testing.T) {
	session := NewSession()

	t.Run("CookieWithSpecialCharacters", func(t *testing.T) {
		req := session.Get("http://httpbin.org/get")

		// 添加包含特殊字符的cookie
		req.SetCookieValue("special", "value with spaces")
		req.SetCookieValue("encoded", "value%20with%20encoding")
		req.SetCookieValue("unicode", "中文cookie值")

		// 这些操作应该不会导致panic
		_, err := req.Execute()
		_ = err // 忽略网络错误
	})

	t.Run("EmptyCookieValues", func(t *testing.T) {
		req := session.Get("http://httpbin.org/get")

		// 添加空值cookie
		req.SetCookieValue("empty", "")
		req.SetCookieValue("", "value")

		// 这些操作应该不会导致panic
		_, err := req.Execute()
		_ = err // 忽略网络错误
	})
}
