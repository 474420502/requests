package requests

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

// Temporary 工作流 设计点: 这个并不影响Session的属性变化 如 NewWorkflow(ses, url).AddHeader() 对ses没影响
type Temporary struct {
	session   *Session
	ParsedURL *url.URL
	Method    string
	Body      IBody
	Header    http.Header
	Cookies   map[string]*http.Cookie
}

// NewTemporary new and init workflow
func NewTemporary(ses *Session, urlstr string) *Temporary {
	wf := &Temporary{}
	wf.SwitchSession(ses)
	wf.SetRawURL(urlstr)

	wf.Body = NewBody()
	wf.Header = make(http.Header)
	wf.Cookies = make(map[string]*http.Cookie)
	return wf
}

// SwitchSession 替换Session
func (wf *Temporary) SwitchSession(ses *Session) {
	wf.session = ses
}

// AddHeader 添加头信息  Get方法从Header参数上获取 必须符合规范 HaHa -> Haha 如果真要HaHa,只能这样 Ha-Ha
func (wf *Temporary) AddHeader(key, value string) *Temporary {
	wf.Header[key] = append(wf.Header[key], value)
	return wf
}

// SetHeader 设置完全替换原有Header 必须符合规范 HaHa -> Haha 如果真要HaHa,只能这样 Ha-Ha
func (wf *Temporary) SetHeader(header http.Header) *Temporary {
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
func (wf *Temporary) GetHeader() http.Header {
	return wf.Header
}

// MergeHeader 合并 Header. 并进 Temporary
func (wf *Temporary) MergeHeader(cheader http.Header) {
	for key, values := range cheader {
		for _, v := range values {
			wf.Header.Add(key, v)
		}
	}
}

// GetCombineHeader 获取后的Header信息
// func (wf *Temporary) GetCombineHeader() http.Header {
// 	if wf.Header != nil {
// 		return wf.Header
// 	}
// 	return zzzzzzz wf.session.Header, wf.Header)
// }

// DelHeader 添加头信息 Get方法从Header参数上获取
func (wf *Temporary) DelHeader(key string) *Temporary {
	wf.Header.Del(key)
	return wf
}

// SetCookie 添加Cookie
func (wf *Temporary) SetCookie(c *http.Cookie) *Temporary {
	wf.Cookies[c.Name] = c
	return wf
}

// AddCookies 添加[]*http.Cookie
func (wf *Temporary) AddCookies(cookies []*http.Cookie) *Temporary {
	for _, c := range cookies {
		wf.SetCookie(c)
	}
	return wf
}

// SetCookieKV 添加 以 key value 的 Cookie
func (wf *Temporary) SetCookieKV(name, value string) *Temporary {
	wf.Cookies[name] = &http.Cookie{Name: name, Value: value}
	return wf
}

// DelCookie 删除Cookie
func (wf *Temporary) DelCookie(name interface{}) *Temporary {
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
func (wf *Temporary) GetParsedURL() *url.URL {
	return wf.ParsedURL
}

// SetParsedURL 获取url的string形式
func (wf *Temporary) SetParsedURL(u *url.URL) *Temporary {
	wf.ParsedURL = u
	return wf
}

// GetRawURL 获取url的string形式
func (wf *Temporary) GetRawURL() string {
	// u := strings.Split(wf.ParsedURL.String(), "?")[0] + "?" + wf.GetCombineQuery().Encode()
	return wf.ParsedURL.String()
}

// SetRawURL 设置 url
func (wf *Temporary) SetRawURL(srcURL string) *Temporary {
	purl, err := url.ParseRequestURI(srcURL)
	if err != nil {
		panic(err)
	}
	wf.ParsedURL = purl
	return wf
}

// GetQuery 获取Query参数
func (wf *Temporary) GetQuery() url.Values {
	return wf.ParsedURL.Query()
}

// GetCombineQuery 获取 与Session合并后的参数
// Query参数 Session 于 Temporary 可能参数设置不一样.
// Temporay修改不影响Session
// func (wf *Temporary) GetCombineQuery() url.Values {
// 	if wf.ParsedURL != nil {
// 		vs := wf.ParsedURL.Query()
// 		return mergeMapList(wf.session.GetQuery(), vs)
// 	}
// 	return nil
// }

// SetQuery 设置Query参数
func (wf *Temporary) SetQuery(query url.Values) *Temporary {
	if query == nil {
		return wf
	}
	// query = (url.Values)(mergeMapList(wf.session.Query, query))
	wf.ParsedURL.RawQuery = query.Encode()
	return wf
}

// MergeQuery 设置Query参数
func (wf *Temporary) MergeQuery(query url.Values) {
	tpquery := wf.ParsedURL.Query()
	for key, values := range query {
		for _, v := range values {
			tpquery.Add(key, v)
		}
	}
	wf.ParsedURL.RawQuery = tpquery.Encode()
}

// QueryParam 设置Query参数
func (wf *Temporary) QueryParam(key string) *Param {
	return &Param{Temp: wf, Key: key}
}

var regexGetPath = regexp.MustCompile("/[^/]*")

// GetURLPath 获取Path参数 http://localhost/anything/user/pwd return [/anything /user /pwd]
func (wf *Temporary) GetURLPath() []string {
	return regexGetPath.FindAllString(wf.ParsedURL.Path, -1)
}

// GetURLRawPath 获取未分解Path参数
func (wf *Temporary) GetURLRawPath() string {
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
func (wf *Temporary) SetURLPath(path []string) *Temporary {
	if path == nil {
		return wf
	}
	wf.ParsedURL.Path = encodePath(path)
	return wf
}

// SetURLRawPath 设置 参数 eg. /get = http:// hostname + /get
func (wf *Temporary) SetURLRawPath(path string) *Temporary {
	if path[0] != '/' {
		wf.ParsedURL.Path = "/" + path
	} else {
		wf.ParsedURL.Path = path
	}
	return wf
}

// SetBody 参数设置
func (wf *Temporary) SetBody(body IBody) *Temporary {
	wf.Body = body
	return wf
}

// GetBody 参数设置
func (wf *Temporary) GetBody() IBody {
	return wf.Body
}

// SetBodyAuto 参数设置
func (wf *Temporary) SetBodyAuto(params ...interface{}) *Temporary {

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

			case string:
				parambytes := []byte(param)
			TOPSTRING:
				for _, c := range parambytes {
					switch c {
					case ' ':
						continue
					case '[', '{':
						if json.Valid(parambytes) {
							wf.Body.SetPrefix(TypeJSON)
							wf.Body.SetIOBody(parambytes)
						} else {
							log.Println("SetBodyAuto -- Param is not json, but like json.\n", string(parambytes))
						}
						break TOPSTRING
					default:
						break TOPSTRING
					}
				}
				wf.Body.SetIOBody(parambytes)
			case []byte:

			TOPBYTES:
				for _, c := range param {
					switch c {
					case ' ':
						continue
					case '[', '{':
						if json.Valid(param) {
							wf.Body.SetPrefix(TypeJSON)
							wf.Body.SetIOBody(param)
						} else {
							log.Println("SetBodyAuto -- Param is not json, but like json.")
						}
						break TOPBYTES
					default:
						break TOPBYTES
					}
				}
				wf.Body.SetIOBody(param)
			case map[string]interface{}, []string, []interface{}:
				paramjson, err := json.Marshal(param)
				if err != nil {
					log.Panic(err)
				}
				wf.Body.SetPrefix(TypeJSON)
				wf.Body.SetIOBody(paramjson)

			case map[string]string:
				values := make(url.Values)
				for k, v := range param {
					values.Set(k, v)
				}
				wf.Body.SetIOBody([]byte(values.Encode()))

			case map[string][]string:
				values = param
				wf.Body.SetIOBody([]byte(values.Encode()))

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

// func mergeMapList(headers ...map[string][]string) map[string][]string {

// 	set := make(map[string]map[string]int)
// 	merged := make(map[string][]string)

// 	for _, header := range headers {
// 		for key, values := range header {

// 			for _, v := range values {
// 				// v := values[0]
// 				if vs, ok := set[key]; ok {
// 					vs[v] = 1
// 				} else {
// 					set[key] = make(map[string]int)
// 					set[key][v] = 1
// 				}
// 			}

// 		}
// 	}

// 	for key, mvalue := range set {
// 		for v := range mvalue {
// 			// merged.Add(key, v)
// 			if mergeValue, ok := merged[key]; ok {
// 				merged[key] = append(mergeValue, v)
// 			} else {
// 				merged[key] = []string{v}
// 			}
// 		}
// 	}

// 	return merged
// }

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

// Execute 执行
func (wf *Temporary) Execute() (IResponse, error) {

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
