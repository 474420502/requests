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

// SetKeepAlives 默认keep alives
func (cfg *Config) SetKeepAlives(is bool) {
	cfg.ses.transport.DisableKeepAlives = !is
}

// SetDecompressNoAccept 设置在没头文件情景下, 接受到压缩数据, 是否要解压. 类型python requests
func (cfg *Config) SetDecompressNoAccept(is bool) {
	cfg.ses.Is.isDecompressNoAccept = is
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

// SetClearBody 设置边界
func (cfg *Config) SetClearBody(is bool) {
	cfg.ses.Is.isClearBodyEvery = is
}

// SetBoundary 设置边界
// func (cfg *Config) SetBoundary(boundary string) {
// 	a := bytes.NewBufferString("")
// 	w := multipart.NewWriter(a)
// 	w.CreateFormFile()
// 	*cfg.ses.boundary = boundary
// }
