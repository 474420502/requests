package requests

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

// BenchmarkBasicRequests 基础请求性能测试
func BenchmarkBasicRequests(b *testing.B) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "hello", "status": "ok"}`))
	}))
	defer server.Close()

	b.Run("Sequential_GET", func(b *testing.B) {
		session := NewSession()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			resp, err := session.Get(server.URL).Execute()
			if err != nil {
				b.Fatal(err)
			}
			resp.Content() // 确保读取响应内容
		}
	})

	b.Run("Parallel_GET", func(b *testing.B) {
		session := NewSession()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				resp, err := session.Get(server.URL).Execute()
				if err != nil {
					b.Fatal(err)
				}
				resp.Content()
			}
		})
	})

	b.Run("Sequential_POST_JSON", func(b *testing.B) {
		session := NewSession()
		data := map[string]interface{}{
			"name":  "test",
			"count": 42,
		}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			resp, err := session.Post(server.URL).SetBodyJson(data).Execute()
			if err != nil {
				b.Fatal(err)
			}
			resp.Content()
		}
	})

	b.Run("Parallel_POST_JSON", func(b *testing.B) {
		session := NewSession()
		data := map[string]interface{}{
			"name":  "test",
			"count": 42,
		}
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				resp, err := session.Post(server.URL).SetBodyJson(data).Execute()
				if err != nil {
					b.Fatal(err)
				}
				resp.Content()
			}
		})
	})
}

// BenchmarkVsStdLib 与标准库性能对比
func BenchmarkVsStdLib(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "hello"}`))
	}))
	defer server.Close()

	b.Run("Requests_Library", func(b *testing.B) {
		session := NewSession()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			resp, err := session.Get(server.URL).Execute()
			if err != nil {
				b.Fatal(err)
			}
			resp.Content()
		}
	})

	b.Run("Standard_Library", func(b *testing.B) {
		client := &http.Client{}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			resp, err := client.Get(server.URL)
			if err != nil {
				b.Fatal(err)
			}
			io.ReadAll(resp.Body)
			resp.Body.Close()
		}
	})
}

// BenchmarkSessionReuse Session复用性能测试
func BenchmarkSessionReuse(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	b.Run("NewSession_PerRequest", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			session := NewSession() // 每次都创建新Session
			resp, err := session.Get(server.URL).Execute()
			if err != nil {
				b.Fatal(err)
			}
			resp.Content()
		}
	})

	b.Run("ReuseSession", func(b *testing.B) {
		session := NewSession() // 复用同一个Session
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			resp, err := session.Get(server.URL).Execute()
			if err != nil {
				b.Fatal(err)
			}
			resp.Content()
		}
	})

	b.Run("ReuseSession_Parallel", func(b *testing.B) {
		session := NewSession()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				resp, err := session.Get(server.URL).Execute()
				if err != nil {
					b.Fatal(err)
				}
				resp.Content()
			}
		})
	})
}

// 简单的日志中间件
type SimpleLoggingMiddleware struct{}

func (m *SimpleLoggingMiddleware) BeforeRequest(req *http.Request) error {
	// 模拟日志记录
	_ = req.URL.String()
	return nil
}

func (m *SimpleLoggingMiddleware) AfterResponse(resp *http.Response) error {
	// 模拟日志记录
	_ = resp.Status
	return nil
}

// 计时中间件
type SimpleTimingMiddleware struct{}

func (m *SimpleTimingMiddleware) BeforeRequest(req *http.Request) error {
	req.Header.Set("X-Start-Time", fmt.Sprintf("%d", time.Now().UnixNano()))
	return nil
}

func (m *SimpleTimingMiddleware) AfterResponse(resp *http.Response) error {
	// 模拟记录耗时
	return nil
}

// BenchmarkMiddleware 中间件性能测试
func BenchmarkMiddleware(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	loggingMw := &SimpleLoggingMiddleware{}
	timingMw := &SimpleTimingMiddleware{}

	b.Run("NoMiddleware", func(b *testing.B) {
		session := NewSession()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			resp, err := session.Get(server.URL).Execute()
			if err != nil {
				b.Fatal(err)
			}
			resp.Content()
		}
	})

	b.Run("SingleMiddleware", func(b *testing.B) {
		session := NewSession()
		session.AddMiddleware(loggingMw)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			resp, err := session.Get(server.URL).Execute()
			if err != nil {
				b.Fatal(err)
			}
			resp.Content()
		}
	})

	b.Run("MultipleMiddleware", func(b *testing.B) {
		session := NewSession()
		session.AddMiddleware(loggingMw)
		session.AddMiddleware(timingMw)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			resp, err := session.Get(server.URL).Execute()
			if err != nil {
				b.Fatal(err)
			}
			resp.Content()
		}
	})
}

// BenchmarkJSONProcessing JSON处理性能测试
func BenchmarkJSONProcessing(b *testing.B) {
	jsonData := `{
		"users": [
			{"id": 1, "name": "Alice", "email": "alice@example.com"},
			{"id": 2, "name": "Bob", "email": "bob@example.com"},
			{"id": 3, "name": "Charlie", "email": "charlie@example.com"}
		],
		"total": 3,
		"page": 1
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonData))
	}))
	defer server.Close()

	b.Run("GetJSONField", func(b *testing.B) {
		session := NewSession()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			resp, err := session.Get(server.URL).Execute()
			if err != nil {
				b.Fatal(err)
			}
			_ = resp.GetJSONField("users.0.name")
		}
	})

	b.Run("GetJSONString", func(b *testing.B) {
		session := NewSession()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			resp, err := session.Get(server.URL).Execute()
			if err != nil {
				b.Fatal(err)
			}
			_, _ = resp.GetJSONString("users.0.name")
		}
	})

	b.Run("DecodeJSON", func(b *testing.B) {
		session := NewSession()
		type User struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		type Response struct {
			Users []User `json:"users"`
			Total int    `json:"total"`
			Page  int    `json:"page"`
		}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			resp, err := session.Get(server.URL).Execute()
			if err != nil {
				b.Fatal(err)
			}
			var result Response
			_ = resp.DecodeJSON(&result)
		}
	})
}

// BenchmarkFileUpload 文件上传性能测试
func BenchmarkFileUpload(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 模拟处理文件上传
		r.ParseMultipartForm(32 << 20) // 32MB
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Upload successful"))
	}))
	defer server.Close()

	// 准备测试数据
	smallData := strings.NewReader(strings.Repeat("a", 1024))      // 1KB
	mediumData := strings.NewReader(strings.Repeat("b", 1024*100)) // 100KB
	largeData := strings.NewReader(strings.Repeat("c", 1024*1024)) // 1MB

	b.Run("SmallFile_1KB", func(b *testing.B) {
		session := NewSession()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			smallData.Seek(0, 0) // 重置reader
			resp, err := session.Post(server.URL).
				AddFormFile("file", "small.txt", smallData).
				Execute()
			if err != nil {
				b.Fatal(err)
			}
			resp.Content()
		}
	})

	b.Run("MediumFile_100KB", func(b *testing.B) {
		session := NewSession()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			mediumData.Seek(0, 0)
			resp, err := session.Post(server.URL).
				AddFormFile("file", "medium.txt", mediumData).
				Execute()
			if err != nil {
				b.Fatal(err)
			}
			resp.Content()
		}
	})

	b.Run("LargeFile_1MB", func(b *testing.B) {
		session := NewSession()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			largeData.Seek(0, 0)
			resp, err := session.Post(server.URL).
				AddFormFile("file", "large.txt", largeData).
				Execute()
			if err != nil {
				b.Fatal(err)
			}
			resp.Content()
		}
	})
}

// BenchmarkConcurrentRequests 并发请求性能测试
func BenchmarkConcurrentRequests(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 模拟一些处理时间
		time.Sleep(time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	concurrencyLevels := []int{1, 10, 50, 100}

	for _, concurrency := range concurrencyLevels {
		b.Run(fmt.Sprintf("Concurrency_%d", concurrency), func(b *testing.B) {
			session := NewSession()

			b.ResetTimer()
			b.SetParallelism(concurrency)

			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					resp, err := session.Get(server.URL).Execute()
					if err != nil {
						b.Fatal(err)
					}
					resp.Content()
				}
			})
		})
	}
}

// BenchmarkMemoryAllocation 内存分配性能测试
func BenchmarkMemoryAllocation(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	b.Run("RequestCreation", func(b *testing.B) {
		session := NewSession()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			req := session.Get(server.URL)
			_ = req
		}
	})

	b.Run("RequestExecution", func(b *testing.B) {
		session := NewSession()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			resp, err := session.Get(server.URL).Execute()
			if err != nil {
				b.Fatal(err)
			}
			_ = resp.Content()
		}
	})
}

// BenchmarkContextUsage Context使用性能测试
func BenchmarkContextUsage(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	b.Run("WithoutContext", func(b *testing.B) {
		session := NewSession()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			resp, err := session.Get(server.URL).Execute()
			if err != nil {
				b.Fatal(err)
			}
			resp.Content()
		}
	})

	b.Run("WithContext", func(b *testing.B) {
		session := NewSession()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			ctx := context.Background()
			resp, err := session.Get(server.URL).WithContext(ctx).Execute()
			if err != nil {
				b.Fatal(err)
			}
			resp.Content()
		}
	})

	b.Run("WithTimeoutContext", func(b *testing.B) {
		session := NewSession()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			resp, err := session.Get(server.URL).WithContext(ctx).Execute()
			cancel()
			if err != nil {
				b.Fatal(err)
			}
			resp.Content()
		}
	})
}

// BenchmarkConnectionPooling 连接池性能测试
func BenchmarkConnectionPooling(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	b.Run("DefaultPooling", func(b *testing.B) {
		session := NewSession()
		var wg sync.WaitGroup
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				resp, err := session.Get(server.URL).Execute()
				if err != nil {
					b.Error(err)
					return
				}
				resp.Content()
			}()
		}
		wg.Wait()
	})

	b.Run("HighPerformanceSession", func(b *testing.B) {
		session, _ := NewHighPerformanceSession()
		var wg sync.WaitGroup
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				resp, err := session.Get(server.URL).Execute()
				if err != nil {
					b.Error(err)
					return
				}
				resp.Content()
			}()
		}
		wg.Wait()
	})
}
