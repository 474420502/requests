package requests

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"runtime"

	"golang.org/x/net/publicsuffix"
)

// BasicAuth 帐号认真结构
type BasicAuth struct {
	// User 帐号
	User string
	// Password 密码
	Password string
}

// IsSetting 是否设置的一些情景
type IsSetting struct {
	isDecompressNoAccept bool
}

type AcceptEncodingType int

const (
	AcceptEncodingNoCompress AcceptEncodingType = 0
	AcceptEncodingGzip       AcceptEncodingType = 1
	// AcceptEncodingCompress   AcceptEncodingType = 2
	AcceptEncodingDeflate AcceptEncodingType = 3
	AcceptEncodingBr      AcceptEncodingType = 4
)

type ContentEncodingType int

const (
	ContentEncodingNoCompress ContentEncodingType = 0
	ContentEncodingGzip       ContentEncodingType = 1
	// ContentEncodingCompress   ContentEncodingType = 2
	ContentEncodingDeflate ContentEncodingType = 3
	ContentEncodingBr      ContentEncodingType = 4
)

// Session 的基本方法
type Session struct {
	auth *BasicAuth

	acceptEncoding  []AcceptEncodingType
	contentEncoding ContentEncodingType

	client    *http.Client
	cookiejar http.CookieJar

	transport *http.Transport

	Header http.Header
	Query  url.Values

	Is IsSetting
}

const (
	// TypeJSON 类型
	TypeJSON = "application/json"

	// TypeXML 类型
	TypeXML = "text/xml"

	// TypePlain 类型
	TypePlain = "text/plain"

	// TypeHTML 类型
	TypeHTML = "text/html"

	// TypeURLENCODED 类型
	TypeURLENCODED = "application/x-www-form-urlencoded"

	// TypeForm PostForm类型
	TypeForm = TypeURLENCODED

	// TypeStream application/octet-stream 只能提交一个二进制流
	TypeStream = "application/octet-stream"

	// TypeFormData 类型 Upload File 支持path(string) 自动转换成UploadFile
	TypeFormData = "multipart/form-data"

	// TypeMixed Mixed类型
	TypeMixed = "multipart/mixed"

	// HeaderKeyHost Host
	HeaderKeyHost = "Host"

	// HeaderKeyUA User-Agent
	HeaderKeyUA = "User-Agent"

	// HeaderKeyContentType Content-Type
	HeaderKeyContentType = "Content-Type"
)

// TypeConfig config type 配置类型
type TypeConfig int

const (
	_ TypeConfig = iota
	// CRequestTimeout request 包括 dial request redirect 总时间超时
	CRequestTimeout // 支持time.Duration 和 int(秒为单位)

	// CDialTimeout 一个Connect过程的Timeout
	CDialTimeout // 支持time.Duration 和 int(秒为单位)

	// CKeepAlives 默认KeepAlives false, 如果默认为true容易被一直KeepAlives, 没关闭链接
	CKeepAlives

	// CProxy 代理链接
	CProxy // http, https, socks5

	// CInsecure InsecureSkipVerify
	CInsecure // true, false

	// CBasicAuth 帐号认证
	CBasicAuth // user pwd

	// CTLS 帐号认证
	CTLS // user pwd

	// CIsWithCookiejar 持久化 CookieJar true or false ; default = true
	CIsWithCookiejar

	// CIsDecompressNoAccept 解压 当response header 不存在 Accept-Encoding
	// 很多特殊情景会不返回Accept-Encoding: Gzip. 如 不按照标准的网站
	CIsDecompressNoAccept
)

// NewSession 创建Session
func NewSession() *Session {
	client := &http.Client{}
	transport := &http.Transport{DisableCompression: true, DisableKeepAlives: true}

	EnsureTransporterFinalized(transport)

	client.Transport = transport
	cjar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		panic(err)
	}

	client.Jar = cjar
	return &Session{client: client, transport: transport, auth: nil, cookiejar: client.Jar, Header: make(http.Header), Is: IsSetting{true}, acceptEncoding: []AcceptEncodingType{}, contentEncoding: ContentEncodingNoCompress}
}

// Config 配置Reqeusts类集合
func (ses *Session) Config() *Config {
	return &Config{ses: ses}
}

// SetQuery 设置url query的持久参数的值
func (ses *Session) SetQuery(values url.Values) {
	ses.Query = values
}

// GetQuery 获取get query的值
func (ses *Session) GetQuery() url.Values {
	return ses.Query
}

// SetContentType 设置set ContentType
func (ses *Session) SetContentType(contentType string) {
	ses.Header.Set(HeaderKeyContentType, contentType)
}

// SetHeader 设置set Header的值, 必须符合规范 HaHa -> Haha 如果真要HaHa,只能这样 Ha-Ha
func (ses *Session) SetHeader(header http.Header) {
	ses.Header = header
}

// AddHeader  添加 Header的值, 必须符合规范 HaHa -> Haha 如果真要HaHa,只能这样 Ha-Ha
func (ses *Session) AddHeader(key, value string) {
	ses.Header.Add(key, value)
}

// GetHeader 获取get Header的值
func (ses *Session) GetHeader() http.Header {
	return ses.Header
}

// SetCookies 设置Cookies 或者添加Cookies Del
func (ses *Session) SetCookies(u *url.URL, cookies []*http.Cookie) {
	ses.cookiejar.SetCookies(u, cookies)
}

// GetCookies 返回 Cookies
func (ses *Session) GetCookies(u *url.URL) []*http.Cookie {
	return ses.cookiejar.Cookies(u)
}

// DelCookies 删除 Cookies
func (ses *Session) DelCookies(u *url.URL, name string) {
	cookies := ses.cookiejar.Cookies(u)
	for _, c := range cookies {
		if c.Name == name {
			c.MaxAge = -1
			break
		}
	}
	ses.SetCookies(u, cookies)
}

// ClearCookies 清楚所有cookiejar上的cookies
func (ses *Session) ClearCookies() {
	cjar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		panic(err)
	}
	ses.cookiejar = cjar
	ses.client.Jar = ses.cookiejar
}

// Head 请求
func (ses *Session) Head(url string) *Temporary {
	wf := NewTemporary(ses, url)
	wf.Method = "HEAD"
	return wf
}

// Get 请求
func (ses *Session) Get(url string) *Temporary {
	wf := NewTemporary(ses, url)
	wf.Method = "GET"
	return wf
}

// Post 请求
func (ses *Session) Post(url string) *Temporary {
	wf := NewTemporary(ses, url)
	wf.Method = "POST"
	return wf
}

// Put 请求
func (ses *Session) Put(url string) *Temporary {
	wf := NewTemporary(ses, url)
	wf.Method = "PUT"
	return wf
}

// Patch 请求
func (ses *Session) Patch(url string) *Temporary {
	wf := NewTemporary(ses, url)
	wf.Method = "PATCH"
	return wf
}

// Options 请求
func (ses *Session) Options(url string) *Temporary {
	wf := NewTemporary(ses, url)
	wf.Method = "OPTIONS"
	return wf
}

// Delete 请求
func (ses *Session) Delete(url string) *Temporary {
	wf := NewTemporary(ses, url)
	wf.Method = "DELETE"
	return wf
}

// Connect 请求
func (ses *Session) Connect(url string) *Temporary {
	wf := NewTemporary(ses, url)
	wf.Method = "CONNECT"
	return wf
}

// Trace 请求
func (ses *Session) Trace(url string) *Temporary {
	wf := NewTemporary(ses, url)
	wf.Method = "TRACE"
	return wf
}

// // CloseIdleConnections  closes the idle connections that a session client may make use of
// // 从levigross/grequests 借鉴
// func (ses *Session) CloseIdleConnections() {
// 	ses.client.Transport.(*http.Transport).CloseIdleConnections()
// }

// EnsureTransporterFinalized will ensure that when the HTTP client is GCed
// the runtime will close the idle connections (so that they won't leak)
// this function was adopted from Hashicorp's go-cleanhttp package
func EnsureTransporterFinalized(httpTransport *http.Transport) {
	runtime.SetFinalizer(&httpTransport, func(transportInt **http.Transport) {
		(*transportInt).CloseIdleConnections()
	})
}
