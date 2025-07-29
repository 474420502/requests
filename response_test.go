package requests

import (
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/tidwall/gjson"
)

type H map[string]interface{}

func TestFromHTTPResponse(t *testing.T) {

	var gresp *http.Response
	var err error
	gresp, err = http.DefaultClient.Get("http://httpbin.org/get")
	if err != nil {
		t.Error(err)
	}
	resp, err := FromHTTPResponse(gresp, false)
	if err != nil {
		t.Error(err)
	}

	if gjson.Get(resp.ContentString(), "headers.Host").String() != "httpbin.org" {
		t.Error("headers.Host != httpbin.org ?")
	}

	if resp.GetStatusCode() != 200 {
		t.Error("StatusCode != 200")
	}

	if len(resp.GetHeader()) == 0 {
		t.Error("esp.GetResponse().Header == nil")
	}

	if resp.GetStatus() != "200 OK" || resp.GetStatusCode() != 200 {
		t.Error(" resp.GetStatus() != 200 OK")
	}

	if len(resp.GetHeader()["Content-Length"]) != 1 {
		t.Error("resp.GetHeader() is error ?")
	}

	if int64(len(resp.ContentString())) != resp.GetContentLength() {
		t.Error("content len is not equal")
	}
}

func TestResponseDeflate(t *testing.T) {
	ses := NewSession()
	if wf := ses.Get("http://httpbin.org/get"); wf != nil {
		wf.AddHeader("accept-encoding", "deflate")
		resp, err := wf.Execute()
		if err != nil {
			t.Error(err)
		} else {
			if gjson.Get(string(resp.Content()), "headers.Accept-Encoding").String() != "deflate" {
				t.Error("Accept-Encoding != deflate ?")
			}
		}
	}

}

func TestAcceptCompressType(t *testing.T) {
	ses := NewSession() //requests.NewSession()
	ses.Config().AddAcceptEncoding(AcceptEncodingGzip)
	tp := ses.Get("http://0.0.0.0/compress")
	resp, err := tp.TestExecute(TestServer)
	if err != nil {
		panic(err)
	}
	if resp.ContentString() != "hello compress" {
		t.Error(resp.ContentString())
	}

	ses = NewSession() //requests.NewSession()
	ses.Config().AddAcceptEncoding(AcceptEncodingDeflate)
	tp = ses.Get("http://0.0.0.0/compress")
	resp, err = tp.TestExecute(TestServer)
	if err != nil {
		panic(err)
	}
	if resp.ContentString() != "hello compress" {
		t.Error(resp.ContentString())
	}

	ses = NewSession() //requests.NewSession()
	tp = ses.Get("http://0.0.0.0/compress")
	resp, err = tp.TestExecute(TestServer)
	if err != nil {
		panic(err)
	}
	if resp.ContentString() != "hello" {
		t.Error(resp.ContentString())
	}

}

func TestCaseAcceptEncoding2(t *testing.T) {
	ses := NewSession() //requests.NewSession()
	ses.Config().AddAcceptEncoding(AcceptEncodingBr)
	tp := ses.Get("http://0.0.0.0/compress")
	tp.SetHeader("Accept-Encoding", "deflate")
	resp, err := tp.TestExecute(TestServer)
	if err != nil {
		panic(err)
	}
	if resp.ContentString() != "hello compress" {
		t.Error(resp.ContentString())
	}
}

func TestContentCompressType(t *testing.T) {
	type H map[string]interface{}
	ses := NewSession() //requests.NewSession()
	ses.Config().SetContentEncoding(ContentEncodingBr)
	tp := ses.Get("http://0.0.0.0/content-compress")
	tp.SetBodyJSON(H{"key": "hello compress"})
	resp, err := tp.TestExecute(TestServer)
	if err != nil {
		panic(err)
	}
	if resp.ContentString() != "hello compress" {
		t.Error(resp.ContentString())
	}

	ses = NewSession() //requests.NewSession()
	ses.Config().SetContentEncoding(ContentEncodingDeflate)
	tp = ses.Get("http://0.0.0.0/content-compress")
	tp.SetBodyJSON(H{"key": "hello compress"})
	resp, err = tp.TestExecute(TestServer)
	if err != nil {
		panic(err)
	}
	if resp.ContentString() != "hello compress" {
		t.Error(resp.ContentString())
	}

	ses = NewSession() //requests.NewSession()
	ses.Config().SetContentEncoding(ContentEncodingGzip)
	tp = ses.Get("http://0.0.0.0/content-compress")
	tp.SetBodyJSON(H{"key": "hello compress"})
	resp, err = tp.TestExecute(TestServer)
	if err != nil {
		panic(err)
	}
	if resp.ContentString() != "hello compress" {
		t.Error(resp.ContentString())
	}

	ses = NewSession() //requests.NewSession()
	tp = ses.Get("http://0.0.0.0/content-compress")
	tp.SetBodyJSON(H{"key": "hello compress"})
	resp, err = tp.TestExecute(TestServer)
	if err != nil {
		panic(err)
	}
	if resp.ContentString() != "error compress" {
		t.Error(resp.ContentString())
	}
}

func TestReadmeEg1_2(t *testing.T) {
	ses := NewSession() //requests.NewSession()
	tp := ses.Get("http://httpbin.org/anything")
	tp.SetBodyJSON(`{"a": 1, "b": 2}`)
	resp, _ := tp.Execute()
	// log.Println(string(resp.Content()))
	// {
	// 	"args": {},
	// 	"data": "{\"a\": 1, \"b\": 2}",
	// 	"files": {},
	// 	"form": {},
	// 	"headers": {
	// 	  "Connection": "close",
	// 	  "Content-Length": "16",
	// 	  "Content-Type": "application/json",
	// 	  "Host": "httpbin.org",
	// 	  "User-Agent": "Go-http-client/1.1"
	// 	},
	// 	"json": {
	// 	  "a": 1,
	// 	  "b": 2
	// 	},
	// 	"method": "GET",
	// 	"origin": "172.17.0.1",
	// 	"url": "http://httpbin.org/anything"
	//   }

	tp = ses.Get("http://httpbin.org/anything")
	tp.SetBodyJSON(map[string]interface{}{"a": "1", "b": 2})
	resp, _ = tp.Execute()
	// log.Println(resp.ContentString())

	tp = ses.Get("http://httpbin.org/anything")
	tp.SetBodyJSON(H{"a": "1", "b": 2})
	resp, _ = tp.Execute()
	log.Println(resp.ContentString())

	tp = ses.Get("http://httpbin.org/anything")
	tp.SetBodyFormData(H{"a": "1", "b": 2})
	resp, _ = tp.Execute()
	log.Println(string(resp.Content()))
	// {
	// 	"args": {},
	// 	"data": "{\"a\":\"1\",\"b\":2}",
	// 	"files": {},
	// 	"form": {},
	// 	"headers": {
	// 	  "Connection": "close",
	// 	  "Content-Length": "15",
	// 	  "Content-Type": "application/json",
	// 	  "Host": "httpbin.org",
	// 	  "User-Agent": "Go-http-client/1.1"
	// 	},
	// 	"json": {
	// 	  "a": "1",
	// 	  "b": 2
	// 	},
	// 	"method": "GET",
	// 	"origin": "172.17.0.1",
	// 	"url": "http://httpbin.org/anything"
	//   }

	tp = ses.Post("http://httpbin.org/anything")
	tp.SetBodyFormData("./tests/learn.js")
	resp, _ = tp.Execute()

	content := `"file0": "learn.js\nfdsfsdavxlearnlearnlearnlearn"`
	if !strings.Contains(resp.ContentString(), content) {
		t.Error(resp.ContentString())
	}

}
