package requests

import (
	"testing"
)

// TestBaseMethods 测试基础HTTP方法
func TestBaseMethods(t *testing.T) {
	// 测试Head方法
	t.Run("Head", func(t *testing.T) {
		req := Head("https://httpbin.org/get")
		if req == nil {
			t.Fatal("Head() returned nil")
		}
		if req.method != "HEAD" {
			t.Errorf("Expected method HEAD, got %s", req.method)
		}
		if req.parsedURL.String() != "https://httpbin.org/get" {
			t.Errorf("Expected URL https://httpbin.org/get, got %s", req.parsedURL.String())
		}
	})

	// 测试Get方法
	t.Run("Get", func(t *testing.T) {
		req := Get("https://httpbin.org/get")
		if req == nil {
			t.Fatal("Get() returned nil")
		}
		if req.method != "GET" {
			t.Errorf("Expected method GET, got %s", req.method)
		}
		if req.parsedURL.String() != "https://httpbin.org/get" {
			t.Errorf("Expected URL https://httpbin.org/get, got %s", req.parsedURL.String())
		}
	})

	// 测试Post方法
	t.Run("Post", func(t *testing.T) {
		req := Post("https://httpbin.org/post")
		if req == nil {
			t.Fatal("Post() returned nil")
		}
		if req.method != "POST" {
			t.Errorf("Expected method POST, got %s", req.method)
		}
		if req.parsedURL.String() != "https://httpbin.org/post" {
			t.Errorf("Expected URL https://httpbin.org/post, got %s", req.parsedURL.String())
		}
	})

	// 测试Put方法
	t.Run("Put", func(t *testing.T) {
		req := Put("https://httpbin.org/put")
		if req == nil {
			t.Fatal("Put() returned nil")
		}
		if req.method != "PUT" {
			t.Errorf("Expected method PUT, got %s", req.method)
		}
		if req.parsedURL.String() != "https://httpbin.org/put" {
			t.Errorf("Expected URL https://httpbin.org/put, got %s", req.parsedURL.String())
		}
	})

	// 测试Patch方法
	t.Run("Patch", func(t *testing.T) {
		req := Patch("https://httpbin.org/patch")
		if req == nil {
			t.Fatal("Patch() returned nil")
		}
		if req.method != "PATCH" {
			t.Errorf("Expected method PATCH, got %s", req.method)
		}
		if req.parsedURL.String() != "https://httpbin.org/patch" {
			t.Errorf("Expected URL https://httpbin.org/patch, got %s", req.parsedURL.String())
		}
	})

	// 测试Delete方法
	t.Run("Delete", func(t *testing.T) {
		req := Delete("https://httpbin.org/delete")
		if req == nil {
			t.Fatal("Delete() returned nil")
		}
		if req.method != "DELETE" {
			t.Errorf("Expected method DELETE, got %s", req.method)
		}
		if req.parsedURL.String() != "https://httpbin.org/delete" {
			t.Errorf("Expected URL https://httpbin.org/delete, got %s", req.parsedURL.String())
		}
	})

	// 测试Connect方法
	t.Run("Connect", func(t *testing.T) {
		req := Connect("https://httpbin.org/")
		if req == nil {
			t.Fatal("Connect() returned nil")
		}
		if req.method != "CONNECT" {
			t.Errorf("Expected method CONNECT, got %s", req.method)
		}
		if req.parsedURL.String() != "https://httpbin.org/" {
			t.Errorf("Expected URL https://httpbin.org/, got %s", req.parsedURL.String())
		}
	})

	// 测试Options方法
	t.Run("Options", func(t *testing.T) {
		req := Options("https://httpbin.org/")
		if req == nil {
			t.Fatal("Options() returned nil")
		}
		if req.method != "OPTIONS" {
			t.Errorf("Expected method OPTIONS, got %s", req.method)
		}
		if req.parsedURL.String() != "https://httpbin.org/" {
			t.Errorf("Expected URL https://httpbin.org/, got %s", req.parsedURL.String())
		}
	})

	// 测试Trace方法
	t.Run("Trace", func(t *testing.T) {
		req := Trace("https://httpbin.org/")
		if req == nil {
			t.Fatal("Trace() returned nil")
		}
		if req.method != "TRACE" {
			t.Errorf("Expected method TRACE, got %s", req.method)
		}
		if req.parsedURL.String() != "https://httpbin.org/" {
			t.Errorf("Expected URL https://httpbin.org/, got %s", req.parsedURL.String())
		}
	})
}
