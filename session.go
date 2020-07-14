package requests

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"runtime"
	"strings"

	"golang.org/x/net/publicsuffix"
)

// Body 相关参数结构
type Body struct {
	// Query       map[string][]string
	ioBody interface{}
	// prefix ContentType 前缀
	prefix string
	// Files       []UploadFile
	contentTypes map[string]int
}

// NewBody new body pointer
func NewBody() *Body {
	b := &Body{}
	b.contentTypes = make(map[string]int)
	return b
}

// SetIOBody 设置IOBody的值
func (body *Body) SetIOBody(iobody interface{}) {
	body.ioBody = iobody
}

// GetIOBody 获取ioBody值
func (body *Body) GetIOBody() interface{} {
	return body.ioBody
}

// ContentType 获取ContentType
func (body *Body) ContentType() string {
	content := body.prefix
	for kvalue := range body.contentTypes {
		content += kvalue + ";"
	}
	return strings.TrimRight(content, ";")
}

// SetPrefix SetPrefix 和 AddContentType的顺序会影响到ContentType()的返回结果
func (body *Body) SetPrefix(ct string) {
	body.prefix = strings.TrimRight(ct, ";") + ";"
}

// AddContentType 添加 Add Type类型
func (body *Body) AddContentType(ct string) {
	for _, v := range strings.Split(ct, ";") {
		v = strings.Trim(v, " ")
		if v != "" {
			if body.prefix != v {
				body.contentTypes[v] = 1
			}
		}
	}

}

// IBody 相关参数结构
type IBody interface {
	// GetIOBody  获取iobody data
	GetIOBody() interface{}
	// SetIOBody  设置iobody data
	SetIOBody(iobody interface{})
	// ContentType      返回包括 Prefix 所有的ContentType
	ContentType() string
	// AppendContent
	AddContentType(ct string)
	// SetPrefix 设置 Prefix;  唯一前缀; 就是ContentType的第一个, ContentType(Prefix);ContentType;ContentType
	SetPrefix(ct string)
}

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

// Session 的基本方法
type Session struct {
	auth *BasicAuth

	body IBody

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

	// TypeStream application/octet-stream 只能提交一个二进制流, 很少用
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

// TypeConfig 配置类型
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
	return &Session{client: client, body: NewBody(), transport: transport, auth: nil, cookiejar: client.Jar, Header: make(http.Header), Is: IsSetting{true}}
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

// SetHeader 设置set Header的值, 必须符合规范 HaHa -> Haha 如果真要HaHa,只能这样 Ha-Ha
func (ses *Session) SetHeader(header http.Header) {
	ses.Header = header
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
