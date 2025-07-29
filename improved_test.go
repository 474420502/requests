package requests

import (
	"context"
	"strings"
	"testing"
	"time"
)

// TestImprovedErrorHandling 测试改进的错误处理
func TestImprovedErrorHandling(t *testing.T) {
	session := NewSession()

	// 测试无效URL错误处理
	resp, err := session.Get("invalid-url").Execute()
	if err == nil {
		t.Errorf("应该返回URL错误，但得到了响应: %v", resp)
	}
	if !strings.Contains(err.Error(), "invalid") {
		t.Errorf("错误信息应该包含'invalid'，实际: %v", err)
	}
}

// TestNewRequestAPI 测试新的Request API
func TestNewRequestAPI(t *testing.T) {
	session, err := NewSessionWithOptions(
		WithTimeout(10*time.Second),
		WithUserAgent("TestAgent/1.0"),
	)
	if err != nil {
		t.Fatalf("创建Session失败: %v", err)
	}

	// 测试链式调用
	req := session.Get("http://httpbin.org/get").
		SetHeader("X-Test", "true").
		AddQuery("test", "value")

	if req.Error() != nil {
		t.Errorf("请求构建失败: %v", req.Error())
	}

	// 测试无效URL的错误处理
	req2 := session.Get("invalid-url")
	if req2.Error() == nil {
		t.Error("应该检测到无效URL错误")
	}
}

// TestContextSupport 测试Context支持
func TestContextSupport(t *testing.T) {
	session, err := NewSessionWithOptions()
	if err != nil {
		t.Fatalf("创建Session失败: %v", err)
	}

	// 测试超时Context
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err = session.Get("http://httpbin.org/delay/1").
		WithContext(ctx).
		Execute()

	duration := time.Since(start)

	if err == nil {
		t.Error("应该因为超时而失败")
	}

	// 应该很快就超时，不应该等待1秒
	if duration > 100*time.Millisecond {
		t.Errorf("超时时间过长: %v", duration)
	}
}

// TestJSONBodySafety 测试JSON Body的类型安全
func TestJSONBodySafety(t *testing.T) {
	session, err := NewSessionWithOptions()
	if err != nil {
		t.Fatalf("创建Session失败: %v", err)
	}

	// 测试各种JSON类型
	testCases := []interface{}{
		"simple string",
		[]byte(`{"key": "value"}`),
		map[string]string{"name": "test"},
		struct {
			Name string `json:"name"`
		}{Name: "test"},
	}

	for _, tc := range testCases {
		req := session.Post("http://httpbin.org/post").SetBodyJSON(tc)
		if req.Error() != nil {
			t.Errorf("JSON序列化失败，类型: %T, 错误: %v", tc, req.Error())
		}
	}
}

// TestMiddlewareBasics 测试基础中间件功能
func TestMiddlewareBasics(t *testing.T) {
	session, err := NewSessionWithOptions()
	if err != nil {
		t.Fatalf("创建Session失败: %v", err)
	}

	// 测试认证中间件
	authMiddleware := &AuthMiddleware{
		TokenProvider: func() (string, error) {
			return "test-token", nil
		},
	}

	req := session.Get("http://httpbin.org/bearer").
		WithMiddlewares(authMiddleware)

	// 构建请求以验证认证头部被添加
	httpReq, err := req.buildHTTPRequest()
	if err != nil {
		t.Fatalf("构建HTTP请求失败: %v", err)
	}

	// 执行中间件
	err = authMiddleware.BeforeRequest(httpReq)
	if err != nil {
		t.Fatalf("中间件执行失败: %v", err)
	}

	authHeader := httpReq.Header.Get("Authorization")
	expected := "Bearer test-token"
	if authHeader != expected {
		t.Errorf("认证头部不正确，期望: %s, 实际: %s", expected, authHeader)
	}
}

// BenchmarkOriginalVsNew 对比原有API和新API的性能
func BenchmarkOriginalAPI(b *testing.B) {
	session := NewSession()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := session.Get("http://httpbin.org/get")
		req.AddHeader("X-Test", "benchmark")
		// 不执行请求，只测试构建过程
		_ = req
	}
}

func BenchmarkNewAPI(b *testing.B) {
	session, _ := NewSessionWithOptions()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := session.Get("http://httpbin.org/get").
			SetHeader("X-Test", "benchmark")
		// 不执行请求，只测试构建过程
		_ = req
	}
}
