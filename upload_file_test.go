package requests

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/tidwall/gjson"
)

func TestUploadFile(t *testing.T) {

	for i := 0; i < 1; i++ {

		ses := NewSession()
		wf := ses.Put("http://httpbin.org/put")

		ufile, err := UploadFileFromPath("tests/json.file")
		if err != nil {
			t.Error(err)
		}
		wf.SetBodyAuto(ufile, TypeFormData)
		resp, _ := wf.Execute()
		if _, ok := gjson.Get(string(resp.Content()), "files").Map()["file0"]; !ok {
			t.Error("file error", string(resp.Content()))
		}

		ses = NewSession()
		wf = ses.Patch("http://httpbin.org/patch")

		wf.SetBodyAuto("tests/json.file", TypeFormData)
		resp, _ = wf.Execute()
		if _, ok := gjson.Get(string(resp.Content()), "files").Map()["file0"]; !ok {
			t.Error("file error", string(resp.Content()))
		}

		ses = NewSession()
		wf = ses.Delete("http://httpbin.org/delete")
		ufile = NewUploadFile()
		ufile.SetFileName("MyFile")
		ufile.SetFieldName("MyField")
		ufile.SetFileReaderCloserFromFile("tests/json.file")
		wf.SetBodyAuto(ufile)
		resp, _ = wf.Execute()
		if _, ok := gjson.Get(string(resp.Content()), "files").Map()["MyField"]; !ok {
			t.Error("file error", string(resp.Content()))
		}

		// ses = NewSession()
		// wf = ses.Put("http://httpbin.org/put")

		ufile.SetFileReaderCloserFromFile("tests/json.file")
		wf.SetBodyAuto(*ufile)
		resp, _ = wf.Execute()
		if _, ok := gjson.Get(string(resp.Content()), "files").Map()["MyField"]; !ok {
			t.Error("file error", string(resp.Content()))
		}

		// ses = NewSession()
		// wf = ses.Put("http://httpbin.org/put")

		ufile = NewUploadFile()
		ufile.SetFileName("MyFile")
		ufile.SetFileReaderCloserFromFile("tests/json.file")
		wf.SetBodyAuto(ufile)
		resp, _ = wf.Execute()
		if _, ok := gjson.Get(string(resp.Content()), "files").Map()["file0"]; !ok {
			t.Error("file error", string(resp.Content()))
		}

		ufile.SetFileReaderCloserFromFile("tests/json.file")
		wf.SetBodyAuto(*ufile)
		resp, _ = wf.Execute()
		if _, ok := gjson.Get(string(resp.Content()), "files").Map()["file0"]; !ok {
			t.Error("file error", string(resp.Content()))
		}

		var ufileList []*UploadFile
		ufile, err = UploadFileFromPath("tests/json.file")
		if err != nil {
			t.Error(err)
		}
		ufileList = append(ufileList, ufile)
		ufile, err = UploadFileFromPath("tests/learn.js")
		if err != nil {
			t.Error(err)
		}
		ufileList = append(ufileList, ufile)
		wf.SetBodyAuto(ufileList)
		resp, _ = wf.Execute()
		if _, ok := gjson.Get(string(resp.Content()), "files").Map()["file1"]; !ok {
			t.Error("file error", string(resp.Content()))
		}

		if wf.GetBody().ContentType() != "" {
			t.Error("Body is not Clear")
		}

		wf.SetBodyAuto([]string{"tests/learn.js", "tests/json.file"}, TypeFormData)
		resp, _ = wf.Execute()
		if _, ok := gjson.Get(string(resp.Content()), "files").Map()["file1_0"]; !ok {
			t.Error("file error", string(resp.Content()))
		}
		if _, ok := gjson.Get(string(resp.Content()), "files").Map()["file0_0"]; !ok {
			t.Error("file error", string(resp.Content()))
		}
	}
}

func TestBoundary(t *testing.T) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//mulitipart/form-data时,需要获取自己关闭的boundary
	boundary := "fsdqwedsads"
	bodyWriter.SetBoundary(boundary)
	closeBuf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	w1, err := bodyWriter.CreateFormField("key1")
	if err != nil {
		panic(err)
	}
	w1.Write([]byte("haha"))
	//建立写入socket的reader对象

	w2, err := bodyWriter.CreateFormField("key2")
	w2.Write([]byte("xixi"))

	requestReader := io.MultiReader(bodyBuf, closeBuf)
	req, err := http.NewRequest("POST", "http://httpbin.org/post", requestReader)
	if err != nil {
		panic(err)
	}

	// body, err := ioutil.ReadAll(req.Body)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(string(body))

	//设置http头
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = int64(bodyBuf.Len()) + int64(closeBuf.Len())
	log.Println(req.ContentLength)
	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		panic(err)
	}
	var buf = &bytes.Buffer{}
	resp.Write(buf)

	t.Error(buf.String())
}
