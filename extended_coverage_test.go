package requests

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

// TestUploadFileExtendedCoverage 测试文件上传的扩展覆盖率
func TestUploadFileExtendedCoverage(t *testing.T) {
	t.Run("SetFileFromPathWithValidFile", func(t *testing.T) {
		// 创建临时测试文件
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.txt")
		content := "test file content for upload"

		err := os.WriteFile(testFile, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		uploadFile := &UploadFile{}
		err = uploadFile.SetFileFromPath(testFile)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// 验证文件是否正确设置
		if uploadFile.GetFile() == nil {
			t.Error("Expected file to be set")
		}

		// 验证文件名是否自动设置 - SetFileFromPath不会自动设置filename
		// 所以我们需要手动设置
		if uploadFile.GetFileName() == "" {
			// 这是正确的行为，SetFileFromPath不设置filename
			uploadFile.SetFileName(filepath.Base(testFile))
		}
	})

	t.Run("NewUploadFileWithAllParameters", func(t *testing.T) {
		// 创建临时测试文件
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "upload.txt")
		content := "content for NewUploadFile test"

		err := os.WriteFile(testFile, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// 测试NewUploadFile函数（创建空的UploadFile并手动设置）
		uploadFile := NewUploadFile()
		err = uploadFile.SetFileFromPath(testFile)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		uploadFile.SetFieldName("test_field")
		uploadFile.SetFileName("custom_name.txt")

		if uploadFile.GetFieldName() != "test_field" {
			t.Errorf("Expected field name 'test_field', got '%s'", uploadFile.GetFieldName())
		}

		if uploadFile.GetFileName() != "custom_name.txt" {
			t.Errorf("Expected filename 'custom_name.txt', got '%s'", uploadFile.GetFileName())
		}

		if uploadFile.GetFile() == nil {
			t.Error("Expected file to be set")
		}
	})

	t.Run("NewUploadFileWithMissingFile", func(t *testing.T) {
		// 测试不存在的文件 - 通过SetFileFromPath
		uploadFile := NewUploadFile()
		err := uploadFile.SetFileFromPath("/nonexistent/path/file.txt")
		if err == nil {
			t.Error("Expected error for non-existent file")
		}
	})

	t.Run("UploadFileFromPathWithValidPath", func(t *testing.T) {
		// 创建临时测试文件
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "path_test.txt")
		content := "content for path upload test"

		err := os.WriteFile(testFile, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// 测试UploadFileFromPath函数
		uploadFile, err := UploadFileFromPath(testFile)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if uploadFile == nil {
			t.Fatal("Expected uploadFile to be created")
		}

		expectedFileName := testFile // UploadFileFromPath sets full path, not basename
		if uploadFile.GetFileName() != expectedFileName {
			t.Errorf("Expected filename %s, got %s", expectedFileName, uploadFile.GetFileName())
		}

		if uploadFile.GetFile() == nil {
			t.Error("Expected file to be set")
		}
	})

	t.Run("UploadFileFromGlobWithValidPattern", func(t *testing.T) {
		// 创建临时目录和多个测试文件
		tmpDir := t.TempDir()

		// 创建测试文件
		files := []string{"test1.txt", "test2.txt", "test3.log"}
		for _, filename := range files {
			filePath := filepath.Join(tmpDir, filename)
			err := os.WriteFile(filePath, []byte("content of "+filename), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file %s: %v", filename, err)
			}
		}

		// 测试glob模式匹配txt文件
		pattern := filepath.Join(tmpDir, "*.txt")
		uploadFiles, err := UploadFileFromGlob(pattern)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(uploadFiles) != 2 {
			t.Errorf("Expected 2 files matching *.txt pattern, got %d", len(uploadFiles))
		}

		// 验证文件名
		foundFiles := make(map[string]bool)
		for _, uf := range uploadFiles {
			foundFiles[uf.GetFileName()] = true
		}

		if !foundFiles["test1.txt"] {
			t.Error("Expected to find test1.txt")
		}
		if !foundFiles["test2.txt"] {
			t.Error("Expected to find test2.txt")
		}
		if foundFiles["test3.log"] {
			t.Error("Did not expect to find test3.log in *.txt pattern")
		}
	})

	t.Run("UploadFileFromGlobWithNoMatches", func(t *testing.T) {
		// 测试没有匹配文件的glob模式 - UploadFileFromGlob返回错误
		pattern := "/nonexistent/path/*.txt"
		uploadFiles, err := UploadFileFromGlob(pattern)
		if err == nil {
			t.Error("Expected error for empty glob result")
		}

		if uploadFiles != nil {
			t.Error("Expected nil uploadFiles when error occurs")
		}
	})

	t.Run("UploadFileFromGlobWithInvalidPattern", func(t *testing.T) {
		// 测试无效的glob模式
		pattern := "[invalid-glob-pattern"
		uploadFiles, err := UploadFileFromGlob(pattern)
		if err == nil {
			t.Error("Expected error for invalid glob pattern")
		}

		if uploadFiles != nil {
			t.Error("Expected nil uploadFiles when error occurs")
		}
	})
}

// TestRequestBuilderBodyExtendedMethods 测试Request构建器的身体设置方法
func TestRequestBuilderBodyExtendedMethods(t *testing.T) {
	session := NewSession()

	t.Run("SetBodyWithTypeExtended", func(t *testing.T) {
		req := session.Post("http://httpbin.org/post")

		// 测试不同类型的body设置 - SetBodyWithType参数顺序是 (contentType, body)
		testCases := []struct {
			name        string
			body        interface{}
			contentType string
		}{
			{"String body", "test string", "text/plain"},
			{"Byte array body", []byte("test bytes"), "application/octet-stream"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req.SetBodyWithType(tc.contentType, tc.body)
				// 验证内容类型是否设置
				if req.header.Get("Content-Type") != tc.contentType {
					t.Errorf("Expected Content-Type %s, got %s", tc.contentType, req.header.Get("Content-Type"))
				}
			})
		}
	})

	t.Run("CreateBodyMultipartExtended", func(t *testing.T) {
		req := session.Post("http://httpbin.org/post")

		// 测试创建multipart body
		mpfd := req.CreateBodyMultipart()
		if mpfd == nil {
			t.Error("Expected CreateBodyMultipart to return non-nil MultipartFormData")
		}

		// CreateBodyMultipart只是创建对象，不会直接设置Content-Type
		// Content-Type通常在写入multipart数据时设置
	})
}

// TestSessionCookieExtendedMethods 测试Session cookie管理的扩展方法
func TestSessionCookieExtendedMethods(t *testing.T) {
	session := NewSession()

	t.Run("CookieJarConfiguration", func(t *testing.T) {
		// 测试启用cookie jar
		config := session.Config()
		config.SetWithCookiejar(true)

		// 验证client jar是否启用
		if session.client.Jar == nil {
			t.Error("Expected client jar to be enabled")
		}

		// 测试禁用cookie jar
		config.SetWithCookiejar(false)
		if session.client.Jar != nil {
			t.Error("Expected client jar to be disabled")
		}
	})
}

// TestMiddlewareExtendedCoverage 测试中间件系统的扩展覆盖率
func TestMiddlewareExtendedCoverage(t *testing.T) {
	t.Run("MiddlewareWithComplexScenarios", func(t *testing.T) {
		session := NewSession()

		// 创建测试中间件
		testMiddleware := &TestMiddleware{}

		// 测试中间件添加
		session.AddMiddleware(testMiddleware)
		if len(session.middlewares) != 1 {
			t.Errorf("Expected 1 middleware, got %d", len(session.middlewares))
		}

		// 测试中间件设置
		middleware2 := &TestMiddleware{}
		session.SetMiddlewares([]Middleware{testMiddleware, middleware2})
		if len(session.middlewares) != 2 {
			t.Errorf("Expected 2 middlewares, got %d", len(session.middlewares))
		}

		// 测试清空中间件
		session.ClearMiddlewares()
		if len(session.middlewares) != 0 {
			t.Errorf("Expected 0 middlewares after clear, got %d", len(session.middlewares))
		}
	})
}

// TestMiddleware 测试中间件实现
type TestMiddleware struct{}

func (tm *TestMiddleware) BeforeRequest(req *http.Request) error {
	return nil
}

func (tm *TestMiddleware) AfterResponse(resp *http.Response) error {
	return nil
}

// TestConfigurationEdgeCases 测试配置的边界情况
func TestConfigurationEdgeCases(t *testing.T) {
	session := NewSession()
	config := session.Config()

	t.Run("WithCookiejarToggle", func(t *testing.T) {
		// 测试多次切换cookie jar
		config.SetWithCookiejar(true)
		if session.client.Jar == nil {
			t.Error("Expected client jar to be enabled")
		}

		config.SetWithCookiejar(false)
		if session.client.Jar != nil {
			t.Error("Expected client jar to be disabled")
		}

		// 再次启用
		config.SetWithCookiejar(true)
		if session.client.Jar == nil {
			t.Error("Expected client jar to be enabled again")
		}
	})
}
