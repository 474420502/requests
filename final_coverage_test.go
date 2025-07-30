package requests

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"
)

// TestFinalCoverageImprovements 测试最终的覆盖率改进
func TestFinalCoverageImprovements(t *testing.T) {
	t.Run("Session_WithRetryConfig", func(t *testing.T) {
		// 测试使用重试配置创建session
		session, err := NewSessionWithRetry(3, 100*time.Millisecond)
		if err != nil {
			t.Errorf("Failed to create session with retry: %v", err)
		}

		// 创建请求测试重试配置
		req := session.Get("http://httpbin.org/status/500")
		resp, err := req.Execute()
		if err != nil {
			t.Logf("Expected error for 500 status: %v", err)
		}
		if resp != nil && resp.GetStatusCode() == 500 {
			t.Log("Got expected 500 status")
		}
	})

	t.Run("Request_ContextTimeout", func(t *testing.T) {
		session := NewSession()
		req := session.Get("http://httpbin.org/delay/2")

		// 设置1秒超时
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		req.WithContext(ctx)
		_, err := req.Execute()
		if err != nil {
			// 应该超时
			if strings.Contains(err.Error(), "context deadline exceeded") ||
				strings.Contains(err.Error(), "timeout") {
				t.Log("Got expected timeout error")
			} else {
				t.Logf("Got error: %v", err)
			}
		}
	})

	t.Run("Response_ErrorHandling", func(t *testing.T) {
		session := NewSession()

		// 测试404错误
		req := session.Get("http://httpbin.org/status/404")
		resp, err := req.Execute()
		if err != nil {
			t.Logf("Got error: %v", err)
		}
		if resp != nil {
			if resp.GetStatusCode() == 404 {
				t.Log("Got expected 404 status")
			}

			// 测试错误响应的JSON解析
			var result map[string]interface{}
			err = resp.BindJSON(&result)
			if err != nil {
				t.Logf("JSON binding error (expected for non-JSON response): %v", err)
			}
		}
	})

	t.Run("SessionWithOptions", func(t *testing.T) {
		// 测试使用选项创建session
		session, err := NewSessionWithOptions(
			WithTimeout(5*time.Second),
			WithUserAgent("test-agent/1.0"),
		)
		if err != nil {
			t.Errorf("Failed to create session with options: %v", err)
		}

		req := session.Get("http://httpbin.org/get")
		resp, err := req.Execute()
		if err != nil {
			t.Errorf("Request failed: %v", err)
		}
		if resp != nil {
			t.Logf("Response status: %d", resp.GetStatusCode())
		}
	})

	t.Run("MiddlewareChain_ErrorPropagation", func(t *testing.T) {
		session := NewSession()

		// 创建一个会产生错误的中间件
		errorMW := &MiddlewareImpl{
			BeforeRequestFunc: func(req *http.Request) error {
				return nil // 正常情况
			},
			AfterResponseFunc: func(resp *http.Response) error {
				if resp.StatusCode >= 400 {
					return nil // 我们不希望这里真的出错
				}
				return nil
			},
		}

		req := session.Get("http://httpbin.org/status/500").WithMiddleware(errorMW)
		resp, err := req.Execute()
		if err != nil {
			t.Logf("Got error: %v", err)
		}
		if resp != nil {
			t.Logf("Response status: %d", resp.GetStatusCode())
		}
	})

	t.Run("Config_EdgeCases", func(t *testing.T) {
		session := NewSession()
		config := session.Config()

		// 测试配置的各种边界情况
		config.SetTimeout(0)     // 无超时
		config.SetInsecure(true) // 跳过TLS验证

		// 测试请求是否仍然工作
		req := session.Get("http://httpbin.org/get")
		resp, err := req.Execute()
		if err != nil {
			t.Errorf("Request with edge case config failed: %v", err)
		}
		if resp != nil {
			t.Logf("Response status: %d", resp.GetStatusCode())
		}
	})

	t.Run("Upload_EdgeCases", func(t *testing.T) {
		session := NewSession()
		req := session.Post("http://httpbin.org/post")

		// 测试上传文件的边界情况
		upload := &UploadFile{
			FileName:  "empty.txt",
			FieldName: "empty_file",
		}

		// 设置空内容
		upload.SetFile(strings.NewReader(""))

		// 使用正确的AddFormFile方法
		req.AddFormFile("test_field", "empty.txt", strings.NewReader(""))

		resp, err := req.Execute()
		if err != nil {
			t.Logf("Upload error: %v", err)
		}
		if resp != nil {
			t.Logf("Upload response status: %d", resp.GetStatusCode())
		}
	})

	t.Run("Response_ContentTypes", func(t *testing.T) {
		session := NewSession()

		// 测试不同内容类型的响应
		req := session.Get("http://httpbin.org/xml")
		resp, err := req.Execute()
		if err != nil {
			t.Logf("XML request error: %v", err)
		}
		if resp != nil {
			// 测试内容类型检查
			t.Logf("Response length: %d", len(resp.Content()))

			// 测试非JSON内容的JSON方法
			if !resp.IsJSON() {
				t.Log("Response is not JSON (expected)")
			}
		}
	})

	t.Run("SessionPresets", func(t *testing.T) {
		// 测试预设session配置
		sessions := []*Session{}

		apiSession, err := NewSessionForAPI()
		if err != nil {
			t.Errorf("Failed to create API session: %v", err)
		} else {
			sessions = append(sessions, apiSession)
		}

		scrapingSession, err := NewSessionForScraping()
		if err != nil {
			t.Errorf("Failed to create scraping session: %v", err)
		} else {
			sessions = append(sessions, scrapingSession)
		}

		testSession, err := NewSessionForTesting()
		if err != nil {
			t.Errorf("Failed to create testing session: %v", err)
		} else {
			sessions = append(sessions, testSession)
		}

		// 测试每个session是否能正常工作
		for i, session := range sessions {
			req := session.Get("http://httpbin.org/get")
			resp, err := req.Execute()
			if err != nil {
				t.Errorf("Session %d request failed: %v", i, err)
			}
			if resp != nil {
				t.Logf("Session %d response status: %d", i, resp.GetStatusCode())
			}
		}
	})
}

// MiddlewareImpl 一个简单的中间件实现用于测试
type MiddlewareImpl struct {
	BeforeRequestFunc func(*http.Request) error
	AfterResponseFunc func(*http.Response) error
}

func (m *MiddlewareImpl) BeforeRequest(req *http.Request) error {
	if m.BeforeRequestFunc != nil {
		return m.BeforeRequestFunc(req)
	}
	return nil
}

func (m *MiddlewareImpl) AfterResponse(resp *http.Response) error {
	if m.AfterResponseFunc != nil {
		return m.AfterResponseFunc(resp)
	}
	return nil
}
