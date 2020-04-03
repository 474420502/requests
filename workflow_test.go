package requests

import (
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"testing"

	"github.com/tidwall/gjson"
)

func TestWorkflow(t *testing.T) {
	ses := NewSession()

	t.Run("set cookie", func(t *testing.T) {
		resp, err := ses.Get("http://httpbin.org/cookies/set").AddKVCookie("a", "1").Execute()
		if err != nil {
			t.Error("cookies set error", err)
		}

		if !regexp.MustCompile(`"a": "1"`).MatchString(resp.Content()) {
			t.Error(resp.Content())
		}

		wf := ses.Get("http://httpbin.org/cookies/set")
		resp, err = wf.AddKVCookie("b", "2").Execute()
		if err != nil {
			t.Error("cookies set error", err)
		}

		result := gjson.Get(resp.Content(), "cookies.a")
		if result.Exists() {
			t.Error(resp.Content())
		}

		result = gjson.Get(resp.Content(), "cookies.b")
		if result.Int() != 2 {
			t.Error(resp.Content())
		}

		resp, err = wf.AddKVCookie("a", "3").Execute()
		results := gjson.GetMany(resp.Content(), "cookies.a", "cookies.b")
		if results[0].Int() != 3 {
			t.Error(resp.Content())
		}

		if results[1].Int() != 2 {
			t.Error(resp.Content())
		}

		resp, err = wf.AddHeader("XX", "123").SetRawURL("http://httpbin.org/headers").Execute()
		if err != nil {
			t.Error("cookies set error", err)
		}

		// headers 只能是String 表示
		result = gjson.Get(resp.Content(), "headers.Xx")
		if result.String() != "123" {
			t.Error(resp.Content())
		}
	})

}

func TestWorkflow_SetHeader(t *testing.T) {
	ses := NewSession()
	wf := ses.Get("http://httpbin.org/headers")
	var header http.Header
	header = make(http.Header)
	header["Eson"] = []string{"Bad"}
	header["HaHa"] = []string{"xixi"}
	wf.SetHeader(header)

	resp, err := wf.Execute()
	if err == nil && gjson.Get(resp.Content(), "headers.Eson").String() != "Bad" {
		t.Error("wf header error", resp.Content())
	}

	if err == nil && gjson.Get(resp.Content(), "headers.Haha").String() != "xixi" {
		t.Error("wf header error", resp.Content())
	}

	// 输入不符合规范不 会自动转换
	if wf.GetHeader()["HaHa"][0] != "xixi" {
		t.Error("Header 错误")
	}

	if len(ses.GetHeader()) != 0 {
		t.Error("session header should be zero")
	}

	delete(header, "HaHa")
	ses.SetHeader(header)
	wf = ses.Get("http://httpbin.org/headers")
	wf.AddHeader("Hello", "Hehe")

	resp, err = wf.Execute()
	if err != nil || gjson.Get(resp.Content(), "headers.Eson").String() != "Bad" {
		t.Error("wf header error", resp.Content())
	}

	if err != nil || gjson.Get(resp.Content(), "headers.Hello").String() != "Hehe" {
		t.Error("wf header error", resp.Content())
	}

	if len(wf.GetHeader()) != 1 || wf.GetHeader()["Hello"][0] != "Hehe" {
		t.Error("session header should be 1")
	}

	cheader := wf.GetCombineHeader()
	if len(cheader) != 2 || cheader["Eson"][0] != "Bad" {
		t.Error("GetCombineHeader error")
	}

	resp, err = wf.DelHeader("Hello").Execute()
	if err != nil {
		t.Error(err, resp.Content())
	}

	if gjson.Get(resp.Content(), "headers.Hello").Exists() {
		t.Error(" wf.DelHeader error")
	}
}

func TestWorkflow_Cookies(t *testing.T) {
	ses := NewSession()
	u, err := url.Parse("http://httpbin.org")
	if err != nil {
		t.Error(err)
	}
	ses.SetCookies(u, []*http.Cookie{&http.Cookie{Name: "Request", Value: "Cookiejar"}})
	wf := ses.Get("http://httpbin.org/cookies")
	wf.AddCookie(&http.Cookie{Name: "eson", Value: "Bad"})

	resp, _ := wf.Execute()
	if gjson.Get(resp.Content(), "cookies.Request").String() != "Cookiejar" {
		t.Error(" wf.AddCookie error")
	}

	if gjson.Get(resp.Content(), "cookies.eson").String() != "Bad" {
		t.Error(" wf.AddCookie error")
	}

	wf.DelCookie("eson")
	resp, _ = wf.Execute()
	if gjson.Get(resp.Content(), "cookies.Request").String() != "Cookiejar" {
		t.Error(" wf.AddCookie error")
	}
	if gjson.Get(resp.Content(), "cookies.eson").Exists() {
		t.Error(" wf.DelCookie error")
	}

	wf.AddCookies([]*http.Cookie{&http.Cookie{Name: "A", Value: "AA"}, &http.Cookie{Name: "B", Value: "BB"}})

	resp, _ = wf.Execute()
	if gjson.Get(resp.Content(), "cookies.Request").String() != "Cookiejar" {
		t.Error(" wf.AddCookie error")
	}
	if gjson.Get(resp.Content(), "cookies.A").String() != "AA" {
		t.Error(" wf.AddCookies error")
	}

	if gjson.Get(resp.Content(), "cookies.B").String() != "BB" {
		t.Error(" wf.AddCookies error")
	}

	wf.DelCookie(&http.Cookie{Name: "A", Value: "AA"})
	resp, _ = wf.Execute()
	if gjson.Get(resp.Content(), "cookies.A").Exists() {
		t.Error(" wf.AddCookies error")
	}

	if gjson.Get(resp.Content(), "cookies.B").String() != "BB" {
		t.Error(" wf.AddCookies error")
	}
}

func TestWorkflow_URL(t *testing.T) {
	ses := NewSession()
	wf := ses.Get("http://httpbin.org/")
	u, err := url.Parse("http://httpbin.org/get")
	if err != nil {
		t.Error(err)
	}
	wf.SetParsedURL(u)
	resp, _ := wf.Execute()
	if gjson.Get(resp.Content(), "url").String() != "http://httpbin.org/get" {
		t.Error("SetParsedURL ", resp.Content())
	}

	if wf.GetParsedURL().String() != "http://httpbin.org/get" {
		t.Error("SetParsedURL ", resp.Content())
	}

	wf = ses.Get("http://httpbin.org/")

	resp, _ = wf.SetURLRawPath("/get").Execute()
	if gjson.Get(resp.Content(), "url").String() != "http://httpbin.org/get" {
		t.Error("SetParsedURL ", resp.Content())
	}

	if wf.GetURLRawPath() != "/get" {
		t.Error("SetParsedURL ", resp.Content())
	}

	resp, _ = wf.SetURLRawPath("anything/user/password").Execute()
	if gjson.Get(resp.Content(), "url").String() != "http://httpbin.org/anything/user/password" {
		t.Error("SetParsedURL ", resp.Content())
	}
	paths := wf.GetURLPath()
	if paths[0] != "/anything" || paths[1] != "/user" || paths[2] != "/password" {
		t.Error("wf.GetURLPath()", paths)
	}

	wf = ses.Get("http://httpbin.org/")
	wf.SetURLPath(paths)
	if gjson.Get(resp.Content(), "url").String() != "http://httpbin.org/anything/user/password" {
		t.Error("SetParsedURL ", resp.Content())
	}
}

func TestWorkflow_Query(t *testing.T) {
	ses := NewSession()
	query := make(url.Values)
	query["session"] = []string{"true"}
	ses.SetQuery(query)
	wf := ses.Get("http://httpbin.org/get")
	wfquery := make(url.Values)
	wfquery["workflow"] = []string{"do", "to"}
	wf.SetQuery(wfquery)

	resp, _ := wf.Execute()
	result := gjson.Get(resp.Content(), "args.workflow")

	for _, r := range result.Array() {
		if !(r.String() == "to" || r.String() == "do") {
			t.Error("workflow SetQuery error")
		}
	}

	if gjson.Get(resp.Content(), "args.session").String() != "true" {
		t.Error("session SetQuery error")
	}

	if v, ok := wf.GetQuery()["workflow"]; ok {
		sort.Slice(v, func(i, j int) bool {
			if v[i] > v[j] {
				return true
			}
			return false
		})
		if !(v[0] == "to" && v[1] == "do") && len(v) != 2 {
			t.Error("workflow GetQuery", v)
		}
	}

	if v, ok := wf.GetQuery()["session"]; ok {
		if v[0] != "true" && len(v) != 1 {
			t.Error("workflow error")
		}
	}
}

func TestWorkflow_Body(t *testing.T) {
	ses := NewSession()
	wf := ses.Post("http://httpbin.org/post")
	body := NewBody()
	body.SetIOBody("a=1&b=2")
	wf.SetBody(body)
	resp, _ := wf.Execute()
	form := gjson.Get(resp.Content(), "form").Map()
	if v, ok := form["a"]; ok {
		if v.String() != "1" {
			t.Error(v)
		}
	}

	if v, ok := form["b"]; ok {
		if v.String() != "2" {
			t.Error(v)
		}
	}

	body.SetPrefix(TypeJSON)
	body.SetIOBody(`{"a": "1",   "b":  "2"}`)
	wf.SetBody(body)
	resp, _ = wf.Execute()
	json := gjson.Get(resp.Content(), "json").Map()
	if v, ok := json["a"]; ok {
		if v.String() != "1" {
			t.Error(v)
		}
	}

	if v, ok := json["b"]; ok {
		if v.String() != "2" {
			t.Error(v)
		}
	}

	// body.SetPrefix(TypeXML)
	// body.SetIOBody(`<root><a>1</a><b>2</b></root>`)
	// wf.SetBody(body)
	// resp, _ = wf.Execute()
}
