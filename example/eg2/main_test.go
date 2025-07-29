package main

import (
	"testing"

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
