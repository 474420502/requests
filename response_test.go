package requests

import (
	"net/http"
	"testing"

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
