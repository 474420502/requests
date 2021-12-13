package requests

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
)

// Temporary    这个并不影响Session的属性变化
type Temporary struct {
	session      *Session
	compressType CompressType
	// mwriter   *MultipartWriter
	mwriter   *multipart.Writer
	ParsedURL *url.URL
	Method    string
	Body      *bytes.Buffer
	Header    http.Header
	Cookies   map[string]*http.Cookie
}

// NewTemporary new and init Temporary
func NewTemporary(ses *Session, urlstr string) *Temporary {
	tp := &Temporary{session: ses}

	tp.SetRawURL(urlstr)

	tp.Header = make(http.Header)
	tp.Cookies = make(map[string]*http.Cookie)
	return tp
}

// SetContentType 设置set ContentType
func (tp *Temporary) SetContentType(contentType string) {
	tp.Header.Set(HeaderKeyContentType, contentType)
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
		newHValues = append(newHValues, HValues...)
		tp.Header[k] = newHValues
	}
	return tp
}

// GetHeader 获取Temporary Header
func (tp *Temporary) GetHeader() http.Header {
	return tp.Header
}

// SetCompress 设置Temporary Compress
func (tp *Temporary) SetCompress(c CompressType) {
	tp.compressType = c
}

// GetCompress 获取Temporary Compress
func (tp *Temporary) GetCompress() CompressType {
	return tp.compressType
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

// GetRawURL get url的string形式
func (tp *Temporary) GetRawURL() string {
	// u := strings.Split(wf.ParsedURL.String(), "?")[0] + "?" + wf.GetCombineQuery().Encode()
	return tp.ParsedURL.String()
}

// SetRawURL set url
func (tp *Temporary) SetRawURL(srcURL string) *Temporary {
	purl, err := url.ParseRequestURI(srcURL)
	if err != nil {
		panic(err)
	}
	tp.ParsedURL = purl
	return tp
}

// GetQuery get Query params
func (tp *Temporary) GetQuery() url.Values {
	return tp.ParsedURL.Query()
}

// SetQuery set Query params
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

// QueryParam Get the Interface of Query Param. never return nil. 不会返回nil
func (tp *Temporary) QueryParam(key string) IParam {
	return &ParamQuery{Temp: tp, Key: key}
}

// PathParam Path param 使用正则匹配路径参数. group为参数 eg. /get?page=1&name=xiaoming
func (tp *Temporary) PathParam(regexpGroup string) IParam {
	return extractorParam(tp, regexpGroup, tp.ParsedURL.Path)
}

// HostParam Host param 使用正则匹配Host参数. group为参数 eg.  httpbin.org
func (tp *Temporary) HostParam(regexpGroup string) IParam {
	return extractorParam(tp, regexpGroup, tp.ParsedURL.Host)
}

var regexGetPath = regexp.MustCompile("/[^/]*")

// GetURLPath get Path param eg: http://localhost/anything/user/pwd return [/anything /user /pwd]
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
func (tp *Temporary) SetBody(body io.Reader) *Temporary {
	var buf = bytes.NewBuffer(nil)
	_, err := io.Copy(buf, body)
	if err != nil {
		panic(err)
	}
	tp.Body = buf

	if tp.Header.Get("Content-Type") == "" {
		tp.Header.Set("Content-Type", TypeStream)
	}

	return tp
}

// // GetBody 参数设置
// func (tp *Temporary) GetBody() IBody {
// 	return tp.Body
// }

// GetBodyMultipart if get multipart, body = NewBody.  使用multipart/form-data. 传递keyvalue. 传递file.
// 每次都需要重置
func (tp *Temporary) CreateBodyMultipart() *multipart.Writer {
	var buf = &bytes.Buffer{}
	tp.mwriter = multipart.NewWriter(buf)
	tp.Header.Set(HeaderKeyContentType, tp.mwriter.FormDataContentType())
	tp.Body = buf
	return tp.mwriter
}

// SetBodyAuto 参数设置
func (tp *Temporary) SetBodyAuto(params ...interface{}) *Temporary {

	if params != nil {

		plen := len(params)
		var defaultContentType = ""
		var mwriter *multipart.Writer
		var err error
		defer func() {
			if mwriter != nil {
				defaultContentType += ";boundary=" + mwriter.Boundary()
			}
			tp.Header.Set(HeaderKeyContentType, defaultContentType)
		}()

		if plen >= 2 {
			t := params[plen-1]
			defaultContentType = t.(string)
		}

		switch defaultContentType {
		case TypeFormData:
			tp.Body, mwriter = createMultipart(params...) // 还存在 Mixed的可能
		case TypeJSON:
			var jsonbytes []byte
			switch param := params[0].(type) {
			case string:
				jsonbytes = []byte(param)
			case []byte:
				jsonbytes = param
			default:
				jsonbytes, err = json.Marshal(param)
				if err != nil {
					log.Panic(err)
				}
			}
			tp.Body = bytes.NewBuffer(jsonbytes)
		case TypeForm:
			fallthrough
		case TypePlain:
			fallthrough
		case TypeStream:
			switch param := params[0].(type) {
			case string:
				parambytes := []byte(param)
				tp.Body = bytes.NewBuffer(parambytes)
			case []byte:
				tp.Body = bytes.NewBuffer(param)
			case []rune:
				tp.Body = bytes.NewBuffer([]byte(string(param)))
			}
		default:
			var values url.Values
			switch param := params[0].(type) {

			case string:
				parambytes := []byte(param)
				tp.Body = bytes.NewBuffer(parambytes)
				defaultContentType = TypePlain
			TOPSTRING:
				for _, c := range parambytes {
					switch c {
					case ' ':
						continue
					case '[', '{':
						if json.Valid(parambytes) {
							defaultContentType = TypeJSON
						}
						break TOPSTRING
					default:
						break TOPSTRING
					}
				}
			case []rune:
				parambytes := []byte(string(param))
				tp.Body = bytes.NewBuffer(parambytes)
				defaultContentType = TypePlain
			TOPRUNE:
				for _, c := range parambytes {
					switch c {
					case ' ':
						continue
					case '[', '{':
						if json.Valid(parambytes) {
							defaultContentType = TypeJSON
						}
						break TOPRUNE
					default:
						break TOPRUNE
					}
				}

			case []byte:
				tp.Body = bytes.NewBuffer(param)
				defaultContentType = TypeStream
			TOPBYTES:
				for _, c := range param {
					switch c {
					case ' ':
						continue
					case '[', '{':
						if json.Valid(param) {
							defaultContentType = TypeJSON
						}
						break TOPBYTES
					default:
						break TOPBYTES
					}
				}

			case map[string]interface{}, []string, []interface{}, map[string]string:
				paramjson, err := json.Marshal(param)
				if err != nil {
					log.Panic(err)
				}
				defaultContentType = TypeJSON
				tp.Body = bytes.NewBuffer(paramjson)
			case url.Values:
				values = param
				tp.Body = bytes.NewBufferString(values.Encode())
				defaultContentType = TypeForm
			case map[string][]string:
				values = param
				tp.Body = bytes.NewBufferString(values.Encode())
				defaultContentType = TypeForm
			case *UploadFile, UploadFile, []*UploadFile, []UploadFile:
				params = append(params, TypeFormData)
				defaultContentType = TypeFormData
				tp.Body, mwriter = createMultipart(params...)
			default:
				var (
					paramjson []byte
					err       error
				)

				pvalue := reflect.ValueOf(param)
				ptype := reflect.TypeOf(param)

				if ptype.ConvertibleTo(compatibleType) {
					cparam := pvalue.Convert(compatibleType)
					paramjson, err = json.Marshal(cparam.Interface())
					if err != nil {
						log.Panic(err)
					}

				} else {
					paramjson, err = json.Marshal(pvalue.Interface())
					if err != nil {
						log.Panic(err)
					}
				}

				defaultContentType = TypeJSON
				tp.Body = bytes.NewBuffer(paramjson)

			}
		}

	}
	return tp
}

// setHeaderRequest 设置request的头
func setHeaderRequest(req *http.Request, wf *Temporary) {
	for key, values := range wf.session.Header {
		req.Header[key] = values
	}

	for key, values := range wf.Header {
		req.Header[key] = values
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
	if tp.mwriter != nil {
		tp.mwriter.Close()
		tp.mwriter = nil
	}
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

	return FromHTTPResponse(resp, tp.session.Is.isDecompressNoAccept)
}
