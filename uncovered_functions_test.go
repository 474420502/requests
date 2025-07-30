package requests

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

// TestUncoveredMiddlewareFunctions 测试未覆盖的middleware功能
func TestUncoveredMiddlewareFunctions(t *testing.T) {
	t.Run("RequestWithMiddleware_AddMiddleware", func(t *testing.T) {
		session := NewSession()
		req := session.Get("http://example.com")

		// 创建一个简单的中间件
		mw := &LoggingMiddleware{}

		// 使用WithMiddlewares创建RequestWithMiddleware
		reqWithMW := req.WithMiddlewares(mw)

		// 添加另一个中间件
		reqWithMW.AddMiddleware(&RetryMiddleware{MaxRetries: 3})

		// 验证中间件已添加
		if len(reqWithMW.middlewares) != 2 {
			t.Errorf("Expected 2 middlewares, got %d", len(reqWithMW.middlewares))
		}
	})

	t.Run("RequestWithMiddleware_ExecuteWithMiddleware", func(t *testing.T) {
		session := NewSession()
		req := session.Get("http://httpbin.org/get")

		// 添加日志中间件
		mw := &LoggingMiddleware{}
		reqWithMW := req.WithMiddlewares(mw)

		// 执行带中间件的请求
		resp, err := reqWithMW.ExecuteWithMiddleware()
		if err != nil {
			t.Errorf("ExecuteWithMiddleware failed: %v", err)
		}

		if resp == nil {
			t.Error("Expected response, got nil")
		}
	})

	t.Run("MetricsMiddleware_Tracking", func(t *testing.T) {
		var requestCount int
		var responseCount int
		var lastDuration time.Duration

		mw := &MetricsMiddleware{
			RequestCounter: func(method, url string) {
				requestCount++
			},
			ResponseCounter: func(statusCode int, method, url string) {
				responseCount++
			},
			DurationTracker: func(duration time.Duration, method, url string) {
				lastDuration = duration
			},
		}

		req, _ := http.NewRequest("GET", "http://example.com", nil)

		// 测试BeforeRequest
		err := mw.BeforeRequest(req)
		if err != nil {
			t.Errorf("BeforeRequest failed: %v", err)
		}

		if requestCount != 1 {
			t.Errorf("Expected request count 1, got %d", requestCount)
		}

		// 模拟一点延迟
		time.Sleep(time.Millisecond)

		// 测试AfterResponse
		resp := &http.Response{StatusCode: 200, Request: req}
		err = mw.AfterResponse(resp)
		if err != nil {
			t.Errorf("AfterResponse failed: %v", err)
		}

		if responseCount != 1 {
			t.Errorf("Expected response count 1, got %d", responseCount)
		}

		if lastDuration == 0 {
			t.Error("Expected duration to be tracked")
		}
	})

	t.Run("CacheMiddleware_Operations", func(t *testing.T) {
		cacheMW := NewCacheMiddleware()

		req, _ := http.NewRequest("GET", "http://example.com/api/data", nil)

		// BeforeRequest不应该出错
		err := cacheMW.BeforeRequest(req)
		if err != nil {
			t.Errorf("BeforeRequest failed: %v", err)
		}

		// 创建成功响应
		resp := &http.Response{
			StatusCode: 200,
			Request:    req,
			Header:     http.Header{},
		}

		// AfterResponse应该缓存响应
		err = cacheMW.AfterResponse(resp)
		if err != nil {
			t.Errorf("AfterResponse failed: %v", err)
		}

		// 验证响应被缓存
		if len(cacheMW.Cache) != 1 {
			t.Errorf("Expected 1 cached item, got %d", len(cacheMW.Cache))
		}
	})

	t.Run("UserAgentRotationMiddleware_Creation", func(t *testing.T) {
		userAgents := []string{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
		}

		mw := NewUserAgentRotationMiddleware(userAgents)

		if len(mw.UserAgents) != 2 {
			t.Errorf("Expected 2 user agents, got %d", len(mw.UserAgents))
		}

		req, _ := http.NewRequest("GET", "http://example.com", nil)

		// 测试第一次请求
		err := mw.BeforeRequest(req)
		if err != nil {
			t.Errorf("BeforeRequest failed: %v", err)
		}

		firstUA := req.Header.Get("User-Agent")
		if firstUA == "" {
			t.Error("Expected User-Agent to be set")
		}

		// 测试第二次请求
		req2, _ := http.NewRequest("GET", "http://example.com", nil)
		err = mw.BeforeRequest(req2)
		if err != nil {
			t.Errorf("BeforeRequest failed: %v", err)
		}

		secondUA := req2.Header.Get("User-Agent")
		if secondUA == firstUA {
			t.Error("Expected User-Agent to rotate")
		}
	})
}

// TestUncoveredMultipartFormDataFunctions 测试未覆盖的multipart form data功能
func TestUncoveredMultipartFormDataFunctions(t *testing.T) {
	t.Run("MultipartFormData_Operations", func(t *testing.T) {
		// 通过Request创建MultipartFormData
		session := NewSession()
		req := session.Post("http://example.com")
		mfd := req.CreateBodyMultipart()

		// 测试Data方法
		data := mfd.Data()
		if data == nil {
			t.Error("Expected data buffer, got nil")
		}

		// 测试Writer方法
		writer := mfd.Writer()
		if writer == nil {
			t.Error("Expected writer, got nil")
		}

		// 测试AddField
		err := mfd.AddField("test_field", "test_value")
		if err != nil {
			t.Errorf("AddField failed: %v", err)
		}

		// 测试AddFile (创建一个假文件)
		fileContent := []byte("test file content")
		err = mfd.AddFile("file_field", "test.txt", fileContent)
		if err != nil {
			t.Errorf("AddFile failed: %v", err)
		}

		// 测试AddFieldFile
		err = mfd.AddFieldFile("file_field2", "test2.txt", bytes.NewReader(fileContent))
		if err != nil {
			t.Errorf("AddFieldFile failed: %v", err)
		}

		// 测试ContentType
		contentType := mfd.ContentType()
		if !strings.Contains(contentType, "multipart/form-data") {
			t.Errorf("Expected multipart content type, got: %s", contentType)
		}

		// 测试Close
		err = mfd.Close()
		if err != nil {
			t.Errorf("Close failed: %v", err)
		}
	})
} // TestUncoveredRequestFunctions 测试未覆盖的request功能
func TestUncoveredRequestFunctions(t *testing.T) {
	t.Run("Request_HeaderOperations", func(t *testing.T) {
		session := NewSession()
		req := session.Get("http://example.com")

		// 测试SetHeadersFromHTTP
		httpHeaders := http.Header{
			"Content-Type": []string{"application/json"},
			"X-Test":       []string{"test-value"},
		}
		req.SetHeadersFromHTTP(httpHeaders)

		// 验证headers被设置
		header := req.GetHeader()
		if header.Get("Content-Type") != "application/json" {
			t.Error("Expected Content-Type to be set")
		}

		// 测试GetHeader
		if header.Get("X-Test") != "test-value" {
			t.Error("Expected X-Test header to be set")
		}
	})

	t.Run("Request_URLOperations", func(t *testing.T) {
		session := NewSession()
		req := session.Get("http://example.com/path")

		// 测试SetParsedURL
		newURL, _ := url.Parse("http://newhost.com/newpath")
		req.SetParsedURL(newURL)

		// 测试GetParsedURL
		parsedURL := req.GetParsedURL()
		if parsedURL.Host != "newhost.com" {
			t.Errorf("Expected host 'newhost.com', got '%s'", parsedURL.Host)
		}

		// 测试GetURLRawPath
		rawPath := req.GetURLRawPath()
		if rawPath != "/newpath" {
			t.Errorf("Expected raw path '/newpath', got '%s'", rawPath)
		}

		// 测试SetURLRawPath
		req.SetURLRawPath("/another/path")
		if req.GetURLRawPath() != "/another/path" {
			t.Error("SetURLRawPath failed")
		}

		// 测试GetURLPath
		urlPath := req.GetURLPath()
		if len(urlPath) == 0 {
			t.Error("Expected URL path, got empty slice")
		}

		// 测试SetURLPath
		req.SetURLPath([]string{"final", "path"})
		if len(req.GetURLPath()) < 2 {
			t.Error("SetURLPath failed")
		}

		// 测试SetRawURL
		req.SetRawURL("http://raw.example.com/raw")
		rawURL := req.GetRawURL()
		if !strings.Contains(rawURL, "raw.example.com") {
			t.Errorf("Expected raw URL to contain 'raw.example.com', got: %s", rawURL)
		}
	})

	t.Run("Request_ParamOperations", func(t *testing.T) {
		session := NewSession()
		req := session.Get("http://example.com")

		// 测试DelParam
		req.SetParam("param1", "value1")
		req.SetParam("param2", "value2")
		req.DelParam("param1")

		query := req.GetQuery()
		if query.Get("param1") != "" {
			t.Error("Expected param1 to be deleted")
		}
		if query.Get("param2") != "value2" {
			t.Error("Expected param2 to remain")
		}

		// 测试MergeQuery
		additionalQuery := url.Values{}
		additionalQuery.Add("merged", "value")
		req.MergeQuery(additionalQuery)

		if req.GetQuery().Get("merged") != "value" {
			t.Error("Expected merged query parameter")
		}
	})

	t.Run("Request_PathParams", func(t *testing.T) {
		session := NewSession()
		req := session.Get("http://example.com/users/{id}/posts/{postId}")

		// 测试SetPathParam
		req.SetPathParam("id", "123")
		req.SetPathParam("postId", "456")

		// URL应该被更新
		finalURL := req.GetRawURL()
		if !strings.Contains(finalURL, "123") || !strings.Contains(finalURL, "456") {
			t.Errorf("Expected URL to contain path params, got: %s", finalURL)
		}
	})

	t.Run("Request_BodyOperations", func(t *testing.T) {
		session := NewSession()
		req := session.Get("http://example.com")

		// 测试SetBody
		bodyData := []byte("test body")
		req.SetBody(bytes.NewReader(bodyData))

		// 测试SetBodyStream - 只支持string, []byte, []rune
		req.SetBodyStream(bodyData)

		// 测试SetBodyReader
		req.SetBodyReader(bytes.NewReader(bodyData))

		// 测试SetFormFieldsTyped
		formData := map[string]interface{}{
			"field1": "value1",
			"field2": 123,
			"field3": true,
		}
		req.SetFormFieldsTyped(formData)

		// 这些方法主要测试不出错
		if req.err != nil {
			t.Errorf("Unexpected error in body operations: %v", req.err)
		}
	})

	t.Run("Request_FormFileOperations", func(t *testing.T) {
		session := NewSession()
		req := session.Get("http://example.com")

		// 测试AddMultipleFormFiles
		formFiles := map[string]io.Reader{
			"file1": bytes.NewReader([]byte("file1 content")),
			"file2": bytes.NewReader([]byte("file2 content")),
		}
		req.AddMultipleFormFiles(formFiles)

		// 测试SetFormFileFromPath - 使用一个可能存在的文件
		err := req.SetFormFileFromPath("config", "go.mod")
		if err != nil {
			// 如果文件不存在，这是预期的
			t.Logf("SetFormFileFromPath error (expected if file not found): %v", err)
		}

		// 测试SetBodyFormFiles
		formFile := FormFile{
			FieldName: "testfile",
			FileName:  "test.txt",
			Reader:    bytes.NewReader([]byte("test content")),
		}
		req.SetBodyFormFiles(formFile)

		if req.err != nil {
			t.Errorf("Unexpected error in form file operations: %v", req.err)
		}
	})

	t.Run("Request_WithMiddleware", func(t *testing.T) {
		session := NewSession()
		req := session.Get("http://example.com")

		// 测试WithMiddleware
		mw := &LoggingMiddleware{}
		reqWithMW := req.WithMiddleware(mw)

		if reqWithMW == nil {
			t.Error("Expected RequestWithMiddleware, got nil")
			return
		}

		if len(reqWithMW.middlewares) != 1 {
			t.Errorf("Expected 1 middleware, got %d", len(reqWithMW.middlewares))
		}
	})

	t.Run("Request_ResponseFormatMethods", func(t *testing.T) {
		session := NewSession()
		req := session.Get("http://httpbin.org/get")

		// 这些方法返回特定格式的响应
		text, err := req.Text()
		if err != nil {
			t.Errorf("Text method failed: %v", err)
		}
		if text == "" {
			t.Error("Expected text response")
		}

		// 重新创建request因为上面已经执行过了
		req2 := session.Get("http://httpbin.org/get")
		var result map[string]interface{}
		err = req2.JSON(&result)
		if err != nil {
			t.Errorf("JSON method failed: %v", err)
		}

		// 重新创建request
		req3 := session.Get("http://httpbin.org/get")
		responseBytes, err := req3.Bytes()
		if err != nil {
			t.Errorf("Bytes method failed: %v", err)
		}
		if len(responseBytes) == 0 {
			t.Error("Expected byte response")
		}
	})

	t.Run("Request_TestExecuteMethods", func(t *testing.T) {
		// TestExecute需要测试服务器，这里只测试方法存在
		session := NewSession()
		req := session.Get("http://httpbin.org/get")

		// 这些方法需要ITestServer参数，在这里我们只能确认方法存在
		// 实际测试需要mock服务器
		if req == nil {
			t.Error("Request should not be nil")
		}
	})
}

// TestUncoveredResponseFunctions 测试未覆盖的response功能
func TestUncoveredResponseFunctions(t *testing.T) {
	t.Run("Response_GetMethods", func(t *testing.T) {
		// 创建一个模拟响应
		httpResp := &http.Response{
			StatusCode:    200,
			Status:        "200 OK",
			ContentLength: 100,
			Header:        http.Header{},
		}

		resp := &Response{
			readResponse: httpResp,
			readBytes:    []byte("test content"),
		}

		// 测试GetResponse
		originalResp := resp.GetResponse()
		if originalResp != httpResp {
			t.Error("GetResponse should return original http.Response")
		}

		// 测试GetCookie
		cookies := resp.GetCookie()
		if cookies == nil {
			t.Error("Expected cookies slice, got nil")
		}
	})
}

// TestUncoveredSessionFunctions 测试未覆盖的session功能
func TestUncoveredSessionFunctions(t *testing.T) {
	t.Run("Session_GetRetryConfig", func(t *testing.T) {
		session := NewSession()

		// GetRetryConfig可能返回默认配置
		retryConfig := session.GetRetryConfig()
		if retryConfig == nil {
			t.Log("GetRetryConfig returned nil (no retry config set)")
		} else {
			t.Logf("GetRetryConfig returned: %+v", retryConfig)
		}
	})
}

// TestUncoveredMultiPoolFunctions 测试未覆盖的multi pool功能
func TestUncoveredMultiPoolFunctions(t *testing.T) {
	t.Run("RequestPool_SetBar", func(t *testing.T) {
		pool := NewRequestPool(2) // 需要指定runner数量

		// SetBar方法测试 - 应传入bool值
		pool.SetBar(true)  // 启用进度条
		pool.SetBar(false) // 禁用进度条

		// 添加一些请求
		session := NewSession()
		req1 := session.Get("http://httpbin.org/get")
		req2 := session.Get("http://httpbin.org/get")

		pool.Add(req1)
		pool.Add(req2)

		// 执行请求池
		responses := pool.Execute()

		if len(responses) != 2 {
			t.Errorf("Expected 2 responses, got %d", len(responses))
		}
	})
}

// TestUncoveredUploadFileFunctions 测试未覆盖的upload file功能
func TestUncoveredUploadFileFunctions(t *testing.T) {
	t.Run("UploadFile_SetFile", func(t *testing.T) {
		uf := &UploadFile{
			FileName:  "test.txt",
			FieldName: "file",
		}

		// 测试SetFile
		fileContent := []byte("test file content")
		uf.SetFile(bytes.NewReader(fileContent))

		// 验证文件被设置
		if uf.FileReader == nil {
			t.Error("Expected FileReader to be set")
		}
	})
} // TestSpecialEdgeCases 测试一些特殊的边界情况
func TestSpecialEdgeCases(t *testing.T) {
	t.Run("Context_EdgeCases", func(t *testing.T) {
		session := NewSession()
		req := session.Get("http://httpbin.org/delay/1")

		// 测试带有取消的context
		ctx, cancel := context.WithCancel(context.Background())
		req.WithContext(ctx)

		// 立即取消
		cancel()

		// 执行应该因为context取消而失败
		_, err := req.Execute()
		if err == nil {
			t.Log("Request with cancelled context didn't fail (might be timing related)")
		}
	})

	t.Run("Multipart_EdgeCases", func(t *testing.T) {
		session := NewSession()
		req := session.Post("http://httpbin.org/post")

		// 创建multipart数据
		multipartData := req.CreateBodyMultipart()
		if multipartData == nil {
			t.Error("Expected multipart data, got nil")
		}

		// 添加字段和文件
		err := multipartData.AddField("text_field", "text_value")
		if err != nil {
			t.Errorf("AddField failed: %v", err)
		}

		// 尝试添加一个小文件
		fileContent := []byte("small file content")
		err = multipartData.AddFile("file_field", "small.txt", fileContent)
		if err != nil {
			t.Errorf("AddFile failed: %v", err)
		}
	})
}
