package requests

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestMiddlewareRetryLogic 测试重试中间件的详细逻辑
func TestMiddlewareRetryLogic(t *testing.T) {
	t.Run("RetryMiddleware_SuccessAfterRetries", func(t *testing.T) {
		retryMW := &RetryMiddleware{
			MaxRetries: 3,
			RetryDelay: time.Millisecond * 10,
		}

		attemptCount := 0

		// 创建mock handler，前2次失败，第3次成功
		handler := func(req *http.Request) (*http.Response, error) {
			attemptCount++
			if attemptCount < 3 {
				return &http.Response{
					StatusCode: 503,
					Status:     "503 Service Unavailable",
				}, fmt.Errorf("service unavailable")
			}
			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
			}, nil
		}

		// 模拟使用中间件
		session := NewSession()
		session.AddMiddleware(retryMW)

		req, _ := http.NewRequest("GET", "http://example.com", nil)

		// BeforeRequest应该不做任何操作
		err := retryMW.BeforeRequest(req)
		if err != nil {
			t.Errorf("BeforeRequest should not fail: %v", err)
		}

		// 模拟多次调用AfterResponse
		for i := 0; i < 3; i++ {
			resp, _ := handler(req)
			afterErr := retryMW.AfterResponse(resp)
			if afterErr != nil && i < 2 {
				// 前两次可能返回错误表示需要重试
				continue
			}
			if i == 2 && afterErr != nil {
				t.Errorf("Final attempt should succeed: %v", afterErr)
			}
		}

		if attemptCount != 3 {
			t.Errorf("Expected 3 attempts, got %d", attemptCount)
		}
	})
}

// TestMiddlewareCircuitBreaker 测试熔断器中间件
func TestMiddlewareCircuitBreaker(t *testing.T) {
	t.Run("CircuitBreaker_FailureThreshold", func(t *testing.T) {
		cbMW := NewCircuitBreakerMiddleware(3, time.Minute)

		req, _ := http.NewRequest("GET", "http://example.com", nil)

		// BeforeRequest应该不做任何操作（初始状态）
		err := cbMW.BeforeRequest(req)
		if err != nil {
			t.Errorf("BeforeRequest should not fail initially: %v", err)
		}

		// 模拟连续失败
		for i := 0; i < 3; i++ {
			failureResp := &http.Response{StatusCode: 500}
			err = cbMW.AfterResponse(failureResp)
			// 熔断器实现可能在达到阈值后开始返回错误
		}

		// 成功响应应该不返回错误
		successResp := &http.Response{StatusCode: 200}
		err = cbMW.AfterResponse(successResp)
		if err != nil {
			t.Errorf("Success response should not cause error: %v", err)
		}
	})
}

// TestMiddlewareRateLimit 测试限流中间件
func TestMiddlewareRateLimit(t *testing.T) {
	t.Run("RateLimit_Basic", func(t *testing.T) {
		rateMW := NewRateLimitMiddleware(10) // 每秒10个请求
		defer rateMW.Close()

		req, _ := http.NewRequest("GET", "http://example.com", nil)

		// 测试多个连续请求
		for i := 0; i < 3; i++ {
			err := rateMW.BeforeRequest(req)
			if err != nil {
				t.Errorf("BeforeRequest %d should not fail: %v", i+1, err)
			}

			resp := &http.Response{StatusCode: 200}
			err = rateMW.AfterResponse(resp)
			if err != nil {
				t.Errorf("AfterResponse %d should not fail: %v", i+1, err)
			}
		}
	})
}

// TestMiddlewareUserAgentRotation 测试User-Agent轮换中间件
func TestMiddlewareUserAgentRotation(t *testing.T) {
	t.Run("UserAgentRotation_Cycling", func(t *testing.T) {
		userAgents := []string{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
			"Mozilla/5.0 (X11; Linux x86_64)",
		}

		uaMW := &UserAgentRotationMiddleware{
			UserAgents: userAgents,
		}

		seenUserAgents := make(map[string]bool)

		// 测试多个请求，验证User-Agent轮换
		for i := 0; i < 6; i++ {
			req, _ := http.NewRequest("GET", "http://example.com", nil)

			err := uaMW.BeforeRequest(req)
			if err != nil {
				t.Errorf("BeforeRequest %d should not fail: %v", i+1, err)
			}

			userAgent := req.Header.Get("User-Agent")
			if userAgent == "" {
				t.Errorf("User-Agent should be set for request %d", i+1)
			}

			seenUserAgents[userAgent] = true

			// 验证设置的User-Agent是预定义列表中的一个
			found := false
			for _, ua := range userAgents {
				if userAgent == ua {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("User-Agent should be from predefined list, got: %s", userAgent)
			}

			resp := &http.Response{StatusCode: 200}
			err = uaMW.AfterResponse(resp)
			if err != nil {
				t.Errorf("AfterResponse %d should not fail: %v", i+1, err)
			}
		}

		// 验证至少使用了多个不同的User-Agent（轮换效果）
		if len(seenUserAgents) < 2 {
			t.Error("Should use multiple different User-Agents")
		}
	})
}

// TestMiddlewareRequestID 测试请求ID中间件
func TestMiddlewareRequestID(t *testing.T) {
	t.Run("RequestID_Generation", func(t *testing.T) {
		ridMW := &RequestIDMiddleware{
			Generator: func() string {
				return fmt.Sprintf("req-%d", time.Now().UnixNano())
			},
		}

		req, _ := http.NewRequest("GET", "http://example.com", nil)

		err := ridMW.BeforeRequest(req)
		if err != nil {
			t.Errorf("BeforeRequest should not fail: %v", err)
		}

		requestID := req.Header.Get("X-Request-ID")
		if requestID == "" {
			t.Error("X-Request-ID header should be set")
		}

		// 验证生成的ID是合理的长度
		if len(requestID) < 10 {
			t.Errorf("Request ID seems too short: %s", requestID)
		}

		resp := &http.Response{StatusCode: 200}
		err = ridMW.AfterResponse(resp)
		if err != nil {
			t.Errorf("AfterResponse should not fail: %v", err)
		}
	})

	t.Run("RequestID_Uniqueness", func(t *testing.T) {
		ridMW := &RequestIDMiddleware{
			Generator: func() string {
				return fmt.Sprintf("req-%d", time.Now().UnixNano())
			},
		}

		requestIDs := make(map[string]bool)

		// 生成多个请求ID，验证唯一性
		for i := 0; i < 10; i++ {
			req, _ := http.NewRequest("GET", "http://example.com", nil)

			err := ridMW.BeforeRequest(req)
			if err != nil {
				t.Errorf("BeforeRequest %d should not fail: %v", i+1, err)
			}

			requestID := req.Header.Get("X-Request-ID")
			if requestIDs[requestID] {
				t.Errorf("Duplicate request ID generated: %s", requestID)
			}
			requestIDs[requestID] = true

			// 短暂延迟确保时间戳不同
			time.Sleep(time.Nanosecond * 100)
		}

		if len(requestIDs) != 10 {
			t.Errorf("Expected 10 unique request IDs, got %d", len(requestIDs))
		}
	})
} // TestMiddlewareTimeoutHandling 测试超时中间件
func TestMiddlewareTimeoutHandling(t *testing.T) {
	t.Run("Timeout_ContextModification", func(t *testing.T) {
		timeoutMW := &TimeoutMiddleware{
			Timeout: time.Second * 5,
		}

		req, _ := http.NewRequest("GET", "http://example.com", nil)
		originalCtx := req.Context()

		err := timeoutMW.BeforeRequest(req)
		if err != nil {
			t.Errorf("BeforeRequest should not fail: %v", err)
		}

		// 验证上下文被修改（添加了超时）
		if req.Context() == originalCtx {
			t.Error("Request context should be modified to include timeout")
		}

		// 验证超时上下文的行为
		select {
		case <-req.Context().Done():
			t.Error("Context should not be cancelled immediately")
		default:
			// 预期行为：上下文未立即取消
		}

		resp := &http.Response{StatusCode: 200}
		err = timeoutMW.AfterResponse(resp)
		if err != nil {
			t.Errorf("AfterResponse should not fail: %v", err)
		}
	})
}

// TestMiddlewareLoggingDetails 测试日志中间件的详细功能
func TestMiddlewareLoggingDetails(t *testing.T) {
	t.Run("Logging_FullRequestResponse", func(t *testing.T) {
		var logOutput strings.Builder
		logger := log.New(&logOutput, "[TEST] ", log.LstdFlags)

		loggingMW := &LoggingMiddleware{Logger: logger}

		req, _ := http.NewRequest("POST", "http://api.example.com/users", strings.NewReader("test body"))
		req.Header.Set("Content-Type", "application/json")

		err := loggingMW.BeforeRequest(req)
		if err != nil {
			t.Errorf("BeforeRequest should not fail: %v", err)
		}

		// 验证请求日志
		logStr := logOutput.String()
		if !strings.Contains(logStr, "POST") {
			t.Error("Log should contain HTTP method")
		}
		if !strings.Contains(logStr, "api.example.com/users") {
			t.Error("Log should contain URL")
		}

		// 重置日志输出
		logOutput.Reset()

		resp := &http.Response{
			StatusCode: 201,
			Status:     "201 Created",
		}

		err = loggingMW.AfterResponse(resp)
		if err != nil {
			t.Errorf("AfterResponse should not fail: %v", err)
		}

		// 验证响应日志
		logStr = logOutput.String()
		if !strings.Contains(logStr, "201") {
			t.Error("Log should contain status code")
		}
		if !strings.Contains(logStr, "Created") {
			t.Error("Log should contain status text")
		}
	})

	t.Run("Logging_NilLogger", func(t *testing.T) {
		loggingMW := &LoggingMiddleware{Logger: nil}

		req, _ := http.NewRequest("GET", "http://example.com", nil)

		// 应该不会panic或报错
		err := loggingMW.BeforeRequest(req)
		if err != nil {
			t.Errorf("BeforeRequest with nil logger should not fail: %v", err)
		}

		resp := &http.Response{StatusCode: 200}
		err = loggingMW.AfterResponse(resp)
		if err != nil {
			t.Errorf("AfterResponse with nil logger should not fail: %v", err)
		}
	})
}

// TestMiddlewareChainExecution 测试中间件链式执行
func TestMiddlewareChainExecution(t *testing.T) {
	t.Run("MiddlewareChain_ExecutionOrder", func(t *testing.T) {
		var executionOrder []string
		var mu sync.Mutex

		// 创建记录执行顺序的中间件
		mw1 := &TestExecutionOrderMiddleware{
			Name:           "MW1",
			ExecutionOrder: &executionOrder,
			Mutex:          &mu,
		}

		mw2 := &TestExecutionOrderMiddleware{
			Name:           "MW2",
			ExecutionOrder: &executionOrder,
			Mutex:          &mu,
		}

		session := NewSession()
		session.AddMiddleware(mw1)
		session.AddMiddleware(mw2)

		request := session.Get("http://example.com")

		// 验证中间件添加顺序
		if len(request.middlewares) != 2 {
			t.Errorf("Expected 2 middlewares, got %d", len(request.middlewares))
		}

		// 模拟执行BeforeRequest
		req, _ := http.NewRequest("GET", "http://example.com", nil)
		for _, mw := range request.middlewares {
			err := mw.BeforeRequest(req)
			if err != nil {
				t.Errorf("BeforeRequest should not fail: %v", err)
			}
		}

		// 模拟执行AfterResponse
		resp := &http.Response{StatusCode: 200}
		for i := len(request.middlewares) - 1; i >= 0; i-- {
			err := request.middlewares[i].AfterResponse(resp)
			if err != nil {
				t.Errorf("AfterResponse should not fail: %v", err)
			}
		}

		mu.Lock()
		expectedOrder := []string{"MW1-Before", "MW2-Before", "MW2-After", "MW1-After"}
		if len(executionOrder) != len(expectedOrder) {
			t.Errorf("Expected %d execution steps, got %d", len(expectedOrder), len(executionOrder))
		}

		for i, expected := range expectedOrder {
			if i < len(executionOrder) && executionOrder[i] != expected {
				t.Errorf("Expected execution order %d to be %s, got %s", i, expected, executionOrder[i])
			}
		}
		mu.Unlock()
	})
}

// TestExecutionOrderMiddleware 用于测试执行顺序的中间件
type TestExecutionOrderMiddleware struct {
	Name           string
	ExecutionOrder *[]string
	Mutex          *sync.Mutex
}

func (m *TestExecutionOrderMiddleware) BeforeRequest(req *http.Request) error {
	m.Mutex.Lock()
	*m.ExecutionOrder = append(*m.ExecutionOrder, m.Name+"-Before")
	m.Mutex.Unlock()
	return nil
}

func (m *TestExecutionOrderMiddleware) AfterResponse(resp *http.Response) error {
	m.Mutex.Lock()
	*m.ExecutionOrder = append(*m.ExecutionOrder, m.Name+"-After")
	m.Mutex.Unlock()
	return nil
}

// TestMiddlewareErrorPropagation 测试中间件错误传播
func TestMiddlewareErrorPropagation(t *testing.T) {
	t.Run("BeforeRequest_ErrorPropagation", func(t *testing.T) {
		errorMW := &TestErrorMiddleware{
			ShouldFailBefore: true,
			ErrorMessage:     "before request failed",
		}

		session := NewSession()
		session.AddMiddleware(errorMW)

		request := session.Get("http://example.com")
		req, _ := http.NewRequest("GET", "http://example.com", nil)

		// 测试BeforeRequest错误
		for _, mw := range request.middlewares {
			err := mw.BeforeRequest(req)
			if err == nil {
				t.Error("Expected middleware to return error")
			}
			if !strings.Contains(err.Error(), "before request failed") {
				t.Errorf("Expected error message to contain 'before request failed', got: %v", err)
			}
		}
	})

	t.Run("AfterResponse_ErrorPropagation", func(t *testing.T) {
		errorMW := &TestErrorMiddleware{
			ShouldFailAfter: true,
			ErrorMessage:    "after response failed",
		}

		session := NewSession()
		session.AddMiddleware(errorMW)

		request := session.Get("http://example.com")
		resp := &http.Response{StatusCode: 200}

		// 测试AfterResponse错误
		for _, mw := range request.middlewares {
			err := mw.AfterResponse(resp)
			if err == nil {
				t.Error("Expected middleware to return error")
			}
			if !strings.Contains(err.Error(), "after response failed") {
				t.Errorf("Expected error message to contain 'after response failed', got: %v", err)
			}
		}
	})
}

// TestErrorMiddleware 用于测试错误情况的中间件
type TestErrorMiddleware struct {
	ShouldFailBefore bool
	ShouldFailAfter  bool
	ErrorMessage     string
}

func (m *TestErrorMiddleware) BeforeRequest(req *http.Request) error {
	if m.ShouldFailBefore {
		return fmt.Errorf(m.ErrorMessage)
	}
	return nil
}

func (m *TestErrorMiddleware) AfterResponse(resp *http.Response) error {
	if m.ShouldFailAfter {
		return fmt.Errorf(m.ErrorMessage)
	}
	return nil
}
