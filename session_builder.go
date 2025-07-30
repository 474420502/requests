package requests

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"golang.org/x/net/publicsuffix"
)

// RetryConfig 重试配置
type RetryConfig struct {
	MaxRetries int
	Backoff    time.Duration
}

// SessionOption 定义Session配置选项
type SessionOption func(*Session) error

// WithTimeout 设置默认超时时间
func WithTimeout(timeout time.Duration) SessionOption {
	return func(s *Session) error {
		if s.client == nil {
			s.client = &http.Client{}
		}
		s.client.Timeout = timeout
		return nil
	}
}

// WithTLSConfig 设置TLS配置
func WithTLSConfig(config *tls.Config) SessionOption {
	return func(s *Session) error {
		if s.transport == nil {
			s.transport = &http.Transport{}
		}
		s.transport.TLSClientConfig = config
		return nil
	}
}

// WithProxy 设置代理
func WithProxy(proxyURL string) SessionOption {
	return func(s *Session) error {
		if s.transport == nil {
			s.transport = &http.Transport{}
		}

		// 检查空字符串
		if proxyURL == "" {
			return fmt.Errorf("proxy URL cannot be empty")
		}

		parsedURL, err := url.Parse(proxyURL)
		if err != nil {
			return err
		}

		// 验证代理URL的scheme必须是支持的协议
		switch parsedURL.Scheme {
		case "http", "https", "socks5":
			// 支持的协议
		case "":
			return fmt.Errorf("proxy URL must have a valid scheme (http, https, or socks5)")
		default:
			return fmt.Errorf("unsupported proxy scheme: %s (supported: http, https, socks5)", parsedURL.Scheme)
		}

		// 验证必须有host
		if parsedURL.Host == "" {
			return fmt.Errorf("proxy URL must have a host")
		}

		s.transport.Proxy = http.ProxyURL(parsedURL)
		return nil
	}
}

// WithBasicAuth 设置基本认证
func WithBasicAuth(username, password string) SessionOption {
	return func(s *Session) error {
		s.auth = &BasicAuth{
			User:     username,
			Password: password,
		}
		return nil
	}
}

// WithHeaders 设置默认请求头
func WithHeaders(headers map[string]string) SessionOption {
	return func(s *Session) error {
		if s.Header == nil {
			s.Header = make(http.Header)
		}
		for k, v := range headers {
			s.Header.Set(k, v)
		}
		return nil
	}
}

// WithUserAgent 设置User-Agent
func WithUserAgent(userAgent string) SessionOption {
	return func(s *Session) error {
		if s.Header == nil {
			s.Header = make(http.Header)
		}
		s.Header.Set("User-Agent", userAgent)
		return nil
	}
}

// WithCookieJar 设置自定义CookieJar
func WithCookieJar(jar http.CookieJar) SessionOption {
	return func(s *Session) error {
		s.cookiejar = jar
		if s.client != nil {
			s.client.Jar = jar
		}
		return nil
	}
}

// WithDisableCookies 禁用Cookie
func WithDisableCookies() SessionOption {
	return func(s *Session) error {
		s.cookiejar = nil
		if s.client != nil {
			s.client.Jar = nil
		}
		return nil
	}
}

// WithKeepAlives 设置是否保持连接
func WithKeepAlives(enabled bool) SessionOption {
	return func(s *Session) error {
		if s.transport == nil {
			s.transport = &http.Transport{}
		}
		s.transport.DisableKeepAlives = !enabled
		return nil
	}
}

// WithCompression 设置是否禁用压缩
func WithCompression(enabled bool) SessionOption {
	return func(s *Session) error {
		if s.transport == nil {
			s.transport = &http.Transport{}
		}
		s.transport.DisableCompression = !enabled
		return nil
	}
}

// WithMaxIdleConns 设置最大空闲连接数
func WithMaxIdleConns(maxConns int) SessionOption {
	return func(s *Session) error {
		if s.transport == nil {
			s.transport = &http.Transport{}
		}
		s.transport.MaxIdleConns = maxConns
		return nil
	}
}

// WithMaxIdleConnsPerHost 设置每个主机的最大空闲连接数
func WithMaxIdleConnsPerHost(maxConns int) SessionOption {
	return func(s *Session) error {
		if s.transport == nil {
			s.transport = &http.Transport{}
		}
		s.transport.MaxIdleConnsPerHost = maxConns
		return nil
	}
}

// WithContext 设置默认上下文
func WithContext(ctx context.Context) SessionOption {
	return func(s *Session) error {
		// 存储默认上下文，在创建Request时使用
		s.defaultContext = ctx
		return nil
	}
}

// WithRetry 设置重试配置
func WithRetry(maxRetries int, backoff time.Duration) SessionOption {
	return func(s *Session) error {
		s.retryConfig = &RetryConfig{
			MaxRetries: maxRetries,
			Backoff:    backoff,
		}
		return nil
	}
}

// WithRedirectPolicy 设置重定向策略
func WithRedirectPolicy(policy func(req *http.Request, via []*http.Request) error) SessionOption {
	return func(s *Session) error {
		if s.client == nil {
			s.client = &http.Client{}
		}
		s.client.CheckRedirect = policy
		return nil
	}
}

// WithDisableRedirects 禁用自动重定向
func WithDisableRedirects() SessionOption {
	return func(s *Session) error {
		if s.client == nil {
			s.client = &http.Client{}
		}
		s.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
		return nil
	}
}

// WithMiddleware 添加中间件
func WithMiddleware(middleware ...Middleware) SessionOption {
	return func(s *Session) error {
		s.middlewares = append(s.middlewares, middleware...)
		return nil
	}
}

// WithInsecureSkipVerify 跳过TLS证书验证
func WithInsecureSkipVerify() SessionOption {
	return func(s *Session) error {
		if s.transport == nil {
			s.transport = &http.Transport{}
		}
		if s.transport.TLSClientConfig == nil {
			s.transport.TLSClientConfig = &tls.Config{}
		}
		s.transport.TLSClientConfig.InsecureSkipVerify = true
		return nil
	}
}

// WithDialTimeout 设置连接超时
func WithDialTimeout(timeout time.Duration) SessionOption {
	return func(s *Session) error {
		if s.transport == nil {
			s.transport = &http.Transport{}
		}
		s.transport.DialContext = (&net.Dialer{
			Timeout: timeout,
		}).DialContext
		return nil
	}
}

// WithMaxConnsPerHost 设置每个主机的最大连接数
func WithMaxConnsPerHost(maxConns int) SessionOption {
	return func(s *Session) error {
		if s.transport == nil {
			s.transport = &http.Transport{}
		}
		s.transport.MaxConnsPerHost = maxConns
		return nil
	}
}

// WithReadBufferSize 设置读缓冲区大小
func WithReadBufferSize(size int) SessionOption {
	return func(s *Session) error {
		if s.transport == nil {
			s.transport = &http.Transport{}
		}
		s.transport.ReadBufferSize = size
		return nil
	}
}

// WithWriteBufferSize 设置写缓冲区大小
func WithWriteBufferSize(size int) SessionOption {
	return func(s *Session) error {
		if s.transport == nil {
			s.transport = &http.Transport{}
		}
		s.transport.WriteBufferSize = size
		return nil
	}
}

// NewSessionWithOptions 使用函数式选项创建Session
func NewSessionWithOptions(opts ...SessionOption) (*Session, error) {
	// 创建默认Session
	session := &Session{
		Header:          make(http.Header),
		Is:              IsSetting{true},
		acceptEncoding:  []AcceptEncodingType{},
		contentEncoding: ContentEncodingNoCompress,
	}

	// 设置默认Transport
	session.transport = &http.Transport{
		DisableCompression: true,
		DisableKeepAlives:  true,
	}

	// 设置默认Client
	session.client = &http.Client{
		Transport: session.transport,
	}

	// 设置默认CookieJar
	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		// 如果无法创建cookie jar，继续但不设置cookies
		session.cookiejar = nil
		session.client.Jar = nil
	} else {
		session.cookiejar = jar
		session.client.Jar = jar
	}

	// 应用选项
	for _, opt := range opts {
		if err := opt(session); err != nil {
			return nil, err
		}
	}

	// 确保Transport被正确设置
	if session.client.Transport == nil {
		session.client.Transport = session.transport
	}

	// 设置finalizer
	EnsureTransporterFinalized(session.transport)

	return session, nil
}

// 提供一些预定义的Session配置

// NewDefaultSession 创建默认Session（向后兼容）
func NewDefaultSession() (*Session, error) {
	return NewSessionWithOptions()
}

// NewSessionForAPI 创建适合API调用的Session
func NewSessionForAPI() (*Session, error) {
	return NewSessionWithOptions(
		WithTimeout(30*time.Second),
		WithUserAgent("Go-Requests/1.0"),
		WithKeepAlives(true),
		WithCompression(true),
		WithMaxIdleConnsPerHost(10),
	)
}

// NewSessionForScraping 创建适合网页抓取的Session
func NewSessionForScraping() (*Session, error) {
	return NewSessionWithOptions(
		WithTimeout(10*time.Second),
		WithUserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
		WithKeepAlives(true),
		WithCompression(true),
		WithMaxIdleConnsPerHost(5),
	)
}

// NewSessionForTesting 创建适合测试的Session
func NewSessionForTesting() (*Session, error) {
	return NewSessionWithOptions(
		WithTimeout(5*time.Second),
		WithUserAgent("Go-Requests-Test/1.0"),
		WithDisableCookies(),
		WithDisableRedirects(),
	)
}

// NewSessionWithDefaults 使用推荐的默认配置创建Session
func NewSessionWithDefaults() (*Session, error) {
	return NewSessionWithOptions(
		WithTimeout(30*time.Second),
		WithUserAgent("Go-Requests/2.0"),
		WithKeepAlives(true),
		WithCompression(true),
		WithMaxIdleConnsPerHost(10),
		WithMaxConnsPerHost(50),
		WithDialTimeout(10*time.Second),
	)
}

// NewSessionWithRetry 创建带重试功能的Session
func NewSessionWithRetry(maxRetries int, backoff time.Duration) (*Session, error) {
	return NewSessionWithOptions(
		WithTimeout(30*time.Second),
		WithUserAgent("Go-Requests-Retry/1.0"),
		WithKeepAlives(true),
		WithCompression(true),
		WithRetry(maxRetries, backoff),
	)
}

// NewSessionWithProxy 创建带代理的Session
func NewSessionWithProxy(proxyURL string) (*Session, error) {
	return NewSessionWithOptions(
		WithTimeout(30*time.Second),
		WithUserAgent("Go-Requests-Proxy/1.0"),
		WithProxy(proxyURL),
		WithKeepAlives(true),
		WithCompression(true),
	)
}

// NewSecureSession 创建安全增强的Session（适用于敏感数据传输）
func NewSecureSession() (*Session, error) {
	return NewSessionWithOptions(
		WithTimeout(30*time.Second),
		WithUserAgent("Go-Requests-Secure/1.0"),
		WithKeepAlives(false), // 不保持连接以提高安全性
		WithCompression(true),
		WithDisableCookies(), // 不使用Cookies
		// 注意：默认不跳过TLS验证以保证安全
	)
}

// NewHighPerformanceSession 创建高性能Session
func NewHighPerformanceSession() (*Session, error) {
	return NewSessionWithOptions(
		WithTimeout(60*time.Second),
		WithUserAgent("Go-Requests-Performance/1.0"),
		WithKeepAlives(true),
		WithCompression(true),
		WithMaxIdleConns(100),
		WithMaxIdleConnsPerHost(20),
		WithMaxConnsPerHost(100),
		WithDialTimeout(5*time.Second),
		WithReadBufferSize(64*1024),  // 64KB
		WithWriteBufferSize(64*1024), // 64KB
	)
}
