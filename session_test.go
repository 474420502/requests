package requests

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/tidwall/gjson"
)

func TestNewSession(t *testing.T) {
	ses := NewSession()
	if ses == nil {
		t.Error("session create fail, value is nil")
	}
}

func TestSession_Get(t *testing.T) {
	type fields struct {
		client *http.Client
	}
	type args struct {
		url string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "Get test",
			fields: fields{client: &http.Client{}},
			args:   args{url: "http://httpbin.org/get"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ses := &Session{
				client: tt.fields.client,
			}
			resp, err := ses.Get(tt.args.url).Execute()
			if err != nil {
				t.Error(err)
			}
			if len(string(resp.Content())) <= 150 {
				t.Error(string(resp.Content()))
			}
		})
	}
}

func TestSession_Post_Urlencoded(t *testing.T) {
	type args struct {
		params interface{}
	}

	tests := []struct {
		name string
		args args
		want *regexp.Regexp
	}{
		{
			name: "Post test",
			args: args{params: nil},
			want: regexp.MustCompile(`"form": \{\}`),
		},
		{
			name: "Post form []byte",
			args: args{params: []byte("a=1&b=2")},
			want: regexp.MustCompile(`"form": [^"]+"a": "1",[^"]+"b": "2"[^"]+`),
		},
		{
			name: "Post form string",
			args: args{params: "a=1&b=3"},
			want: regexp.MustCompile(`"form": [^"]+"a": "1",[^"]+"b": "3"[^"]+`),
		},

		{
			name: "Post form map[string]string",
			args: args{params: map[string]string{"a": "1", "b": "4"}},
			want: regexp.MustCompile(`"form": [^"]+"a": "1",[^"]+"b": "4"[^"]+`),
		},
		{
			name: "Post form map[string]int",
			args: args{params: map[string]int{"a": 1, "b": 4}},
			want: regexp.MustCompile(`"form": [^"]+"a": "1",[^"]+"b": "4"[^"]+`),
		},

		{
			name: "Post form map[string]int64",
			args: args{params: map[string]int64{"a": 1, "b": 4}},
			want: regexp.MustCompile(`"form": [^"]+"a": "1",[^"]+"b": "4"[^"]+`),
		},
		{
			name: "Post form map[string]uint",
			args: args{params: map[string]uint{"a": 1, "b": 4}},
			want: regexp.MustCompile(`"form": [^"]+"a": "1",[^"]+"b": "4"[^"]+`),
		},

		{
			name: "Post form map[string]uint64",
			args: args{params: map[string]uint64{"a": 1, "b": 4}},
			want: regexp.MustCompile(`"form": [^"]+"a": "1",[^"]+"b": "4"[^"]+`),
		},
		{
			name: "Post form map[string]float64",
			args: args{params: map[string]float64{"a": 1.23, "b": 4.543}},
			want: regexp.MustCompile(`"form": [^"]+"a": "1.23",[^"]+"b": "4.54"[^"]+`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ses := NewSession()
			got, err := ses.Post("http://httpbin.org/post").SetBodyUrlencoded(tt.args.params).Execute()

			if err != nil {
				t.Errorf("Metchod error = %v", err)
				return
			}
			// result := gjson.Parse(got.ContentString())
			if tt.want.MatchString(got.ContentString()) == false {
				t.Errorf("Metchod = %v \n want %v", got.ContentString(), tt.want)
			}

		})
	}
}

func TestSession_SetParams(t *testing.T) {

}

func TestSession_PostUploadFile_2(t *testing.T) {
	type args struct {
		params interface{}
	}

	tests := []struct {
		name string
		args args
		want *regexp.Regexp
	}{
		{
			name: "test post uploadfile glob",
			args: args{params: "tests/*.js"},
			want: regexp.MustCompile(`"file0": "tests/\*\.js"`),
		},
		{
			name: "test post uploadfile only one file",
			args: args{params: "tests/json.file"},
			want: regexp.MustCompile(`"file0": "tests/json\.file"`),
		},
		{
			name: "test post uploadfile key values",
			args: args{params: map[string]string{"a": "32"}},
			want: regexp.MustCompile(`"a": "32"`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ses := NewSession()
			got, err := ses.Post("http://httpbin.org/post").SetBodyFormData(tt.args.params).Execute()

			if err != nil {
				t.Errorf("Metchod error = %v", err)
				return
			}

			if tt.want.MatchString(string(got.Content())) == false {
				t.Errorf("Metchod = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestSession_Put(t *testing.T) {
	type args struct {
		params interface{}
	}

	tests := []struct {
		name string
		args args
		want *regexp.Regexp
	}{
		{
			name: "test post uploadfile glob",
			args: args{params: "tests/*.js"},
			want: regexp.MustCompile(`"file0": "tests/\*\.js"`),
		},
		{
			name: "test post uploadfile only one file",
			args: args{params: "tests/json.file"},
			want: regexp.MustCompile(`"file0": "tests/json\.file"`),
		},
		{
			name: "test post uploadfile key values",
			args: args{params: map[string]string{"a": "32"}},
			want: regexp.MustCompile(`"a": "32"`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ses := NewSession()
			got, err := ses.Put("http://httpbin.org/put").SetBodyFormData(tt.args.params).Execute()

			if err != nil {
				t.Errorf("Metchod error = %v", err)
				return
			}

			if tt.want.MatchString(string(got.Content())) == false {
				t.Errorf("Metchod = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestSession_Patch(t *testing.T) {
	type args struct {
		params interface{}
	}

	tests := []struct {
		name string
		args args
		want *regexp.Regexp
	}{
		{
			name: "test post uploadfile glob",
			args: args{params: "tests/*.js"},
			want: regexp.MustCompile(`"file0": "tests/\*\.js"`),
		},
		{
			name: "test post uploadfile only one file",
			args: args{params: "tests/json.file"},
			want: regexp.MustCompile(`"file0": "tests/json\.file"`),
		},
		{
			name: "test post uploadfile key values",
			args: args{params: map[string]string{"a": "32"}},
			want: regexp.MustCompile(`"a": "32"`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ses := NewSession()
			got, err := ses.Patch("http://httpbin.org/patch").SetBodyFormData(tt.args.params).Execute()

			if err != nil {
				t.Errorf("Metchod error = %v", err)
				return
			}

			if tt.want.MatchString(string(got.Content())) == false {
				t.Errorf("Metchod = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestSession_SetConfig(t *testing.T) {

	type args struct {
		typeConfig TypeConfig
		values     interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "test timeout",
			args:    args{typeConfig: CRequestTimeout, values: 0.000000001},
			wantErr: true,
		},

		{
			name:    "test not timeout",
			args:    args{typeConfig: CRequestTimeout, values: 5},
			wantErr: false,
		},
		{
			name:    "test proxy",
			args:    args{typeConfig: CProxy, values: "http://" + ProxyAddress},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ses := NewSession()

			switch tt.args.typeConfig {
			case CRequestTimeout:
				// 将interface{}转换为time.Duration
				switch v := tt.args.values.(type) {
				case float64:
					ses.Config().SetTimeout(time.Duration(v * float64(time.Second)))
				case int:
					ses.Config().SetTimeout(time.Duration(v) * time.Second)
				case int64:
					ses.Config().SetTimeout(time.Duration(v) * time.Second)
				case time.Duration:
					ses.Config().SetTimeout(v)
				}
			case CProxy:
				// 将interface{}转换为string
				if str, ok := tt.args.values.(string); ok {
					ses.Config().SetProxy(str)
				}
			}

			_, err := ses.Get("http://httpbin.org/get").Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("%v Metchod error = %v", tt.name, err)
				return
			}

		})
	}

	ProxyCloseChan <- 1
}

func TestSession_SetConfigInsecure(t *testing.T) {

	ses := NewSession()

	ses.Config().SetInsecure(true)
	for _, badSSL := range []string{
		"https://self-signed.badssl.com/",
	} {
		resp, err := ses.Get(badSSL).Execute()
		if err != nil {
			t.Error("Unable to make request", err)
		}
		if resp.GetStatusCode() != 200 {
			t.Error("Request did not return OK, is ", resp.GetStatusCode())
		}
	}

}

func TestSession_Cookies(t *testing.T) {
	ses := NewSession()

	t.Run("set cookie", func(t *testing.T) {
		resp, err := ses.Get("http://httpbin.org/cookies/set").SetCookieValue("a", "1").Execute()
		if err != nil {
			t.Error("cookies set error", err)
		}

		if !regexp.MustCompile(`"a": "1"`).MatchString(string(resp.Content())) {
			t.Error(string(resp.Content()))
		}
	})
}

func TestSession_Header(t *testing.T) {
	chromeua := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36"
	ses := NewSession()

	t.Run("ua header test", func(t *testing.T) {

		ses.Header.Add(HeaderKeyUA, chromeua)
		resp, err := ses.Get("https://www.baidu.com").Execute()
		if err != nil {
			t.Error("cookies set error", err)
		}

		if len(string(resp.Content())) <= 5000 {
			t.Error(string(resp.Content()), len(string(resp.Content())))
		}

		ses = NewSession()
		resp, err = ses.Get("https://www.baidu.com").AddHeader(HeaderKeyUA, chromeua).Execute()
		if err != nil {
			t.Error("cookies set error", err)
		}

		if len(string(resp.Content())) <= 5000 {
			t.Error(string(resp.Content()), len(string(resp.Content())))
		}
	})
}

func TestSession_ConfigEx(t *testing.T) {
	ses := NewSession()
	ses.Config().SetTimeout(time.Microsecond)
	resp, err := ses.Get("http://httpbin.org/get").Execute()
	if err == nil {
		t.Error(resp)
	} else {
		if strings.LastIndex(err.Error(), "Client.Timeout exceeded while awaiting headers") < 0 {
			t.Error(err)
		}
	}

	ses.Config().SetTimeout(time.Duration(float32(0.0000001) * float32(time.Second)))

	resp, err = ses.Get("http://httpbin.org/get").Execute()
	if err == nil {
		t.Error(resp)
	} else {
		if strings.LastIndex(err.Error(), "Client.Timeout exceeded while awaiting headers") < 0 {
			t.Error(err)
		}
	}

	ses.Config().SetKeepAlives(true)
	ses.Config().SetTimeout(time.Duration(int64(5)) * time.Second)
	// jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	u, err := url.Parse("http://httpbin.org")
	if err != nil {
		t.Error(err)
	} else {
		// jar.SetCookies(u, []*http.Cookie{&http.Cookie{Name: "Request", Value: "Cookiejar"}})

		cfg := ses.Config()
		cfg.SetWithCookiejar(false)
		cfg.SetWithCookiejar(true)

		ses.SetCookies(u, []*http.Cookie{{Name: "Request", Value: "Cookiejar"}, &http.Cookie{Name: "eson", Value: "bad"}})
		resp, err = ses.Get("http://httpbin.org/get").Execute()
		if err != nil {
			t.Error(err)
		}

		if gjson.Get(string(resp.Content()), "headers.Cookie").String() != "Request=Cookiejar; eson=bad" {
			t.Error(string(resp.Content()))
		}

		if resp.GetHeader()["Connection"][0] != "keep-alive" {
			t.Error("CKeepAlive is error")
		}
	}

	ses.Config().SetProxy("") // 使用空字符串代替nil来清除代理
	if u, err := url.Parse("http://" + ProxyAddress); err != nil {
		t.Error(err)
	} else {
		ses.Config().SetProxy(u.String()) // 转换为字符串
	}

	resp, err = ses.Get("http://httpbin.org/get").Execute()
	if err != nil {
		t.Error(err)
	}
	ProxyCloseChan <- 1

	if !regexp.MustCompile("eson=bad").Match(resp.Content()) {
		t.Error(string(resp.Content()))
	}

	ses.DelCookies(u, "eson")
	resp, err = ses.Get("http://httpbin.org/cookies").Execute()
	if err != nil {
		t.Error(err)
	}
	if regexp.MustCompile("eson=bad").Match(resp.Content()) {
		t.Error(string(resp.Content()))
	}

	cookies := ses.GetCookies(u)
	if len(cookies) != 1 && cookies[0].String() != "Request=Cookiejar" {
		t.Error("cookies del get error please check it")
	}

	ses.ClearCookies()
	resp, err = ses.Get("http://httpbin.org/cookies").Execute()
	if err != nil {
		t.Error(err)
	}
	if gjson.Get(string(resp.Content()), "cookies").String() != "{}" {
		t.Error(string(resp.Content()))
	}
}

func TestSession_SetQuery(t *testing.T) {
	ses := NewSession()
	ses.SetQuery(url.Values{"query": []string{"a", "b"}})
	resp, err := ses.Get("http://httpbin.org/get").Execute()
	if err != nil {
		t.Error(err)
	}
	query := gjson.Get(string(resp.Content()), "args.query").Array()
	for _, q := range query {
		if !(q.String() == "a" || q.String() == "b") {
			t.Error("query error, ", string(resp.Content()))
		}
	}
}

func TestSession_SetHeader(t *testing.T) {
	ses := NewSession()
	var header http.Header = make(http.Header)
	header["xx-xx"] = []string{"Header"}
	ses.SetHeader(header)

	resp, err := ses.Get("http://httpbin.org/headers").Execute()
	if err != nil {
		t.Error(err)
	}

	if gjson.Get(string(resp.Content()), "headers.Xx-Xx").String() != "Header" {
		t.Error("Xx-Xx is not exists", string(resp.Content()))
	}

	var m = map[string][]string(ses.GetHeader())
	if m["xx-xx"][0] != "Header" {
		t.Error("header error")
	}
}

func TestSession_SetBasicAuth(t *testing.T) {
	ses := NewSession()
	err := ses.Config().SetBasicAuth("eson", "123456")
	if err != nil {
		t.Error("SetBasicAuth failed:", err)
	}
	resp, err := ses.Get("http://httpbin.org/basic-auth/eson/123456").Execute()
	if err != nil {
		t.Error(err)
	}
	if resp.GetStatusCode() != 200 {
		t.Error("code != 200, code = ", resp.GetStatus())
	}

	err = ses.Config().SetBasicAuth("eson", "12345")
	resp, err = ses.Get("http://httpbin.org/basic-auth/eson/123456").Execute()
	if err != nil {
		t.Error(err)
	}

	if resp.GetStatusCode() != 401 {
		t.Error("code != 401, code = ", resp.GetStatus())
	}

	resp, err = ses.Get("http://httpbin.org/basic-auth/eson/123456").Execute()
	if err != nil {
		t.Error(err)
	}

	if resp.GetStatusCode() != 401 {
		t.Error("code != 401, code = ", resp.GetStatus())
	}

	ses.Config().SetBasicAuth("son", "123456")
	err = ses.Config().SetBasicAuth("eson", "12345")
	if err != nil {
		t.Error("SetBasicAuth failed:", err)
	}
	resp, err = ses.Get("http://httpbin.org/basic-auth/eson/123456").Execute()
	if err != nil {
		t.Error(err)
	}
	if resp.GetStatusCode() != 401 {
		t.Error("code != 401, code = ", resp.GetStatus())
	}

	ses.Config().ClearBasicAuth() // 使用ClearBasicAuth代替nil
	resp, err = ses.Get("http://httpbin.org/basic-auth/eson/123456").Execute()
	if err != nil {
		t.Error(err)
	}
	if resp.GetStatusCode() != 401 {
		t.Error("code != 401, code = ", resp.GetStatus())
	}
}
