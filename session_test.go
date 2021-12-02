package requests

// func TestNewSession(t *testing.T) {
// 	ses := NewSession()
// 	if ses == nil {
// 		t.Error("session create fail, value is nil")
// 	}
// }

// func TestSession_Get(t *testing.T) {
// 	type fields struct {
// 		client *http.Client
// 	}
// 	type args struct {
// 		url string
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 	}{
// 		{
// 			name:   "Get test",
// 			fields: fields{client: &http.Client{}},
// 			args:   args{url: "http://httpbin.org/get"},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ses := &Session{
// 				client: tt.fields.client,
// 			}
// 			resp, err := ses.Get(tt.args.url).Execute()
// 			if err != nil {
// 				t.Error(err)
// 			}
// 			if len(string(resp.Content())) <= 150 {
// 				t.Error(string(resp.Content()))
// 			}
// 		})
// 	}
// }

// func TestSession_Post(t *testing.T) {
// 	type args struct {
// 		params []interface{}
// 	}

// 	tests := []struct {
// 		name string
// 		args args
// 		want *regexp.Regexp
// 	}{
// 		{
// 			name: "Post test",
// 			args: args{params: nil},
// 			want: regexp.MustCompile(`"form": \{\}`),
// 		},
// 		{
// 			name: "Post data",
// 			args: args{params: []interface{}{[]byte("a=1&b=2")}},
// 			want: regexp.MustCompile(`"form": \{[^"]+"a": "1"[^"]+"b": "2"[^\}]+\}`),
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ses := NewSession()
// 			got, err := ses.Post("http://httpbin.org/post").SetBodyAuto(tt.args.params...).Execute()

// 			if err != nil {
// 				t.Errorf("Metchod error = %v", err)
// 				return
// 			}

// 			if tt.want.MatchString(string(got.Content())) == false {
// 				t.Errorf("Metchod = %v, want %v", got, tt.want)
// 			}

// 		})
// 	}
// }

// func TestSession_Setparams(t *testing.T) {
// 	type fields struct {
// 		client *http.Client
// 		params *Body
// 	}
// 	type args struct {
// 		params []interface{}
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    *regexp.Regexp
// 		wantErr bool
// 	}{
// 		{
// 			name: "test Setparams",
// 			args: args{params: []interface{}{map[string]string{"a": "1", "b": "2"}}},
// 			want: regexp.MustCompile(`"form": \{[^"]+"a": "1"[^"]+"b": "2"[^\}]+\}`),
// 		},
// 		{
// 			name: "test json",
// 			args: args{params: []interface{}{`{"a":"1","b":"2"}`, TypeJSON}},
// 			want: regexp.MustCompile(`"json": \{[^"]+"a": "1"[^"]+"b": "2"[^\}]+\}`),
// 		},
// 		{
// 			name:   "test xml",
// 			fields: fields{client: &http.Client{}, params: NewBody()},
// 			args:   args{params: []interface{}{`<request><parameters><password>test</password></parameters></request>`, TypeXML}},
// 			want:   regexp.MustCompile(`"data": "<request><parameters><password>test</password></parameters></request>"`),
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ses := NewSession()

// 			got, err := ses.Post("http://httpbin.org/post").SetBodyAuto(tt.args.params...).Execute()
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Metchod error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}

// 			if tt.want.MatchString(string(got.Content())) == false {
// 				t.Errorf("Metchod = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestSession_PostUploadFile(t *testing.T) {
// 	type args struct {
// 		params interface{}
// 	}

// 	tests := []struct {
// 		name string
// 		args args
// 		want *regexp.Regexp
// 	}{
// 		{
// 			name: "test post uploadfile glob",
// 			args: args{params: "tests/*.js"},
// 			want: regexp.MustCompile(`"file0": "data:application/octet-stream;base64`),
// 		},
// 		{
// 			name: "test post uploadfile only one file",
// 			args: args{params: "tests/json.file"},
// 			want: regexp.MustCompile(`"file0": "json.file.+jsonjsonjsonjson"`),
// 		},
// 		{
// 			name: "test post uploadfile key values",
// 			args: args{params: map[string]string{"a": "32"}},
// 			want: regexp.MustCompile(`"a": "32"`),
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ses := NewSession()
// 			got, err := ses.Post("http://httpbin.org/post").SetBodyAuto(tt.args.params, TypeFormData).Execute()

// 			if err != nil {
// 				t.Errorf("Metchod error = %v", err)
// 				return
// 			}

// 			if tt.want.MatchString(string(got.Content())) == false {
// 				t.Errorf("Metchod = %v, want %v", got, tt.want)
// 			}

// 		})
// 	}
// }

// func TestSession_Put(t *testing.T) {
// 	type args struct {
// 		params interface{}
// 	}

// 	tests := []struct {
// 		name string
// 		args args
// 		want *regexp.Regexp
// 	}{
// 		{
// 			name: "test post uploadfile glob",
// 			args: args{params: "tests/*.js"},
// 			want: regexp.MustCompile(`"file0": "data:application/octet-stream;base64`),
// 		},
// 		{
// 			name: "test post uploadfile only one file",
// 			args: args{params: "tests/json.file"},
// 			want: regexp.MustCompile(`"file0": "json.file.+jsonjsonjsonjson"`),
// 		},
// 		{
// 			name: "test post uploadfile key values",
// 			args: args{params: map[string]string{"a": "32"}},
// 			want: regexp.MustCompile(`"a": "32"`),
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ses := NewSession()
// 			got, err := ses.Put("http://httpbin.org/put").SetBodyAuto(tt.args.params, TypeFormData).Execute()

// 			if err != nil {
// 				t.Errorf("Metchod error = %v", err)
// 				return
// 			}

// 			if tt.want.MatchString(string(got.Content())) == false {
// 				t.Errorf("Metchod = %v, want %v", got, tt.want)
// 			}

// 		})
// 	}
// }

// func TestSession_Patch(t *testing.T) {
// 	type args struct {
// 		params interface{}
// 	}

// 	tests := []struct {
// 		name string
// 		args args
// 		want *regexp.Regexp
// 	}{
// 		{
// 			name: "test post uploadfile glob",
// 			args: args{params: "tests/*.js"},
// 			want: regexp.MustCompile(`"file0": "data:application/octet-stream;base64`),
// 		},
// 		{
// 			name: "test post uploadfile only one file",
// 			args: args{params: "tests/json.file"},
// 			want: regexp.MustCompile(`"file0": "json.file.+jsonjsonjsonjson"`),
// 		},
// 		{
// 			name: "test post uploadfile key values",
// 			args: args{params: map[string]string{"a": "32"}},
// 			want: regexp.MustCompile(`"a": "32"`),
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ses := NewSession()
// 			got, err := ses.Patch("http://httpbin.org/patch").SetBodyAuto(tt.args.params, TypeFormData).Execute()

// 			if err != nil {
// 				t.Errorf("Metchod error = %v", err)
// 				return
// 			}

// 			if tt.want.MatchString(string(got.Content())) == false {
// 				t.Errorf("Metchod = %v, want %v", got, tt.want)
// 			}

// 		})
// 	}
// }

// func TestSession_SetConfig(t *testing.T) {

// 	type args struct {
// 		typeConfig TypeConfig
// 		values     interface{}
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name:    "test timeout",
// 			args:    args{typeConfig: CRequestTimeout, values: 0.0001},
// 			wantErr: true,
// 		},

// 		{
// 			name:    "test not timeout",
// 			args:    args{typeConfig: CRequestTimeout, values: 5},
// 			wantErr: false,
// 		},

// 		{
// 			name:    "test proxy",
// 			args:    args{typeConfig: CProxy, values: "http://" + ProxyAddress},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ses := NewSession()

// 			switch tt.args.typeConfig {
// 			case CRequestTimeout:
// 				ses.Config().SetTimeout(tt.args.values)
// 			case CProxy:
// 				ses.Config().SetProxy(tt.args.values)
// 			}

// 			_, err := ses.Get("http://httpbin.org/get").Execute()

// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Metchod error = %v", err)
// 				return
// 			}

// 		})
// 	}
// }

// func TestSession_SetConfigInsecure(t *testing.T) {

// 	ses := NewSession()

// 	ses.Config().SetInsecure(true)
// 	for _, badSSL := range []string{
// 		"https://self-signed.badssl.com/",
// 		"https://expired.badssl.com/",
// 		"https://wrong.host.badssl.com/",
// 	} {
// 		resp, err := ses.Get(badSSL).Execute()
// 		if err != nil {
// 			t.Error("Unable to make request", err)
// 		}
// 		if resp.GetStatusCode() != 200 {
// 			t.Error("Request did not return OK, is ", resp.GetStatusCode())
// 		}
// 	}

// }

// func TestSession_Cookies(t *testing.T) {
// 	ses := NewSession()

// 	t.Run("set cookie", func(t *testing.T) {
// 		resp, err := ses.Get("http://httpbin.org/cookies/set").SetCookieKV("a", "1").Execute()
// 		if err != nil {
// 			t.Error("cookies set error", err)
// 		}

// 		if !regexp.MustCompile(`"a": "1"`).MatchString(string(resp.Content())) {
// 			t.Error(string(resp.Content()))
// 		}
// 	})
// }

// func TestSession_Header(t *testing.T) {
// 	chromeua := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36"
// 	ses := NewSession()

// 	t.Run("ua header test", func(t *testing.T) {

// 		ses.Header.Add(HeaderKeyUA, chromeua)
// 		resp, err := ses.Get("https://www.baidu.com").Execute()
// 		if err != nil {
// 			t.Error("cookies set error", err)
// 		}

// 		if len(string(resp.Content())) <= 5000 {
// 			t.Error(string(resp.Content()), len(string(resp.Content())))
// 		}

// 		ses = NewSession()
// 		resp, err = ses.Get("https://www.baidu.com").AddHeader(HeaderKeyUA, chromeua).Execute()
// 		if err != nil {
// 			t.Error("cookies set error", err)
// 		}

// 		if len(string(resp.Content())) <= 5000 {
// 			t.Error(string(resp.Content()), len(string(resp.Content())))
// 		}
// 	})
// }

// func TestSession_ConfigEx(t *testing.T) {
// 	ses := NewSession()
// 	ses.Config().SetTimeout(time.Microsecond)
// 	resp, err := ses.Get("http://httpbin.org/get").Execute()
// 	if err == nil {
// 		t.Error(resp)
// 	} else {
// 		if strings.LastIndex(err.Error(), "Client.Timeout exceeded while awaiting headers") < 0 {
// 			t.Error(err)
// 		}
// 	}

// 	ses.Config().SetTimeout(float32(0.0000001))

// 	resp, err = ses.Get("http://httpbin.org/get").Execute()
// 	if err == nil {
// 		t.Error(resp)
// 	} else {
// 		if strings.LastIndex(err.Error(), "Client.Timeout exceeded while awaiting headers") < 0 {
// 			t.Error(err)
// 		}
// 	}

// 	ses.Config().SetKeepAlives(true)
// 	ses.Config().SetTimeout(int64(5))
// 	// jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
// 	u, err := url.Parse("http://httpbin.org")
// 	if err != nil {
// 		t.Error(err)
// 	} else {
// 		// jar.SetCookies(u, []*http.Cookie{&http.Cookie{Name: "Request", Value: "Cookiejar"}})

// 		cfg := ses.Config()
// 		cfg.SetWithCookiejar(false)
// 		cfg.SetWithCookiejar(true)

// 		ses.SetCookies(u, []*http.Cookie{&http.Cookie{Name: "Request", Value: "Cookiejar"}, &http.Cookie{Name: "eson", Value: "bad"}})
// 		resp, err = ses.Get("http://httpbin.org/get").Execute()
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		if gjson.Get(string(resp.Content()), "headers.Cookie").String() != "Request=Cookiejar; eson=bad" {
// 			t.Error(string(resp.Content()))
// 		}

// 		if resp.GetHeader()["Connection"][0] != "keep-alive" {
// 			t.Error("CKeepAlive is error")
// 		}
// 	}

// 	ses.Config().SetProxy(nil)
// 	if u, err := url.Parse("http://" + ProxyAddress); err != nil {
// 		t.Error(err)
// 	} else {
// 		ses.Config().SetProxy(u)
// 	}

// 	resp, err = ses.Get("http://httpbin.org/get").Execute()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if !regexp.MustCompile("eson=bad").Match(resp.Content()) {
// 		t.Error(string(resp.Content()))
// 	}

// 	ses.DelCookies(u, "eson")
// 	resp, err = ses.Get("http://httpbin.org/cookies").Execute()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if regexp.MustCompile("eson=bad").Match(resp.Content()) {
// 		t.Error(string(resp.Content()))
// 	}

// 	cookies := ses.GetCookies(u)
// 	if len(cookies) != 1 && cookies[0].String() != "Request=Cookiejar" {
// 		t.Error("cookies del get error please check it")
// 	}

// 	ses.ClearCookies()
// 	resp, err = ses.Get("http://httpbin.org/cookies").Execute()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if gjson.Get(string(resp.Content()), "cookies").String() != "{}" {
// 		t.Error(string(resp.Content()))
// 	}
// }

// func TestSession_SetQuery(t *testing.T) {
// 	ses := NewSession()
// 	ses.SetQuery(url.Values{"query": []string{"a", "b"}})
// 	resp, err := ses.Get("http://httpbin.org/get").Execute()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	query := gjson.Get(string(resp.Content()), "args.query").Array()
// 	for _, q := range query {
// 		if !(q.String() == "a" || q.String() == "b") {
// 			t.Error("query error, ", string(resp.Content()))
// 		}
// 	}
// }

// func TestSession_SetHeader(t *testing.T) {
// 	ses := NewSession()
// 	var header http.Header = make(http.Header)
// 	header["xx-xx"] = []string{"Header"}
// 	ses.SetHeader(header)

// 	resp, err := ses.Get("http://httpbin.org/headers").Execute()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if gjson.Get(string(resp.Content()), "headers.Xx-Xx").String() != "Header" {
// 		t.Error("Xx-Xx is not exists", string(resp.Content()))
// 	}

// 	var m = map[string][]string(ses.GetHeader())
// 	if m["xx-xx"][0] != "Header" {
// 		t.Error("header error")
// 	}
// }

// func TestSession_SetBasicAuth(t *testing.T) {
// 	ses := NewSession()
// 	ses.Config().SetBasicAuth(&BasicAuth{User: "eson", Password: "123456"})
// 	resp, err := ses.Get("http://httpbin.org/basic-auth/eson/123456").Execute()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if resp.GetStatusCode() != 200 {
// 		t.Error("code != 200, code = ", resp.GetStatus())
// 	}

// 	ses.Config().SetBasicAuth(&BasicAuth{User: "eson", Password: "12345"})
// 	resp, err = ses.Get("http://httpbin.org/basic-auth/eson/123456").Execute()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if resp.GetStatusCode() != 401 {
// 		t.Error("code != 401, code = ", resp.GetStatus())
// 	}

// 	resp, err = ses.Get("http://httpbin.org/basic-auth/eson/123456").Execute()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if resp.GetStatusCode() != 401 {
// 		t.Error("code != 401, code = ", resp.GetStatus())
// 	}

// 	ses.Config().SetBasicAuth("son", "123456")
// 	resp, err = ses.Get("http://httpbin.org/basic-auth/eson/123456").Execute()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if resp.GetStatusCode() != 401 {
// 		t.Error("code != 401, code = ", resp.GetStatus())
// 	}

// 	ses.Config().SetBasicAuth(nil)
// 	resp, err = ses.Get("http://httpbin.org/basic-auth/eson/123456").Execute()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if resp.GetStatusCode() != 401 {
// 		t.Error("code != 401, code = ", resp.GetStatus())
// 	}
// }
