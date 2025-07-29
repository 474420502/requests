package requests

import (
	"testing"
	"time"
)

// TestUnifiedAPI 测试统一的API
func TestUnifiedAPI(t *testing.T) {
	// 创建Session
	session := NewSession()

	// 测试Session方法返回Request而不是Temporary
	req := session.Get("http://httpbin.org/get")
	if req == nil {
		t.Fatal("session.Get() should return a Request")
	}

	// 测试Request的链式调用
	req.AddParam("test", "value").
		SetTimeout(10*time.Second).
		AddHeader("User-Agent", "TestAgent")

	// 测试错误累积
	if err := req.Error(); err != nil {
		t.Logf("Request has error: %v", err)
	}

	t.Log("统一API测试通过")
}

// TestMiddlewareIntegration 测试中间件集成
func TestMiddlewareIntegration(t *testing.T) {
	session := NewSession()

	// 添加中间件
	session.AddMiddleware(&LoggingMiddleware{})

	// 创建请求
	req := session.Post("http://httpbin.org/post")
	req.SetBodyJSON(map[string]string{"test": "data"})

	// 验证中间件被正确集成
	if len(req.middlewares) == 0 {
		t.Error("中间件应该被传递给Request")
	}

	t.Log("中间件集成测试通过")
}

// TestRequestBuilderAPI 测试请求构建器API完整性
func TestRequestBuilderAPI(t *testing.T) {
	session := NewSession()
	req := session.Post("http://httpbin.org/post")

	// 测试头部管理
	req.AddHeader("Custom-Header", "value")
	req.SetHeader("Content-Type", "application/json")
	req.DelHeader("User-Agent")

	// 测试参数管理
	req.AddParam("param1", "value1")
	req.SetParam("param2", "value2")

	// 测试Body设置
	req.SetBodyJSON(map[string]string{"key": "value"})

	// 测试错误处理
	if err := req.Error(); err != nil {
		t.Logf("Request error: %v", err)
	}

	t.Log("请求构建器API测试通过")
}
