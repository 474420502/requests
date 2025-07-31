package requests

import (
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

// FileUploadTestServer 模拟文件上传服务器
type FileUploadTestServer struct{}

func (s *FileUploadTestServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 解析multipart表单
	err := r.ParseMultipartForm(32 << 20) // 32MB max
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"method": r.Method,
		"url":    r.URL.String(),
		"files":  make(map[string]interface{}),
		"form":   make(map[string]interface{}),
	}

	// 处理文件
	if r.MultipartForm != nil && r.MultipartForm.File != nil {
		for fieldName, fileHeaders := range r.MultipartForm.File {
			if len(fileHeaders) > 0 {
				fileHeader := fileHeaders[0]
				file, err := fileHeader.Open()
				if err != nil {
					continue
				}
				defer file.Close()

				// 读取文件内容
				content, err := io.ReadAll(file)
				if err != nil {
					continue
				}

				response["files"].(map[string]interface{})[fieldName] = string(content)
			}
		}
	}

	// 处理表单字段
	if r.MultipartForm != nil && r.MultipartForm.Value != nil {
		for key, values := range r.MultipartForm.Value {
			if len(values) > 0 {
				response["form"].(map[string]interface{})[key] = values[0]
			}
		}
	}

	// 返回JSON响应
	w.WriteHeader(http.StatusOK)
	responseJSON := `{
"method": "` + r.Method + `",
"url": "` + r.URL.String() + `",
"files": {`

	filesParts := []string{}
	for key, value := range response["files"].(map[string]interface{}) {
		filesParts = append(filesParts, `"`+key+`": "`+value.(string)+`"`)
	}
	responseJSON += strings.Join(filesParts, ", ")

	responseJSON += `},
"form": {`

	formParts := []string{}
	for key, value := range response["form"].(map[string]interface{}) {
		formParts = append(formParts, `"`+key+`": "`+value.(string)+`"`)
	}
	responseJSON += strings.Join(formParts, ", ")

	responseJSON += `}
}`

	w.Write([]byte(responseJSON))
}

func TestUploadFile(t *testing.T) {
	server := &FileUploadTestServer{}

	// 测试UploadFile值类型
	t.Run("UploadFile value test", func(t *testing.T) {
		ses := NewSession()
		req := ses.Post("/upload")

		ufile := NewUploadFile()
		ufile.SetFileName("test.txt")
		ufile.SetFieldName("testfield")
		ufile.SetFile(strings.NewReader("test content"))

		resp, err := req.AddFormFile(ufile.FieldName, ufile.FileName, ufile.FileReader).TestExecute(server)
		if err != nil {
			t.Fatalf("TestExecute failed: %v", err)
		}

		content := resp.ContentString()
		if !strings.Contains(content, "test content") {
			t.Errorf("Expected 'test content' not found in response: %s", content)
		}
	})

	// 测试多个文件上传
	t.Run("Multiple files test", func(t *testing.T) {
		ses := NewSession()
		req := ses.Post("/upload")

		// 第一个文件
		ufile1 := NewUploadFile()
		ufile1.SetFileName("file1.txt")
		ufile1.SetFieldName("file1")
		ufile1.SetFile(strings.NewReader("first file content"))

		// 第二个文件
		ufile2 := NewUploadFile()
		ufile2.SetFileName("file2.txt")
		ufile2.SetFieldName("file2")
		ufile2.SetFile(strings.NewReader("second file content"))

		// 传递多个文件
		resp, err := req.
			AddFormFile(ufile1.FieldName, ufile1.FileName, ufile1.FileReader).
			AddFormFile(ufile2.FieldName, ufile2.FileName, ufile2.FileReader).
			TestExecute(server)
		if err != nil {
			t.Fatalf("TestExecute failed: %v", err)
		}

		content := resp.ContentString()
		// 验证两个文件都存在
		if !strings.Contains(content, "first file content") {
			t.Errorf("First file content not found: %s", content)
		}
		if !strings.Contains(content, "second file content") {
			t.Errorf("Second file content not found: %s", content)
		}
	})
}

func TestBoundary(t *testing.T) {
	server := &FileUploadTestServer{}

	t.Run("Multipart boundary test", func(t *testing.T) {
		ses := NewSession()
		req := ses.Post("/boundary-test")

		// 使用现代API创建multipart数据
		file := strings.NewReader("test file content")
		resp, err := req.
			SetFormFields(map[string]string{
				"key1": "value1",
				"key2": "value2",
			}).
			AddFormFile("testfile", "test.txt", file).
			TestExecute(server)
		if err != nil {
			t.Fatalf("TestExecute failed: %v", err)
		}

		content := resp.ContentString()

		// 验证表单字段
		if !strings.Contains(content, "value1") {
			t.Errorf("Form field 'key1: value1' not found: %s", content)
		}
		if !strings.Contains(content, "value2") {
			t.Errorf("Form field 'key2: value2' not found: %s", content)
		}

		// 验证文件内容
		if !strings.Contains(content, "test file content") {
			t.Errorf("File content not found: %s", content)
		}
	})
}

func TestUploadFileFromPath(t *testing.T) {
	server := &FileUploadTestServer{}

	t.Run("Upload from file path", func(t *testing.T) {
		ses := NewSession()
		req := ses.Post("/upload-path")

		ufile, err := UploadFileFromPath("tests/json.file")
		if err != nil {
			t.Fatalf("UploadFileFromPath failed: %v", err)
		}
		defer ufile.fileCloser.Close()

		// 验证UploadFile属性
		if ufile.FileName != "tests/json.file" {
			t.Errorf("Expected FileName 'tests/json.file', got: %s", ufile.FileName)
		}

		resp, err := req.AddFormFile(ufile.FieldName, ufile.FileName, ufile.FileReader).TestExecute(server)
		if err != nil {
			t.Fatalf("TestExecute failed: %v", err)
		}

		content := resp.ContentString()
		// 验证文件内容包含预期的字符串
		if !strings.Contains(content, "jsonjsonjsonjson") {
			t.Errorf("Expected file content not found: %s", content)
		}
	})

	t.Run("Upload with custom field name", func(t *testing.T) {
		ses := NewSession()
		req := ses.Post("/upload-custom")

		ufile, err := UploadFileFromPath("tests/json.file")
		if err != nil {
			t.Fatalf("UploadFileFromPath failed: %v", err)
		}
		defer ufile.fileCloser.Close()

		ufile.SetFieldName("customfile")
		ufile.SetFileName("custom.json")

		resp, err := req.AddFormFile(ufile.FieldName, ufile.FileName, ufile.FileReader).TestExecute(server)
		if err != nil {
			t.Fatalf("TestExecute failed: %v", err)
		}

		content := resp.ContentString()
		if !strings.Contains(content, "jsonjsonjsonjson") {
			t.Errorf("File content not found: %s", content)
		}
	})
}

func TestUploadFileFromGlob(t *testing.T) {
	server := &FileUploadTestServer{}

	t.Run("Upload from glob pattern", func(t *testing.T) {
		// 创建测试文件
		testContent := "glob test content"
		testFile := "test_glob.txt"
		err := os.WriteFile(testFile, []byte(testContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		defer os.Remove(testFile)

		files, err := UploadFileFromGlob("test_*.txt")
		if err != nil {
			t.Fatalf("UploadFileFromGlob failed: %v", err)
		}

		if len(files) == 0 {
			t.Fatal("No files found from glob pattern")
		}

		// 验证第一个文件
		file := files[0]
		defer file.fileCloser.Close()

		if !strings.Contains(file.FileName, "test_glob.txt") {
			t.Errorf("Expected filename to contain 'test_glob.txt', got: %s", file.FileName)
		}

		// 测试上传单个文件
		ses := NewSession()
		req := ses.Post("/upload-glob")

		resp, err := req.AddFormFile(file.FieldName, file.FileName, file.FileReader).TestExecute(server)
		if err != nil {
			t.Fatalf("TestExecute failed: %v", err)
		}

		content := resp.ContentString()
		if !strings.Contains(content, testContent) {
			t.Errorf("Expected content '%s' not found: %s", testContent, content)
		}
	})

	t.Run("Upload multiple files from glob", func(t *testing.T) {
		// 创建多个测试文件
		testFiles := []string{"multi1.test", "multi2.test"}
		testContents := []string{"content1", "content2"}

		for i, filename := range testFiles {
			err := os.WriteFile(filename, []byte(testContents[i]), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file %s: %v", filename, err)
			}
			defer os.Remove(filename)
		}

		files, err := UploadFileFromGlob("multi*.test")
		if err != nil {
			t.Fatalf("UploadFileFromGlob failed: %v", err)
		}

		if len(files) < 2 {
			t.Fatalf("Expected at least 2 files, got: %d", len(files))
		}

		// 关闭所有文件
		for _, file := range files {
			defer file.fileCloser.Close()
		}

		// 测试上传 - 传递多个文件
		ses := NewSession()
		req := ses.Post("/upload-multi-glob")

		resp, err := req.
			AddFormFile(files[0].FieldName, files[0].FileName, files[0].FileReader).
			AddFormFile(files[1].FieldName, files[1].FileName, files[1].FileReader).
			TestExecute(server)
		if err != nil {
			t.Fatalf("TestExecute failed: %v", err)
		}

		content := resp.ContentString()

		// 验证所有内容都存在
		for _, expectedContent := range testContents {
			if !strings.Contains(content, expectedContent) {
				t.Errorf("Expected content '%s' not found: %s", expectedContent, content)
			}
		}
	})
}
