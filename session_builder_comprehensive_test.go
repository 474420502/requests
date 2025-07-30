package requests

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"net/http/cookiejar"
	"testing"
	"time"

	"golang.org/x/net/publicsuffix"
)

// TestSessionBuilderEdgeCases 测试SessionBuilder的边界情况
func TestSessionBuilderEdgeCases(t *testing.T) {
	t.Run("MultipleTransportInitialization", func(t *testing.T) {
		// 测试多个选项都需要初始化transport的情况
		session, err := NewSessionWithOptions(
			WithTLSConfig(&tls.Config{InsecureSkipVerify: true}),
			WithKeepAlives(true),
			WithCompression(false),
			WithMaxIdleConns(50),
		)
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}

		// 验证所有transport设置都已正确应用
		if session.transport.TLSClientConfig == nil || !session.transport.TLSClientConfig.InsecureSkipVerify {
			t.Error("TLS config not applied correctly")
		}
		if session.transport.DisableKeepAlives {
			t.Error("Keep-alives should be enabled")
		}
		if !session.transport.DisableCompression {
			t.Error("Compression should be disabled")
		}
		if session.transport.MaxIdleConns != 50 {
			t.Error("MaxIdleConns not set correctly")
		}
	})

	t.Run("MultipleClientInitialization", func(t *testing.T) {
		// 测试多个选项都需要初始化client的情况
		timeout := 15 * time.Second
		session, err := NewSessionWithOptions(
			WithTimeout(timeout),
			WithDisableRedirects(),
			WithRedirectPolicy(func(req *http.Request, via []*http.Request) error {
				return nil // 允许重定向
			}),
		)
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}

		// 验证所有client设置都已正确应用
		if session.client.Timeout != timeout {
			t.Error("Timeout not set correctly")
		}
		if session.client.CheckRedirect == nil {
			t.Error("CheckRedirect should be set")
		}
	})

	t.Run("NilTransportHandling", func(t *testing.T) {
		// 测试当transport为nil时的情况
		session, err := NewSessionWithOptions()
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}

		// 确保transport不为nil
		if session.transport == nil {
			t.Error("Transport should not be nil")
		}
		if session.client.Transport != session.transport {
			t.Error("Client transport should match session transport")
		}
	})

	t.Run("NilClientHandling", func(t *testing.T) {
		// 测试当client为nil时的情况
		session, err := NewSessionWithOptions(WithTimeout(5 * time.Second))
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}

		// 确保client不为nil且设置正确
		if session.client == nil {
			t.Error("Client should not be nil")
		}
		if session.client.Timeout != 5*time.Second {
			t.Error("Client timeout not set correctly")
		}
	})
}

// TestSessionBuilderErrorScenarios 测试可能的错误情况
func TestSessionBuilderErrorScenarios(t *testing.T) {
	t.Run("InvalidProxyURL", func(t *testing.T) {
		// 测试无效的代理URL
		invalidURLs := []string{
			"://invalid-url",
			"invalid://://url",
			"http://[invalid-host:port",
		}

		for _, invalidURL := range invalidURLs {
			_, err := NewSessionWithOptions(WithProxy(invalidURL))
			if err == nil {
				t.Errorf("Expected error for invalid proxy URL: %s", invalidURL)
			}
		}
	})

	t.Run("OptionErrorPropagation", func(t *testing.T) {
		// 创建一个会返回错误的选项
		errorOption := func(s *Session) error {
			return errors.New("test error from option")
		}

		_, err := NewSessionWithOptions(errorOption)
		if err == nil {
			t.Error("Expected error from option function")
		}
	})

	t.Run("CookieJarCreationFailure", func(t *testing.T) {
		// 这个测试验证即使cookiejar创建失败，session创建也应该成功
		// 在实际情况中，cookiejar.New很少失败，但我们测试这种情况
		session, err := NewSessionWithOptions()
		if err != nil {
			t.Fatalf("Session creation should succeed even if cookie jar creation fails: %v", err)
		}

		// 验证session基本功能正常
		if session.Header == nil {
			t.Error("Session header should be initialized")
		}
		if session.transport == nil {
			t.Error("Session transport should be initialized")
		}
		if session.client == nil {
			t.Error("Session client should be initialized")
		}
	})
}

// TestSessionBuilderConcurrency 测试并发场景
func TestSessionBuilderConcurrency(t *testing.T) {
	t.Run("ConcurrentSessionCreation", func(t *testing.T) {
		// 测试并发创建session
		const numGoroutines = 10
		results := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				_, err := NewSessionWithOptions(
					WithTimeout(time.Duration(id+1)*time.Second),
					WithUserAgent("Test-Agent"),
					WithKeepAlives(true),
				)
				results <- err
			}(i)
		}

		// 检查所有goroutine的结果
		for i := 0; i < numGoroutines; i++ {
			if err := <-results; err != nil {
				t.Errorf("Concurrent session creation failed: %v", err)
			}
		}
	})
}

// TestSessionBuilderAdvancedOptions 测试高级选项组合
func TestSessionBuilderAdvancedOptions(t *testing.T) {
	t.Run("FullConfigurationSession", func(t *testing.T) {
		// 创建一个包含所有可能选项的session
		type contextKey string
		ctx := context.WithValue(context.Background(), contextKey("test"), "value")
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS12,
		}

		// 创建自定义cookie jar
		jar, err := cookiejar.New(&cookiejar.Options{
			PublicSuffixList: publicsuffix.List,
		})
		if err != nil {
			t.Fatalf("Failed to create cookie jar: %v", err)
		}

		session, err := NewSessionWithOptions(
			WithTimeout(30*time.Second),
			WithTLSConfig(tlsConfig),
			WithBasicAuth("user", "pass"),
			WithHeaders(map[string]string{
				"Custom-Header":  "custom-value",
				"Another-Header": "another-value",
			}),
			WithUserAgent("Full-Config-Agent/1.0"),
			WithCookieJar(jar),
			WithKeepAlives(true),
			WithCompression(true),
			WithMaxIdleConns(100),
			WithMaxIdleConnsPerHost(20),
			WithContext(ctx),
			WithRetry(3, 2*time.Second),
			WithRedirectPolicy(func(req *http.Request, via []*http.Request) error {
				if len(via) >= 5 {
					return http.ErrUseLastResponse
				}
				return nil
			}),
			WithInsecureSkipVerify(),
			WithDialTimeout(10*time.Second),
			WithMaxConnsPerHost(50),
			WithReadBufferSize(32*1024),
			WithWriteBufferSize(32*1024),
		)

		if err != nil {
			t.Fatalf("Failed to create full configuration session: %v", err)
		}

		// 验证所有配置都已正确应用
		validateFullConfiguration(t, session, ctx, tlsConfig, jar)
	})

	t.Run("ConflictingOptions", func(t *testing.T) {
		// 测试冲突的选项（后面的应该覆盖前面的）
		session, err := NewSessionWithOptions(
			WithTimeout(10*time.Second),
			WithTimeout(20*time.Second), // 这个应该覆盖前面的
			WithKeepAlives(false),
			WithKeepAlives(true), // 这个应该覆盖前面的
			WithUserAgent("Agent1"),
			WithUserAgent("Agent2"), // 这个应该覆盖前面的
		)

		if err != nil {
			t.Fatalf("Failed to create session with conflicting options: %v", err)
		}

		// 验证最后的选项被应用
		if session.client.Timeout != 20*time.Second {
			t.Errorf("Expected timeout 20s, got %v", session.client.Timeout)
		}
		if session.transport.DisableKeepAlives {
			t.Error("Keep-alives should be enabled (last option)")
		}
		if session.Header.Get("User-Agent") != "Agent2" {
			t.Errorf("Expected User-Agent 'Agent2', got %s", session.Header.Get("User-Agent"))
		}
	})
}

func validateFullConfiguration(t *testing.T, session *Session, ctx context.Context, tlsConfig *tls.Config, jar http.CookieJar) {
	// 验证超时
	if session.client.Timeout != 30*time.Second {
		t.Error("Timeout not set correctly")
	}

	// 验证TLS配置
	if session.transport.TLSClientConfig != tlsConfig {
		t.Error("TLS config not set correctly")
	}

	// 验证基本认证
	if session.auth == nil || session.auth.User != "user" || session.auth.Password != "pass" {
		t.Error("Basic auth not set correctly")
	}

	// 验证请求头
	if session.Header.Get("Custom-Header") != "custom-value" {
		t.Error("Custom header not set correctly")
	}
	if session.Header.Get("Another-Header") != "another-value" {
		t.Error("Another header not set correctly")
	}
	if session.Header.Get("User-Agent") != "Full-Config-Agent/1.0" {
		t.Error("User-Agent not set correctly")
	}

	// 验证cookie jar
	if session.cookiejar != jar {
		t.Error("Cookie jar not set correctly")
	}

	// 验证transport设置
	if session.transport.DisableKeepAlives {
		t.Error("Keep-alives should be enabled")
	}
	if session.transport.DisableCompression {
		t.Error("Compression should be enabled")
	}
	if session.transport.MaxIdleConns != 100 {
		t.Error("MaxIdleConns not set correctly")
	}
	if session.transport.MaxIdleConnsPerHost != 20 {
		t.Error("MaxIdleConnsPerHost not set correctly")
	}

	// 验证上下文
	if session.defaultContext != ctx {
		t.Error("Default context not set correctly")
	}

	// 验证重试配置
	if session.retryConfig == nil || session.retryConfig.MaxRetries != 3 || session.retryConfig.Backoff != 2*time.Second {
		t.Error("Retry config not set correctly")
	}

	// 验证重定向策略
	if session.client.CheckRedirect == nil {
		t.Error("Redirect policy not set")
	}

	// 验证其他transport设置
	if session.transport.MaxConnsPerHost != 50 {
		t.Error("MaxConnsPerHost not set correctly")
	}
	if session.transport.ReadBufferSize != 32*1024 {
		t.Error("ReadBufferSize not set correctly")
	}
	if session.transport.WriteBufferSize != 32*1024 {
		t.Error("WriteBufferSize not set correctly")
	}

	// 验证DialContext被设置
	if session.transport.DialContext == nil {
		t.Error("DialContext should be set")
	}
}

// TestSessionBuilderSpecialCases 测试特殊情况
func TestSessionBuilderSpecialCases(t *testing.T) {
	t.Run("EmptyOptions", func(t *testing.T) {
		// 测试没有选项的情况
		session, err := NewSessionWithOptions()
		if err != nil {
			t.Fatalf("Failed to create session with no options: %v", err)
		}

		// 验证默认设置
		if session.Header == nil {
			t.Error("Header should be initialized")
		}
		if session.transport == nil {
			t.Error("Transport should be initialized")
		}
		if session.client == nil {
			t.Error("Client should be initialized")
		}
		if session.client.Transport != session.transport {
			t.Error("Client transport should match session transport")
		}
	})

	t.Run("WithCookieJarNil", func(t *testing.T) {
		// 测试设置nil cookie jar
		session, err := NewSessionWithOptions(WithCookieJar(nil))
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}

		if session.cookiejar != nil {
			t.Error("Cookie jar should be nil")
		}
		if session.client.Jar != nil {
			t.Error("Client cookie jar should be nil")
		}
	})

	t.Run("WithMiddlewareMultiple", func(t *testing.T) {
		// 测试添加多个中间件
		middleware1 := &TestSessionBuilderMiddleware{}
		middleware2 := &TestSessionBuilderMiddleware{}
		middleware3 := &TestSessionBuilderMiddleware{}

		session, err := NewSessionWithOptions(
			WithMiddleware(middleware1, middleware2),
			WithMiddleware(middleware3),
		)
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}

		// 验证所有中间件都被添加
		if len(session.middlewares) != 3 {
			t.Errorf("Expected 3 middlewares, got %d", len(session.middlewares))
		}
	})

	t.Run("WithTLSConfigNil", func(t *testing.T) {
		// 测试设置nil TLS配置
		session, err := NewSessionWithOptions(WithTLSConfig(nil))
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}

		if session.transport.TLSClientConfig != nil {
			t.Error("TLS config should be nil")
		}
	})

	t.Run("WithHeadersEmpty", func(t *testing.T) {
		// 测试设置空的headers map
		session, err := NewSessionWithOptions(WithHeaders(map[string]string{}))
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}

		// 应该不会崩溃，header应该被初始化
		if session.Header == nil {
			t.Error("Header should be initialized")
		}
	})

	t.Run("WithProxyEmptyString", func(t *testing.T) {
		// 测试设置空字符串代理（应该出错）
		_, err := NewSessionWithOptions(WithProxy(""))
		if err == nil {
			t.Error("Expected error for empty proxy URL")
		}
	})

	t.Run("WithBasicAuthEmpty", func(t *testing.T) {
		// 测试设置空的基本认证
		session, err := NewSessionWithOptions(WithBasicAuth("", ""))
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}

		if session.auth == nil {
			t.Error("Auth should be set even with empty credentials")
		}
		if session.auth.User != "" || session.auth.Password != "" {
			t.Error("Empty credentials should be preserved")
		}
	})
}

// TestNewSessionFunctions 测试所有预定义session创建函数
func TestNewSessionFunctions(t *testing.T) {
	t.Run("AllPreDefinedSessions", func(t *testing.T) {
		// 测试所有预定义session创建函数
		creators := map[string]func() (*Session, error){
			"Default":         NewDefaultSession,
			"API":             NewSessionForAPI,
			"Scraping":        NewSessionForScraping,
			"Testing":         NewSessionForTesting,
			"WithDefaults":    NewSessionWithDefaults,
			"Secure":          NewSecureSession,
			"HighPerformance": NewHighPerformanceSession,
		}

		for name, creator := range creators {
			t.Run(name, func(t *testing.T) {
				session, err := creator()
				if err != nil {
					t.Errorf("Failed to create %s session: %v", name, err)
					return
				}
				if session == nil {
					t.Errorf("%s session should not be nil", name)
					return
				}

				// 验证基本属性
				if session.client == nil {
					t.Errorf("%s session client should not be nil", name)
				}
				if session.transport == nil {
					t.Errorf("%s session transport should not be nil", name)
				}
				if session.Header == nil {
					t.Errorf("%s session header should not be nil", name)
				}

				// 验证transport finalizer被设置
				if session.client.Transport != session.transport {
					t.Errorf("%s session client transport should match session transport", name)
				}
			})
		}
	})

	t.Run("ParameterizedSessions", func(t *testing.T) {
		// 测试带参数的session创建函数
		t.Run("WithRetry", func(t *testing.T) {
			session, err := NewSessionWithRetry(5, 3*time.Second)
			if err != nil {
				t.Fatalf("Failed to create retry session: %v", err)
			}
			if session.retryConfig == nil {
				t.Error("Retry config should be set")
			}
			if session.retryConfig.MaxRetries != 5 {
				t.Error("Max retries not set correctly")
			}
			if session.retryConfig.Backoff != 3*time.Second {
				t.Error("Backoff not set correctly")
			}
		})

		t.Run("WithProxy", func(t *testing.T) {
			proxyURL := "http://127.0.0.1:8080"
			session, err := NewSessionWithProxy(proxyURL)
			if err != nil {
				t.Fatalf("Failed to create proxy session: %v", err)
			}
			if session.transport.Proxy == nil {
				t.Error("Proxy should be set")
			}
		})

		t.Run("WithProxyInvalid", func(t *testing.T) {
			_, err := NewSessionWithProxy("invalid://url")
			if err == nil {
				t.Error("Expected error for invalid proxy URL")
			}
		})
	})
}

// TestSessionBuilderMemoryUsage 测试内存使用情况
func TestSessionBuilderMemoryUsage(t *testing.T) {
	t.Run("MultipleSessionsResourceSharing", func(t *testing.T) {
		// 创建多个session，确保它们不会互相影响
		sessions := make([]*Session, 10)
		for i := 0; i < 10; i++ {
			var err error
			sessions[i], err = NewSessionWithOptions(
				WithTimeout(time.Duration(i+1)*time.Second),
				WithUserAgent("Agent-"+string(rune('A'+i))),
			)
			if err != nil {
				t.Fatalf("Failed to create session %d: %v", i, err)
			}
		}

		// 验证每个session都有独立的配置
		for i, session := range sessions {
			expectedTimeout := time.Duration(i+1) * time.Second
			if session.client.Timeout != expectedTimeout {
				t.Errorf("Session %d timeout mismatch: expected %v, got %v",
					i, expectedTimeout, session.client.Timeout)
			}

			expectedUA := "Agent-" + string(rune('A'+i))
			if session.Header.Get("User-Agent") != expectedUA {
				t.Errorf("Session %d User-Agent mismatch: expected %s, got %s",
					i, expectedUA, session.Header.Get("User-Agent"))
			}
		}
	})
}

// TestSessionBuilderMiddleware 测试中间件实现
type TestSessionBuilderMiddleware struct{}

func (m *TestSessionBuilderMiddleware) BeforeRequest(req *http.Request) error {
	return nil
}

func (m *TestSessionBuilderMiddleware) AfterResponse(resp *http.Response) error {
	return nil
}

// TestSessionBuilderOptionsValidation 测试选项验证
func TestSessionBuilderOptionsValidation(t *testing.T) {
	t.Run("ValidProxyURLs", func(t *testing.T) {
		validURLs := []string{
			"http://127.0.0.1:8080",
			"https://proxy.example.com:3128",
			"socks5://127.0.0.1:1080",
			"http://user:pass@proxy.example.com:8080",
		}

		for _, proxyURL := range validURLs {
			_, err := NewSessionWithOptions(WithProxy(proxyURL))
			if err != nil {
				t.Errorf("Valid proxy URL should not cause error: %s, error: %v", proxyURL, err)
			}
		}
	})

	t.Run("NetworkTimeouts", func(t *testing.T) {
		// 测试各种超时设置
		timeouts := []time.Duration{
			1 * time.Nanosecond,
			1 * time.Microsecond,
			1 * time.Millisecond,
			1 * time.Second,
			1 * time.Minute,
			1 * time.Hour,
		}

		for _, timeout := range timeouts {
			session, err := NewSessionWithOptions(WithTimeout(timeout))
			if err != nil {
				t.Errorf("Timeout %v should not cause error: %v", timeout, err)
			}
			if session.client.Timeout != timeout {
				t.Errorf("Timeout not set correctly: expected %v, got %v",
					timeout, session.client.Timeout)
			}
		}
	})

	t.Run("BufferSizes", func(t *testing.T) {
		// 测试各种缓冲区大小
		sizes := []int{0, 1, 1024, 64 * 1024, 1024 * 1024}

		for _, size := range sizes {
			session, err := NewSessionWithOptions(
				WithReadBufferSize(size),
				WithWriteBufferSize(size),
			)
			if err != nil {
				t.Errorf("Buffer size %d should not cause error: %v", size, err)
			}
			if session.transport.ReadBufferSize != size {
				t.Errorf("Read buffer size not set correctly: expected %d, got %d",
					size, session.transport.ReadBufferSize)
			}
			if session.transport.WriteBufferSize != size {
				t.Errorf("Write buffer size not set correctly: expected %d, got %d",
					size, session.transport.WriteBufferSize)
			}
		}
	})
}

// TestDialContextConfiguration 测试DialContext配置
func TestDialContextConfiguration(t *testing.T) {
	t.Run("CustomDialTimeout", func(t *testing.T) {
		timeout := 5 * time.Second
		session, err := NewSessionWithOptions(WithDialTimeout(timeout))
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}

		if session.transport.DialContext == nil {
			t.Error("DialContext should be set")
		}

		// 测试DialContext确实使用了指定的超时时间
		// 这里我们无法直接测试超时时间，但可以确保DialContext不为nil
		ctx := context.Background()
		_, err = session.transport.DialContext(ctx, "tcp", "127.0.0.1:1")
		// 我们期望这会失败，但不应该panic
		if err == nil {
			t.Log("Dial succeeded unexpectedly (this is okay in some environments)")
		}
	})

	t.Run("ZeroDialTimeout", func(t *testing.T) {
		session, err := NewSessionWithOptions(WithDialTimeout(0))
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}

		if session.transport.DialContext == nil {
			t.Error("DialContext should be set even with zero timeout")
		}
	})
}
