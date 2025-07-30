package requests

import (
	"net/http"
	"net/url"
	"testing"
)

// TestMiddlewareMethods 测试中间件相关方法
func TestMiddlewareMethods(t *testing.T) {
	t.Run("AddMiddleware", func(t *testing.T) {
		session := NewSession()

		// 创建一个测试中间件
		middleware := &TestLoggingMiddleware{}

		session.AddMiddleware(middleware)

		middlewares := session.GetMiddlewares()
		if len(middlewares) != 1 {
			t.Errorf("Expected 1 middleware, got %d", len(middlewares))
		}

		// 检查中间件是否正确添加（类型断言）
		if _, ok := middlewares[0].(*TestLoggingMiddleware); !ok {
			t.Error("Middleware not added correctly")
		}
	})

	t.Run("SetMiddlewares", func(t *testing.T) {
		session := NewSession()

		middleware1 := &TestLoggingMiddleware{}
		middleware2 := &TestLoggingMiddleware{}
		middlewares := []Middleware{middleware1, middleware2}

		session.SetMiddlewares(middlewares)

		sessionMiddlewares := session.GetMiddlewares()
		if len(sessionMiddlewares) != 2 {
			t.Errorf("Expected 2 middlewares, got %d", len(sessionMiddlewares))
		}
	})

	t.Run("ClearMiddlewares", func(t *testing.T) {
		session := NewSession()
		session.AddMiddleware(&TestLoggingMiddleware{})

		// 确认有中间件
		if len(session.GetMiddlewares()) == 0 {
			t.Error("Should have middleware before clearing")
		}

		session.ClearMiddlewares()

		// 确认已清除
		if len(session.GetMiddlewares()) != 0 {
			t.Error("Middlewares should be cleared")
		}
	})
}

// TestSessionQueryMethods 测试Session的查询参数方法
func TestSessionQueryMethods(t *testing.T) {
	t.Run("SetQuery", func(t *testing.T) {
		session := NewSession()
		values := url.Values{}
		values.Add("key1", "value1")
		values.Add("key2", "value2")

		session.SetQuery(values)

		queryValues := session.GetQuery()
		if queryValues.Get("key1") != "value1" {
			t.Error("Query parameter not set correctly")
		}
		if queryValues.Get("key2") != "value2" {
			t.Error("Query parameter not set correctly")
		}
	})

	t.Run("GetQuery", func(t *testing.T) {
		session := NewSession()
		queryValues := session.GetQuery()

		// 初始应该为空
		if len(queryValues) != 0 {
			t.Error("Initial query should be empty")
		}
	})
}

// TestSessionHeaderMethods 测试Session的Header方法
func TestSessionHeaderMethods(t *testing.T) {
	t.Run("SetHeader", func(t *testing.T) {
		session := NewSession()
		headers := make(map[string][]string)
		headers["Content-Type"] = []string{"application/json"}
		headers["User-Agent"] = []string{"Test-Agent"}

		session.SetHeader(headers)

		sessionHeaders := session.GetHeader()
		if sessionHeaders.Get("Content-Type") != "application/json" {
			t.Error("Header not set correctly")
		}
		if sessionHeaders.Get("User-Agent") != "Test-Agent" {
			t.Error("Header not set correctly")
		}
	})

	t.Run("AddHeader", func(t *testing.T) {
		session := NewSession()
		session.AddHeader("X-Custom-Header", "custom-value")

		headers := session.GetHeader()
		if headers.Get("X-Custom-Header") != "custom-value" {
			t.Error("Header not added correctly")
		}
	})

	t.Run("SetContentType", func(t *testing.T) {
		session := NewSession()
		session.SetContentType("application/xml")

		headers := session.GetHeader()
		if headers.Get("Content-Type") != "application/xml" {
			t.Error("Content-Type not set correctly")
		}
	})
}

// TestSessionCookieMethods 测试Session的Cookie方法
func TestSessionCookieMethods(t *testing.T) {
	t.Run("SetCookies", func(t *testing.T) {
		session := NewSession()
		testURL, _ := url.Parse("https://example.com")

		cookies := []*http.Cookie{
			{Name: "cookie1", Value: "value1"},
			{Name: "cookie2", Value: "value2"},
		}

		session.SetCookies(testURL, cookies)

		retrievedCookies := session.GetCookies(testURL)
		if len(retrievedCookies) != 2 {
			t.Errorf("Expected 2 cookies, got %d", len(retrievedCookies))
		}
	})

	t.Run("DelCookies", func(t *testing.T) {
		session := NewSession()
		testURL, _ := url.Parse("https://example.com")

		// 先设置一个cookie
		cookies := []*http.Cookie{
			{Name: "test-cookie", Value: "test-value"},
		}
		session.SetCookies(testURL, cookies)

		// 删除cookie
		session.DelCookies(testURL, "test-cookie")

		// 验证cookie已删除
		remainingCookies := session.GetCookies(testURL)
		for _, cookie := range remainingCookies {
			if cookie.Name == "test-cookie" {
				t.Error("Cookie should be deleted")
			}
		}
	})

	t.Run("ClearCookies", func(t *testing.T) {
		session := NewSession()
		testURL, _ := url.Parse("https://example.com")

		// 先设置一些cookies
		cookies := []*http.Cookie{
			{Name: "cookie1", Value: "value1"},
			{Name: "cookie2", Value: "value2"},
		}
		session.SetCookies(testURL, cookies)

		// 清除所有cookies
		err := session.ClearCookies()
		if err != nil {
			t.Errorf("ClearCookies should not return error: %v", err)
		}

		// 验证cookies已清除
		remainingCookies := session.GetCookies(testURL)
		if len(remainingCookies) != 0 {
			t.Error("All cookies should be cleared")
		}
	})
}

// TestSessionHTTPMethods 测试Session的HTTP方法
func TestSessionHTTPMethods(t *testing.T) {
	session := NewSession()
	testURL := "https://httpbin.org/get"

	t.Run("Head", func(t *testing.T) {
		req := session.Head(testURL)
		if req == nil {
			t.Fatal("Head request should not be nil")
		}
		if req.method != "HEAD" {
			t.Errorf("Expected HEAD method, got %s", req.method)
		}
	})

	t.Run("Get", func(t *testing.T) {
		req := session.Get(testURL)
		if req == nil {
			t.Fatal("Get request should not be nil")
		}
		if req.method != "GET" {
			t.Errorf("Expected GET method, got %s", req.method)
		}
	})

	t.Run("Post", func(t *testing.T) {
		req := session.Post(testURL)
		if req == nil {
			t.Fatal("Post request should not be nil")
		}
		if req.method != "POST" {
			t.Errorf("Expected POST method, got %s", req.method)
		}
	})

	t.Run("Put", func(t *testing.T) {
		req := session.Put(testURL)
		if req == nil {
			t.Fatal("Put request should not be nil")
		}
		if req.method != "PUT" {
			t.Errorf("Expected PUT method, got %s", req.method)
		}
	})

	t.Run("Patch", func(t *testing.T) {
		req := session.Patch(testURL)
		if req == nil {
			t.Fatal("Patch request should not be nil")
		}
		if req.method != "PATCH" {
			t.Errorf("Expected PATCH method, got %s", req.method)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		req := session.Delete(testURL)
		if req == nil {
			t.Fatal("Delete request should not be nil")
		}
		if req.method != "DELETE" {
			t.Errorf("Expected DELETE method, got %s", req.method)
		}
	})

	t.Run("Options", func(t *testing.T) {
		req := session.Options(testURL)
		if req == nil {
			t.Fatal("Options request should not be nil")
		}
		if req.method != "OPTIONS" {
			t.Errorf("Expected OPTIONS method, got %s", req.method)
		}
	})
}

// TestLoggingMiddleware 用于测试的简单中间件
type TestLoggingMiddleware struct{}

func (m *TestLoggingMiddleware) BeforeRequest(req *http.Request) error {
	return nil
}

func (m *TestLoggingMiddleware) AfterResponse(resp *http.Response) error {
	return nil
}
