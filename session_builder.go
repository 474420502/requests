package requests

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"golang.org/x/net/publicsuffix"
)

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

		parsedURL, err := url.Parse(proxyURL)
		if err != nil {
			return err
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
		// 这里可以存储默认上下文，在Request中使用
		// 暂时不实现，因为需要修改Request结构
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
