package requests

import (
	"log"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

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

func TestReadmeEg1(t *testing.T) {
	ses := NewSession() //requests.NewSession()
	tp := ses.Get("http://httpbin.org/anything")
	tp.SetBodyAuto(`{"a": 1, "b": 2}`)
	resp, _ := tp.Execute()
	log.Println(string(resp.Content()))
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
	tp.SetBodyAuto(map[string]interface{}{"a": "1", "b": 2})
	resp, _ = tp.Execute()
	log.Println(string(resp.Content()))

	tp = ses.Get("http://httpbin.org/anything")
	tp.SetBodyAuto(gin.H{"a": "1", "b": 2})
	resp, _ = tp.Execute()
	log.Println(string(resp.Content()))

	tp = ses.Get("http://httpbin.org/anything")
	tp.SetBodyAuto(gin.H{"a": "1", "b": 2}, TypeFormData)
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
	tp.SetBodyAuto("./tests/learn.js", TypeFormData)
	resp, _ = tp.Execute()
	// log.Println(string(resp.Content()))
	// {
	// 	"args": {},
	// 	"data": "",
	// 	"files": {
	// 	  "file0": "learn.js\nfdsfsdavxlearnlearnlearnlearn"
	// 	},
	// 	"form": {},
	// 	"headers": {
	// 	  "Connection": "close",
	// 	  "Content-Length": "279",
	// 	  "Content-Type": "multipart/form-data; boundary=1b8ffe52a1241b6caa93af8d5d2c3b6172eb650224ad959c69ea8df7c04d",
	// 	  "Host": "httpbin.org",
	// 	  "User-Agent": "Go-http-client/1.1"
	// 	},
	// 	"json": null,
	// 	"method": "POST",
	// 	"origin": "172.17.0.1",
	// 	"url": "http://httpbin.org/anything"
	//   }
}
