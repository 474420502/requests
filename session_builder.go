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

// SessionOption 定义Session配置选项的函数类型。
// 每个配置选项都是一个函数，接收一个Session指针并返回error。
// 这种模式允许在创建Session时进行灵活的配置组合。
//
// 示例:
//
//	session := NewSessionBuilder().
//	  WithTimeout(30*time.Second).
//	  WithUserAgent("MyApp/1.0").
//	  Build()
type SessionOption func(*Session) error

// WithTimeout 设置HTTP客户端的默认超时时间。
//
// 该超时时间应用于整个请求-响应周期，包括连接建立、
// 请求发送、响应接收的全过程。
//
// 参数:
//
//	timeout - 超时时间，如果设置为0则表示无超时限制
//
// 示例:
//
//	session := NewSessionBuilder().
//	  WithTimeout(30*time.Second).  // 30秒超时
//	  Build()
//
// 注意: 这个超时时间会覆盖任何之前设置的超时时间。
func WithTimeout(timeout time.Duration) SessionOption {
	return func(s *Session) error {
		if s.client == nil {
			s.client = &http.Client{}
		}
		s.client.Timeout = timeout
		return nil
	}
}

// WithTLSConfig 设置TLS/SSL连接的自定义配置。
//
// 通过此函数可以配置客户端证书、根证书、加密套件等TLS相关设置。
// 常用于需要客户端证书认证或自定义根证书的场景。
//
// 参数:
//
//	config - TLS配置对象，包含证书、密钥、根证书等信息
//
// 示例:
//
//	tlsConfig := &tls.Config{
//	  InsecureSkipVerify: false,
//	  ClientAuth: tls.RequireAndVerifyClientCert,
//	}
//	session := NewSessionBuilder().
//	  WithTLSConfig(tlsConfig).
//	  Build()
//
// 安全提示: 避免在生产环境中使用 InsecureSkipVerify: true
func WithTLSConfig(config *tls.Config) SessionOption {
	return func(s *Session) error {
		if s.transport == nil {
			s.transport = &http.Transport{}
		}
		s.transport.TLSClientConfig = config
		return nil
	}
}

// WithProxy 设置HTTP代理服务器。
//
// 支持HTTP、HTTPS和SOCKS5代理协议。代理URL应包含完整的协议和地址信息。
//
// 参数:
//
//	proxyURL - 代理服务器地址，格式为 "protocol://host:port"
//
// 支持的代理协议:
//   - HTTP:  "http://proxy.example.com:8080"
//   - HTTPS: "https://proxy.example.com:8080"
//   - SOCKS5: "socks5://proxy.example.com:1080"
//
// 示例:
//
//	session := NewSessionBuilder().
//	  WithProxy("http://proxy.company.com:8080").
//	  Build()
//
// 错误处理: 如果代理URL格式无效，将返回错误
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

// WithBasicAuth 设置HTTP基本认证（Basic Authentication）。
//
// 基本认证会在每个请求的Authorization头中发送用户名和密码，
// 密码会被Base64编码（注意：不是加密，只是编码）。
//
// 参数:
//
//	username - 用户名
//	password - 密码
//
// 示例:
//
//	session := NewSessionBuilder().
//	  WithBasicAuth("admin", "secret123").
//	  Build()
//
// 安全提示:
//   - 基本认证不提供加密，建议仅在HTTPS连接中使用
//   - 避免在日志中记录包含认证信息的请求
func WithBasicAuth(username, password string) SessionOption {
	return func(s *Session) error {
		s.auth = &BasicAuth{
			User:     username,
			Password: password,
		}
		return nil
	}
}

// WithHeaders 设置默认的HTTP请求头。
//
// 这些头部会应用到所有通过该Session发送的请求中。
// 如果单个请求设置了相同名称的头部，则请求级别的设置会覆盖默认设置。
//
// 参数:
//
//	headers - 包含头部名称和值的映射表
//
// 示例:
//
//	headers := map[string]string{
//	  "User-Agent": "MyApp/1.0",
//	  "Accept": "application/json",
//	  "X-API-Key": "your-api-key",
//	}
//	session := NewSessionBuilder().
//	  WithHeaders(headers).
//	  Build()
//
// 注意:
//   - 头部名称不区分大小写，但建议使用标准格式
//   - 多次调用会合并头部，相同名称的会被覆盖
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

// WithUserAgent 设置HTTP请求的User-Agent头部。
//
// User-Agent头部用于标识客户端应用程序、版本、操作系统等信息。
// 许多网站和API会根据User-Agent来返回不同的内容或实施不同的策略。
//
// 参数:
//
//	userAgent - User-Agent字符串，建议包含应用名称和版本
//
// 示例:
//
//	session := NewSessionBuilder().
//	  WithUserAgent("MyApp/1.0 (https://example.com)").
//	  Build()
//
//	// 移动设备模拟
//	session := NewSessionBuilder().
//	  WithUserAgent("Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X)").
//	  Build()
//
// 最佳实践:
//   - 使用描述性的User-Agent，包含应用名称和版本
//   - 遵守目标网站的robots.txt和使用条款
//   - 避免伪装成浏览器进行恶意行为
func WithUserAgent(userAgent string) SessionOption {
	return func(s *Session) error {
		if s.Header == nil {
			s.Header = make(http.Header)
		}
		s.Header.Set("User-Agent", userAgent)
		return nil
	}
}

// WithCookieJar 设置自定义的Cookie管理器。
//
// Cookie jar负责在请求之间自动管理Cookie的存储和发送。
// 默认情况下，Session使用内置的Cookie jar，但可以通过此选项自定义。
//
// 参数:
//
//	jar - 实现http.CookieJar接口的Cookie管理器
//
// 示例:
//
//	// 使用内存Cookie jar
//	jar, _ := cookiejar.New(nil)
//	session := NewSessionBuilder().
//	  WithCookieJar(jar).
//	  Build()
//
//	// 使用自定义Cookie jar（如持久化存储）
//	customJar := &MyPersistentCookieJar{}
//	session := NewSessionBuilder().
//	  WithCookieJar(customJar).
//	  Build()
//
// 用例:
//   - 跨会话持久化Cookie
//   - 多个Session之间共享Cookie
//   - 实现自定义Cookie策略
func WithCookieJar(jar http.CookieJar) SessionOption {
	return func(s *Session) error {
		s.cookiejar = jar
		if s.client != nil {
			s.client.Jar = jar
		}
		return nil
	}
}

// WithDisableCookies 完全禁用Cookie功能。
//
// 禁用Cookie后，Session将不会发送或接收任何Cookie。
// 这对于某些无状态的API调用或有特殊安全要求的场景很有用。
//
// 示例:
//
//	// 创建不使用Cookie的Session
//	session := NewSessionBuilder().
//	  WithDisableCookies().
//	  Build()
//
//	// 适用于RESTful API调用
//	resp, err := session.Get("https://api.example.com/data").
//	  SetBearerToken("your-token").
//	  Execute()
//
// 注意:
//   - 禁用Cookie后无法进行基于Session的身份验证
//   - 某些网站可能需要Cookie才能正常工作
//   - 可以与基于Token的认证方案配合使用
func WithDisableCookies() SessionOption {
	return func(s *Session) error {
		s.cookiejar = nil
		if s.client != nil {
			s.client.Jar = nil
		}
		return nil
	}
}

// WithKeepAlives 设置是否启用HTTP Keep-Alive连接复用。
//
// Keep-Alive允许在同一个TCP连接上发送多个HTTP请求，
// 这可以显著提高性能，特别是在需要向同一服务器发送多个请求时。
//
// 参数:
//
//	enabled - true启用Keep-Alive，false禁用
//
// 示例:
//
//	// 启用Keep-Alive（推荐用于生产环境）
//	session := NewSessionBuilder().
//	  WithKeepAlives(true).
//	  Build()
//
//	// 禁用Keep-Alive（用于特殊场景）
//	session := NewSessionBuilder().
//	  WithKeepAlives(false).
//	  Build()
//
// 性能影响:
//   - 启用: 减少TCP握手开销，提高并发性能
//   - 禁用: 每个请求都建立新连接，增加延迟但减少资源占用
//
// 建议:
//   - 生产环境建议启用
//   - 短生命周期应用或有连接限制时可考虑禁用
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
