package requests

import (
	"net/http"
	"net/http/httptest"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestDefaultConnectionPoolConfig(t *testing.T) {
	config := DefaultConnectionPoolConfig()

	if config.MaxIdleConns != 100 {
		t.Errorf("Expected MaxIdleConns to be 100, got %d", config.MaxIdleConns)
	}

	if config.MaxIdleConnsPerHost != 10 {
		t.Errorf("Expected MaxIdleConnsPerHost to be 10, got %d", config.MaxIdleConnsPerHost)
	}

	if config.IdleConnTimeout != 90*time.Second {
		t.Errorf("Expected IdleConnTimeout to be 90s, got %v", config.IdleConnTimeout)
	}
}

func TestHighPerformanceConnectionPoolConfig(t *testing.T) {
	config := HighPerformanceConnectionPoolConfig()
	cpuCount := runtime.NumCPU()

	expectedMaxIdle := cpuCount * 20
	if config.MaxIdleConns != expectedMaxIdle {
		t.Errorf("Expected MaxIdleConns to be %d, got %d", expectedMaxIdle, config.MaxIdleConns)
	}

	if !config.EnableAdaptive {
		t.Error("Expected EnableAdaptive to be true for high performance config")
	}

	if config.ScaleUpThreshold != 0.7 {
		t.Errorf("Expected ScaleUpThreshold to be 0.7, got %f", config.ScaleUpThreshold)
	}
}

func TestLowLatencyConnectionPoolConfig(t *testing.T) {
	config := LowLatencyConnectionPoolConfig()

	if config.MaxIdleConns != 200 {
		t.Errorf("Expected MaxIdleConns to be 200, got %d", config.MaxIdleConns)
	}

	if config.DialTimeout != 10*time.Second {
		t.Errorf("Expected DialTimeout to be 10s, got %v", config.DialTimeout)
	}

	if config.TLSHandshakeTimeout != 3*time.Second {
		t.Errorf("Expected TLSHandshakeTimeout to be 3s, got %v", config.TLSHandshakeTimeout)
	}
}

func TestResourceConstrainedConnectionPoolConfig(t *testing.T) {
	config := ResourceConstrainedConnectionPoolConfig()

	if config.MaxIdleConns != 20 {
		t.Errorf("Expected MaxIdleConns to be 20, got %d", config.MaxIdleConns)
	}

	if config.MaxIdleConnsPerHost != 2 {
		t.Errorf("Expected MaxIdleConnsPerHost to be 2, got %d", config.MaxIdleConnsPerHost)
	}

	if config.MaxConnsPerHost != 5 {
		t.Errorf("Expected MaxConnsPerHost to be 5, got %d", config.MaxConnsPerHost)
	}
}

func TestNewAdaptiveConnectionPool(t *testing.T) {
	config := DefaultConnectionPoolConfig()
	pool := NewAdaptiveConnectionPool(config)
	defer pool.Close()

	if pool.config != config {
		t.Error("Config not set correctly")
	}

	if pool.transport == nil {
		t.Error("Transport not initialized")
	}

	if pool.stats == nil {
		t.Error("Stats not initialized")
	}
}

func TestAdaptiveConnectionPoolWithNilConfig(t *testing.T) {
	pool := NewAdaptiveConnectionPool(nil)
	defer pool.Close()

	if pool.config == nil {
		t.Error("Expected default config when nil is passed")
	}

	// 验证默认配置值
	if pool.config.MaxIdleConns != 100 {
		t.Errorf("Expected default MaxIdleConns to be 100, got %d", pool.config.MaxIdleConns)
	}
}

func TestCreateOptimizedTransport(t *testing.T) {
	config := DefaultConnectionPoolConfig()
	transport := createOptimizedTransport(config)

	if transport.MaxIdleConns != config.MaxIdleConns {
		t.Errorf("Expected MaxIdleConns to be %d, got %d",
			config.MaxIdleConns, transport.MaxIdleConns)
	}

	if transport.MaxIdleConnsPerHost != config.MaxIdleConnsPerHost {
		t.Errorf("Expected MaxIdleConnsPerHost to be %d, got %d",
			config.MaxIdleConnsPerHost, transport.MaxIdleConnsPerHost)
	}

	if transport.IdleConnTimeout != config.IdleConnTimeout {
		t.Errorf("Expected IdleConnTimeout to be %v, got %v",
			config.IdleConnTimeout, transport.IdleConnTimeout)
	}

	if transport.TLSHandshakeTimeout != config.TLSHandshakeTimeout {
		t.Errorf("Expected TLSHandshakeTimeout to be %v, got %v",
			config.TLSHandshakeTimeout, transport.TLSHandshakeTimeout)
	}

	if transport.ForceAttemptHTTP2 != true {
		t.Error("Expected ForceAttemptHTTP2 to be true")
	}

	if transport.DisableKeepAlives != false {
		t.Error("Expected DisableKeepAlives to be false")
	}
}

func TestConnectionPoolStats(t *testing.T) {
	pool := NewAdaptiveConnectionPool(DefaultConnectionPoolConfig())
	defer pool.Close()

	// 测试记录请求
	pool.RecordRequest(true, 100*time.Millisecond)
	pool.RecordRequest(false, 200*time.Millisecond)

	stats := pool.GetStats()

	if stats.TotalRequests != 2 {
		t.Errorf("Expected TotalRequests to be 2, got %d", stats.TotalRequests)
	}

	if stats.FailedConnections != 1 {
		t.Errorf("Expected FailedConnections to be 1, got %d", stats.FailedConnections)
	}

	if stats.AverageConnTime == 0 {
		t.Error("Expected AverageConnTime to be set")
	}
}

func TestConnectionPoolWithSession(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	// 测试使用高性能连接池
	session, err := NewSessionWithOptions(WithHighPerformanceConnectionPool())

	if err != nil {
		t.Fatalf("Failed to create session with high performance pool: %v", err)
	}

	// 验证传输层配置
	if session.transport == nil {
		t.Error("Transport not set")
	}

	cpuCount := runtime.NumCPU()
	expectedMaxIdle := cpuCount * 20
	if session.transport.MaxIdleConns != expectedMaxIdle {
		t.Errorf("Expected MaxIdleConns to be %d, got %d",
			expectedMaxIdle, session.transport.MaxIdleConns)
	}

	// 测试请求
	resp, err := session.Get(server.URL).Execute()
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.GetStatusCode() != 200 {
		t.Errorf("Expected status 200, got %d", resp.GetStatusCode())
	}
}

func TestLowLatencyConnectionPool(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	session, err := NewSessionWithOptions(WithLowLatencyConnectionPool())

	if err != nil {
		t.Fatalf("Failed to create session with low latency pool: %v", err)
	}

	// 验证配置
	if session.transport.MaxIdleConns != 200 {
		t.Errorf("Expected MaxIdleConns to be 200, got %d", session.transport.MaxIdleConns)
	}

	if session.transport.TLSHandshakeTimeout != 3*time.Second {
		t.Errorf("Expected TLSHandshakeTimeout to be 3s, got %v",
			session.transport.TLSHandshakeTimeout)
	}
}

func TestResourceConstrainedConnectionPool(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	session, err := NewSessionWithOptions(WithResourceConstrainedConnectionPool())

	if err != nil {
		t.Fatalf("Failed to create session with resource constrained pool: %v", err)
	}

	// 验证配置适合资源受限环境
	if session.transport.MaxIdleConns != 20 {
		t.Errorf("Expected MaxIdleConns to be 20, got %d", session.transport.MaxIdleConns)
	}

	if session.transport.MaxIdleConnsPerHost != 2 {
		t.Errorf("Expected MaxIdleConnsPerHost to be 2, got %d",
			session.transport.MaxIdleConnsPerHost)
	}
}

func TestCustomConnectionPoolConfig(t *testing.T) {
	customConfig := &ConnectionPoolConfig{
		MaxIdleConns:        50,
		MaxIdleConnsPerHost: 5,
		MaxConnsPerHost:     10,
		IdleConnTimeout:     60 * time.Second,
		DialTimeout:         20 * time.Second,
		KeepAliveTimeout:    20 * time.Second,
		TLSHandshakeTimeout: 8 * time.Second,
	}

	session, err := NewSessionWithOptions(WithConnectionPool(customConfig))

	if err != nil {
		t.Fatalf("Failed to create session with custom pool: %v", err)
	}

	// 验证自定义配置
	if session.transport.MaxIdleConns != 50 {
		t.Errorf("Expected MaxIdleConns to be 50, got %d", session.transport.MaxIdleConns)
	}

	if session.transport.MaxIdleConnsPerHost != 5 {
		t.Errorf("Expected MaxIdleConnsPerHost to be 5, got %d",
			session.transport.MaxIdleConnsPerHost)
	}
}

func TestConcurrentConnectionPoolUsage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond) // 模拟处理时间
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	session, err := NewSessionWithOptions(WithHighPerformanceConnectionPool())

	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// 并发发送请求测试连接池
	const numRequests = 50
	var wg sync.WaitGroup
	errors := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			resp, err := session.Get(server.URL).Execute()
			if err != nil {
				errors <- err
				return
			}

			if resp.GetStatusCode() != 200 {
				errors <- err
				return
			}
		}()
	}

	wg.Wait()
	close(errors)

	// 检查是否有错误
	for err := range errors {
		t.Errorf("Request failed: %v", err)
	}
}

func TestAdaptivePoolAdjustment(t *testing.T) {
	config := DefaultConnectionPoolConfig()
	config.EnableAdaptive = true
	config.ScaleUpThreshold = 0.5
	config.ScaleDownThreshold = 0.2

	pool := NewAdaptiveConnectionPool(config)
	defer pool.Close()

	// 模拟高使用率触发扩容
	originalMax := pool.config.MaxIdleConnsPerHost

	// 手动调用调整逻辑（实际中由监控协程处理）
	pool.stats.mutex.Lock()
	pool.stats.ActiveConnections = int64(float64(originalMax) * 0.8) // 80%使用率
	pool.stats.mutex.Unlock()

	pool.adjustConnectionPool()

	if pool.config.MaxIdleConnsPerHost <= originalMax {
		t.Error("Expected MaxIdleConnsPerHost to increase due to high utilization")
	}
}

func BenchmarkConnectionPoolPerformance(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	b.Run("DefaultPool", func(b *testing.B) {
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

	b.Run("HighPerformancePool", func(b *testing.B) {
		session, _ := NewSessionWithOptions(WithHighPerformanceConnectionPool())
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			resp, err := session.Get(server.URL).Execute()
			if err != nil {
				b.Fatal(err)
			}
			resp.Content()
		}
	})

	b.Run("LowLatencyPool", func(b *testing.B) {
		session, _ := NewSessionWithOptions(WithLowLatencyConnectionPool())
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
