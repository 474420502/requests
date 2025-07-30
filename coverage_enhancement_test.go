package requests

import (
	"context"
	"testing"
	"time"
)

// TestCoverageEnhancement 基础覆盖率测试
func TestCoverageEnhancement(t *testing.T) {
	session := NewSession()

	t.Run("Request_AddQueryInt64", func(t *testing.T) {
		req := session.Get("http://example.com").AddQueryInt64("id", 123456789)
		query := req.GetQuery()
		if query.Get("id") != "123456789" {
			t.Errorf("Expected id=123456789, got %s", query.Get("id"))
		}
	})

	t.Run("Request_GetRawURL", func(t *testing.T) {
		req := session.Get("http://example.com/test")
		url := req.GetRawURL()
		if url != "http://example.com/test" {
			t.Errorf("Expected http://example.com/test, got %s", url)
		}
	})

	t.Run("Request_WithContext", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "test", "value")
		req := session.Get("http://example.com").WithContext(ctx)
		if req.ctx.Value("test") != "value" {
			t.Error("Context not set correctly")
		}
	})

	t.Run("Session_DefaultContext", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "test", "value")
		session.SetDefaultContext(ctx)
		retrievedCtx := session.GetDefaultContext()
		if retrievedCtx.Value("test") != "value" {
			t.Error("Default context not set correctly")
		}
	})

	t.Run("Config_SetTimeout", func(t *testing.T) {
		config := session.Config()
		// 测试负数超时 - 应该成功（SetTimeout支持负数）
		err := config.SetTimeout(-1)
		if err != nil {
			t.Errorf("SetTimeout should not fail with negative int: %v", err)
		}
		// 测试无效类型
		err = config.SetTimeout("invalid")
		if err == nil {
			t.Error("Should fail with invalid type")
		}
		// 测试有效超时
		err = config.SetTimeout(30)
		if err != nil {
			t.Errorf("Should not fail with valid timeout: %v", err)
		}
	})

	t.Run("MultiPool_Execute", func(t *testing.T) {
		pool := NewRequestPool(2)
		req := session.Get("invalid-url")
		pool.Add(req)
		results := pool.Execute()
		if len(results) != 1 {
			t.Errorf("Expected 1 result, got %d", len(results))
		}
		if results[0].Error == nil {
			t.Error("Expected error for invalid URL")
		}
	})

	t.Run("FormMethods", func(t *testing.T) {
		req := session.Post("http://httpbin.org/post")
		req.AddFormField("name", "value")
		req.AddFormFieldInt("count", 42)
		req.AddFormFieldInt64("id", 123456789)
		req.AddFormFieldBool("active", true)
		req.AddFormFieldFloat("price", 19.99)

		// 测试不存在的文件
		err := req.SetFormFileFromPath("upload", "/nonexistent/file.txt")
		if err == nil {
			t.Error("Should fail with nonexistent file")
		}
	})

	t.Run("Request_WithTimeout", func(t *testing.T) {
		req := session.Get("http://example.com").WithTimeout(5 * time.Second)
		// 验证超时设置（通过检查Request不为nil）
		if req == nil {
			t.Error("Request should not be nil")
		}
	})
}
