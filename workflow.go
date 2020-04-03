package requests

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// Workflow 工作流 设计点: 这个并不影响Session的属性变化 如 NewWorkflow(ses, url).AddHeader() 对ses没影响
type Workflow struct {
	session   *Session
	ParsedURL *url.URL
	Method    string
	Body      IBody
	Header    http.Header
	Cookies   map[string]*http.Cookie
}

// NewWorkflow new and init workflow
func NewWorkflow(ses *Session, urlstr string) *Workflow {
	wf := &Workflow{}
	wf.SwitchSession(ses)
	wf.SetRawURL(urlstr)

	wf.Body = NewBody()
	wf.Header = make(http.Header)
	wf.Cookies = make(map[string]*http.Cookie)
	return wf
}

// SwitchSession 替换Session
func (wf *Workflow) SwitchSession(ses *Session) {
	wf.session = ses
}

// AddHeader 添加头信息  Get方法从Header参数上获取 必须符合规范 HaHa -> Haha 如果真要HaHa,只能这样 Ha-Ha
func (wf *Workflow) AddHeader(key, value string) *Workflow {
	wf.Header[key] = append(wf.Header[key], value)
	return wf
}

// SetHeader 设置完全替换原有Header 必须符合规范 HaHa -> Haha 如果真要HaHa,只能这样 Ha-Ha
func (wf *Workflow) SetHeader(header http.Header) *Workflow {
	wf.Header = make(http.Header)
	for k, HValues := range header {
		var newHValues []string
		for _, value := range HValues {
			newHValues = append(newHValues, value)
		}
		wf.Header[k] = newHValues
	}
	return wf
}

// GetHeader 获取Workflow Header
func (wf *Workflow) GetHeader() http.Header {
	return wf.Header
}

// GetCombineHeader 获取后的Header信息
func (wf *Workflow) GetCombineHeader() http.Header {
	return mergeMapList(wf.session.Header, wf.Header)
}

// DelHeader 添加头信息 Get方法从Header参数上获取
func (wf *Workflow) DelHeader(key string) *Workflow {
	wf.Header.Del(key)
	return wf
}

// AddCookie 添加Cookie
func (wf *Workflow) AddCookie(c *http.Cookie) *Workflow {
	wf.Cookies[c.Name] = c
	return wf
}

// AddCookies 添加[]*http.Cookie
func (wf *Workflow) AddCookies(cookies []*http.Cookie) *Workflow {
	for _, c := range cookies {
		wf.AddCookie(c)
	}
	return wf
}

// AddKVCookie 添加 以 key value 的 Cookie
func (wf *Workflow) AddKVCookie(name, value string) *Workflow {
	wf.Cookies[name] = &http.Cookie{Name: name, Value: value}
	return wf
}

// DelCookie 删除Cookie
func (wf *Workflow) DelCookie(name interface{}) *Workflow {
	switch n := name.(type) {
	case string:
		if _, ok := wf.Cookies[n]; ok {
			delete(wf.Cookies, n)
			return wf
		}
	case *http.Cookie:
		if _, ok := wf.Cookies[n.Name]; ok {
			delete(wf.Cookies, n.Name)
			return wf
		}
	default:
		panic("name type is not support")
	}
	return nil
}

// GetParsedURL 获取url的string形式
func (wf *Workflow) GetParsedURL() *url.URL {
	return wf.ParsedURL
}

// SetParsedURL 获取url的string形式
func (wf *Workflow) SetParsedURL(u *url.URL) *Workflow {
	wf.ParsedURL = u
	return wf
}

// GetRawURL 获取url的string形式
func (wf *Workflow) GetRawURL() string {
	u := strings.Split(wf.ParsedURL.String(), "?")[0] + "?" + wf.GetCombineQuery().Encode()
	return u
}

// SetRawURL 设置 url
func (wf *Workflow) SetRawURL(srcURL string) *Workflow {
	purl, err := url.ParseRequestURI(srcURL)
	if err != nil {
		panic(err)
	}
	wf.ParsedURL = purl
	return wf
}

// GetQuery 获取Query参数
func (wf *Workflow) GetQuery() url.Values {
	return wf.ParsedURL.Query()
}

// GetCombineQuery 获取Query参数
func (wf *Workflow) GetCombineQuery() url.Values {
	if wf.ParsedURL != nil {
		vs := wf.ParsedURL.Query()
		return mergeMapList(wf.session.GetQuery(), vs)
	}
	return nil
}

// SetQuery 设置Query参数
func (wf *Workflow) SetQuery(query url.Values) *Workflow {
	if query == nil {
		return wf
	}
	query = (url.Values)(mergeMapList(wf.session.Query, query))
	wf.ParsedURL.RawQuery = query.Encode()
	return wf
}

var regexGetPath = regexp.MustCompile("/[^/]*")

// GetURLPath 获取Path参数 http://localhost/anything/user/pwd return [/anything /user /pwd]
func (wf *Workflow) GetURLPath() []string {
	return regexGetPath.FindAllString(wf.ParsedURL.Path, -1)
}

// GetURLRawPath 获取未分解Path参数
func (wf *Workflow) GetURLRawPath() string {
	return wf.ParsedURL.Path
}

// encodePath path格式每个item都必须以/开头
func encodePath(path []string) string {
	rawpath := ""
	for _, p := range path {
		if p[0] != '/' {
			p = "/" + p
		}
		rawpath += p
	}
	return rawpath
}

// SetURLPath 设置Path参数 对应 GetURLPath
func (wf *Workflow) SetURLPath(path []string) *Workflow {
	if path == nil {
		return wf
	}
	wf.ParsedURL.Path = encodePath(path)
	return wf
}

// SetURLRawPath 设置 参数 eg. /get = http:// hostname + /get
func (wf *Workflow) SetURLRawPath(path string) *Workflow {
	if path[0] != '/' {
		wf.ParsedURL.Path = "/" + path
	} else {
		wf.ParsedURL.Path = path
	}
	return wf
}

// SetBody 参数设置
func (wf *Workflow) SetBody(body IBody) *Workflow {
	wf.Body = body
	return wf
}

// GetBody 参数设置
func (wf *Workflow) GetBody() IBody {
	return wf.Body
}

// SetBodyAuto 参数设置
func (wf *Workflow) SetBodyAuto(params ...interface{}) *Workflow {

	if params != nil {
		plen := len(params)
		defaultContentType := TypeURLENCODED

		if plen >= 2 {
			t := params[plen-1]
			defaultContentType = t.(string)
		}

		wf.Body.SetPrefix(defaultContentType)

		switch defaultContentType {
		case TypeFormData:
			createMultipart(wf.Body, params) // 还存在 Mixed的可能
		default:
			var values url.Values
			switch param := params[0].(type) {
			case map[string]string:
				values := make(url.Values)
				for k, v := range param {
					values.Set(k, v)
				}
				wf.Body.SetIOBody([]byte(values.Encode()))
			case map[string][]string:
				values = param
				wf.Body.SetIOBody([]byte(values.Encode()))
			case string:
				wf.Body.SetIOBody([]byte(param))
			case []byte:
				wf.Body.SetIOBody(param)

			case *UploadFile:
				params = append(params, TypeFormData)
				wf.Body.SetPrefix(TypeFormData)
				createMultipart(wf.Body, params)
			case UploadFile:
				params = append(params, TypeFormData)
				wf.Body.SetPrefix(TypeFormData)
				createMultipart(wf.Body, params)
			case []*UploadFile:
				params = append(params, TypeFormData)
				wf.Body.SetPrefix(TypeFormData)
				createMultipart(wf.Body, params)
			case []UploadFile:
				params = append(params, TypeFormData)
				wf.Body.SetPrefix(TypeFormData)
				createMultipart(wf.Body, params)
			}
		}

	}
	return wf
}

func mergeMapList(headers ...map[string][]string) map[string][]string {

	set := make(map[string]map[string]int)
	merged := make(map[string][]string)

	for _, header := range headers {
		for key, values := range header {
			for _, v := range values {
				if vs, ok := set[key]; ok {
					vs[v] = 1
				} else {
					set[key] = make(map[string]int)
					set[key][v] = 1
				}
			}
		}
	}

	for key, mvalue := range set {
		for v := range mvalue {
			// merged.Add(key, v)
			if mergeValue, ok := merged[key]; ok {
				merged[key] = append(mergeValue, v)
			} else {
				merged[key] = []string{v}
			}
		}
	}

	return merged
}

// setHeaderRequest 设置request的头
func setHeaderRequest(req *http.Request, wf *Workflow) {
	req.Header = mergeMapList(req.Header, wf.session.Header, wf.Header)
}

// setHeaderRequest 设置request的临时Cookie, 永久需要在session上设置cookie
func setTempCookieRequest(req *http.Request, wf *Workflow) {
	if wf.Cookies != nil {
		for _, c := range wf.Cookies {
			req.AddCookie(c)
		}
	}
}

// Execute 执行
func (wf *Workflow) Execute() (*Response, error) {

	req := buildBodyRequest(wf)

	setHeaderRequest(req, wf)
	setTempCookieRequest(req, wf)

	if wf.session.auth != nil {
		req.SetBasicAuth(wf.session.auth.User, wf.session.auth.Password)
	}

	resp, err := wf.session.client.Do(req)
	if err != nil {
		return nil, err
	}

	wf.Body = NewBody()
	return FromHTTPResponse(resp, wf.session.Is.isDecompressNoAccept)
}
