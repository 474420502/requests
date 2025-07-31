package requests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestResponse_JSONMethods 测试Response的JSON相关方法
func TestResponse_JSONMethods(t *testing.T) {
	// 创建返回JSON的测试服务器
	jsonServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{
			"name": "test",
			"age": 25,
			"active": true,
			"score": 98.5,
			"nested": {
				"field": "value"
			},
			"items": ["item1", "item2"]
		}`))
	}))
	defer jsonServer.Close()

	// 创建返回非JSON的测试服务器
	textServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(200)
		w.Write([]byte("plain text response"))
	}))
	defer textServer.Close()

	session := NewSession()

	t.Run("IsJSON_True", func(t *testing.T) {
		resp, err := session.Get(jsonServer.URL).Execute()
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		if !resp.IsJSON() {
			t.Error("Expected response to be JSON")
		}
	})

	t.Run("IsJSON_False", func(t *testing.T) {
		resp, err := session.Get(textServer.URL).Execute()
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		if resp.IsJSON() {
			t.Error("Expected response to not be JSON")
		}
	})

	t.Run("Json", func(t *testing.T) {
		resp, err := session.Get(jsonServer.URL).Execute()
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		jsonResult := resp.Json()
		if !jsonResult.Exists() {
			t.Error("Expected JSON result to exist")
		}

		name := jsonResult.Get("name").String()
		if name != "test" {
			t.Errorf("Expected name 'test', got '%s'", name)
		}
	})

	t.Run("GetJSONField", func(t *testing.T) {
		resp, err := session.Get(jsonServer.URL).Execute()
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		result := resp.GetJSONField("name")
		if !result.Exists() {
			t.Error("Expected field 'name' to exist")
		}
		if result.String() != "test" {
			t.Errorf("Expected name 'test', got '%s'", result.String())
		}

		// 测试嵌套字段
		nestedResult := resp.GetJSONField("nested.field")
		if nestedResult.String() != "value" {
			t.Errorf("Expected nested field 'value', got '%s'", nestedResult.String())
		}

		// 测试不存在的字段
		nonExistentResult := resp.GetJSONField("nonexistent")
		if nonExistentResult.Exists() {
			t.Error("Expected non-existent field to not exist")
		}
	})

	t.Run("GetJSONString", func(t *testing.T) {
		resp, err := session.Get(jsonServer.URL).Execute()
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		name, err := resp.GetJSONString("name")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if name != "test" {
			t.Errorf("Expected name 'test', got '%s'", name)
		}

		// 测试不存在的字段
		_, err = resp.GetJSONString("nonexistent")
		if err == nil {
			t.Error("Expected error for non-existent field")
		}
	})

	t.Run("GetJSONInt", func(t *testing.T) {
		resp, err := session.Get(jsonServer.URL).Execute()
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		age, err := resp.GetJSONInt("age")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if age != 25 {
			t.Errorf("Expected age 25, got %d", age)
		}

		// 测试不存在的字段
		_, err = resp.GetJSONInt("nonexistent")
		if err == nil {
			t.Error("Expected error for non-existent field")
		}

		// 测试类型不匹配
		_, err = resp.GetJSONInt("name")
		if err == nil {
			t.Error("Expected error for type mismatch")
		}
	})

	t.Run("GetJSONFloat", func(t *testing.T) {
		resp, err := session.Get(jsonServer.URL).Execute()
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		score, err := resp.GetJSONFloat("score")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if score != 98.5 {
			t.Errorf("Expected score 98.5, got %f", score)
		}

		// 测试不存在的字段
		_, err = resp.GetJSONFloat("nonexistent")
		if err == nil {
			t.Error("Expected error for non-existent field")
		}

		// 测试类型不匹配
		_, err = resp.GetJSONFloat("name")
		if err == nil {
			t.Error("Expected error for type mismatch")
		}
	})

	t.Run("GetJSONBool", func(t *testing.T) {
		resp, err := session.Get(jsonServer.URL).Execute()
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		active, err := resp.GetJSONBool("active")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !active {
			t.Error("Expected active to be true")
		}

		// 测试不存在的字段
		_, err = resp.GetJSONBool("nonexistent")
		if err == nil {
			t.Error("Expected error for non-existent field")
		}

		// 测试类型不匹配
		_, err = resp.GetJSONBool("name")
		if err == nil {
			t.Error("Expected error for type mismatch")
		}
	})

	t.Run("DecodeJSON", func(t *testing.T) {
		resp, err := session.Get(jsonServer.URL).Execute()
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		var data map[string]interface{}
		err = resp.DecodeJSON(&data)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if data["name"] != "test" {
			t.Errorf("Expected name 'test', got %v", data["name"])
		}

		// 测试类型转换
		age, ok := data["age"].(float64) // JSON数字解析为float64
		if !ok || age != 25 {
			t.Errorf("Expected age 25, got %v", data["age"])
		}
	})

	t.Run("BindJSON", func(t *testing.T) {
		resp, err := session.Get(jsonServer.URL).Execute()
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		type TestStruct struct {
			Name   string  `json:"name"`
			Age    int     `json:"age"`
			Active bool    `json:"active"`
			Score  float64 `json:"score"`
		}

		var data TestStruct
		err = resp.BindJSON(&data)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if data.Name != "test" {
			t.Errorf("Expected name 'test', got '%s'", data.Name)
		}
		if data.Age != 25 {
			t.Errorf("Expected age 25, got %d", data.Age)
		}
		if !data.Active {
			t.Error("Expected active to be true")
		}
		if data.Score != 98.5 {
			t.Errorf("Expected score 98.5, got %f", data.Score)
		}
	})
}

// TestResponse_StatusMethods 测试Response的状态相关方法
func TestResponse_StatusMethods(t *testing.T) {
	// 创建不同状态码的测试服务器
	statusServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		status := r.URL.Query().Get("status")
		switch status {
		case "200":
			w.WriteHeader(200)
		case "404":
			w.WriteHeader(404)
		case "500":
			w.WriteHeader(500)
		default:
			w.WriteHeader(200)
		}
		w.Write([]byte("response"))
	}))
	defer statusServer.Close()

	session := NewSession()

	t.Run("GetStatusCode", func(t *testing.T) {
		resp, err := session.Get(statusServer.URL + "?status=404").Execute()
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		if resp.GetStatusCode() != 404 {
			t.Errorf("Expected status code 404, got %d", resp.GetStatusCode())
		}
	})

	t.Run("GetStatus", func(t *testing.T) {
		resp, err := session.Get(statusServer.URL + "?status=200").Execute()
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		status := resp.GetStatus()
		if !strings.Contains(status, "200") {
			t.Errorf("Expected status to contain '200', got '%s'", status)
		}
	})
}

// TestResponse_HeaderMethods 测试Response的Header相关方法
func TestResponse_HeaderMethods(t *testing.T) {
	headerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Custom-Header", "custom-value")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"message": "test"}`))
	}))
	defer headerServer.Close()

	session := NewSession()

	t.Run("GetHeader", func(t *testing.T) {
		resp, err := session.Get(headerServer.URL).Execute()
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		headers := resp.GetHeader()
		if headers == nil {
			t.Error("Expected headers to not be nil")
		}

		customHeader := headers.Get("X-Custom-Header")
		if customHeader != "custom-value" {
			t.Errorf("Expected custom header 'custom-value', got '%s'", customHeader)
		}

		contentType := headers.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			t.Errorf("Expected Content-Type to contain 'application/json', got '%s'", contentType)
		}

		nonExistentHeader := headers.Get("X-Non-Existent")
		if nonExistentHeader != "" {
			t.Errorf("Expected empty header for non-existent header, got '%s'", nonExistentHeader)
		}
	})
}

// TestResponse_BodyMethods 测试Response的Body相关方法
func TestResponse_BodyMethods(t *testing.T) {
	bodyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(200)
		w.Write([]byte("test response body"))
	}))
	defer bodyServer.Close()

	session := NewSession()

	t.Run("ContentString", func(t *testing.T) {
		resp, err := session.Get(bodyServer.URL).Execute()
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		text := resp.ContentString()
		if text != "test response body" {
			t.Errorf("Expected text 'test response body', got '%s'", text)
		}
	})

	t.Run("Content", func(t *testing.T) {
		resp, err := session.Get(bodyServer.URL).Execute()
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		bytes := resp.Content()
		if string(bytes) != "test response body" {
			t.Errorf("Expected bytes 'test response body', got '%s'", string(bytes))
		}
	})
}

// TestResponse_EdgeCases 测试Response的边界情况
func TestResponse_EdgeCases(t *testing.T) {
	// 创建返回空响应的服务器
	emptyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204) // No Content
	}))
	defer emptyServer.Close()

	// 创建返回无效JSON的服务器
	invalidJSONServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"invalid": json}`)) // 故意的无效JSON
	}))
	defer invalidJSONServer.Close()

	session := NewSession()

	t.Run("EmptyResponse", func(t *testing.T) {
		resp, err := session.Get(emptyServer.URL).Execute()
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		text := resp.ContentString()
		if text != "" {
			t.Errorf("Expected empty text, got '%s'", text)
		}

		bytes := resp.Content()
		if len(bytes) != 0 {
			t.Errorf("Expected empty bytes, got %d bytes", len(bytes))
		}
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		resp, err := session.Get(invalidJSONServer.URL).Execute()
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		var data map[string]interface{}
		err = resp.DecodeJSON(&data)
		if err == nil {
			t.Error("Expected error for invalid JSON")
		}

		err = resp.BindJSON(&data)
		if err == nil {
			t.Error("Expected error for invalid JSON")
		}
	})
}
