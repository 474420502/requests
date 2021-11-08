package requests

import (
	"bytes"
	"encoding/json"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
)

// Temporary 工作流 设计点: 这个并不影响Session的属性变化 如 NewWorkflow(ses, url).AddHeader() 对ses没影响
type Temporary struct {
	session   *Session
	mwriter   *MultipartWriter
	ParsedURL *url.URL
	Method    string
	Body      IBody
	Header    http.Header
	Cookies   map[string]*http.Cookie
}

// NewTemporary new and init workflow
func NewTemporary(ses *Session, urlstr string) *Temporary {
	tp := &Temporary{}
	tp.SwitchSession(ses)
	tp.SetRawURL(urlstr)

	tp.Body = NewBody()
	tp.Header = make(http.Header)
	tp.Cookies = make(map[string]*http.Cookie)
	return tp
}

// SwitchSession 替换Session
func (tp *Temporary) SwitchSession(ses *Session) {
	tp.session = ses
}

// AddHeader 添加头信息  Get方法从Header参数上获取 必须符合规范 HaHa -> Haha 如果真要HaHa,只能这样 Ha-Ha
func (tp *Temporary) AddHeader(key, value string) *Temporary {
	tp.Header[key] = append(tp.Header[key], value)
	return tp
}

// SetHeader 设置完全替换原有Header 必须符合规范 HaHa -> Haha 如果真要HaHa,只能这样 Ha-Ha
func (tp *Temporary) SetHeader(header http.Header) *Temporary {
	tp.Header = make(http.Header)
	for k, HValues := range header {
		var newHValues []string
		for _, value := range HValues {
			newHValues = append(newHValues, value)
		}
		tp.Header[k] = newHValues
	}
	return tp
}

// GetHeader 获取Workflow Header
func (tp *Temporary) GetHeader() http.Header {
	return tp.Header
}

// MergeHeader 合并 Header. 并进 Temporary
func (tp *Temporary) MergeHeader(cheader http.Header) {
	for key, values := range cheader {
		for _, v := range values {
			tp.Header.Add(key, v)
		}
	}
}

// DelHeader 添加头信息 Get方法从Header参数上获取
func (tp *Temporary) DelHeader(key string) *Temporary {
	tp.Header.Del(key)
	return tp
}

// SetCookie 添加Cookie
func (tp *Temporary) SetCookie(c *http.Cookie) *Temporary {
	tp.Cookies[c.Name] = c
	return tp
}

// AddCookies 添加[]*http.Cookie
func (tp *Temporary) AddCookies(cookies []*http.Cookie) *Temporary {
	for _, c := range cookies {
		tp.SetCookie(c)
	}
	return tp
}

// SetCookieKV 添加 以 key value 的 Cookie
func (tp *Temporary) SetCookieKV(name, value string) *Temporary {
	tp.Cookies[name] = &http.Cookie{Name: name, Value: value}
	return tp
}

// DelCookie 删除Cookie
func (tp *Temporary) DelCookie(name interface{}) *Temporary {
	switch n := name.(type) {
	case string:
		if _, ok := tp.Cookies[n]; ok {
			delete(tp.Cookies, n)
			return tp
		}
	case *http.Cookie:
		if _, ok := tp.Cookies[n.Name]; ok {
			delete(tp.Cookies, n.Name)
			return tp
		}
	default:
		panic("name type is not support")
	}
	return nil
}

// GetParsedURL 获取url的string形式
func (tp *Temporary) GetParsedURL() *url.URL {
	return tp.ParsedURL
}

// SetParsedURL 获取url的string形式
func (tp *Temporary) SetParsedURL(u *url.URL) *Temporary {
	tp.ParsedURL = u
	return tp
}

// GetRawURL 获取url的string形式
func (tp *Temporary) GetRawURL() string {
	// u := strings.Split(wf.ParsedURL.String(), "?")[0] + "?" + wf.GetCombineQuery().Encode()
	return tp.ParsedURL.String()
}

// SetRawURL 设置 url
func (tp *Temporary) SetRawURL(srcURL string) *Temporary {
	purl, err := url.ParseRequestURI(srcURL)
	if err != nil {
		panic(err)
	}
	tp.ParsedURL = purl
	return tp
}

// GetQuery 获取Query参数
func (tp *Temporary) GetQuery() url.Values {
	return tp.ParsedURL.Query()
}

// SetQuery 设置Query参数
func (tp *Temporary) SetQuery(query url.Values) *Temporary {
	if query == nil {
		return tp
	}
	// query = (url.Values)(mergeMapList(wf.session.Query, query))
	tp.ParsedURL.RawQuery = query.Encode()
	return tp
}

// MergeQuery 设置Query参数
func (tp *Temporary) MergeQuery(query url.Values) {
	tpquery := tp.ParsedURL.Query()
	for key, values := range query {
		for _, v := range values {
			tpquery.Add(key, v)
		}
	}
	tp.ParsedURL.RawQuery = tpquery.Encode()
}

// QueryParam 设置Query参数 不会返回nil
func (tp *Temporary) QueryParam(key string) IParam {
	return &ParamQuery{Temp: tp, Key: key}
}

// PathParam Path参数 使用正则匹配路径参数. group为参数 eg. /get?page=1&name=xiaoming
func (tp *Temporary) PathParam(regexpGroup string) IParam {
	return extractorParam(tp, regexpGroup, tp.ParsedURL.Path)
}

// HostParam Host参数 使用正则匹配Host参数. group为参数 eg.  httpbin.org
func (tp *Temporary) HostParam(regexpGroup string) IParam {
	return extractorParam(tp, regexpGroup, tp.ParsedURL.Host)
}

var regexGetPath = regexp.MustCompile("/[^/]*")

// GetURLPath 获取Path参数 http://localhost/anything/user/pwd return [/anything /user /pwd]
func (tp *Temporary) GetURLPath() []string {
	return regexGetPath.FindAllString(tp.ParsedURL.Path, -1)
}

// GetURLRawPath 获取未分解Path参数
func (tp *Temporary) GetURLRawPath() string {
	return tp.ParsedURL.Path
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
func (tp *Temporary) SetURLPath(path []string) *Temporary {
	if path == nil {
		return tp
	}
	tp.ParsedURL.Path = encodePath(path)
	return tp
}

// SetURLRawPath 设置 参数 eg. /get = http:// hostname + /get
func (tp *Temporary) SetURLRawPath(path string) *Temporary {
	if path[0] != '/' {
		tp.ParsedURL.Path = "/" + path
	} else {
		tp.ParsedURL.Path = path
	}
	return tp
}

// SetBody 参数设置
func (tp *Temporary) SetBody(body IBody) *Temporary {
	tp.mwriter = nil
	tp.Body = body
	return tp
}

// GetBody 参数设置
func (tp *Temporary) GetBody() IBody {
	return tp.Body
}

// GetBodyMultipart if get multipart, body = NewBody.  使用multipart/form-data. 传递keyvalue. 传递file.
// 每次都需要重置
func (tp *Temporary) GetBodyMultipart() *MultipartWriter {
	mw := &MultipartWriter{}
	var buf = &bytes.Buffer{}
	mwriter := multipart.NewWriter(buf)
	tp.Body.SetIOBody(buf)
	mw.mwriter = mwriter
	tp.mwriter = mw
	return mw
}

// SetBodyAuto 参数设置
func (tp *Temporary) SetBodyAuto(params ...interface{}) *Temporary {
	tp.Body = NewBody()
	if params != nil {
		tp.mwriter = nil

		plen := len(params)
		defaultContentType := TypeURLENCODED

		if plen >= 2 {
			t := params[plen-1]
			defaultContentType = t.(string)
		}

		tp.Body.SetPrefix(defaultContentType)

		switch defaultContentType {
		case TypeFormData:
			createMultipart(tp.Body, params) // 还存在 Mixed的可能
		default:
			var values url.Values
			switch param := params[0].(type) {

			case string:
				parambytes := []byte(param)
			TOPSTRING:
				for _, c := range parambytes {
					switch c {
					case ' ':
						continue
					case '[', '{':
						if json.Valid(parambytes) {
							tp.Body.SetPrefix(TypeJSON)
							tp.Body.SetIOBody(parambytes)
						} else {
							log.Println("SetBodyAuto -- Param is not json, but like json.\n", string(parambytes))
						}
						break TOPSTRING
					default:
						break TOPSTRING
					}
				}
				tp.Body.SetIOBody(parambytes)
			case []byte:

			TOPBYTES:
				for _, c := range param {
					switch c {
					case ' ':
						continue
					case '[', '{':
						if json.Valid(param) {
							tp.Body.SetPrefix(TypeJSON)
							tp.Body.SetIOBody(param)
						} else {
							log.Println("SetBodyAuto -- Param is not json, but like json.")
						}
						break TOPBYTES
					default:
						break TOPBYTES
					}
				}
				tp.Body.SetIOBody(param)
			case map[string]interface{}, []string, []interface{}:
				paramjson, err := json.Marshal(param)
				if err != nil {
					log.Panic(err)
				}
				tp.Body.SetPrefix(TypeJSON)
				tp.Body.SetIOBody(paramjson)

			case map[string]string:
				values := make(url.Values)
				for k, v := range param {
					values.Set(k, v)
				}
				tp.Body.SetIOBody([]byte(values.Encode()))

			case map[string][]string:
				values = param
				tp.Body.SetIOBody([]byte(values.Encode()))

			case *UploadFile:
				params = append(params, TypeFormData)
				tp.Body.SetPrefix(TypeFormData)
				createMultipart(tp.Body, params)
			case UploadFile:
				params = append(params, TypeFormData)
				tp.Body.SetPrefix(TypeFormData)
				createMultipart(tp.Body, params)
			case []*UploadFile:
				params = append(params, TypeFormData)
				tp.Body.SetPrefix(TypeFormData)
				createMultipart(tp.Body, params)
			case []UploadFile:
				params = append(params, TypeFormData)
				tp.Body.SetPrefix(TypeFormData)
				createMultipart(tp.Body, params)
			default:

				pvalue := reflect.ValueOf(param)
				ptype := reflect.TypeOf(param)

				if ptype.ConvertibleTo(compatibleType) {
					cparam := pvalue.Convert(compatibleType)
					paramjson, err := json.Marshal(cparam.Interface())
					if err != nil {
						log.Panic(err)
					}
					tp.Body.SetPrefix(TypeJSON)
					tp.Body.SetIOBody(paramjson)
				} else {
					paramjson, err := json.Marshal(pvalue.Interface())
					if err != nil {
						log.Panic(err)
					}
					tp.Body.SetPrefix(TypeJSON)
					tp.Body.SetIOBody(paramjson)
				}

			}
		}

	}
	return tp
}

// setHeaderRequest 设置request的头
func setHeaderRequest(req *http.Request, wf *Temporary) {
	var header http.Header
	if len(wf.Header) != 0 {
		header = wf.Header
	} else {
		header = wf.session.Header
	}
	for key, values := range header {
		for _, v := range values {
			req.Header.Add(key, v)
		}
	}
}

// setHeaderRequest 设置request的临时Cookie, 永久需要在session上设置cookie
func setTempCookieRequest(req *http.Request, wf *Temporary) {
	if wf.Cookies != nil {
		for _, c := range wf.Cookies {
			req.AddCookie(c)
		}
	}
}

// Execute 执行. 请求后会清楚Body的内容. 需要重新
func (tp *Temporary) Execute() (IResponse, error) {

	req := buildBodyRequest(tp)
	setHeaderRequest(req, tp)
	setTempCookieRequest(req, tp)

	if tp.session.auth != nil {
		req.SetBasicAuth(tp.session.auth.User, tp.session.auth.Password)
	}

	resp, err := tp.session.client.Do(req)
	if err != nil {
		return nil, err
	}

	if tp.session.Is.isClearBodyEvery {
		// tp.Body = NewBody()
	}

	return FromHTTPResponse(resp, tp.session.Is.isDecompressNoAccept)
}
