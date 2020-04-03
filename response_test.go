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

	if gjson.Get(resp.Content(), "headers.Host").String() != "httpbin.org" {
		t.Error("headers.Host != httpbin.org ?")
	}

	if resp.GetSrcResponse().StatusCode != 200 {
		t.Error("StatusCode != 200")
	}

	if len(resp.GetSrcResponse().Header) == 0 {
		t.Error("esp.GetSrcResponse().Header == nil")
	}

	if resp.GetStatue() != "200 OK" || resp.GetStatueCode() != 200 {
		t.Error(" resp.GetStatue() != 200 OK")
	}

	if len(resp.GetHeader()["Content-Length"]) != 1 {
		t.Error("resp.GetHeader() is error ?")
	}

	if int64(len(resp.Content())) != resp.GetContentLength() {
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
			if gjson.Get(resp.Content(), "headers.Accept-Encoding").String() != "deflate" {
				t.Error("Accept-Encoding != deflate ?")
			}
		}
	}

}
