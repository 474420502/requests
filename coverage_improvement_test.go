package requests

import (
	"context"
	"crypto/tls"
	"net/http"
	"strings"
	"testing"
	"time"
)

// TestConfigLowCoverageMethods 测试配置方法中覆盖率较低的函数
func TestConfigLowCoverageMethods(t *testing.T) {
	session := NewSession()
	config := session.Config()

	t.Run("SetDecompressNoAccept", func(t *testing.T) {
		// 测试设置在没有Accept-Encoding头的情况下是否解压
		config.SetDecompressNoAccept(true)
		if !session.Is.isDecompressNoAccept {
			t.Error("Expected isDecompressNoAccept to be true")
		}

		config.SetDecompressNoAccept(false)
		if session.Is.isDecompressNoAccept {
			t.Error("Expected isDecompressNoAccept to be false")
		}
	})

	t.Run("GetAcceptEncoding", func(t *testing.T) {
		// 清空现有编码
		session.acceptEncoding = []AcceptEncodingType{}

		// 添加编码类型
		config.AddAcceptEncoding(AcceptEncodingGzip)
		config.AddAcceptEncoding(AcceptEncodingDeflate)

		// 测试获取编码类型
		encodings := config.GetAcceptEncoding(AcceptEncodingGzip)
		if len(encodings) != 2 {
			t.Errorf("Expected 2 encodings, got %d", len(encodings))
		}

		// 验证编码类型
		found := false
		for _, enc := range encodings {
			if enc == AcceptEncodingGzip {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected to find AcceptEncodingGzip")
		}
	})

	t.Run("SetContentEncoding", func(t *testing.T) {
		// 测试设置内容编码
		config.SetContentEncoding(ContentEncodingGzip)
		if session.contentEncoding != ContentEncodingGzip {
			t.Error("Expected contentEncoding to be ContentEncodingGzip")
		}

		config.SetContentEncoding(ContentEncodingDeflate)
		if session.contentEncoding != ContentEncodingDeflate {
			t.Error("Expected contentEncoding to be ContentEncodingDeflate")
		}
	})

	t.Run("SetHeaderAuthorization", func(t *testing.T) {
		// 测试设置JWT token Authorization头
		token := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
		config.SetHeaderAuthorization(token)

		// 验证头部是否正确设置
		authHeader := session.Header.Get("Authorization")
		if authHeader != token {
			t.Errorf("Expected Authorization header to be %s, got %s", token, authHeader)
		}

		// 测试多次添加Authorization头
		token2 := "Bearer another.jwt.token"
		config.SetHeaderAuthorization(token2)

		// 应该有两个Authorization头
		authHeaders := session.Header.Values("Authorization")
		if len(authHeaders) != 2 {
			t.Errorf("Expected 2 Authorization headers, got %d", len(authHeaders))
		}
	})
}

// TestSetBasicAuthLegacyEdgeCases 测试传统BasicAuth方法的边界情况
func TestSetBasicAuthLegacyEdgeCases(t *testing.T) {
	session := NewSession()
	config := session.Config()

	t.Run("SetBasicAuthLegacyWithComplexCredentials", func(t *testing.T) {
		// 测试包含特殊字符的用户名和密码
		username := "user@domain.com"
		password := "pass:word!@#$%"

		err := config.SetBasicAuthLegacy(username, password)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// 验证认证是否正确设置
		if session.auth == nil {
			t.Error("Expected auth to be set")
		} else {
			if session.auth.User != username {
				t.Errorf("Expected username %s, got %s", username, session.auth.User)
			}
			if session.auth.Password != password {
				t.Errorf("Expected password %s, got %s", password, session.auth.Password)
			}
		}
	})

	t.Run("SetBasicAuthLegacyWithEmptyCredentials", func(t *testing.T) {
		// 测试空凭证
		err := config.SetBasicAuthLegacy("", "")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// 即使是空凭证也应该被设置
		if session.auth == nil {
			t.Error("Expected auth to be set even with empty credentials")
		}
	})

	t.Run("SetBasicAuthLegacyWithUnicodeCredentials", func(t *testing.T) {
		// 测试Unicode字符
		username := "用户名"
		password := "密码123"

		err := config.SetBasicAuthLegacy(username, password)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if session.auth == nil {
			t.Error("Expected auth to be set")
		} else {
			if session.auth.User != username {
				t.Errorf("Expected username %s, got %s", username, session.auth.User)
			}
		}
	})
}

// TestSetProxyEdgeCases 测试代理设置的边界情况
func TestSetProxyEdgeCases(t *testing.T) {
	session := NewSession()
	config := session.Config()

	t.Run("SetProxyWithDifferentSchemes", func(t *testing.T) {
		testCases := []struct {
			name     string
			proxyURL string
			wantErr  bool
		}{
			{"HTTP proxy", "http://proxy.example.com:8080", false},
			{"HTTPS proxy", "https://proxy.example.com:8080", false},
			{"SOCKS5 proxy", "socks5://proxy.example.com:1080", false},
			// These URL formats are actually accepted by the URL parser
			{"Invalid scheme", "invalid://proxy.example.com:8080", false}, // URL parser accepts any scheme
			{"No scheme", "proxy.example.com:8080", false},                // This becomes a relative URL
			{"Empty URL", "", false},                                      // Empty string clears proxy, no error
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := config.SetProxyString(tc.proxyURL)
				if tc.wantErr && err == nil {
					t.Error("Expected error but got none")
				}
				if !tc.wantErr && err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			})
		}
	})

	t.Run("SetProxyWithAuthentication", func(t *testing.T) {
		// 测试包含认证信息的代理
		proxyURL := "http://user:pass@proxy.example.com:8080"
		err := config.SetProxyString(proxyURL)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// 验证代理是否正确设置
		transport := session.client.Transport.(*http.Transport)
		if transport.Proxy == nil {
			t.Error("Expected proxy to be set")
		}
	})
}

// TestRequestCompatibilityMethods 测试Request的兼容性方法
func TestRequestCompatibilityMethods(t *testing.T) {
	session := NewSession()

	t.Run("RequestBasicMethods", func(t *testing.T) {
		req := session.Get("http://httpbin.org/get")
		if req == nil {
			t.Fatal("Expected Request object to be created")
		}

		// 测试Error方法
		err := req.Error()
		// 在没有错误的情况下应该返回nil
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("RequestBodyMethods", func(t *testing.T) {
		req := session.Post("http://httpbin.org/post")

		// 测试SetBodyUrlencoded
		data := map[string]string{"key": "value"}
		req.SetBodyUrlencoded(data)

		// 测试SetBodyPlain
		plainData := "plain text data"
		req.SetBodyPlain(plainData)

		// 测试SetBodyStream
		reader := strings.NewReader("stream data")
		req.SetBodyStream(reader)

		// 这些方法应该不会导致panic
	})

	t.Run("RequestWithInvalidURL", func(t *testing.T) {
		// 测试无效URL
		req := NewRequest(session, "GET", "://invalid-url")
		if req.Error() == nil {
			t.Error("Expected error for invalid URL")
		}
	})

	t.Run("RequestExecuteMethods", func(t *testing.T) {
		req := session.Get("http://httpbin.org/get")

		// 测试BuildRequest (不执行实际请求)
		httpReq, err := req.buildHTTPRequest()
		if err != nil {
			t.Errorf("buildHTTPRequest returned error: %v", err)
		}
		if httpReq == nil {
			t.Error("Expected buildHTTPRequest to return non-nil request")
		}

		// 测试TestExecute - 需要一个测试服务器接口
		// 由于TestExecute需要ITestServer参数，我们跳过这个测试
		// 或者创建一个模拟的测试服务器
	})
}

// TestUploadFileMethods 测试文件上传相关方法
func TestUploadFileMethods(t *testing.T) {
	t.Run("UploadFileGetMethods", func(t *testing.T) {
		uploadFile := &UploadFile{
			FileName:  "test.txt",
			FieldName: "file",
		}

		// 测试GetFileName
		fileName := uploadFile.GetFileName()
		if fileName != "test.txt" {
			t.Errorf("Expected filename 'test.txt', got '%s'", fileName)
		}

		// 测试GetFieldName
		fieldName := uploadFile.GetFieldName()
		if fieldName != "file" {
			t.Errorf("Expected field name 'file', got '%s'", fieldName)
		}

		// 测试GetFile (当没有设置File时)
		file := uploadFile.GetFile()
		if file != nil {
			t.Error("Expected GetFile to return nil when no file is set")
		}
	})

	t.Run("UploadFileSetFromPath", func(t *testing.T) {
		uploadFile := &UploadFile{}

		// 测试SetFileFromPath - 使用不存在的文件路径
		err := uploadFile.SetFileFromPath("/nonexistent/file.txt")
		if err == nil {
			t.Error("Expected error when setting file from non-existent path")
		}

		// 测试使用现有文件
		// 创建临时文件用于测试
		testContent := "test file content"
		testFile := strings.NewReader(testContent)
		uploadFile.SetFile(testFile)

		file := uploadFile.GetFile()
		if file == nil {
			t.Error("Expected file to be set")
		}
	})
}

// TestRequestBuilderExtendedMethods 测试Request构建器的扩展方法
func TestRequestBuilderExtendedMethods(t *testing.T) {
	session := NewSession()

	t.Run("RequestWithComplexChaining", func(t *testing.T) {
		// 测试复杂的方法链式调用
		req := session.Get("http://httpbin.org/get").
			SetHeader("X-Custom-Header", "test-value").
			AddQuery("param1", "value1").
			AddQuery("param2", "value2").
			WithTimeout(30 * time.Second).
			WithContext(context.Background())

		if req == nil {
			t.Fatal("Expected request to be created")
		}

		// 验证headers是否正确设置
		if req.header.Get("X-Custom-Header") != "test-value" {
			t.Error("Expected custom header to be set")
		}
	})

	t.Run("RequestBodyMethods", func(t *testing.T) {
		req := session.Post("http://httpbin.org/post")

		// 测试SetBodyWithType方法
		req.SetBodyWithType("test data", "text/plain")

		// 测试CreateBodyMultipart方法
		req.CreateBodyMultipart()

		// 这些方法应该不会导致panic
	})
}

// TestCompressionHandling 测试压缩处理
func TestCompressionHandling(t *testing.T) {
	session := NewSession()

	t.Run("AcceptEncodingConfiguration", func(t *testing.T) {
		config := session.Config()

		// 清空现有编码设置
		session.acceptEncoding = []AcceptEncodingType{}

		// 添加多种编码类型
		config.AddAcceptEncoding(AcceptEncodingGzip)
		config.AddAcceptEncoding(AcceptEncodingDeflate)
		config.AddAcceptEncoding(AcceptEncodingBr)

		// 验证编码类型数量
		if len(session.acceptEncoding) != 3 {
			t.Errorf("Expected 3 encoding types, got %d", len(session.acceptEncoding))
		}

		// 验证具体编码类型
		encodingMap := make(map[AcceptEncodingType]bool)
		for _, enc := range session.acceptEncoding {
			encodingMap[enc] = true
		}

		if !encodingMap[AcceptEncodingGzip] {
			t.Error("Expected AcceptEncodingGzip to be set")
		}
		if !encodingMap[AcceptEncodingDeflate] {
			t.Error("Expected AcceptEncodingDeflate to be set")
		}
		if !encodingMap[AcceptEncodingBr] {
			t.Error("Expected AcceptEncodingBr to be set")
		}
	})

	t.Run("ContentEncodingConfiguration", func(t *testing.T) {
		config := session.Config()

		// 测试设置不同的内容编码
		testCases := []ContentEncodingType{
			ContentEncodingGzip,
			ContentEncodingDeflate,
		}

		for _, encoding := range testCases {
			config.SetContentEncoding(encoding)
			if session.contentEncoding != encoding {
				t.Errorf("Expected content encoding %v, got %v", encoding, session.contentEncoding)
			}
		}
	})
}

// TestTimeoutConfiguration 测试超时配置
func TestTimeoutConfiguration(t *testing.T) {
	session := NewSession()
	config := session.Config()

	t.Run("SetTimeoutSeconds", func(t *testing.T) {
		// 测试设置超时秒数
		config.SetTimeoutSeconds(30)
		expectedTimeout := 30 * time.Second
		if session.client.Timeout != expectedTimeout {
			t.Errorf("Expected timeout %v, got %v", expectedTimeout, session.client.Timeout)
		}

		// 测试设置不同的超时时间
		config.SetTimeoutSeconds(60)
		expectedTimeout = 60 * time.Second
		if session.client.Timeout != expectedTimeout {
			t.Errorf("Expected timeout %v, got %v", expectedTimeout, session.client.Timeout)
		}
	})

	t.Run("SetTimeoutDuration", func(t *testing.T) {
		// 测试设置Duration类型的超时
		duration := 45 * time.Second
		config.SetTimeoutDuration(duration)
		if session.client.Timeout != duration {
			t.Errorf("Expected timeout %v, got %v", duration, session.client.Timeout)
		}
	})
}

// TestProxyAuthenticationScenarios 测试代理认证场景
func TestProxyAuthenticationScenarios(t *testing.T) {
	session := NewSession()
	config := session.Config()

	t.Run("ProxyWithSpecialCharacters", func(t *testing.T) {
		// 测试包含特殊字符的代理认证
		proxyURL := "http://user%40domain:pass%21word@proxy.example.com:8080"
		err := config.SetProxyString(proxyURL)
		if err != nil {
			t.Errorf("Unexpected error with URL-encoded proxy credentials: %v", err)
		}
	})

	t.Run("ProxyWithIPAddress", func(t *testing.T) {
		// 测试使用IP地址的代理
		proxyURL := "http://192.168.1.100:3128"
		err := config.SetProxyString(proxyURL)
		if err != nil {
			t.Errorf("Unexpected error with IP proxy: %v", err)
		}
	})

	t.Run("ProxyWithNonStandardPort", func(t *testing.T) {
		// 测试非标准端口的代理
		proxyURL := "http://proxy.example.com:9999"
		err := config.SetProxyString(proxyURL)
		if err != nil {
			t.Errorf("Unexpected error with non-standard port proxy: %v", err)
		}
	})
}

// TestTLSConfigurationExtended 测试扩展的TLS配置
func TestTLSConfigurationExtended(t *testing.T) {
	session := NewSession()
	config := session.Config()

	t.Run("CustomTLSConfig", func(t *testing.T) {
		// 创建自定义TLS配置
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS12,
			MaxVersion:         tls.VersionTLS13,
		}

		config.SetTLSConfig(tlsConfig)

		// 验证TLS配置是否正确设置
		transport := session.client.Transport.(*http.Transport)
		if transport.TLSClientConfig != tlsConfig {
			t.Error("Expected TLS config to be set correctly")
		}
	})

	t.Run("InsecureConfiguration", func(t *testing.T) {
		// 测试设置不安全连接
		config.SetInsecure(true)

		transport := session.client.Transport.(*http.Transport)
		if !transport.TLSClientConfig.InsecureSkipVerify {
			t.Error("Expected InsecureSkipVerify to be true")
		}

		// 测试恢复安全连接
		config.SetInsecure(false)
		if transport.TLSClientConfig.InsecureSkipVerify {
			t.Error("Expected InsecureSkipVerify to be false")
		}
	})
}

// TestErrorHandlingExtended 测试扩展错误处理
func TestErrorHandlingExtended(t *testing.T) {
	t.Run("RequestWithError", func(t *testing.T) {
		session := NewSession()

		// 创建带有错误的Request对象
		req := NewRequest(session, "GET", "://malformed-url")
		if req.Error() == nil {
			t.Error("Expected error for malformed URL")
		}

		// 测试Error方法返回错误
		err := req.Error()
		if err == nil {
			t.Error("Expected Error() to return non-nil error")
		}
	})

	t.Run("RequestExecutionError", func(t *testing.T) {
		session := NewSession()

		// 创建一个会导致错误的请求 (无效的URL)
		req := session.Get("://invalid-url-scheme")
		if req.err == nil {
			t.Error("Expected error for invalid URL")
		}
	})
}

// TestKeepAliveConfiguration 测试Keep-Alive配置
func TestKeepAliveConfiguration(t *testing.T) {
	session := NewSession()
	config := session.Config()

	t.Run("EnableKeepAlive", func(t *testing.T) {
		// 启用Keep-Alive
		config.SetKeepAlives(true)

		transport := session.client.Transport.(*http.Transport)
		if transport.DisableKeepAlives {
			t.Error("Expected DisableKeepAlives to be false when keep-alive is enabled")
		}
	})

	t.Run("DisableKeepAlive", func(t *testing.T) {
		// 禁用Keep-Alive
		config.SetKeepAlives(false)

		transport := session.client.Transport.(*http.Transport)
		if !transport.DisableKeepAlives {
			t.Error("Expected DisableKeepAlives to be true when keep-alive is disabled")
		}
	})
}
