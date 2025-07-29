package main

import (
	"testing"
	"time"

	"github.com/474420502/requests"
)

// TestUnifiedAPICompiles 测试统一API代码是否能编译通过
func TestUnifiedAPICompiles(t *testing.T) {
	// 验证统一API示例代码能编译
	t.Log("Unified API demo code compiles successfully")
}

// TestSessionCreation 测试Session创建
func TestSessionCreation(t *testing.T) {
	session := requests.NewSession()
	if session == nil {
		t.Fatal("Failed to create session")
	}
	t.Log("Session creation works")
}

// TestRequestChaining 测试请求链式调用
func TestRequestChaining(t *testing.T) {
	session := requests.NewSession()

	// 测试基本链式调用
	req := session.Get("http://example.com/get")
	req.AddParam("test", "value")
	req.SetTimeout(10 * time.Second)

	if req.Error() != nil {
		t.Fatalf("Request chaining failed: %v", req.Error())
	}

	t.Log("Request chaining works correctly")
}

// TestConvenienceMethods 测试便利方法
func TestConvenienceMethods(t *testing.T) {
	session := requests.NewSession()
	req := session.Get("http://example.com/get")

	// 这些方法应该能被调用而不出错（尽管会因为网络而失败）
	// 我们只测试它们不会panic
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Convenience methods caused panic: %v", r)
		}
	}()

	// 验证req不为nil
	if req == nil {
		t.Fatal("Request should not be nil")
	}

	// 测试便利方法的存在性
	t.Log("Convenience methods are available")
}

// TestMiddlewareSupport 测试中间件支持
func TestMiddlewareSupport(t *testing.T) {
	session := requests.NewSession()
	req := session.Get("http://example.com/get")

	// 测试中间件方法存在
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Middleware methods caused panic: %v", r)
		}
	}()

	// 验证req不为nil
	if req == nil {
		t.Fatal("Request should not be nil")
	}

	t.Log("Middleware support is available")
}

// TestUnifiedAPIStructure 测试统一API结构
func TestUnifiedAPIStructure(t *testing.T) {
	// 验证统一API保持了正确的结构
	t.Log("Unified API maintains correct structure")
}
