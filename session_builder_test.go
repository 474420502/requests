package requests

import (
	"context"
	"crypto/tls"
	"net/http"
	"testing"
	"time"
)

// TestSessionBuilderOptions 测试SessionBuilder选项功能
func TestSessionBuilderOptions(t *testing.T) {
	t.Run("WithTimeout", func(t *testing.T) {
		timeout := 10 * time.Second
		session, err := NewSessionWithOptions(WithTimeout(timeout))
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}
		if session.client.Timeout != timeout {
			t.Errorf("Expected timeout %v, got %v", timeout, session.client.Timeout)
		}
	})

	t.Run("WithTLSConfig", func(t *testing.T) {
		tlsConfig := &tls.Config{InsecureSkipVerify: true}
		session, err := NewSessionWithOptions(WithTLSConfig(tlsConfig))
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}
		if session.transport.TLSClientConfig != tlsConfig {
			t.Error("TLS config was not set correctly")
		}
	})

	t.Run("WithProxy", func(t *testing.T) {
		proxyURL := "http://127.0.0.1:1080"
		session, err := NewSessionWithOptions(WithProxy(proxyURL))
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}
		if session.transport.Proxy == nil {
			t.Error("Proxy was not set")
		}
	})

	t.Run("WithProxyInvalidURL", func(t *testing.T) {
		_, err := NewSessionWithOptions(WithProxy("://invalid-scheme-url"))
		if err == nil {
			t.Error("Expected error for invalid proxy URL")
		}
	})

	t.Run("WithBasicAuth", func(t *testing.T) {
		username, password := "testuser", "testpass"
		session, err := NewSessionWithOptions(WithBasicAuth(username, password))
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}
		if session.auth == nil {
			t.Fatal("Basic auth was not set")
		}
		if session.auth.User != username || session.auth.Password != password {
			t.Error("Basic auth credentials not set correctly")
		}
	})

	t.Run("WithHeaders", func(t *testing.T) {
		headers := map[string]string{
			"Custom-Header": "test-value",
			"Another":       "value",
		}
		session, err := NewSessionWithOptions(WithHeaders(headers))
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}
		for k, v := range headers {
			if session.Header.Get(k) != v {
				t.Errorf("Header %s: expected %s, got %s", k, v, session.Header.Get(k))
			}
		}
	})

	t.Run("WithUserAgent", func(t *testing.T) {
		userAgent := "Test-Agent/1.0"
		session, err := NewSessionWithOptions(WithUserAgent(userAgent))
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}
		if session.Header.Get("User-Agent") != userAgent {
			t.Errorf("Expected User-Agent %s, got %s", userAgent, session.Header.Get("User-Agent"))
		}
	})

	t.Run("WithDisableCookies", func(t *testing.T) {
		session, err := NewSessionWithOptions(WithDisableCookies())
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}
		if session.cookiejar != nil {
			t.Error("Cookie jar should be nil when cookies are disabled")
		}
		if session.client.Jar != nil {
			t.Error("Client cookie jar should be nil when cookies are disabled")
		}
	})

	t.Run("WithKeepAlives", func(t *testing.T) {
		session, err := NewSessionWithOptions(WithKeepAlives(false))
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}
		if !session.transport.DisableKeepAlives {
			t.Error("Keep-alives should be disabled")
		}

		session2, err := NewSessionWithOptions(WithKeepAlives(true))
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}
		if session2.transport.DisableKeepAlives {
			t.Error("Keep-alives should be enabled")
		}
	})

	t.Run("WithCompression", func(t *testing.T) {
		session, err := NewSessionWithOptions(WithCompression(false))
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}
		if !session.transport.DisableCompression {
			t.Error("Compression should be disabled")
		}

		session2, err := NewSessionWithOptions(WithCompression(true))
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}
		if session2.transport.DisableCompression {
			t.Error("Compression should be enabled")
		}
	})

	t.Run("WithMaxIdleConns", func(t *testing.T) {
		maxConns := 50
		session, err := NewSessionWithOptions(WithMaxIdleConns(maxConns))
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}
		if session.transport.MaxIdleConns != maxConns {
			t.Errorf("Expected MaxIdleConns %d, got %d", maxConns, session.transport.MaxIdleConns)
		}
	})

	t.Run("WithMaxIdleConnsPerHost", func(t *testing.T) {
		maxConns := 10
		session, err := NewSessionWithOptions(WithMaxIdleConnsPerHost(maxConns))
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}
		if session.transport.MaxIdleConnsPerHost != maxConns {
			t.Errorf("Expected MaxIdleConnsPerHost %d, got %d", maxConns, session.transport.MaxIdleConnsPerHost)
		}
	})

	t.Run("WithContext", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "test", "value")
		session, err := NewSessionWithOptions(WithContext(ctx))
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}
		if session.defaultContext != ctx {
			t.Error("Default context was not set correctly")
		}
	})

	t.Run("WithRetry", func(t *testing.T) {
		maxRetries := 3
		backoff := time.Second
		session, err := NewSessionWithOptions(WithRetry(maxRetries, backoff))
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}
		if session.retryConfig == nil {
			t.Fatal("Retry config was not set")
		}
		if session.retryConfig.MaxRetries != maxRetries {
			t.Errorf("Expected MaxRetries %d, got %d", maxRetries, session.retryConfig.MaxRetries)
		}
		if session.retryConfig.Backoff != backoff {
			t.Errorf("Expected Backoff %v, got %v", backoff, session.retryConfig.Backoff)
		}
	})

	t.Run("WithRedirectPolicy", func(t *testing.T) {
		policy := func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
		session, err := NewSessionWithOptions(WithRedirectPolicy(policy))
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}
		if session.client.CheckRedirect == nil {
			t.Error("Redirect policy was not set")
		}
	})

	t.Run("WithDisableRedirects", func(t *testing.T) {
		session, err := NewSessionWithOptions(WithDisableRedirects())
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}
		if session.client.CheckRedirect == nil {
			t.Error("Redirect policy should be set to disable redirects")
		}
	})

	t.Run("WithInsecureSkipVerify", func(t *testing.T) {
		session, err := NewSessionWithOptions(WithInsecureSkipVerify())
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}
		if session.transport.TLSClientConfig == nil {
			t.Fatal("TLS config should be set")
		}
		if !session.transport.TLSClientConfig.InsecureSkipVerify {
			t.Error("InsecureSkipVerify should be true")
		}
	})

	t.Run("WithDialTimeout", func(t *testing.T) {
		timeout := 5 * time.Second
		session, err := NewSessionWithOptions(WithDialTimeout(timeout))
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}
		if session.transport.DialContext == nil {
			t.Error("DialContext should be set")
		}
	})

	t.Run("WithMaxConnsPerHost", func(t *testing.T) {
		maxConns := 100
		session, err := NewSessionWithOptions(WithMaxConnsPerHost(maxConns))
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}
		if session.transport.MaxConnsPerHost != maxConns {
			t.Errorf("Expected MaxConnsPerHost %d, got %d", maxConns, session.transport.MaxConnsPerHost)
		}
	})

	t.Run("WithReadBufferSize", func(t *testing.T) {
		size := 64 * 1024
		session, err := NewSessionWithOptions(WithReadBufferSize(size))
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}
		if session.transport.ReadBufferSize != size {
			t.Errorf("Expected ReadBufferSize %d, got %d", size, session.transport.ReadBufferSize)
		}
	})

	t.Run("WithWriteBufferSize", func(t *testing.T) {
		size := 64 * 1024
		session, err := NewSessionWithOptions(WithWriteBufferSize(size))
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}
		if session.transport.WriteBufferSize != size {
			t.Errorf("Expected WriteBufferSize %d, got %d", size, session.transport.WriteBufferSize)
		}
	})
}

// TestPreDefinedSessions 测试预定义的Session配置
func TestPreDefinedSessions(t *testing.T) {
	t.Run("NewDefaultSession", func(t *testing.T) {
		session, err := NewDefaultSession()
		if err != nil {
			t.Fatalf("Failed to create default session: %v", err)
		}
		if session == nil {
			t.Error("Default session should not be nil")
		}
	})

	t.Run("NewSessionForAPI", func(t *testing.T) {
		session, err := NewSessionForAPI()
		if err != nil {
			t.Fatalf("Failed to create API session: %v", err)
		}
		if session.client.Timeout != 30*time.Second {
			t.Error("API session should have 30 second timeout")
		}
		if session.Header.Get("User-Agent") != "Go-Requests/1.0" {
			t.Error("API session should have correct User-Agent")
		}
	})

	t.Run("NewSessionForScraping", func(t *testing.T) {
		session, err := NewSessionForScraping()
		if err != nil {
			t.Fatalf("Failed to create scraping session: %v", err)
		}
		if session.client.Timeout != 10*time.Second {
			t.Error("Scraping session should have 10 second timeout")
		}
		expectedUA := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
		if session.Header.Get("User-Agent") != expectedUA {
			t.Error("Scraping session should have browser-like User-Agent")
		}
	})

	t.Run("NewSessionForTesting", func(t *testing.T) {
		session, err := NewSessionForTesting()
		if err != nil {
			t.Fatalf("Failed to create testing session: %v", err)
		}
		if session.client.Timeout != 5*time.Second {
			t.Error("Testing session should have 5 second timeout")
		}
		if session.cookiejar != nil {
			t.Error("Testing session should have cookies disabled")
		}
	})

	t.Run("NewSessionWithDefaults", func(t *testing.T) {
		session, err := NewSessionWithDefaults()
		if err != nil {
			t.Fatalf("Failed to create session with defaults: %v", err)
		}
		if session.client.Timeout != 30*time.Second {
			t.Error("Default session should have 30 second timeout")
		}
		if session.Header.Get("User-Agent") != "Go-Requests/2.0" {
			t.Error("Default session should have correct User-Agent")
		}
	})

	t.Run("NewSessionWithRetry", func(t *testing.T) {
		maxRetries := 5
		backoff := 2 * time.Second
		session, err := NewSessionWithRetry(maxRetries, backoff)
		if err != nil {
			t.Fatalf("Failed to create retry session: %v", err)
		}
		if session.retryConfig == nil {
			t.Fatal("Retry session should have retry config")
		}
		if session.retryConfig.MaxRetries != maxRetries {
			t.Errorf("Expected MaxRetries %d, got %d", maxRetries, session.retryConfig.MaxRetries)
		}
		if session.retryConfig.Backoff != backoff {
			t.Errorf("Expected Backoff %v, got %v", backoff, session.retryConfig.Backoff)
		}
	})

	t.Run("NewSessionWithProxy", func(t *testing.T) {
		proxyURL := "http://127.0.0.1:8080"
		session, err := NewSessionWithProxy(proxyURL)
		if err != nil {
			t.Fatalf("Failed to create proxy session: %v", err)
		}
		if session.transport.Proxy == nil {
			t.Error("Proxy session should have proxy set")
		}
	})

	t.Run("NewSecureSession", func(t *testing.T) {
		session, err := NewSecureSession()
		if err != nil {
			t.Fatalf("Failed to create secure session: %v", err)
		}
		if session.cookiejar != nil {
			t.Error("Secure session should have cookies disabled")
		}
		if !session.transport.DisableKeepAlives {
			t.Error("Secure session should have keep-alives disabled")
		}
	})

	t.Run("NewHighPerformanceSession", func(t *testing.T) {
		session, err := NewHighPerformanceSession()
		if err != nil {
			t.Fatalf("Failed to create high performance session: %v", err)
		}
		if session.client.Timeout != 60*time.Second {
			t.Error("High performance session should have 60 second timeout")
		}
		if session.transport.MaxIdleConns != 100 {
			t.Error("High performance session should have 100 max idle connections")
		}
	})
}

// TestSessionOptionsChaining 测试多个选项的链式组合
func TestSessionOptionsChaining(t *testing.T) {
	session, err := NewSessionWithOptions(
		WithTimeout(15*time.Second),
		WithUserAgent("Custom-Agent/1.0"),
		WithKeepAlives(true),
		WithCompression(true),
		WithMaxIdleConnsPerHost(20),
		WithRetry(2, time.Second),
	)
	if err != nil {
		t.Fatalf("Failed to create session with multiple options: %v", err)
	}

	// 验证所有选项都已正确应用
	if session.client.Timeout != 15*time.Second {
		t.Error("Timeout not set correctly")
	}
	if session.Header.Get("User-Agent") != "Custom-Agent/1.0" {
		t.Error("User-Agent not set correctly")
	}
	if session.transport.DisableKeepAlives {
		t.Error("Keep-alives should be enabled")
	}
	if session.transport.DisableCompression {
		t.Error("Compression should be enabled")
	}
	if session.transport.MaxIdleConnsPerHost != 20 {
		t.Error("MaxIdleConnsPerHost not set correctly")
	}
	if session.retryConfig == nil || session.retryConfig.MaxRetries != 2 {
		t.Error("Retry config not set correctly")
	}
}
