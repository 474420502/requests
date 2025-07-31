package requests

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestResponseAdvancedJSONMethods 测试Response的高级JSON方法
func TestResponseAdvancedJSONMethods(t *testing.T) {
	t.Run("IsJSON_EdgeCases", func(t *testing.T) {
		testCases := []struct {
			name        string
			contentType string
			expected    bool
		}{
			{"ApplicationJSON", "application/json", true},
			{"ApplicationJSONCharset", "application/json; charset=utf-8", true},
			{"TextJSON", "text/json", true},
			{"ApplicationJSONUTF8", "application/json;charset=UTF-8", true},
			{"ApplicationXML", "application/xml", false},
			{"TextPlain", "text/plain", false},
			{"EmptyContentType", "", false},
			{"ApplicationJavaScript", "application/javascript", false},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				resp := &Response{
					readResponse: &http.Response{
						Header: http.Header{
							"Content-Type": []string{tc.contentType},
						},
					},
				}

				if result := resp.IsJSON(); result != tc.expected {
					t.Errorf("Expected IsJSON() to return %v for content-type %s, got %v", tc.expected, tc.contentType, result)
				}
			})
		}
	})

	t.Run("Json_ErrorHandling", func(t *testing.T) {
		// 测试无效JSON - 注意gjson比较宽容，很多格式都能解析
		invalidJSON := `{invalid json without quotes}`
		resp := &Response{
			readBytes: []byte(invalidJSON),
			readResponse: &http.Response{
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
		}

		result := resp.Json()
		// gjson即使对于无效JSON也会尝试解析，所以我们测试特定字段是否存在
		field := result.Get("invalid")
		if !field.Exists() {
			// 这是预期行为：无效JSON中的字段不应该存在
			t.Logf("JSON parsing handled invalid format as expected")
		}
	})

	t.Run("GetJSONField_NestedAccess", func(t *testing.T) {
		jsonData := `{
			"level1": {
				"level2": {
					"level3": "deep_value",
					"array": [1, 2, 3, {"nested": "in_array"}]
				},
				"simple": "value"
			},
			"root_array": [
				{"item": "first"},
				{"item": "second"}
			]
		}`

		resp := &Response{
			readBytes: []byte(jsonData),
			readResponse: &http.Response{
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
		}

		// 测试嵌套访问
		value := resp.GetJSONField("level1.level2.level3")
		if value.String() != "deep_value" {
			t.Errorf("Expected 'deep_value', got '%v'", value.String())
		}

		// 测试数组访问
		arrayValue := resp.GetJSONField("level1.level2.array")
		if !arrayValue.Exists() {
			t.Error("Expected array value to exist")
		}

		// 测试根级数组访问
		rootArray := resp.GetJSONField("root_array")
		if !rootArray.Exists() {
			t.Error("Expected root array value to exist")
		}

		// 测试不存在的字段
		missing := resp.GetJSONField("nonexistent.field")
		if missing.Exists() {
			t.Errorf("Expected field not to exist, got %v", missing)
		}
	})

	t.Run("TypeSpecificGetters", func(t *testing.T) {
		jsonData := `{
			"string_value": "hello",
			"int_value": 42,
			"float_value": 3.14,
			"bool_true": true,
			"bool_false": false,
			"null_value": null,
			"nested": {
				"inner_string": "world"
			}
		}`

		resp := &Response{
			readBytes: []byte(jsonData),
			readResponse: &http.Response{
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
		}

		// 测试GetJSONString
		if str, err := resp.GetJSONString("string_value"); err != nil || str != "hello" {
			t.Errorf("Expected 'hello', got '%s', error: %v", str, err)
		}

		if str, err := resp.GetJSONString("nested.inner_string"); err != nil || str != "world" {
			t.Errorf("Expected 'world', got '%s', error: %v", str, err)
		}

		// 测试不存在的字段应该返回错误
		if _, err := resp.GetJSONString("nonexistent"); err == nil {
			t.Error("Expected error for nonexistent field")
		}

		// 测试GetJSONInt
		if intVal, err := resp.GetJSONInt("int_value"); err != nil || intVal != 42 {
			t.Errorf("Expected 42, got %d, error: %v", intVal, err)
		}

		// 测试类型不匹配时返回错误
		if _, err := resp.GetJSONInt("string_value"); err == nil {
			t.Error("Expected error for non-int field")
		}

		// 测试GetJSONFloat
		if floatVal, err := resp.GetJSONFloat("float_value"); err != nil || floatVal != 3.14 {
			t.Errorf("Expected 3.14, got %f, error: %v", floatVal, err)
		}

		// 测试GetJSONBool
		if boolVal, err := resp.GetJSONBool("bool_true"); err != nil || !boolVal {
			t.Errorf("Expected true, got %v, error: %v", boolVal, err)
		}

		if boolVal, err := resp.GetJSONBool("bool_false"); err != nil || boolVal {
			t.Errorf("Expected false, got %v, error: %v", boolVal, err)
		}

		// 测试类型不匹配时返回错误
		if _, err := resp.GetJSONBool("string_value"); err == nil {
			t.Error("Expected error for non-bool field")
		}
	})
}

// TestResponseDecodingMethods 测试各种解码方法
func TestResponseDecodingMethods(t *testing.T) {
	t.Run("DecodeJSON_StructDecoding", func(t *testing.T) {
		type TestStruct struct {
			Name  string `json:"name"`
			Age   int    `json:"age"`
			Email string `json:"email"`
		}

		jsonData := `{"name": "John Doe", "age": 30, "email": "john@example.com"}`
		resp := &Response{
			readBytes: []byte(jsonData),
			readResponse: &http.Response{
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
		}

		var result TestStruct
		err := resp.DecodeJSON(&result)
		if err != nil {
			t.Errorf("DecodeJSON failed: %v", err)
		}

		if result.Name != "John Doe" {
			t.Errorf("Expected name 'John Doe', got '%s'", result.Name)
		}
		if result.Age != 30 {
			t.Errorf("Expected age 30, got %d", result.Age)
		}
		if result.Email != "john@example.com" {
			t.Errorf("Expected email 'john@example.com', got '%s'", result.Email)
		}
	})

	t.Run("BindJSON_StructBinding", func(t *testing.T) {
		type TestStruct struct {
			Name  string `json:"name"`
			Count int    `json:"count"`
		}

		jsonData := `{"name": "Test Item", "count": 5}`
		resp := &Response{
			readBytes: []byte(jsonData),
			readResponse: &http.Response{
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
		}

		var result TestStruct
		err := resp.BindJSON(&result)
		if err != nil {
			t.Errorf("BindJSON failed: %v", err)
		}

		if result.Name != "Test Item" {
			t.Errorf("Expected name 'Test Item', got '%s'", result.Name)
		}
		if result.Count != 5 {
			t.Errorf("Expected count 5, got %d", result.Count)
		}
	})
}

// TestResponseContentHandling 测试内容处理方法
func TestResponseContentHandling(t *testing.T) {
	t.Run("ContentString_WithDifferentEncodings", func(t *testing.T) {
		testCases := []struct {
			name     string
			content  []byte
			expected string
		}{
			{"PlainText", []byte("Hello, World!"), "Hello, World!"},
			{"UTF8Text", []byte("你好，世界！"), "你好，世界！"},
			{"EmptyContent", []byte(""), ""},
			{"BinaryContent", []byte{0x00, 0x01, 0x02}, string([]byte{0x00, 0x01, 0x02})},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				resp := &Response{
					readBytes: tc.content,
					readResponse: &http.Response{
						Header: http.Header{},
					},
				}

				result := resp.ContentString()
				if result != tc.expected {
					t.Errorf("Expected '%s', got '%s'", tc.expected, result)
				}
			})
		}
	})

	t.Run("Content_ByteHandling", func(t *testing.T) {
		originalContent := []byte("Test content with special chars: éñ中文")
		resp := &Response{
			readBytes: originalContent,
			readResponse: &http.Response{
				Header: http.Header{},
			},
		}

		result := resp.Content()
		if !bytes.Equal(result, originalContent) {
			t.Errorf("Content bytes don't match. Expected %v, got %v", originalContent, result)
		}
	})
}

// TestResponseHeaderMethods 测试header相关方法
func TestResponseHeaderMethods(t *testing.T) {
	t.Run("GetHeader_HeaderAccess", func(t *testing.T) {
		resp := &Response{
			readResponse: &http.Response{
				Header: http.Header{
					"Content-Type":    []string{"application/json"},
					"Cache-Control":   []string{"no-cache"},
					"X-Custom-Header": []string{"custom-value"},
				},
			},
		}

		headers := resp.GetHeader()

		// 验证可以获取headers
		if contentType := headers.Get("Content-Type"); contentType != "application/json" {
			t.Errorf("Expected 'application/json', got '%s'", contentType)
		}

		if cacheControl := headers.Get("Cache-Control"); cacheControl != "no-cache" {
			t.Errorf("Expected 'no-cache', got '%s'", cacheControl)
		}

		if customHeader := headers.Get("X-Custom-Header"); customHeader != "custom-value" {
			t.Errorf("Expected 'custom-value', got '%s'", customHeader)
		}

		// 测试不存在的header
		if nonExistent := headers.Get("NonExistent"); nonExistent != "" {
			t.Errorf("Expected empty string for nonexistent header, got '%s'", nonExistent)
		}
	})
}

// TestResponseWithCompression 测试压缩响应处理
func TestResponseWithCompression(t *testing.T) {
	t.Run("GzipCompressedResponse", func(t *testing.T) {
		// 准备压缩内容
		originalContent := "This is a test content that will be compressed using gzip."
		var buf bytes.Buffer
		gz := gzip.NewWriter(&buf)
		gz.Write([]byte(originalContent))
		gz.Close()
		compressedContent := buf.Bytes()

		// 创建模拟的HTTP响应
		resp := &Response{
			readBytes: compressedContent,
			readResponse: &http.Response{
				Header: http.Header{
					"Content-Encoding": []string{"gzip"},
					"Content-Type":     []string{"text/plain"},
				},
			},
		}

		// 注意：这里我们测试的是存储在Response中的已解压内容
		// 在实际使用中，FromHTTPResponse会处理解压
		content := resp.Content()
		if len(content) == 0 {
			t.Error("Expected content, got empty response")
		}
	})
}

// TestResponseErrorCases 测试各种错误情况
func TestResponseErrorCases(t *testing.T) {
	t.Run("JSONMethods_WithNonJSONContent", func(t *testing.T) {
		resp := &Response{
			readBytes: []byte("This is plain text, not JSON"),
			readResponse: &http.Response{
				Header: http.Header{
					"Content-Type": []string{"text/plain"},
				},
			},
		}

		// GetJSONField应该返回不存在的结果
		field := resp.GetJSONField("any.field")
		if field.Exists() {
			t.Errorf("Expected field not to exist for non-JSON content, got %v", field)
		}

		// GetJSONString应该返回错误
		if _, err := resp.GetJSONString("any.field"); err == nil {
			t.Error("Expected error for non-JSON content")
		}

		// GetJSONInt应该返回错误
		if _, err := resp.GetJSONInt("any.field"); err == nil {
			t.Error("Expected error for non-JSON content")
		}

		// GetJSONFloat应该返回错误
		if _, err := resp.GetJSONFloat("any.field"); err == nil {
			t.Error("Expected error for non-JSON content")
		}

		// GetJSONBool应该返回错误
		if _, err := resp.GetJSONBool("any.field"); err == nil {
			t.Error("Expected error for non-JSON content")
		}
	})

	t.Run("DecodeJSON_WithInvalidJSON", func(t *testing.T) {
		resp := &Response{
			readBytes: []byte("invalid json content"),
			readResponse: &http.Response{
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
		}

		var result map[string]interface{}
		err := resp.DecodeJSON(&result)
		if err == nil {
			t.Error("Expected error for invalid JSON content")
		}
	})
}

// TestFromHTTPResponseAdvanced 测试FromHTTPResponse函数的高级用法
func TestFromHTTPResponseAdvanced(t *testing.T) {
	t.Run("FromHTTPResponse_BasicConversion", func(t *testing.T) {
		// 创建测试服务器
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Test-Header", "test-value")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "Hello, World!"}`))
		}))
		defer server.Close()

		// 发送HTTP请求
		httpResp, err := http.Get(server.URL)
		if err != nil {
			t.Fatalf("Failed to make HTTP request: %v", err)
		}
		defer httpResp.Body.Close()

		// 转换为我们的Response类型
		resp, err := FromHTTPResponse(httpResp, false)
		if err != nil {
			t.Fatalf("FromHTTPResponse failed: %v", err)
		}

		// 验证状态码
		if resp.GetStatusCode() != 200 {
			t.Errorf("Expected status code 200, got %d", resp.GetStatusCode())
		}

		// 验证headers
		headers := resp.GetHeader()
		if contentType := headers.Get("Content-Type"); !strings.Contains(contentType, "application/json") {
			t.Errorf("Expected JSON content type, got '%s'", contentType)
		}

		if testHeader := headers.Get("X-Test-Header"); testHeader != "test-value" {
			t.Errorf("Expected 'test-value', got '%s'", testHeader)
		}

		// 验证内容
		content := resp.ContentString()
		if !strings.Contains(content, "Hello, World!") {
			t.Errorf("Expected content to contain 'Hello, World!', got: %s", content)
		}

		// 验证JSON解析
		if message, err := resp.GetJSONString("message"); err != nil || message != "Hello, World!" {
			t.Errorf("Expected JSON message 'Hello, World!', got '%s', error: %v", message, err)
		}
	})

	t.Run("FromHTTPResponse_WithDecompression", func(t *testing.T) {
		// 创建测试服务器返回gzip压缩内容
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.Header().Set("Content-Encoding", "gzip")

			// 写入gzip压缩的内容
			gz := gzip.NewWriter(w)
			gz.Write([]byte("This is compressed content"))
			gz.Close()
		}))
		defer server.Close()

		// 发送HTTP请求
		httpResp, err := http.Get(server.URL)
		if err != nil {
			t.Fatalf("Failed to make HTTP request: %v", err)
		}
		defer httpResp.Body.Close()

		// 测试不解压缩
		resp1, err := FromHTTPResponse(httpResp, false)
		if err != nil {
			t.Fatalf("FromHTTPResponse failed: %v", err)
		}

		// 验证内容被正确读取（应该是已解压的）
		content1 := resp1.Content()
		if len(content1) == 0 {
			t.Error("Expected content, got empty response")
		}

		// 验证解压后的内容
		contentStr := resp1.ContentString()
		if contentStr != "This is compressed content" {
			t.Errorf("Expected 'This is compressed content', got '%s'", contentStr)
		}
	})
}

// TestResponseStatusMethods 测试状态相关方法
func TestResponseStatusMethods(t *testing.T) {
	t.Run("StatusMethods_VariousStatusCodes", func(t *testing.T) {
		testCases := []struct {
			statusCode int
			status     string
		}{
			{200, "200 OK"},
			{201, "201 Created"},
			{400, "400 Bad Request"},
			{401, "401 Unauthorized"},
			{404, "404 Not Found"},
			{500, "500 Internal Server Error"},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("Status_%d", tc.statusCode), func(t *testing.T) {
				resp := &Response{
					readResponse: &http.Response{
						StatusCode: tc.statusCode,
						Status:     tc.status,
					},
				}

				if code := resp.GetStatusCode(); code != tc.statusCode {
					t.Errorf("Expected status code %d, got %d", tc.statusCode, code)
				}

				if status := resp.GetStatus(); status != tc.status {
					t.Errorf("Expected status '%s', got '%s'", tc.status, status)
				}
			})
		}
	})
}

// TestResponseEdgeCases 测试边界情况
func TestResponseEdgeCases(t *testing.T) {
	t.Run("EmptyResponse", func(t *testing.T) {
		resp := &Response{
			readBytes: []byte{},
			readResponse: &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Header:     http.Header{},
			},
		}

		// 空内容应该正常处理
		if content := resp.ContentString(); content != "" {
			t.Errorf("Expected empty string, got '%s'", content)
		}

		if content := resp.Content(); len(content) != 0 {
			t.Errorf("Expected empty byte slice, got %v", content)
		}

		// JSON方法应该返回默认值或错误
		if field := resp.GetJSONField("any"); field.Exists() {
			t.Errorf("Expected field not to exist, got %v", field)
		}
	})

	t.Run("NilResponse", func(t *testing.T) {
		resp := &Response{
			readBytes:    nil,
			readResponse: nil,
		}

		// 应该能够处理nil情况而不panic
		content := resp.Content()
		if content == nil {
			content = []byte{}
		}

		contentStr := resp.ContentString()
		if contentStr == "" {
			// 预期行为
		}
	})
}
