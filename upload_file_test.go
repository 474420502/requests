package requests

import (
	"io/ioutil"
	"os"
	"path/filepath"
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
		wf.SetBodyFormData(ufile)
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

		wf.SetBodyFormData("tests/json.file")
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
		wf.SetBodyFormData(ufile)
		resp, _ = wf.Execute()
		if _, ok := gjson.Get(string(resp.Content()), "files").Map()["MyField"]; !ok {
			t.Error("file error", string(resp.Content()))
		}

		// ses = NewSession()
		// wf = ses.Put("http://httpbin.org/put")

		ufile.SetFileFromPath("tests/json.file")
		wf.SetBodyFormData(*ufile)
		resp, _ = wf.Execute()
		if _, ok := gjson.Get(string(resp.Content()), "files").Map()["MyField"]; !ok {
			t.Error("file error", string(resp.Content()))
		}

		// ses = NewSession()
		// wf = ses.Put("http://httpbin.org/put")

		ufile = NewUploadFile()
		ufile.SetFileName("MyFile")
		ufile.SetFileFromPath("tests/json.file")
		wf.SetBodyFormData(ufile)
		resp, _ = wf.Execute()
		if _, ok := gjson.Get(string(resp.Content()), "files").Map()["file0"]; !ok {
			t.Error("file error", string(resp.Content()))
		}

		ufile.SetFileFromPath("tests/json.file")
		wf.SetBodyFormData(*ufile)
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
		wf.SetBodyFormData(ufileList)
		resp, _ = wf.Execute()
		if _, ok := gjson.Get(string(resp.Content()), "files").Map()["file1"]; !ok {
			t.Error("file error", string(resp.Content()))
		}

		// if wf.GetBody().ContentType() != "" {
		// 	t.Error("Body is not Clear")
		// }

		wf.SetBodyFormData([]string{"tests/learn.js", "tests/json.file"})

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
	mw.AddField("key1", "haha")
	mw.AddField("key2", "xixi")

	// mw.AddField("key2", "xixi")
	// data, err := ioutil.ReadAll(tp.Body)
	// log.Println(string(data))
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	tp.SetBodyFormData(mw)
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
	mw.AddField("key1", "haha")
	mw.AddField("key2", "xixi")

	f, err := os.Open("./tests/learn.js")
	if err != nil {
		t.Error(err)
		return
	}

	err = mw.AddFieldFile("filekey", "file0", f)
	if err != nil {
		t.Error(err)
		return
	}

	tp.SetBodyFormData(mw)

	resp, err = tp.Execute()
	if err != nil {
		t.Error(err)
		return
	}

	if _, ok := gjson.Get(string(resp.Content()), "files").Map()["filekey"]; !ok {
		t.Error("file error", string(resp.Content()))
	}
}

func TestUploadFileFromPath(t *testing.T) {
	// 创建一个临时文件
	tempFile, err := ioutil.TempFile("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	// 写入数据
	content := []byte("Hello, World!")
	if _, err := tempFile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatal(err)
	}

	// 使用 UploadFileFromPath 测试
	ufile, err := UploadFileFromPath(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	fileContent, err := ioutil.ReadAll(ufile.GetFile())
	if err != nil {
		t.Fatal(err)
	}
	if string(fileContent) != string(content) {
		t.Errorf("Expected file content %s, but got %s", string(content), string(fileContent))
	}
}

func TestUploadFileFromGlob(t *testing.T) {
	// 创建一个临时目录
	tempDir, err := ioutil.TempDir("", "testdir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// 在临时目录中创建两个文件
	fileNames := []string{"file1.txt", "file2.txt"}
	for _, fileName := range fileNames {
		filePath := filepath.Join(tempDir, fileName)
		if err := ioutil.WriteFile(filePath, []byte("Hello, World!"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// 使用 UploadFileFromGlob 测试
	globPattern := filepath.Join(tempDir, "*.txt")
	ufiles, err := UploadFileFromGlob(globPattern)
	if err != nil {
		t.Fatal(err)
	}

	// 检查返回的文件数
	if len(ufiles) != len(fileNames) {
		t.Errorf("Expected %d files, but got %d", len(fileNames), len(ufiles))
	}

	// 检查返回的文件名
	for i, ufile := range ufiles {
		expectedFileName := filepath.Base(fileNames[i])
		if ufile.GetFileName() != expectedFileName {
			t.Errorf("Expected file name %s, but got %s", expectedFileName, ufile.GetFileName())
		}
	}
}
