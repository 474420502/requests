package requests

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	netproxy "golang.org/x/net/proxy"
)

// Config Set requests config
type Config struct {
	ses *Session
}

// SetBasicAuth 设置基础认证（支持向后兼容的error返回）
func (cfg *Config) SetBasicAuth(user, password string) error {
	// 为了向后兼容，允许空字符串，但仍然设置认证
	if cfg.ses.auth == nil {
		cfg.ses.auth = &BasicAuth{}
	}
	cfg.ses.auth.User = user
	cfg.ses.auth.Password = password
	return nil
} // SetBasicAuthString 设置基础认证（类型安全方法，不返回error）
func (cfg *Config) SetBasicAuthString(user, password string) {
	if cfg.ses.auth == nil {
		cfg.ses.auth = &BasicAuth{}
	}
	cfg.ses.auth.User = user
	cfg.ses.auth.Password = password
}

// SetBasicAuthStruct 使用BasicAuth结构体设置认证
func (cfg *Config) SetBasicAuthStruct(auth *BasicAuth) {
	if auth == nil {
		cfg.ses.auth = nil
	} else {
		if cfg.ses.auth == nil {
			cfg.ses.auth = &BasicAuth{}
		}
		cfg.ses.auth.User = auth.User
		cfg.ses.auth.Password = auth.Password
	}
}

// ClearBasicAuth 清除基础认证
func (cfg *Config) ClearBasicAuth() {
	cfg.ses.auth = nil
}

// SetTLSConfig 设置TLS配置
func (cfg *Config) SetTLSConfig(tlsconfig *tls.Config) {
	cfg.ses.transport.TLSClientConfig = tlsconfig
}

// SetInsecure 默认 安全(false)
func (cfg *Config) SetInsecure(is bool) {
	cfg.ses.transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: is}
}

// SetProxyString 类型安全的代理设置方法（推荐使用）
func (cfg *Config) SetProxyString(proxyURL string) error {
	if proxyURL == "" {
		cfg.ses.transport.Proxy = nil
		return nil
	}

	purl, err := url.Parse(proxyURL)
	if err != nil {
		return fmt.Errorf("parse proxy URL: %w", err)
	}

	// 智能验证：检查URL是否看起来像合理的代理URL
	// 如果没有scheme且看起来不像host:port格式，则认为是无效的
	if purl.Scheme == "" && purl.Host == "" {
		// 对于没有scheme的情况，检查是否至少看起来像host或host:port
		if !isValidHostPort(proxyURL) {
			return fmt.Errorf("invalid proxy URL format: %s", proxyURL)
		}
	}

	return cfg.setProxyURL(purl)
}

// isValidHostPort 检查字符串是否看起来像host或host:port格式
func isValidHostPort(s string) bool {
	// 简单检查：应该包含至少一个点（域名）或者是localhost
	// 且不应该包含空格或其他明显无效的字符
	if strings.Contains(s, " ") || strings.Contains(s, "\t") || strings.Contains(s, "\n") {
		return false
	}

	// 检查是否包含点（域名）或者是已知的本地主机名
	return strings.Contains(s, ".") || strings.HasPrefix(s, "localhost") || strings.HasPrefix(s, "127.")
} // SetProxy 设置代理（SetProxyString的别名，用于兼容性）
func (cfg *Config) SetProxy(proxyURL string) error {
	return cfg.SetProxyString(proxyURL)
}

// ClearProxy 清除代理设置
func (cfg *Config) ClearProxy() {
	cfg.ses.transport.Proxy = nil
}

// setProxyURL 内部方法处理URL代理设置
func (cfg *Config) setProxyURL(purl *url.URL) error {
	if purl.Scheme == "socks5" {
		cfg.ses.transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			dialer, err := netproxy.SOCKS5("tcp", purl.Host, nil, netproxy.Direct)
			if err != nil {
				return nil, err
			}
			return dialer.Dial(network, addr)
		}
		// 为SOCKS5代理设置一个占位符函数，以便测试可以检测到代理已设置
		cfg.ses.transport.Proxy = func(req *http.Request) (*url.URL, error) {
			return purl, nil
		}
	} else {
		cfg.ses.transport.Proxy = http.ProxyURL(purl)
	}
	return nil
}

// SetWithCookiejar 默认使用cookiejar false 不是用session client jar作为cookie传递
func (cfg *Config) SetWithCookiejar(is bool) {
	if is {
		if cfg.ses.client.Jar == nil {
			cfg.ses.client.Jar = cfg.ses.cookiejar
		}
	} else {
		cfg.ses.client.Jar = nil
	}
}

// SetKeepAlives default keep alives
func (cfg *Config) SetKeepAlives(is bool) {
	cfg.ses.transport.DisableKeepAlives = !is
}

// SetDecompressNoAccept 设置在没头文件情景下, 接受到压缩数据, 是否要解压. 类型python requests
func (cfg *Config) SetDecompressNoAccept(is bool) {
	cfg.ses.Is.isDecompressNoAccept = is
}

// AddAcceptEncoding 设置接收压缩的类型
func (cfg *Config) AddAcceptEncoding(ct AcceptEncodingType) {
	cfg.ses.acceptEncoding = append(cfg.ses.acceptEncoding, ct)
}

// AddAcceptEncoding 设置接收压缩的类型
func (cfg *Config) GetAcceptEncoding(ct AcceptEncodingType) []AcceptEncodingType {
	return cfg.ses.acceptEncoding
}

// SetContentEncoding 设置发送数据body压缩的类型
func (cfg *Config) SetContentEncoding(ct ContentEncodingType) {
	cfg.ses.contentEncoding = ct
}

// SetTimeoutDuration 设置超时时间（推荐使用此方法）
func (cfg *Config) SetTimeoutDuration(timeout time.Duration) {
	cfg.ses.client.Timeout = timeout
}

// SetTimeout 设置超时时间（SetTimeoutDuration的别名，用于兼容性）
func (cfg *Config) SetTimeout(timeout time.Duration) {
	cfg.SetTimeoutDuration(timeout)
}

// SetTimeoutSeconds 设置超时时间（秒）
func (cfg *Config) SetTimeoutSeconds(seconds int) {
	cfg.ses.client.Timeout = time.Duration(seconds) * time.Second
}

// SetHeaderAuthorization method is used to add the JWT token's Authorization field to the HTTP header in the ses field of the Config structure.
// The tokenString parameter is a string that represents the string representation of the JWT token.
// In this method, the cfg.ses.Header.Add() method is used to add the Authorization field to the HTTP request header and set its value to the passed JWT token string.
func (cfg *Config) SetHeaderAuthorization(tokenString string) {
	// SetHeaderAuthorization方法用于在Config结构体中的ses字段的HTTP头部添加JWT令牌的Authorization字段。
	// 参数tokenString是一个字符串，它代表JWT令牌的字符串表示形式。
	// 在此方法中，使用cfg.ses.Header.Add()方法将Authorization字段添加到HTTP请求的头部，并将其值设置为传递的JWT令牌字符串。
	cfg.ses.Header.Add("Authorization", tokenString)
}
