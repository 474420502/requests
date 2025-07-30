package requests

import (
	"context"
	"crypto/tls"
	"errors"
	"strings"
	"testing"
	"time"
)

// TestSessionBuilderEdgeCasesExtended 测试Session构建器的更多边界情况
func TestSessionBuilderEdgeCasesExtended(t *testing.T) {
	t.Run("WithProxyEdgeCases", func(t *testing.T) {
		// 测试各种边界情况的代理URL
		testCases := []struct {
			name      string
			proxyURL  string
			shouldErr bool
		}{
			{"EmptyProxy", "", true},
			{"NoScheme", "proxy.example.com:8080", true},
			{"InvalidScheme", "invalid://proxy.example.com:8080", true},
			{"NoHost", "http://", true},
			{"ValidHTTP", "http://proxy.example.com:8080", false},
			{"ValidHTTPS", "https://proxy.example.com:8080", false},
			{"ValidSOCKS5", "socks5://proxy.example.com:1080", false},
			{"WithAuth", "http://user:pass@proxy.example.com:8080", false},
			{"WithEncodedAuth", "http://user%40domain:pass%40word@proxy.example.com:8080", false},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := NewSessionWithOptions(WithProxy(tc.proxyURL))
				if tc.shouldErr && err == nil {
					t.Errorf("Expected error for proxy URL '%s', but got nil", tc.proxyURL)
				}
				if !tc.shouldErr && err != nil {
					t.Errorf("Unexpected error for proxy URL '%s': %v", tc.proxyURL, err)
				}
			})
		}
	})

	t.Run("WithTLSConfigNilHandling", func(t *testing.T) {
		// 测试传入nil TLS配置
		session, err := NewSessionWithOptions(WithTLSConfig(nil))
		if err != nil {
			t.Fatalf("WithTLSConfig(nil) should not return error: %v", err)
		}

		if session.transport.TLSClientConfig != nil {
			t.Error("TLS config should be nil when nil is passed")
		}
	})

	t.Run("WithHeadersEmptyMap", func(t *testing.T) {
		// 测试传入空的headers map
		emptyHeaders := make(map[string]string)
		session, err := NewSessionWithOptions(WithHeaders(emptyHeaders))
		if err != nil {
			t.Fatalf("WithHeaders with empty map should not return error: %v", err)
		}

		if session.Header == nil {
			t.Error("Header should not be nil even with empty input")
		}
	})

	t.Run("WithTimeoutZero", func(t *testing.T) {
		// 测试设置零超时
		session, err := NewSessionWithOptions(WithTimeout(0))
		if err != nil {
			t.Fatalf("WithTimeout(0) should not return error: %v", err)
		}

		if session.client.Timeout != 0 {
			t.Errorf("Expected timeout 0, got %v", session.client.Timeout)
		}
	})

	t.Run("WithBasicAuthEmptyCredentials", func(t *testing.T) {
		// 测试空的认证凭据
		session, err := NewSessionWithOptions(WithBasicAuth("", ""))
		if err != nil {
			t.Fatalf("WithBasicAuth with empty credentials should not return error: %v", err)
		}

		if session.auth == nil {
			t.Fatal("Auth should not be nil")
		}

		if session.auth.User != "" || session.auth.Password != "" {
			t.Error("Empty credentials should be preserved")
		}
	})
}

// TestSessionBuilderCombinations 测试多个选项的复杂组合
func TestSessionBuilderCombinations(t *testing.T) {
	t.Run("ComplexConfiguration", func(t *testing.T) {
		// 创建一个包含多种配置的复杂Session
		ctx := context.Background()
		tlsConfig := &tls.Config{InsecureSkipVerify: true}
		headers := map[string]string{
			"Custom-Header":  "test-value",
			"Another-Header": "another-value",
		}

		session, err := NewSessionWithOptions(
			WithTimeout(30*time.Second),
			WithTLSConfig(tlsConfig),
			WithBasicAuth("user", "pass"),
			WithHeaders(headers),
			WithUserAgent("TestAgent/1.0"),
			WithKeepAlives(true),
			WithCompression(true),
			WithMaxIdleConns(100),
			WithMaxIdleConnsPerHost(20),
			WithContext(ctx),
			WithRetry(3, time.Second),
			WithInsecureSkipVerify(),
		)

		if err != nil {
			t.Fatalf("Complex configuration should not return error: %v", err)
		}

		// 验证所有配置都被正确应用
		if session.client.Timeout != 30*time.Second {
			t.Error("Timeout not applied correctly")
		}

		if session.auth == nil || session.auth.User != "user" {
			t.Error("BasicAuth not applied correctly")
		}

		if session.Header.Get("Custom-Header") != "test-value" {
			t.Error("Headers not applied correctly")
		}

		if session.Header.Get("User-Agent") != "TestAgent/1.0" {
			t.Error("User-Agent not applied correctly")
		}

		if session.transport.DisableKeepAlives {
			t.Error("KeepAlives should be enabled")
		}

		if session.transport.DisableCompression {
			t.Error("Compression should be enabled")
		}

		if session.transport.MaxIdleConns != 100 {
			t.Error("MaxIdleConns not applied correctly")
		}

		if session.retryConfig == nil || session.retryConfig.MaxRetries != 3 {
			t.Error("Retry config not applied correctly")
		}
	})

	t.Run("ConflictingTLSConfigs", func(t *testing.T) {
		// 测试冲突的TLS配置（先设置自定义，再设置InsecureSkipVerify）
		customTLS := &tls.Config{ServerName: "custom.example.com"}

		session, err := NewSessionWithOptions(
			WithTLSConfig(customTLS),
			WithInsecureSkipVerify(),
		)

		if err != nil {
			t.Fatalf("Conflicting TLS configs should not return error: %v", err)
		}

		// WithInsecureSkipVerify 应该覆盖之前的配置
		if !session.transport.TLSClientConfig.InsecureSkipVerify {
			t.Error("InsecureSkipVerify should override previous TLS config")
		}
	})
}

// TestSessionBuilderErrorPropagation 测试错误传播
func TestSessionBuilderErrorPropagation(t *testing.T) {
	t.Run("ProxyErrorPropagation", func(t *testing.T) {
		// 测试代理错误是否正确传播
		_, err := NewSessionWithOptions(WithProxy("://invalid-url"))
		if err == nil {
			t.Error("Expected error for invalid proxy URL")
		}

		// 检查错误信息是否与URL解析相关
		errorMsg := strings.ToLower(err.Error())
		if !strings.Contains(errorMsg, "proxy") &&
			!strings.Contains(errorMsg, "url") &&
			!strings.Contains(errorMsg, "scheme") &&
			!strings.Contains(errorMsg, "parse") {
			t.Errorf("Error should mention URL parsing issue: %v", err)
		}
	})

	t.Run("MultipleErrorsPropagation", func(t *testing.T) {
		// 使用mock的SessionOption来测试错误传播
		errorOption := func(*Session) error {
			return errors.New("test error")
		}

		_, err := NewSessionWithOptions(
			WithTimeout(10*time.Second), // 正常选项
			errorOption,                 // 错误选项
			WithUserAgent("Test"),       // 这个不应该被执行
		)

		if err == nil {
			t.Error("Expected error from errorOption")
		}

		if !strings.Contains(err.Error(), "test error") {
			t.Errorf("Expected 'test error', got: %v", err)
		}
	})
}

// TestPreDefinedSessionsValidation 验证预定义Session的配置正确性
func TestPreDefinedSessionsValidation(t *testing.T) {
	t.Run("NewSessionForAPIValidation", func(t *testing.T) {
		session, err := NewSessionForAPI()
		if err != nil {
			t.Fatalf("NewSessionForAPI failed: %v", err)
		}

		// 验证API Session的具体配置
		if session.client.Timeout != 30*time.Second {
			t.Errorf("API session timeout should be 30s, got %v", session.client.Timeout)
		}

		userAgent := session.Header.Get("User-Agent")
		if userAgent == "" {
			t.Error("API session should have User-Agent set")
		}

		if !strings.Contains(userAgent, "Go-Requests") {
			t.Errorf("API session User-Agent should contain 'Go-Requests', got: %s", userAgent)
		}
	})

	t.Run("NewSessionForScrapingValidation", func(t *testing.T) {
		session, err := NewSessionForScraping()
		if err != nil {
			t.Fatalf("NewSessionForScraping failed: %v", err)
		}

		// 验证Scraping Session的具体配置
		if session.client.Timeout != 10*time.Second {
			t.Errorf("Scraping session timeout should be 10s, got %v", session.client.Timeout)
		}

		userAgent := session.Header.Get("User-Agent")
		if !strings.Contains(userAgent, "Mozilla") {
			t.Errorf("Scraping session should have browser-like User-Agent, got: %s", userAgent)
		}
	})

	t.Run("NewSessionForTestingValidation", func(t *testing.T) {
		session, err := NewSessionForTesting()
		if err != nil {
			t.Fatalf("NewSessionForTesting failed: %v", err)
		}

		// 验证Testing Session的具体配置
		if session.client.Timeout != 5*time.Second {
			t.Errorf("Testing session timeout should be 5s, got %v", session.client.Timeout)
		}

		if session.cookiejar != nil {
			t.Error("Testing session should have cookies disabled")
		}
	})

	t.Run("NewHighPerformanceSessionValidation", func(t *testing.T) {
		session, err := NewHighPerformanceSession()
		if err != nil {
			t.Fatalf("NewHighPerformanceSession failed: %v", err)
		}

		// 验证High Performance Session的具体配置
		if session.client.Timeout != 60*time.Second {
			t.Errorf("High performance session timeout should be 60s, got %v", session.client.Timeout)
		}

		if session.transport.MaxIdleConns != 100 {
			t.Errorf("High performance session MaxIdleConns should be 100, got %d", session.transport.MaxIdleConns)
		}

		if session.transport.MaxIdleConnsPerHost != 20 {
			t.Errorf("High performance session MaxIdleConnsPerHost should be 20, got %d", session.transport.MaxIdleConnsPerHost)
		}
	})

	t.Run("NewSecureSessionValidation", func(t *testing.T) {
		session, err := NewSecureSession()
		if err != nil {
			t.Fatalf("NewSecureSession failed: %v", err)
		}

		// 验证Secure Session的具体配置
		if session.cookiejar != nil {
			t.Error("Secure session should have cookies disabled")
		}

		if !session.transport.DisableKeepAlives {
			t.Error("Secure session should have keep-alives disabled")
		}

		// Secure session 使用默认的安全TLS设置，TLSClientConfig可以为nil
		// 这是正确的安全行为，使用系统默认的TLS配置
		if session.transport.TLSClientConfig != nil {
			// 如果有TLS配置，确保它是安全的
			if session.transport.TLSClientConfig.InsecureSkipVerify {
				t.Error("Secure session should not skip TLS verification")
			}
		}
	})
}

// TestSessionBuilderMemoryManagement 测试内存管理
func TestSessionBuilderMemoryManagement(t *testing.T) {
	t.Run("MultipleSessionsIndependence", func(t *testing.T) {
		// 创建多个Session，确保它们是独立的
		session1, err1 := NewSessionWithOptions(WithUserAgent("Agent1"))
		session2, err2 := NewSessionWithOptions(WithUserAgent("Agent2"))

		if err1 != nil || err2 != nil {
			t.Fatalf("Failed to create sessions: %v, %v", err1, err2)
		}

		// 验证它们有不同的配置
		if session1.Header.Get("User-Agent") == session2.Header.Get("User-Agent") {
			t.Error("Sessions should have independent configurations")
		}

		// 修改一个session不应该影响另一个
		session1.Header.Set("Custom", "value1")
		session2.Header.Set("Custom", "value2")

		if session1.Header.Get("Custom") == session2.Header.Get("Custom") {
			t.Error("Sessions should maintain independent headers")
		}
	})

	t.Run("SessionCreationPerformance", func(t *testing.T) {
		// 测试Session创建的性能
		start := time.Now()
		const numSessions = 100

		for i := 0; i < numSessions; i++ {
			_, err := NewSessionWithOptions(
				WithTimeout(30*time.Second),
				WithMaxIdleConns(50),
				WithUserAgent("PerfTest"),
			)
			if err != nil {
				t.Fatalf("Failed to create session %d: %v", i, err)
			}
		}

		elapsed := time.Since(start)
		avgTime := elapsed / numSessions

		// 平均每个Session创建时间不应超过1ms
		if avgTime > time.Millisecond {
			t.Errorf("Session creation too slow: average %v per session", avgTime)
		}

		t.Logf("Created %d sessions in %v (avg: %v per session)", numSessions, elapsed, avgTime)
	})
}
