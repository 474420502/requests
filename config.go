package requests

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

// Config Set requests config
type Config struct {
	ses *Session
}

// SetBasicAuth 可以 User, Password or *BasicAuth or BasicAuth or nil(clear)
func (cfg *Config) SetBasicAuth(values ...interface{}) {
	if cfg.ses.auth == nil {
		cfg.ses.auth = &BasicAuth{}
	}

	switch len(values) {
	case 1:
		switch v := values[0].(type) {
		case *BasicAuth:
			cfg.ses.auth.User = v.User
			cfg.ses.auth.Password = v.Password
		case BasicAuth:
			cfg.ses.auth.User = v.User
			cfg.ses.auth.Password = v.Password
		case nil:
			cfg.ses.auth = nil
		default:
			panic(errors.New("error type " + reflect.TypeOf(v).String()))
		}
	case 2:
		cfg.ses.auth.User = values[0].(string)
		cfg.ses.auth.Password = values[1].(string)
	}
}

// SetTLSConfig 默认 string or *url.URL or nil
func (cfg *Config) SetTLSConfig(tlsconfig *tls.Config) {
	cfg.ses.transport.TLSClientConfig = tlsconfig
}

// SetInsecure 默认 安全(false)
func (cfg *Config) SetInsecure(is bool) {
	cfg.ses.transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: is}
}

// SetProxy 默认 string or *url.URL or nil
func (cfg *Config) SetProxy(proxy interface{}) {
	switch v := proxy.(type) {
	case string:
		purl, err := (url.Parse(v))
		if err != nil {
			panic(err)
		}
		cfg.ses.transport.Proxy = http.ProxyURL(purl)
	case *url.URL:
		cfg.ses.transport.Proxy = http.ProxyURL(v)
	case nil:
		cfg.ses.transport.Proxy = nil
	default:
		panic(errors.New("error type " + reflect.TypeOf(v).String()))
	}
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

// SetConfig 设置配置
func (cfg *Config) SetTimeout(t interface{}) {
	switch v := t.(type) {
	case time.Duration:
		cfg.ses.client.Timeout = v
	case int:
		cfg.ses.client.Timeout = time.Duration(v * int(time.Second))
	case int64:
		cfg.ses.client.Timeout = time.Duration(v * int64(time.Second))
	case float32:
		cfg.ses.client.Timeout = time.Duration(v * float32(time.Second))
	case float64:
		cfg.ses.client.Timeout = time.Duration(v * float64(time.Second))
	default:
		panic(errors.New("error type " + reflect.TypeOf(v).String()))
	}

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
