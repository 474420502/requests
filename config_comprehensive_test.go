package requests

import (
	"crypto/tls"
	"net/url"
	"testing"
	"time"
)

// TestConfig_SetBasicAuthLegacy 测试遗留的基础认证方法的各种情况
func TestConfig_SetBasicAuthLegacy(t *testing.T) {
	session := NewSession()
	config := session.Config()

	t.Run("BasicAuthStruct", func(t *testing.T) {
		auth := &BasicAuth{User: "user1", Password: "pass1"}
		err := config.SetBasicAuthLegacy(auth)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if session.auth == nil || session.auth.User != "user1" || session.auth.Password != "pass1" {
			t.Error("BasicAuth not set correctly")
		}
	})

	t.Run("BasicAuthValue", func(t *testing.T) {
		auth := BasicAuth{User: "user2", Password: "pass2"}
		err := config.SetBasicAuthLegacy(auth)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if session.auth == nil || session.auth.User != "user2" || session.auth.Password != "pass2" {
			t.Error("BasicAuth not set correctly")
		}
	})

	t.Run("NilAuth", func(t *testing.T) {
		err := config.SetBasicAuthLegacy(nil)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if session.auth != nil {
			t.Error("BasicAuth should be cleared")
		}
	})

	t.Run("TwoStringArgs", func(t *testing.T) {
		err := config.SetBasicAuthLegacy("user3", "pass3")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if session.auth == nil || session.auth.User != "user3" || session.auth.Password != "pass3" {
			t.Error("BasicAuth not set correctly")
		}
	})

	t.Run("InvalidSingleArg", func(t *testing.T) {
		err := config.SetBasicAuthLegacy(123)
		if err == nil {
			t.Error("Expected error for invalid type")
		}
	})

	t.Run("InvalidFirstArg", func(t *testing.T) {
		err := config.SetBasicAuthLegacy(123, "password")
		if err == nil {
			t.Error("Expected error for invalid first argument type")
		}
	})

	t.Run("InvalidSecondArg", func(t *testing.T) {
		err := config.SetBasicAuthLegacy("username", 123)
		if err == nil {
			t.Error("Expected error for invalid second argument type")
		}
	})

	t.Run("InvalidArgCount", func(t *testing.T) {
		err := config.SetBasicAuthLegacy("user1", "pass1", "extra")
		if err == nil {
			t.Error("Expected error for too many arguments")
		}

		err = config.SetBasicAuthLegacy()
		if err == nil {
			t.Error("Expected error for no arguments")
		}
	})
}

// TestConfig_SetProxy 测试代理设置的各种情况
func TestConfig_SetProxy(t *testing.T) {
	session := NewSession()
	config := session.Config()

	t.Run("StringProxy", func(t *testing.T) {
		err := config.SetProxy("http://proxy.example.com:8080")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if session.transport.Proxy == nil {
			t.Error("Proxy should be set")
		}
	})

	t.Run("URLProxy", func(t *testing.T) {
		proxyURL, _ := url.Parse("http://proxy2.example.com:8080")
		err := config.SetProxy(proxyURL)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if session.transport.Proxy == nil {
			t.Error("Proxy should be set")
		}
	})

	t.Run("NilProxy", func(t *testing.T) {
		err := config.SetProxy(nil)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if session.transport.Proxy != nil {
			t.Error("Proxy should be cleared")
		}
	})

	t.Run("InvalidProxyType", func(t *testing.T) {
		err := config.SetProxy(123)
		if err == nil {
			t.Error("Expected error for invalid proxy type")
		}
	})
}

// TestConfig_setProxyURL 测试内部代理URL设置方法
func TestConfig_setProxyURL(t *testing.T) {
	session := NewSession()
	config := session.Config()

	t.Run("HTTPProxy", func(t *testing.T) {
		proxyURL, _ := url.Parse("http://proxy.example.com:8080")
		err := config.setProxyURL(proxyURL)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if session.transport.Proxy == nil {
			t.Error("HTTP proxy should be set")
		}
	})

	t.Run("SOCKS5Proxy", func(t *testing.T) {
		proxyURL, _ := url.Parse("socks5://127.0.0.1:1080")
		err := config.setProxyURL(proxyURL)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if session.transport.DialContext == nil {
			t.Error("SOCKS5 dial context should be set")
		}
		if session.transport.Proxy == nil {
			t.Error("SOCKS5 proxy placeholder should be set")
		}
	})
}

// TestConfig_SetTimeout 测试超时设置的各种类型
func TestConfig_SetTimeout(t *testing.T) {
	session := NewSession()
	config := session.Config()

	t.Run("Duration", func(t *testing.T) {
		timeout := 30 * time.Second
		err := config.SetTimeout(timeout)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if session.client.Timeout != timeout {
			t.Error("Duration timeout not set correctly")
		}
	})

	t.Run("IntSeconds", func(t *testing.T) {
		err := config.SetTimeout(45)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if session.client.Timeout != 45*time.Second {
			t.Error("Int timeout not set correctly")
		}
	})

	t.Run("Int64Seconds", func(t *testing.T) {
		err := config.SetTimeout(int64(60))
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if session.client.Timeout != 60*time.Second {
			t.Error("Int64 timeout not set correctly")
		}
	})

	t.Run("Float32Seconds", func(t *testing.T) {
		err := config.SetTimeout(float32(1.5))
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		expected := time.Duration(1.5 * float32(time.Second))
		if session.client.Timeout != expected {
			t.Error("Float32 timeout not set correctly")
		}
	})

	t.Run("Float64Seconds", func(t *testing.T) {
		err := config.SetTimeout(2.5)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		expected := time.Duration(2.5 * float64(time.Second))
		if session.client.Timeout != expected {
			t.Error("Float64 timeout not set correctly")
		}
	})

	t.Run("UnsupportedType", func(t *testing.T) {
		err := config.SetTimeout("invalid")
		if err == nil {
			t.Error("Expected error for unsupported timeout type")
		}
	})
}

// TestConfig_TLSConfig 测试TLS配置
func TestConfig_TLSConfig(t *testing.T) {
	session := NewSession()
	config := session.Config()

	t.Run("SetTLSConfig", func(t *testing.T) {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS12,
		}
		config.SetTLSConfig(tlsConfig)
		if session.transport.TLSClientConfig != tlsConfig {
			t.Error("TLS config not set correctly")
		}
	})

	t.Run("SetInsecure", func(t *testing.T) {
		config.SetInsecure(true)
		if session.transport.TLSClientConfig == nil || !session.transport.TLSClientConfig.InsecureSkipVerify {
			t.Error("Insecure setting not applied correctly")
		}

		config.SetInsecure(false)
		if session.transport.TLSClientConfig == nil || session.transport.TLSClientConfig.InsecureSkipVerify {
			t.Error("Secure setting not applied correctly")
		}
	})
}

// TestConfig_CookieJar 测试Cookie Jar配置
func TestConfig_CookieJar(t *testing.T) {
	session := NewSession()
	config := session.Config()

	t.Run("EnableCookieJar", func(t *testing.T) {
		config.SetWithCookiejar(true)
		if session.client.Jar == nil {
			t.Error("Cookie jar should be enabled")
		}
	})

	t.Run("DisableCookieJar", func(t *testing.T) {
		config.SetWithCookiejar(false)
		if session.client.Jar != nil {
			t.Error("Cookie jar should be disabled")
		}
	})
}

// TestConfig_KeepAlives 测试Keep-Alive配置
func TestConfig_KeepAlives(t *testing.T) {
	session := NewSession()
	config := session.Config()

	t.Run("EnableKeepAlives", func(t *testing.T) {
		config.SetKeepAlives(true)
		if session.transport.DisableKeepAlives {
			t.Error("Keep-alives should be enabled")
		}
	})

	t.Run("DisableKeepAlives", func(t *testing.T) {
		config.SetKeepAlives(false)
		if !session.transport.DisableKeepAlives {
			t.Error("Keep-alives should be disabled")
		}
	})
}

// TestConfig_CompressionAndEncoding 测试压缩和编码配置
func TestConfig_CompressionAndEncoding(t *testing.T) {
	session := NewSession()
	config := session.Config()

	t.Run("DecompressNoAccept", func(t *testing.T) {
		config.SetDecompressNoAccept(true)
		if !session.Is.isDecompressNoAccept {
			t.Error("DecompressNoAccept should be enabled")
		}

		config.SetDecompressNoAccept(false)
		if session.Is.isDecompressNoAccept {
			t.Error("DecompressNoAccept should be disabled")
		}
	})

	t.Run("AcceptEncoding", func(t *testing.T) {
		initialCount := len(session.acceptEncoding)
		config.AddAcceptEncoding(AcceptEncodingGzip)
		if len(session.acceptEncoding) != initialCount+1 {
			t.Error("Accept encoding not added correctly")
		}

		encodings := config.GetAcceptEncoding(AcceptEncodingGzip)
		if len(encodings) == 0 {
			t.Error("Should return accept encodings")
		}
	})

	t.Run("ContentEncoding", func(t *testing.T) {
		config.SetContentEncoding(ContentEncodingGzip)
		if session.contentEncoding != ContentEncodingGzip {
			t.Error("Content encoding not set correctly")
		}
	})
}

// TestConfig_TimeoutMethods 测试超时相关方法
func TestConfig_TimeoutMethods(t *testing.T) {
	session := NewSession()
	config := session.Config()

	t.Run("SetTimeoutDuration", func(t *testing.T) {
		timeout := 25 * time.Second
		config.SetTimeoutDuration(timeout)
		if session.client.Timeout != timeout {
			t.Error("Timeout duration not set correctly")
		}
	})

	t.Run("SetTimeoutSeconds", func(t *testing.T) {
		config.SetTimeoutSeconds(35)
		if session.client.Timeout != 35*time.Second {
			t.Error("Timeout seconds not set correctly")
		}
	})
}

// TestConfig_Authorization 测试认证相关方法
func TestConfig_Authorization(t *testing.T) {
	session := NewSession()
	config := session.Config()

	t.Run("SetHeaderAuthorization", func(t *testing.T) {
		token := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
		config.SetHeaderAuthorization(token)
		if session.Header.Get("Authorization") != token {
			t.Error("Authorization header not set correctly")
		}
	})

	t.Run("SetBasicAuthStruct", func(t *testing.T) {
		auth := &BasicAuth{User: "testuser", Password: "testpass"}
		config.SetBasicAuthStruct(auth)
		if session.auth == nil || session.auth.User != "testuser" || session.auth.Password != "testpass" {
			t.Error("BasicAuth struct not set correctly")
		}
	})

	t.Run("SetBasicAuthStructNil", func(t *testing.T) {
		config.SetBasicAuthStruct(nil)
		if session.auth != nil {
			t.Error("BasicAuth should be cleared when set to nil")
		}
	})

	t.Run("ClearBasicAuth", func(t *testing.T) {
		// First set some auth
		config.SetBasicAuth("user", "pass")
		if session.auth == nil {
			t.Error("Auth should be set")
		}

		// Then clear it
		config.ClearBasicAuth()
		if session.auth != nil {
			t.Error("Auth should be cleared")
		}
	})
}

// TestConfig_ProxyMethods 测试代理相关方法
func TestConfig_ProxyMethods(t *testing.T) {
	session := NewSession()
	config := session.Config()

	t.Run("SetProxyString", func(t *testing.T) {
		proxyURL := "http://proxy.test.com:8080"
		err := config.SetProxyString(proxyURL)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if session.transport.Proxy == nil {
			t.Error("Proxy should be set")
		}
	})

	t.Run("SetProxyStringEmpty", func(t *testing.T) {
		err := config.SetProxyString("")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if session.transport.Proxy != nil {
			t.Error("Proxy should be cleared for empty string")
		}
	})

	t.Run("SetProxyStringInvalid", func(t *testing.T) {
		err := config.SetProxyString("://invalid-url")
		if err == nil {
			t.Error("Expected error for invalid proxy URL")
		}
	})

	t.Run("ClearProxy", func(t *testing.T) {
		// First set a proxy
		config.SetProxyString("http://proxy.test.com:8080")
		if session.transport.Proxy == nil {
			t.Error("Proxy should be set")
		}

		// Then clear it
		config.ClearProxy()
		if session.transport.Proxy != nil {
			t.Error("Proxy should be cleared")
		}
	})
}

// TestConfig_EdgeCases 测试边界情况
func TestConfig_EdgeCases(t *testing.T) {
	session := NewSession()
	config := session.Config()

	t.Run("MultipleAuthOperations", func(t *testing.T) {
		// Test multiple auth operations in sequence
		config.SetBasicAuth("user1", "pass1")
		config.SetBasicAuthString("user2", "pass2")
		config.SetBasicAuthStruct(&BasicAuth{User: "user3", Password: "pass3"})

		if session.auth == nil || session.auth.User != "user3" || session.auth.Password != "pass3" {
			t.Error("Final auth state not correct")
		}
	})

	t.Run("MultipleProxyOperations", func(t *testing.T) {
		// Test multiple proxy operations
		config.SetProxyString("http://proxy1.test.com:8080")
		config.SetProxyString("http://proxy2.test.com:8080")
		config.ClearProxy()

		if session.transport.Proxy != nil {
			t.Error("Proxy should be cleared")
		}
	})

	t.Run("MultipleTimeoutOperations", func(t *testing.T) {
		// Test multiple timeout operations
		config.SetTimeoutDuration(10 * time.Second)
		config.SetTimeoutSeconds(20)
		config.SetTimeout(30 * time.Second)

		if session.client.Timeout != 30*time.Second {
			t.Error("Final timeout not correct")
		}
	})
}
