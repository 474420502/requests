package requests

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	netproxy "golang.org/x/net/proxy"
)

// Config Set requests config
type Config struct {
	ses *Session
}

// SetBasicAuth 设置基础认证（支持向后兼容的error返回）
func (cfg *Config) SetBasicAuth(user, password string) error {
	if cfg.ses.auth == nil {
		cfg.ses.auth = &BasicAuth{}
	}
	cfg.ses.auth.User = user
	cfg.ses.auth.Password = password
	return nil
}

// SetBasicAuthString 设置基础认证（类型安全方法，不返回error）
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

// SetBasicAuthLegacy 支持多种参数形式的遗留方法
// Deprecated: 使用 SetBasicAuth(user, password string) 或 SetBasicAuthStruct(*BasicAuth) 代替
func (cfg *Config) SetBasicAuthLegacy(args ...interface{}) error {
	if len(args) == 1 {
		switch v := args[0].(type) {
		case *BasicAuth:
			cfg.SetBasicAuthStruct(v)
			return nil
		case BasicAuth:
			cfg.SetBasicAuthStruct(&v)
			return nil
		case nil:
			cfg.ClearBasicAuth()
			return nil
		default:
			return fmt.Errorf("unsupported basic auth type: %T", v)
		}
	} else if len(args) == 2 {
		user, ok := args[0].(string)
		if !ok {
			return fmt.Errorf("first argument must be string, got %T", args[0])
		}
		password, ok := args[1].(string)
		if !ok {
			return fmt.Errorf("second argument must be string, got %T", args[1])
		}
		return cfg.SetBasicAuth(user, password)
	}
	return fmt.Errorf("invalid number of arguments: %d", len(args))
}

// SetTLSConfig 设置TLS配置
func (cfg *Config) SetTLSConfig(tlsconfig *tls.Config) {
	cfg.ses.transport.TLSClientConfig = tlsconfig
}

// SetInsecure 默认 安全(false)
func (cfg *Config) SetInsecure(is bool) {
	cfg.ses.transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: is}
}

// SetProxy 设置代理（支持多种类型以保持向后兼容）
// Deprecated: 使用 SetProxyString(proxyURL string) 代替以获得更好的类型安全性
func (cfg *Config) SetProxy(proxy interface{}) error {
	switch v := proxy.(type) {
	case string:
		return cfg.SetProxyString(v)
	case *url.URL:
		return cfg.setProxyURL(v)
	case nil:
		cfg.ClearProxy()
		return nil
	default:
		return fmt.Errorf("unsupported proxy type: %T", v)
	}
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
	return cfg.setProxyURL(purl)
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

// SetTimeout 设置超时时间（支持多种类型以保持向后兼容）
// Deprecated: 使用 SetTimeoutDuration(time.Duration) 或 SetTimeoutSeconds(int) 代替以获得更好的类型安全性
func (cfg *Config) SetTimeout(t interface{}) error {
	switch v := t.(type) {
	case time.Duration:
		cfg.SetTimeoutDuration(v)
	case int:
		cfg.SetTimeoutSeconds(v)
	case int64:
		cfg.SetTimeoutSeconds(int(v))
	case float32:
		cfg.ses.client.Timeout = time.Duration(v * float32(time.Second))
	case float64:
		cfg.ses.client.Timeout = time.Duration(v * float64(time.Second))
	default:
		return fmt.Errorf("unsupported timeout type: %T", v)
	}
	return nil
}

// SetTimeoutDuration 设置超时时间（推荐使用此方法）
func (cfg *Config) SetTimeoutDuration(timeout time.Duration) {
	cfg.ses.client.Timeout = timeout
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
