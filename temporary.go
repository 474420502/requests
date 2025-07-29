package requests

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Temporary 兼容性结构，现在内部使用Request实现
// Deprecated: 使用 Request 代替。Temporary将在未来版本中移除。
type Temporary struct {
	req *Request // 内部使用Request实现

	// 向后兼容性字段 - 这些字段现在主要是占位符，用于保持API兼容
	session   *Session
	ParsedURL *url.URL
	Method    string
	Body      *bytes.Buffer
	Header    http.Header
	Cookies   map[string]*http.Cookie
	FloatPrec int
	err       error

	// 兼容性字段
	acceptEncoding  []AcceptEncodingType
	contentEncoding ContentEncodingType
	mwriter         interface{} // 不再使用，保留以避免编译错误
}

// NewTemporary 创建新的Temporary对象，内部使用Request
// Deprecated: 使用 session.Get/Post/etc 方法代替，它们返回 *Request
func NewTemporary(ses *Session, urlstr string) *Temporary {
	req := NewRequest(ses, "GET", urlstr) // 默认GET，后续会被SetMethod覆盖

	tp := &Temporary{
		req:       req,
		session:   ses,
		Header:    make(http.Header),
		Cookies:   make(map[string]*http.Cookie),
		FloatPrec: 2,
	}

	// 如果URL解析失败，记录错误
	if req.err != nil {
		tp.err = req.err
		return tp
	}

	tp.ParsedURL = req.parsedURL
	tp.Method = req.method

	return tp
}

// 将Temporary的方法重定向到内部的Request实现
// 这样可以保持API兼容性，同时内部使用统一的Request实现

func (tp *Temporary) SetContentType(contentType string) {
	if tp.req != nil {
		tp.req.SetContentType(contentType)
	}
	tp.Header.Set("Content-Type", contentType) // 兼容性
}

func (tp *Temporary) AddHeader(key, value string) *Temporary {
	if tp.req != nil {
		tp.req.AddHeader(key, value)
	}
	tp.Header.Add(key, value) // 兼容性
	return tp
}

func (tp *Temporary) SetHeader(header http.Header) *Temporary {
	if tp.req != nil {
		tp.req.SetHeadersFromHTTP(header)
	}
	// 兼容性：同步到local header
	tp.Header = make(http.Header)
	for k, values := range header {
		tp.Header[k] = append([]string(nil), values...)
	}
	return tp
}

func (tp *Temporary) GetHeader() http.Header {
	if tp.req != nil {
		return tp.req.GetHeader()
	}
	return tp.Header
}

func (tp *Temporary) AddAcceptEncoding(c AcceptEncodingType) {
	tp.acceptEncoding = append(tp.acceptEncoding, c)
}

func (tp *Temporary) GetAcceptEncoding() []AcceptEncodingType {
	return tp.acceptEncoding
}

func (tp *Temporary) MergeHeader(cheader http.Header) {
	if tp.req != nil {
		tp.req.MergeHeader(cheader)
	}
	// 兼容性同步
	for key, values := range cheader {
		for _, v := range values {
			tp.Header.Add(key, v)
		}
	}
}

func (tp *Temporary) DelHeader(key string) *Temporary {
	if tp.req != nil {
		tp.req.DelHeader(key)
	}
	tp.Header.Del(key) // 兼容性
	return tp
}

func (tp *Temporary) SetCookie(cookie *http.Cookie) *Temporary {
	if tp.req != nil {
		tp.req.SetCookie(cookie)
	}
	tp.Cookies[cookie.Name] = cookie // 兼容性
	return tp
}

func (tp *Temporary) SetCookieValue(name, value string) *Temporary {
	if tp.req != nil {
		tp.req.SetCookieValue(name, value)
	}
	tp.Cookies[name] = &http.Cookie{Name: name, Value: value} // 兼容性
	return tp
}

func (tp *Temporary) DelCookie(name interface{}) *Temporary {
	if tp.req != nil {
		tp.req.DelCookie(name)
	}
	// 兼容性处理
	switch v := name.(type) {
	case string:
		delete(tp.Cookies, v)
	case *http.Cookie:
		delete(tp.Cookies, v.Name)
	}
	return tp
}

func (tp *Temporary) SetParsedURL(u *url.URL) *Temporary {
	if tp.req != nil {
		tp.req.SetParsedURL(u)
	}
	tp.ParsedURL = u // 兼容性
	return tp
}

func (tp *Temporary) GetParsedURL() *url.URL {
	if tp.req != nil {
		return tp.req.GetParsedURL()
	}
	return tp.ParsedURL
}

func (tp *Temporary) GetRawURL() string {
	if tp.req != nil {
		return tp.req.GetRawURL()
	}
	if tp.ParsedURL != nil {
		return tp.ParsedURL.String()
	}
	return ""
}

func (tp *Temporary) SetQuery(params url.Values) *Temporary {
	if tp.req != nil {
		tp.req.SetQuery(params)
	}
	return tp
}

func (tp *Temporary) GetQuery() url.Values {
	if tp.req != nil {
		return tp.req.GetQuery()
	}
	if tp.ParsedURL != nil {
		return tp.ParsedURL.Query()
	}
	return make(url.Values)
}

func (tp *Temporary) MergeQuery(query url.Values) *Temporary {
	if tp.req != nil {
		tp.req.MergeQuery(query)
	}
	return tp
}

func (tp *Temporary) GetURLRawPath() string {
	if tp.req != nil {
		return tp.req.GetURLRawPath()
	}
	if tp.ParsedURL != nil {
		return tp.ParsedURL.Path
	}
	return ""
}

func (tp *Temporary) SetURLRawPath(path string) *Temporary {
	if tp.req != nil {
		tp.req.SetURLRawPath(path)
	}
	// 兼容性：同步到ParsedURL
	if tp.ParsedURL != nil {
		if path[0] != '/' {
			tp.ParsedURL.Path = "/" + path
		} else {
			tp.ParsedURL.Path = path
		}
	}
	return tp
}

func (tp *Temporary) GetURLPath() []string {
	if tp.req != nil {
		return tp.req.GetURLPath()
	}
	// 兼容性实现
	if tp.ParsedURL != nil {
		// 这里需要实现与Request相同的逻辑
		path := tp.ParsedURL.Path
		// 简化实现，使用strings.Split
		if path == "" || path == "/" {
			return []string{}
		}
		// 移除开头的/并split
		if path[0] == '/' {
			path = path[1:]
		}
		parts := make([]string, 0)
		for _, part := range strings.Split(path, "/") {
			if part != "" {
				parts = append(parts, "/"+part)
			}
		}
		return parts
	}
	return nil
}

func (tp *Temporary) SetURLPath(path []string) *Temporary {
	if tp.req != nil {
		tp.req.SetURLPath(path)
	}
	// 兼容性实现
	if tp.ParsedURL != nil && path != nil {
		rawpath := ""
		for _, p := range path {
			if p[0] != '/' {
				p = "/" + p
			}
			rawpath += p
		}
		tp.ParsedURL.Path = rawpath
	}
	return tp
}

func (tp *Temporary) QueryParam(key string) IParam {
	if tp.req != nil {
		return tp.req.QueryParam(key)
	}
	// 兼容性：创建一个简单的参数处理器
	return &ParamQuery{req: tp.req, Key: key}
}

func (tp *Temporary) PathParam(regexpGroup string) IParam {
	if tp.req != nil {
		return tp.req.PathParam(regexpGroup)
	}
	// 兼容性实现
	return extractorParam(tp.req, regexpGroup, tp.GetURLRawPath())
}

func (tp *Temporary) HostParam(regexpGroup string) IParam {
	if tp.req != nil {
		return tp.req.HostParam(regexpGroup)
	}
	// 兼容性实现
	host := ""
	if tp.ParsedURL != nil {
		host = tp.ParsedURL.Host
	}
	return extractorParam(tp.req, regexpGroup, host)
}

func (tp *Temporary) SetBody(body io.Reader) *Temporary {
	if tp.req != nil {
		tp.req.SetBodyReader(body)
	}
	// 兼容性：尝试读取到bytes.Buffer中
	if buf, ok := body.(*bytes.Buffer); ok {
		tp.Body = buf
	} else {
		// 创建新的buffer并复制数据
		tp.Body = &bytes.Buffer{}
		io.Copy(tp.Body, body)
	}
	return tp
}

func (tp *Temporary) SetBodyJson(v interface{}) *Temporary {
	if tp.req != nil {
		tp.req.SetBodyJSON(v)
	}
	return tp
}

func (tp *Temporary) SetBodyWithType(contentType string, params interface{}) *Temporary {
	if tp.req != nil {
		tp.req.SetBodyWithType(contentType, params)
	}
	return tp
}

func (tp *Temporary) SetBodyFormData(params ...interface{}) *Temporary {
	if tp.req != nil {
		tp.req.SetBodyFormData(params...)
	}
	return tp
}

func (tp *Temporary) SetBodyUrlencoded(data interface{}) *Temporary {
	if tp.req != nil {
		tp.req.SetBodyUrlencoded(data)
	}
	return tp
}

func (tp *Temporary) SetBodyPlain(params interface{}) *Temporary {
	if tp.req != nil {
		tp.req.SetBodyPlain(params)
	}
	return tp
}

func (tp *Temporary) SetBodyStream(params interface{}) *Temporary {
	if tp.req != nil {
		tp.req.SetBodyStream(params)
	}
	return tp
}

// CreateBodyMultipart 返回一个兼容性的multipart构建器
func (tp *Temporary) CreateBodyMultipart() *MultipartFormData {
	if tp.req != nil {
		return tp.req.CreateBodyMultipart()
	}
	// 兼容性实现
	return &MultipartFormData{}
}

func (tp *Temporary) Error() error {
	if tp.err != nil {
		return tp.err
	}
	if tp.req != nil {
		return tp.req.Error()
	}
	return nil
}

func (tp *Temporary) Execute() (*Response, error) {
	if tp.err != nil {
		return nil, tp.err
	}

	if tp.req == nil {
		return nil, errors.New("internal Request is nil")
	}

	// 确保method被正确设置
	if tp.Method != "" && tp.Method != tp.req.method {
		tp.req.method = tp.Method
	}

	return tp.req.Execute()
}

func (tp *Temporary) BuildRequest() (*http.Request, error) {
	if tp.err != nil {
		return nil, tp.err
	}

	if tp.req == nil {
		return nil, errors.New("internal Request is nil")
	}

	// 使用Request的buildHTTPRequest方法
	return tp.req.buildHTTPRequest()
}

func (tp *Temporary) TestExecute(server ITestServer) (*Response, error) {
	if tp.err != nil {
		return nil, tp.err
	}

	if tp.req == nil {
		return nil, errors.New("internal Request is nil")
	}

	return tp.req.TestExecute(server)
}

// 辅助方法用于设置Method（兼容性需要）
func (tp *Temporary) SetMethod(method string) *Temporary {
	tp.Method = method
	if tp.req != nil {
		tp.req.method = method
	}
	return tp
}
