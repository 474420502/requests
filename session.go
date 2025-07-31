package requests

import (
	"context"
	"fmt"
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

	// middlewares 中间件列表
	middlewares []Middleware

	// 第三阶段增强功能
	defaultContext context.Context // 默认上下文
	retryConfig    *RetryConfig    // 重试配置
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

// NewSession 创建一个新的HTTP会话实例。
//
// 该函数创建一个配置了默认设置的Session对象，适用于大多数HTTP请求场景。
// 默认配置包括：
//   - 禁用压缩处理（可通过SetCompression手动启用）
//   - 禁用keep-alive连接（可通过SetKeepAlives手动启用）
//   - 启用Cookie jar，自动处理Cookie
//   - 使用系统默认的TLS配置
//
// 返回值:
//
//	*Session - 新创建的Session实例
//
// 示例:
//
//	// 基本用法
//	session := NewSession()
//	resp, err := session.Get("https://api.example.com/users").Execute()
//
//	// 链式调用
//	content, err := NewSession().
//	  Get("https://httpbin.org/json").
//	  Execute().
//	  Content()
//
// 注意:
//   - 如果无法创建Cookie jar（极少见情况），Session仍可正常工作，只是不支持Cookie
//   - 建议在应用中复用Session实例以获得更好的性能
//   - 对于更复杂的配置需求，推荐使用 NewSessionBuilder()
func NewSession() *Session {
	client := &http.Client{}
	transport := &http.Transport{DisableCompression: true, DisableKeepAlives: true}

	EnsureTransporterFinalized(transport)

	client.Transport = transport
	cjar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		// 如果无法创建cookie jar，我们仍然可以创建一个可用的Session，只是没有cookie支持
		// 这比直接panic更好
		client.Jar = nil
	} else {
		client.Jar = cjar
	}

	return &Session{
		client:          client,
		transport:       transport,
		auth:            nil,
		cookiejar:       client.Jar,
		Header:          make(http.Header),
		Is:              IsSetting{true},
		acceptEncoding:  []AcceptEncodingType{},
		contentEncoding: ContentEncodingNoCompress,
		defaultContext:  context.Background(),
		middlewares:     []Middleware{},
	}
}

// GetDefaultContext 获取默认上下文
func (ses *Session) GetDefaultContext() context.Context {
	if ses.defaultContext == nil {
		return context.Background()
	}
	return ses.defaultContext
}

// SetDefaultContext 设置默认上下文
func (ses *Session) SetDefaultContext(ctx context.Context) {
	ses.defaultContext = ctx
}

// GetRetryConfig 获取重试配置
func (ses *Session) GetRetryConfig() *RetryConfig {
	return ses.retryConfig
}

// ClearMiddlewares 清除所有中间件
func (ses *Session) ClearMiddlewares() {
	ses.middlewares = []Middleware{}
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

// ClearCookies 清除所有cookiejar上的cookies
func (ses *Session) ClearCookies() error {
	cjar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return fmt.Errorf("failed to create new cookie jar: %w", err)
	}
	ses.cookiejar = cjar
	ses.client.Jar = ses.cookiejar
	return nil
}

// AddMiddleware 添加中间件到Session
func (ses *Session) AddMiddleware(middleware Middleware) {
	ses.middlewares = append(ses.middlewares, middleware)
}

// SetMiddlewares 设置Session的中间件列表
func (ses *Session) SetMiddlewares(middlewares []Middleware) {
	ses.middlewares = middlewares
}

// GetMiddlewares 获取Session的中间件列表
func (ses *Session) GetMiddlewares() []Middleware {
	return ses.middlewares
}

// Head 请求 - 统一返回 Request 对象
// Head 创建一个HTTP HEAD请求。
//
// HEAD请求获取资源的响应头信息，但不返回响应体内容。
// 常用于检查资源是否存在、获取文件大小、检查最后修改时间等。
//
// 参数:
//
//	url - 请求的目标URL
//
// 返回值:
//
//	*Request - 可进一步配置和执行的请求对象
//
// 示例:
//
//	// 检查文件是否存在
//	resp, err := session.Head("https://example.com/file.pdf").Execute()
//	if err == nil && resp.GetStatusCode() == 200 {
//	  fmt.Println("文件存在")
//	}
//
//	// 获取Content-Length头部
//	resp, _ := session.Head("https://example.com/download").Execute()
//	size := resp.GetHeader("Content-Length")
func (ses *Session) Head(url string) *Request {
	return NewRequest(ses, "HEAD", url)
}

// Get 创建一个HTTP GET请求。
//
// GET是最常用的HTTP方法，用于从服务器获取资源。
// GET请求应该是安全的（不改变服务器状态）和幂等的（多次请求结果相同）。
//
// 参数:
//
//	url - 请求的目标URL
//
// 返回值:
//
//	*Request - 可进一步配置和执行的请求对象
//
// 示例:
//
//	// 简单GET请求
//	resp, err := session.Get("https://api.example.com/users").Execute()
//
//	// 带查询参数的GET请求
//	resp, err := session.Get("https://api.example.com/search").
//	  SetQuery(url.Values{"q": {"golang"}, "limit": {"10"}}).
//	  Execute()
//
//	// 获取JSON响应
//	var result map[string]interface{}
//	err := session.Get("https://api.example.com/data").
//	  Execute().
//	  JSON(&result)
func (ses *Session) Get(url string) *Request {
	return NewRequest(ses, "GET", url)
}

// Post 创建一个HTTP POST请求。
//
// POST方法用于向服务器提交数据，常用于创建资源、提交表单、上传文件等。
// POST请求不是幂等的，多次执行可能产生不同的结果。
//
// 参数:
//
//	url - 请求的目标URL
//
// 返回值:
//
//	*Request - 可进一步配置和执行的请求对象
//
// 示例:
//
//	// 提交JSON数据
//	data := map[string]interface{}{"name": "John", "age": 30}
//	resp, err := session.Post("https://api.example.com/users").
//	  SetJSON(data).
//	  Execute()
//
//	// 提交表单数据
//	form := url.Values{"username": {"john"}, "password": {"secret"}}
//	resp, err := session.Post("https://api.example.com/login").
//	  SetForm(form).
//	  Execute()
//
//	// 上传文件
//	resp, err := session.Post("https://api.example.com/upload").
//	  SetFile("file", "document.pdf").
//	  Execute()
func (ses *Session) Post(url string) *Request {
	return NewRequest(ses, "POST", url)
}

// Put 创建一个HTTP PUT请求。
//
// PUT方法用于更新或创建指定资源的完整表示。
// PUT请求是幂等的，多次执行相同的PUT请求应该产生相同的结果。
//
// 参数:
//
//	url - 请求的目标URL
//
// 返回值:
//
//	*Request - 可进一步配置和执行的请求对象
//
// 示例:
//
//	// 更新用户信息
//	user := map[string]interface{}{
//	  "id": 123,
//	  "name": "John Doe",
//	  "email": "john@example.com",
//	}
//	resp, err := session.Put("https://api.example.com/users/123").
//	  SetJSON(user).
//	  Execute()
//
//	// 替换整个资源
//	resp, err := session.Put("https://api.example.com/documents/456").
//	  SetContentType("text/plain").
//	  SetBody("新的文档内容").
//	  Execute()
func (ses *Session) Put(url string) *Request {
	return NewRequest(ses, "PUT", url)
}

// Patch 创建一个HTTP PATCH请求。
//
// PATCH方法用于对资源进行部分更新，只修改指定的字段。
// 与PUT不同，PATCH不需要提供资源的完整表示。
//
// 参数:
//
//	url - 请求的目标URL
//
// 返回值:
//
//	*Request - 可进一步配置和执行的请求对象
//
// 示例:
//
//	// 部分更新用户信息
//	updates := map[string]interface{}{"email": "newemail@example.com"}
//	resp, err := session.Patch("https://api.example.com/users/123").
//	  SetJSON(updates).
//	  Execute()
//
//	// JSON Patch格式
//	patches := []map[string]interface{}{
//	  {"op": "replace", "path": "/email", "value": "new@example.com"},
//	  {"op": "add", "path": "/phone", "value": "123-456-7890"},
//	}
//	resp, err := session.Patch("https://api.example.com/users/123").
//	  SetContentType("application/json-patch+json").
//	  SetJSON(patches).
//	  Execute()
func (ses *Session) Patch(url string) *Request {
	return NewRequest(ses, "PATCH", url)
}

// Options 请求 - 统一返回 Request 对象
func (ses *Session) Options(url string) *Request {
	return NewRequest(ses, "OPTIONS", url)
}

// Delete 请求 - 统一返回 Request 对象
func (ses *Session) Delete(url string) *Request {
	return NewRequest(ses, "DELETE", url)
}

// Connect 请求 - 统一返回 Request 对象
func (ses *Session) Connect(url string) *Request {
	return NewRequest(ses, "CONNECT", url)
}

// Trace 请求 - 统一返回 Request 对象
func (ses *Session) Trace(url string) *Request {
	return NewRequest(ses, "TRACE", url)
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
