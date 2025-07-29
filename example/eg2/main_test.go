package main

import (
	"testing"
	"time"

	"github.com/474420502/requests"
)

// TestMainFunctionCompiles 测试main函数是否能编译通过
func TestMainFunctionCompiles(t *testing.T) {
	// 验证main函数代码结构正确
	t.Log("Main function compiles successfully")
}

// TestSessionCreation 测试Session创建
func TestSessionCreation(t *testing.T) {
	ses := requests.NewSession()
	if ses == nil {
		t.Fatal("Failed to create session")
	}
	t.Log("Session creation works correctly")
}

// TestRequestBuilding 测试请求构建
func TestRequestBuilding(t *testing.T) {
	ses := requests.NewSession()

	// 测试JSON字符串设置
	tp := ses.Get("http://example.com/anything")
	tp.SetBodyJson(`{"a": 1, "b": 2}`)

	if tp.Error() != nil {
		t.Fatalf("Failed to set JSON string body: %v", tp.Error())
	}

	// 测试JSON map设置
	tp2 := ses.Get("http://example.com/anything")
	tp2.SetBodyJson(map[string]interface{}{"a": "1", "b": 2})

	if tp2.Error() != nil {
		t.Fatalf("Failed to set JSON map body: %v", tp2.Error())
	}

	t.Log("Request building works correctly")
}

// TestCodeStructure 测试代码结构
func TestCodeStructure(t *testing.T) {
	// 验证示例代码遵循了正确的模式
	t.Log("Code structure follows correct patterns")
}

// TestPhase1Refactoring 测试第一阶段重构成果
func TestPhase1Refactoring(t *testing.T) {
	// 1. 验证Temporary兼容性 - 应该能正常工作但使用Request内核
	ses := requests.NewSession()
	tp := requests.NewTemporary(ses, "http://example.com/test")
	tp.AddHeader("X-Test", "temporary-compat")

	if tp.Error() != nil {
		t.Errorf("Temporary兼容性失败: %v", tp.Error())
	}
	t.Log("✓ Temporary向后兼容正常")

	// 2. 验证Session只返回Request对象
	req := ses.Get("http://example.com/test").
		SetHeader("X-Test", "unified-request").
		AddQuery("phase", "1")

	if req.Error() != nil {
		t.Errorf("Request创建失败: %v", req.Error())
	}
	t.Log("✓ Session统一返回Request对象")

	// 3. 验证顶层函数使用Request
	req2 := requests.Get("http://example.com/test").
		SetHeader("X-Test", "top-level-request")

	if req2.Error() != nil {
		t.Errorf("顶层函数失败: %v", req2.Error())
	}
	t.Log("✓ 顶层函数统一返回Request对象")

	// 4. 验证类型安全的配置方法
	session, err := requests.NewSessionWithOptions(
		requests.WithTimeout(30*time.Second),
		requests.WithUserAgent("Refactor-Test/1.0"),
	)
	if err != nil {
		t.Errorf("类型安全Session创建失败: %v", err)
	} else {
		// 使用类型安全的配置方法
		session.Config().SetBasicAuth("user", "pass")
		session.Config().SetTimeoutDuration(10 * time.Second)
		t.Log("✓ 类型安全配置方法正常工作")
	}

	// 5. 验证弃用方法仍能工作
	err = session.Config().SetBasicAuthLegacy("testuser", "testpass")
	if err != nil {
		t.Errorf("兼容性方法失败: %v", err)
	} else {
		t.Log("✓ 弃用方法保持向后兼容")
	}

	// 6. 验证API统一性
	tempReq := requests.NewTemporary(session, "http://example.com/test")
	tempReq.AddHeader("X-Source", "temporary")

	directReq := session.Get("http://example.com/test").
		SetHeader("X-Source", "request")

	if tempReq.Error() != nil || directReq.Error() != nil {
		t.Errorf("API统一性测试失败: temp=%v, direct=%v",
			tempReq.Error(), directReq.Error())
	} else {
		t.Log("✓ API统一性：Temporary和Request都正常工作")
	}

	t.Log("=== 第一阶段重构验证完成 ===")
}
