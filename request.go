package requests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// Request 替代 Temporary，更符合直觉的命名
// 这是一个请求构建器，支持链式调用和健壮的错误处理
type Request struct {
	session *Session
	ctx     context.Context

	method    string
	parsedURL *url.URL
	header    http.Header
	cookies   map[string]*http.Cookie
	body      *bytes.Buffer

	// 错误处理：在链式调用中累积错误
	err error

	// 配置选项
	timeout time.Duration

	// 中间件支持
	middlewares []Middleware
}

// NewRequest 创建一个新的请求构建器
func NewRequest(session *Session, method, urlStr string) *Request {
	req := &Request{
		session: session,
		method:  method,
		header:  make(http.Header),
		cookies: make(map[string]*http.Cookie),
		ctx:     context.Background(),
		// 继承Session的中间件
		middlewares: append([]Middleware{}, session.middlewares...),
	}

	parsedURL, err := url.ParseRequestURI(urlStr)
	if err != nil {
		req.err = fmt.Errorf("invalid URL: %w", err)
		return req
	}
	req.parsedURL = parsedURL

	return req
}

// WithContext 设置请求的上下文，支持超时和取消
func (r *Request) WithContext(ctx context.Context) *Request {
	if r.err != nil {
		return r
	}
	r.ctx = ctx
	return r
}

// WithTimeout 设置请求超时
func (r *Request) WithTimeout(timeout time.Duration) *Request {
	if r.err != nil {
		return r
	}
	r.timeout = timeout
	return r
}

// SetHeader 设置单个请求头
func (r *Request) SetHeader(key, value string) *Request {
	if r.err != nil {
		return r
	}
	r.header.Set(key, value)
	return r
}

// SetHeaders 批量设置请求头
func (r *Request) SetHeaders(headers map[string]string) *Request {
	if r.err != nil {
		return r
	}
	for k, v := range headers {
		r.header.Set(k, v)
	}
	return r
}

// SetHeadersFromHTTP 从http.Header设置请求头
func (r *Request) SetHeadersFromHTTP(headers http.Header) *Request {
	if r.err != nil {
		return r
	}
	// 清空现有头部
	r.header = make(http.Header)
	// 复制新头部
	for k, values := range headers {
		r.header[k] = append([]string(nil), values...)
	}
	return r
}

// AddHeader 添加请求头（不覆盖已有值）
func (r *Request) AddHeader(key, value string) *Request {
	if r.err != nil {
		return r
	}
	r.header.Add(key, value)
	return r
}

// DelHeader 删除请求头
func (r *Request) DelHeader(key string) *Request {
	if r.err != nil {
		return r
	}
	r.header.Del(key)
	return r
}

// GetHeader 获取请求头
func (r *Request) GetHeader() http.Header {
	return r.header
}

// SetContentType 设置 Content-Type
func (r *Request) SetContentType(contentType string) *Request {
	return r.SetHeader("Content-Type", contentType)
}

// SetCookie 设置Cookie
func (r *Request) SetCookie(cookie *http.Cookie) *Request {
	if r.err != nil {
		return r
	}
	r.cookies[cookie.Name] = cookie
	return r
}

// SetCookieValue 设置Cookie键值对
func (r *Request) SetCookieValue(name, value string) *Request {
	if r.err != nil {
		return r
	}
	r.cookies[name] = &http.Cookie{Name: name, Value: value}
	return r
}

// AddCookies 批量添加Cookie
func (r *Request) AddCookies(cookies []*http.Cookie) *Request {
	if r.err != nil {
		return r
	}
	for _, cookie := range cookies {
		r.cookies[cookie.Name] = cookie
	}
	return r
}

// DelCookie 删除Cookie
func (r *Request) DelCookie(name interface{}) *Request {
	if r.err != nil {
		return r
	}
	switch v := name.(type) {
	case string:
		delete(r.cookies, v)
	case *http.Cookie:
		delete(r.cookies, v.Name)
	}
	return r
}

// SetParsedURL 设置解析后的URL
func (r *Request) SetParsedURL(u *url.URL) *Request {
	if r.err != nil {
		return r
	}
	r.parsedURL = u
	return r
}

// MergeHeader 合并Header
func (r *Request) MergeHeader(header http.Header) *Request {
	if r.err != nil {
		return r
	}
	if r.header == nil {
		r.header = make(http.Header)
	}
	for key, values := range header {
		for _, value := range values {
			r.header.Add(key, value)
		}
	}
	return r
}

// SetBody 设置请求体 (兼容性方法)
func (r *Request) SetBody(body io.Reader) *Request {
	return r.SetBodyReader(body)
}

// SetBodyJson 设置JSON请求体 (兼容性方法)
func (r *Request) SetBodyJson(v interface{}) *Request {
	return r.SetBodyJSON(v)
}

// SetBodyWithType 设置指定类型的请求体 (兼容性方法)
func (r *Request) SetBodyWithType(contentType string, params interface{}) *Request {
	if r.err != nil {
		return r
	}

	r.SetContentType(contentType)

	if params == nil {
		return r
	}

	switch param := params.(type) {
	case string:
		return r.SetBodyString(param)
	case []byte:
		return r.SetBodyBytes(param)
	case []rune:
		return r.SetBodyString(string(param))
	default:
		r.err = fmt.Errorf("SetBodyWithType only supports string, []byte, []rune, got %T", params)
		return r
	}
}

// CreateBodyMultipart 创建multipart表单数据 (兼容性方法)
func (r *Request) CreateBodyMultipart() *MultipartFormData {
	mpfd := &MultipartFormData{}
	mpfd.writer = multipart.NewWriter(&mpfd.data)
	return mpfd
}

// SetQuery 设置查询参数
func (r *Request) SetQuery(params url.Values) *Request {
	if r.err != nil {
		return r
	}
	r.parsedURL.RawQuery = params.Encode()
	return r
}

// AddQuery 添加查询参数
func (r *Request) AddQuery(key, value string) *Request {
	if r.err != nil {
		return r
	}
	q := r.parsedURL.Query()
	q.Add(key, value)
	r.parsedURL.RawQuery = q.Encode()
	return r
}

// AddParam 添加URL参数 (AddQuery的别名)
func (r *Request) AddParam(key, value string) *Request {
	return r.AddQuery(key, value)
}

// SetParam 设置URL参数（会覆盖已存在的）
func (r *Request) SetParam(key, value string) *Request {
	if r.err != nil {
		return r
	}
	q := r.parsedURL.Query()
	q.Set(key, value)
	r.parsedURL.RawQuery = q.Encode()
	return r
}

// DelParam 删除URL参数
func (r *Request) DelParam(key string) *Request {
	if r.err != nil {
		return r
	}
	q := r.parsedURL.Query()
	q.Del(key)
	r.parsedURL.RawQuery = q.Encode()
	return r
}

// SetTimeout 设置请求超时时间 (WithTimeout的别名)
func (r *Request) SetTimeout(timeout time.Duration) *Request {
	return r.WithTimeout(timeout)
}

// GetQuery 获取查询参数
func (r *Request) GetQuery() url.Values {
	if r.parsedURL == nil {
		return make(url.Values)
	}
	return r.parsedURL.Query()
}

// MergeQuery 合并查询参数
func (r *Request) MergeQuery(query url.Values) *Request {
	if r.err != nil {
		return r
	}
	q := r.parsedURL.Query()
	for key, values := range query {
		for _, value := range values {
			q.Add(key, value)
		}
	}
	r.parsedURL.RawQuery = q.Encode()
	return r
}

// GetRawURL 获取原始URL字符串
func (r *Request) GetRawURL() string {
	if r.parsedURL == nil {
		return ""
	}
	return r.parsedURL.String()
}

// SetRawURL 设置原始URL字符串
func (r *Request) SetRawURL(srcURL string) *Request {
	if r.err != nil {
		return r
	}

	parsedURL, err := url.ParseRequestURI(srcURL)
	if err != nil {
		r.err = fmt.Errorf("invalid URL: %w", err)
		return r
	}
	r.parsedURL = parsedURL
	return r
}

// GetParsedURL 获取解析后的URL
func (r *Request) GetParsedURL() *url.URL {
	return r.parsedURL
}

// GetURLRawPath 获取URL路径
func (r *Request) GetURLRawPath() string {
	if r.parsedURL == nil {
		return ""
	}
	return r.parsedURL.Path
}

// SetURLRawPath 设置URL路径
func (r *Request) SetURLRawPath(path string) *Request {
	if r.err != nil {
		return r
	}
	if r.parsedURL == nil {
		r.err = fmt.Errorf("URL not initialized")
		return r
	}
	if path[0] != '/' {
		r.parsedURL.Path = "/" + path
	} else {
		r.parsedURL.Path = path
	}
	return r
}

// GetURLPath 获得URL路径分段
func (r *Request) GetURLPath() []string {
	if r.parsedURL == nil {
		return nil
	}
	regexGetPath := regexp.MustCompile("/[^/]*")
	return regexGetPath.FindAllString(r.parsedURL.Path, -1)
}

// SetURLPath 设置URL路径分段
func (r *Request) SetURLPath(path []string) *Request {
	if r.err != nil {
		return r
	}
	if r.parsedURL == nil {
		r.err = fmt.Errorf("URL not initialized")
		return r
	}
	if path == nil {
		return r
	}
	rawpath := ""
	for _, p := range path {
		if p[0] != '/' {
			p = "/" + p
		}
		rawpath += p
	}
	r.parsedURL.Path = rawpath
	return r
}

// QueryParam 获取查询参数处理器（为向后兼容性提供）
func (r *Request) QueryParam(key string) IParam {
	// 创建一个临时的Temporary对象用于参数处理
	temp := &Temporary{
		ParsedURL: r.parsedURL,
	}
	return &ParamQuery{Temp: temp, Key: key}
}

// PathParam 路径参数处理器（为向后兼容性提供）
func (r *Request) PathParam(regexpGroup string) IParam {
	temp := &Temporary{
		ParsedURL: r.parsedURL,
	}
	return extractorParam(temp, regexpGroup, r.parsedURL.Path)
}

// HostParam 主机参数处理器（为向后兼容性提供）
func (r *Request) HostParam(regexpGroup string) IParam {
	temp := &Temporary{
		ParsedURL: r.parsedURL,
	}
	return extractorParam(temp, regexpGroup, r.parsedURL.Host)
}

// WithMiddleware 添加中间件
func (r *Request) WithMiddleware(middleware ...Middleware) *Request {
	if r.err != nil {
		return r
	}
	r.middlewares = append(r.middlewares, middleware...)
	return r
}

// SetBodyReader 设置请求体
func (r *Request) SetBodyReader(body io.Reader) *Request {
	if r.err != nil {
		return r
	}

	buf := bytes.NewBuffer(nil)
	_, err := io.Copy(buf, body)
	if err != nil {
		r.err = fmt.Errorf("failed to read body: %w", err)
		return r
	}
	r.body = buf
	return r
}

// SetBodyString 设置字符串请求体
func (r *Request) SetBodyString(body string) *Request {
	if r.err != nil {
		return r
	}
	r.body = bytes.NewBufferString(body)
	return r
}

// SetBodyBytes 设置字节请求体
func (r *Request) SetBodyBytes(body []byte) *Request {
	if r.err != nil {
		return r
	}
	r.body = bytes.NewBuffer(body)
	return r
}

// SetBodyJSON 设置JSON请求体
func (r *Request) SetBodyJSON(v interface{}) *Request {
	if r.err != nil {
		return r
	}

	if v == nil {
		r.body = nil
		return r.SetHeader("Content-Type", "application/json")
	}

	// 处理已经是字符串或字节的情况
	switch data := v.(type) {
	case string:
		r.body = bytes.NewBufferString(data)
	case []byte:
		r.body = bytes.NewBuffer(data)
	case []rune:
		r.body = bytes.NewBuffer([]byte(string(data)))
	default:
		// 需要序列化的类型
		jsonData, err := json.Marshal(v)
		if err != nil {
			r.err = fmt.Errorf("failed to marshal JSON: %w", err)
			return r
		}
		r.body = bytes.NewBuffer(jsonData)
	}

	return r.SetHeader("Content-Type", "application/json")
}

// SetBodyForm 设置表单请求体
func (r *Request) SetBodyForm(values url.Values) *Request {
	if r.err != nil {
		return r
	}

	r.body = bytes.NewBufferString(values.Encode())
	return r.SetHeader("Content-Type", "application/x-www-form-urlencoded")
}

// SetBodyFormValues 设置表单键值对
func (r *Request) SetBodyFormValues(values map[string]string) *Request {
	if r.err != nil {
		return r
	}

	form := make(url.Values)
	for k, v := range values {
		form.Set(k, v)
	}
	return r.SetBodyForm(form)
}

// SetBodyPlain 设置纯文本请求体
func (r *Request) SetBodyPlain(data interface{}) *Request {
	if r.err != nil {
		return r
	}

	var content string
	switch v := data.(type) {
	case string:
		content = v
	case []byte:
		content = string(v)
	case []rune:
		content = string(v)
	default:
		r.err = fmt.Errorf("SetBodyPlain only supports string, []byte, []rune, got %T", data)
		return r
	}

	r.body = bytes.NewBufferString(content)
	return r.SetHeader("Content-Type", "text/plain")
}

// SetBodyStream 设置流式请求体
func (r *Request) SetBodyStream(data interface{}) *Request {
	if r.err != nil {
		return r
	}

	var content []byte
	switch v := data.(type) {
	case string:
		content = []byte(v)
	case []byte:
		content = v
	case []rune:
		content = []byte(string(v))
	default:
		r.err = fmt.Errorf("SetBodyStream only supports string, []byte, []rune, got %T", data)
		return r
	}

	r.body = bytes.NewBuffer(content)
	return r.SetHeader("Content-Type", "application/octet-stream")
}

// SetBodyUrlencoded 设置URL编码请求体（兼容Temporary的方法名）
func (r *Request) SetBodyUrlencoded(data interface{}) *Request {
	if r.err != nil {
		return r
	}

	switch v := data.(type) {
	case url.Values:
		return r.SetBodyForm(v)
	case map[string][]string:
		return r.SetBodyForm(url.Values(v))
	case map[string]string:
		return r.SetBodyFormValues(v)
	case string:
		r.body = bytes.NewBufferString(v)
		return r.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	case []byte:
		r.body = bytes.NewBuffer(v)
		return r.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	case []rune:
		r.body = bytes.NewBuffer([]byte(string(v)))
		return r.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	default:
		r.err = fmt.Errorf("SetBodyUrlencoded supports url.Values, map[string][]string, map[string]string, string, []byte, []rune, got %T", data)
		return r
	}
}

// SetBodyFormData 设置多部分表单数据（兼容Temporary，但简化实现）
func (r *Request) SetBodyFormData(params ...interface{}) *Request {
	if r.err != nil {
		return r
	}

	// 简化实现：转换为FormFile切片
	var files []FormFile
	var hasData bool

	for i, param := range params {
		switch v := param.(type) {
		case map[string]string:
			// 将键值对转换为表单字段
			for key, value := range v {
				files = append(files, FormFile{
					FieldName: key,
					FileName:  "", // 空文件名表示这是一个表单字段
					Reader:    strings.NewReader(value),
				})
			}
			hasData = true
		case string:
			// 假设是文件路径或文件内容
			files = append(files, FormFile{
				FieldName: fmt.Sprintf("file%d", i),
				FileName:  "file.txt",
				Reader:    strings.NewReader(v),
			})
			hasData = true
		default:
			r.err = fmt.Errorf("SetBodyFormData unsupported parameter type: %T", v)
			return r
		}
	}

	if !hasData {
		r.err = fmt.Errorf("SetBodyFormData requires at least one parameter")
		return r
	}

	return r.SetBodyFormFiles(files...)
}

// SetBodyFormFiles 设置多部分表单（文件上传）
func (r *Request) SetBodyFormFiles(files ...FormFile) *Request {
	if r.err != nil {
		return r
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for _, file := range files {
		err := r.writeFormFile(writer, file)
		if err != nil {
			r.err = fmt.Errorf("failed to write form file: %w", err)
			return r
		}
	}

	err := writer.Close()
	if err != nil {
		r.err = fmt.Errorf("failed to close multipart writer: %w", err)
		return r
	}

	r.body = body
	return r.SetHeader("Content-Type", "multipart/form-data; boundary="+writer.Boundary())
}

// FormFile 表示一个表单文件
type FormFile struct {
	FieldName string
	FileName  string
	Reader    io.Reader
}

// writeFormFile 写入表单文件
func (r *Request) writeFormFile(writer *multipart.Writer, file FormFile) error {
	part, err := writer.CreateFormFile(file.FieldName, file.FileName)
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file.Reader)
	return err
}

// Execute 执行请求，返回响应和错误。统一处理中间件逻辑
func (r *Request) Execute() (*Response, error) {
	if r.err != nil {
		return nil, fmt.Errorf("request validation failed: %w", r.err)
	}

	// 构建HTTP请求
	httpReq, err := r.buildHTTPRequest()
	if err != nil {
		return nil, fmt.Errorf("failed to build HTTP request: %w", err)
	}

	// 执行BeforeRequest中间件
	for _, middleware := range r.middlewares {
		if err := middleware.BeforeRequest(httpReq); err != nil {
			return nil, fmt.Errorf("middleware BeforeRequest failed: %w", err)
		}
	}

	// 应用超时
	if r.timeout > 0 {
		ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
		defer cancel()
		httpReq = httpReq.WithContext(ctx)
	} else if r.ctx != context.Background() {
		httpReq = httpReq.WithContext(r.ctx)
	}

	// 执行请求
	resp, err := r.session.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request execution failed: %w", err)
	}

	// 执行AfterResponse中间件
	for _, middleware := range r.middlewares {
		if err := middleware.AfterResponse(resp); err != nil {
			return nil, fmt.Errorf("middleware AfterResponse failed: %w", err)
		}
	}

	// 转换为我们的Response类型
	myResponse, err := FromHTTPResponse(resp, r.session.Is.isDecompressNoAccept)
	if err != nil {
		return nil, fmt.Errorf("failed to process response: %w", err)
	}

	myResponse.readResponse = resp
	return myResponse, nil
}

// buildHTTPRequest 构建标准库的HTTP请求
func (r *Request) buildHTTPRequest() (*http.Request, error) {
	var bodyReader io.Reader
	if r.body != nil {
		bodyReader = r.body
	}

	req, err := http.NewRequest(r.method, r.parsedURL.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	// 设置会话级别的头部
	for key, values := range r.session.Header {
		req.Header[key] = values
	}

	// 设置请求级别的头部
	for key, values := range r.header {
		req.Header[key] = values
	}

	// 设置Cookie
	for _, cookie := range r.cookies {
		req.AddCookie(cookie)
	}

	// 设置基本认证
	if r.session.auth != nil {
		req.SetBasicAuth(r.session.auth.User, r.session.auth.Password)
	}

	return req, nil
}

// Error 返回累积的错误（如果有）
func (r *Request) Error() error {
	return r.err
}

// Text 发送请求并返回响应内容的字符串形式
func (r *Request) Text() (string, error) {
	resp, err := r.Execute()
	if err != nil {
		return "", err
	}
	defer resp.readResponse.Body.Close()
	return resp.ContentString(), nil
}

// JSON 发送请求并将响应内容解析为JSON
func (r *Request) JSON(v interface{}) error {
	resp, err := r.Execute()
	if err != nil {
		return err
	}
	defer resp.readResponse.Body.Close()
	return json.Unmarshal(resp.Content(), v)
}

// Bytes 发送请求并返回响应内容的字节切片
func (r *Request) Bytes() ([]byte, error) {
	resp, err := r.Execute()
	if err != nil {
		return nil, err
	}
	defer resp.readResponse.Body.Close()
	return resp.Content(), nil
}

// TestExecute 使用测试服务器执行请求（用于测试）
func (r *Request) TestExecute(server ITestServer) (*Response, error) {
	if r.err != nil {
		return nil, r.err
	}

	req, err := r.buildHTTPRequest()
	if err != nil {
		return nil, fmt.Errorf("failed to build HTTP request: %w", err)
	}

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	resp, err := FromHTTPResponse(w.Result(), false)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// TestExecuteWithDecompress 使用测试服务器执行请求并自动解压（用于测试）
func (r *Request) TestExecuteWithDecompress(server ITestServer) (*Response, error) {
	if r.err != nil {
		return nil, r.err
	}

	req, err := r.buildHTTPRequest()
	if err != nil {
		return nil, fmt.Errorf("failed to build HTTP request: %w", err)
	}

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	resp, err := FromHTTPResponse(w.Result(), true)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
