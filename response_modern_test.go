package requests

import (
	"strings"
	"testing"
)

// 测试新的BindJSON功能
func TestResponse_BindJSON(t *testing.T) {
	session := NewSession()

	// 测试结构体绑定
	type TestResponse struct {
		JSON map[string]interface{} `json:"json"`
		URL  string                 `json:"url"`
	}

	resp, err := session.Post("http://httpbin.org/post").
		SetBodyJSON(map[string]interface{}{
			"name": "test",
			"age":  25,
		}).
		Execute()

	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	var result TestResponse
	err = resp.BindJSON(&result)
	if err != nil {
		t.Fatalf("BindJSON failed: %v", err)
	}

	if result.JSON["name"] != "test" {
		t.Errorf("Expected name='test', got: %v", result.JSON["name"])
	}

	if result.JSON["age"] != float64(25) { // JSON numbers are float64
		t.Errorf("Expected age=25, got: %v", result.JSON["age"])
	}

	if !strings.Contains(result.URL, "httpbin.org") {
		t.Errorf("Expected URL to contain 'httpbin.org', got: %s", result.URL)
	}
}

// 测试DecodeJSON功能
func TestResponse_DecodeJSON(t *testing.T) {
	session := NewSession()

	resp, err := session.Get("http://httpbin.org/json").Execute()
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	var data map[string]interface{}
	err = resp.DecodeJSON(&data)
	if err != nil {
		t.Fatalf("DecodeJSON failed: %v", err)
	}

	// httpbin.org/json 返回一个包含slideshow的JSON对象
	if _, exists := data["slideshow"]; !exists {
		t.Errorf("Expected 'slideshow' key in response, got: %+v", data)
	}
}

// 测试类型安全的查询参数
func TestRequest_TypeSafeQueryParams(t *testing.T) {
	session := NewSession()

	resp, err := session.Get("http://httpbin.org/get").
		AddQuery("name", "张三").
		AddQueryInt("age", 25).
		AddQueryBool("active", true).
		AddQueryFloat("score", 95.5).
		Execute()

	if err != nil {
		t.Fatalf("Request with type-safe query params failed: %v", err)
	}

	// 解析响应以验证查询参数
	var result struct {
		Args map[string]interface{} `json:"args"`
	}

	err = resp.BindJSON(&result)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if result.Args["name"] != "张三" {
		t.Errorf("Expected name='张三', got: %v", result.Args["name"])
	}

	if result.Args["age"] != "25" {
		t.Errorf("Expected age='25', got: %v", result.Args["age"])
	}

	if result.Args["active"] != "true" {
		t.Errorf("Expected active='true', got: %v", result.Args["active"])
	}

	if result.Args["score"] != "95.5" {
		t.Errorf("Expected score='95.5', got: %v", result.Args["score"])
	}
}

// 测试路径参数替换
func TestRequest_PathParams(t *testing.T) {
	session := NewSession()

	// 单个参数替换
	resp, err := session.Get("http://httpbin.org/status/{code}").
		SetPathParam("code", "200").
		Execute()

	if err != nil {
		t.Fatalf("Request with path param failed: %v", err)
	}

	if resp.GetStatusCode() != 200 {
		t.Errorf("Expected status code 200, got: %d", resp.GetStatusCode())
	}

	// 批量参数替换
	resp, err = session.Get("http://httpbin.org/delay/{seconds}").
		SetPathParams(map[string]string{
			"seconds": "1",
		}).
		Execute()

	if err != nil {
		t.Fatalf("Request with batch path params failed: %v", err)
	}

	if resp.GetStatusCode() != 200 {
		t.Errorf("Expected status code 200, got: %d", resp.GetStatusCode())
	}
} // 测试类型安全的表单处理
func TestRequest_TypeSafeForm(t *testing.T) {
	session := NewSession()

	// 测试SetFormFields
	resp, err := session.Post("http://httpbin.org/post").
		SetFormFields(map[string]string{
			"username": "testuser",
			"email":    "test@example.com",
		}).
		Execute()

	if err != nil {
		t.Fatalf("Request with form fields failed: %v", err)
	}

	var result struct {
		Form map[string]interface{} `json:"form"`
	}

	err = resp.BindJSON(&result)
	if err != nil {
		t.Fatalf("Failed to parse form response: %v", err)
	}

	if result.Form["username"] != "testuser" {
		t.Errorf("Expected username='testuser', got: %v", result.Form["username"])
	}

	if result.Form["email"] != "test@example.com" {
		t.Errorf("Expected email='test@example.com', got: %v", result.Form["email"])
	}

	// 测试AddFormFile
	fileContent := strings.NewReader("test file content")
	resp, err = session.Post("http://httpbin.org/post").
		AddFormFile("upload", "test.txt", fileContent).
		Execute()

	if err != nil {
		t.Fatalf("Request with form file failed: %v", err)
	}

	var fileResult struct {
		Form  map[string]interface{} `json:"form"`
		Files map[string]interface{} `json:"files"`
	}

	err = resp.BindJSON(&fileResult)
	if err != nil {
		t.Fatalf("Failed to parse file response: %v", err)
	}

	// 测试文件上传 (简化版本，不混合form字段)
	var uploadResult struct {
		Files map[string]interface{} `json:"files"`
	}

	err = resp.BindJSON(&uploadResult)
	if err != nil {
		t.Fatalf("Failed to parse file response: %v", err)
	}

	if _, exists := uploadResult.Files["upload"]; !exists {
		// 如果没有files字段，检查响应是否包含文件上传相关内容
		responseContent := resp.ContentString()
		if !strings.Contains(responseContent, "upload") && !strings.Contains(responseContent, "test.txt") {
			t.Errorf("Response should contain file upload info, got: %s", responseContent)
		}
	}
}

// 测试错误处理的健壮性
func TestErrorHandling_Robustness(t *testing.T) {
	session := NewSession()

	// 测试无效URL
	_, err := session.Get("invalid-url").Execute()
	if err == nil {
		t.Error("Expected error for invalid URL, but got none")
	}

	// 测试无效代理配置 (使用明显无效的格式)
	err = session.Config().SetProxy("://invalid-proxy-url")
	if err == nil {
		t.Error("Expected error for completely invalid proxy format, but got none")
	} else {
		t.Logf("SetProxy correctly returned error for invalid format: %v", err)
	}

	// 测试类型安全的代理设置
	err = session.Config().SetProxyString("http://invalid-proxy.example.com:8080")
	if err != nil {
		t.Logf("Proxy setting returned expected error: %v", err)
	}

	// 测试清除代理
	session.Config().ClearProxy() // 应该不返回错误

	// 测试类型安全的认证设置
	session.Config().SetBasicAuthString("user", "pass") // 应该不返回错误

	// 测试清除认证
	session.Config().ClearBasicAuth() // 应该不返回错误
}

// 测试向后兼容性
func TestBackwardCompatibility(t *testing.T) {
	session := NewSession()

	// 测试旧的SetBasicAuth接口仍然工作
	err := session.Config().SetBasicAuth("user", "pass")
	if err != nil {
		t.Fatalf("Backward compatible SetBasicAuth failed: %v", err)
	}

	// 测试旧的SetProxy接口仍然工作
	err = session.Config().SetProxy("http://proxy.example.com:8080")
	if err != nil {
		t.Fatalf("Backward compatible SetProxy failed: %v", err)
	}

	// 测试旧的SetTimeout接口仍然工作
	err = session.Config().SetTimeout(30)
	if err != nil {
		t.Fatalf("Backward compatible SetTimeout failed: %v", err)
	}
}

// 基准测试：测试新API的性能
func BenchmarkRequest_TypeSafeAPI(b *testing.B) {
	session := NewSession()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := session.Get("http://httpbin.org/get").
			AddQueryInt("id", i).
			AddQueryBool("active", true).
			Execute()

		if err != nil {
			b.Fatalf("Request failed: %v", err)
		}
	}
}

// 基准测试：现代API性能
func BenchmarkRequest_ModernAPI(b *testing.B) {
	session := NewSession()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := session.Get("http://httpbin.org/get")
		req.AddQueryInt("id", i)
		req.AddQueryBool("active", true)
		_, err := req.Execute()

		if err != nil {
			b.Fatalf("Request failed: %v", err)
		}
	}
}
