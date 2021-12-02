package requests

import (
	"io"
	"os"
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
			panic("")
		}
		wf.SetBodyAuto(ufile, TypeFormData)
		resp, err := wf.Execute()
		if err != nil {
			panic(err)
		}
		if _, ok := gjson.Get(string(resp.Content()), "files").Map()["file0"]; !ok {
			t.Error("file error", string(resp.Content()))
			panic("")
		}

		ses = NewSession()
		wf = ses.Patch("http://httpbin.org/patch")

		wf.SetBodyAuto("tests/json.file", TypeFormData)
		resp, _ = wf.Execute()
		if _, ok := gjson.Get(string(resp.Content()), "files").Map()["file0"]; !ok {
			t.Error("file error", string(resp.Content()))
			panic("")
		}

		ses = NewSession()
		wf = ses.Delete("http://httpbin.org/delete")
		ufile = NewUploadFile()
		ufile.SetFileName("MyFile")
		ufile.SetFieldName("MyField")
		ufile.SetFileFromPath("tests/json.file")
		wf.SetBodyAuto(ufile)
		resp, _ = wf.Execute()
		if _, ok := gjson.Get(string(resp.Content()), "files").Map()["MyField"]; !ok {
			t.Error("file error", string(resp.Content()))
		}

		// ses = NewSession()
		// wf = ses.Put("http://httpbin.org/put")

		ufile.SetFileFromPath("tests/json.file")
		wf.SetBodyAuto(*ufile)
		resp, _ = wf.Execute()
		if _, ok := gjson.Get(string(resp.Content()), "files").Map()["MyField"]; !ok {
			t.Error("file error", string(resp.Content()))
		}

		// ses = NewSession()
		// wf = ses.Put("http://httpbin.org/put")

		ufile = NewUploadFile()
		ufile.SetFileName("MyFile")
		ufile.SetFileFromPath("tests/json.file")
		wf.SetBodyAuto(ufile)
		resp, _ = wf.Execute()
		if _, ok := gjson.Get(string(resp.Content()), "files").Map()["file0"]; !ok {
			t.Error("file error", string(resp.Content()))
		}

		ufile.SetFileFromPath("tests/json.file")
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

		// if wf.GetBody().ContentType() != "" {
		// 	t.Error("Body is not Clear")
		// }

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

	ses := NewSession()
	tp := ses.Post("http://httpbin.org/post")

	mw := tp.CreateBodyMultipart()
	mw.WriteField("key1", "haha")
	mw.WriteField("key2", "xixi")

	// mw.AddField("key2", "xixi")

	resp, err := tp.Execute()
	if err != nil {
		t.Error(err)
		return
	}

	if v, ok := gjson.Get(string(resp.Content()), "form").Map()["key2"]; !ok || v.String() != "xixi" {
		t.Error("file error", string(resp.Content()))
	}

	resp, err = tp.Execute()
	if err != nil {
		t.Error(err)
		return
	}

	if v, ok := gjson.Get(string(resp.Content()), "form").Map()["key1"]; !ok || v.String() != "haha" {
		t.Error("file error", string(resp.Content()))
	}

	mw = tp.CreateBodyMultipart()
	mw.WriteField("key1", "haha")
	mw.WriteField("key2", "xixi")

	f, err := os.Open("./tests/learn.js")
	if err != nil {
		t.Error(err)
		return
	}

	writer, err := mw.CreateFormFile("filekey", "file0")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = io.Copy(writer, f)
	if err != nil {
		t.Error(err)
		return
	}

	resp, err = tp.Execute()
	if err != nil {
		t.Error(err)
		return
	}

	if _, ok := gjson.Get(string(resp.Content()), "files").Map()["file0"]; !ok {
		t.Error("file error", string(resp.Content()))
	}
}

func TestCaseCreateMultiPart(t *testing.T) {
	// ufile, err := UploadFileFromPath("tests/json.file")
	// if err != nil {
	// 	t.Error(err)
	// }

	// createMultipart(ufile, TypeFormData)

}
