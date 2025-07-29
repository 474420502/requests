package requests

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Middleware 定义中间件接口
type Middleware interface {
	// BeforeRequest 在请求发送前调用
	BeforeRequest(req *http.Request) error
	// AfterResponse 在收到响应后调用
	AfterResponse(resp *http.Response) error
}

// LoggingMiddleware 日志记录中间件
type LoggingMiddleware struct {
	Logger *log.Logger
}

func (m *LoggingMiddleware) BeforeRequest(req *http.Request) error {
	if m.Logger != nil {
		m.Logger.Printf("发送请求: %s %s", req.Method, req.URL.String())
	}
	return nil
}

func (m *LoggingMiddleware) AfterResponse(resp *http.Response) error {
	if m.Logger != nil {
		m.Logger.Printf("收到响应: %d %s", resp.StatusCode, resp.Status)
	}
	return nil
}

// RetryMiddleware 重试中间件
type RetryMiddleware struct {
	MaxRetries int
	RetryDelay time.Duration
}

func (m *RetryMiddleware) BeforeRequest(req *http.Request) error {
	// 在请求前不需要做什么
	return nil
}

func (m *RetryMiddleware) AfterResponse(resp *http.Response) error {
	// 可以在这里决定是否需要重试
	// 实际的重试逻辑需要在更高层实现
	return nil
}

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	TokenProvider func() (string, error)
}

func (m *AuthMiddleware) BeforeRequest(req *http.Request) error {
	if m.TokenProvider != nil {
		token, err := m.TokenProvider()
		if err != nil {
			return fmt.Errorf("获取认证令牌失败: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return nil
}

func (m *AuthMiddleware) AfterResponse(resp *http.Response) error {
	// 可以在这里处理401响应，刷新token等
	return nil
}

// RateLimitMiddleware 限流中间件
type RateLimitMiddleware struct {
	requests chan struct{}
	ticker   *time.Ticker
}

func NewRateLimitMiddleware(requestsPerSecond int) *RateLimitMiddleware {
	requests := make(chan struct{}, requestsPerSecond)
	ticker := time.NewTicker(time.Second / time.Duration(requestsPerSecond))

	// 启动令牌桶填充
	go func() {
		for range ticker.C {
			select {
			case requests <- struct{}{}:
			default:
				// 桶已满，丢弃令牌
			}
		}
	}()

	return &RateLimitMiddleware{
		requests: requests,
		ticker:   ticker,
	}
}

func (m *RateLimitMiddleware) BeforeRequest(req *http.Request) error {
	// 等待令牌
	<-m.requests
	return nil
}

func (m *RateLimitMiddleware) AfterResponse(resp *http.Response) error {
	return nil
}

func (m *RateLimitMiddleware) Close() {
	m.ticker.Stop()
	close(m.requests)
}

// RequestWithMiddleware 支持中间件的请求结构
type RequestWithMiddleware struct {
	*Request
	middlewares []Middleware
}

// AddMiddleware 添加中间件
func (r *RequestWithMiddleware) AddMiddleware(middleware Middleware) *RequestWithMiddleware {
	r.middlewares = append(r.middlewares, middleware)
	return r
}

// ExecuteWithMiddleware 执行带中间件的请求
func (r *RequestWithMiddleware) ExecuteWithMiddleware() (*Response, error) {
	if r.err != nil {
		return nil, fmt.Errorf("request validation failed: %w", r.err)
	}

	// 构建HTTP请求
	httpReq, err := r.buildHTTPRequest()
	if err != nil {
		return nil, fmt.Errorf("failed to build HTTP request: %w", err)
	}

	// 执行BeforeRequest中间件
	for _, middleware := range r.middlewares {
		if err := middleware.BeforeRequest(httpReq); err != nil {
			return nil, fmt.Errorf("middleware BeforeRequest failed: %w", err)
		}
	}

	// 应用超时
	if r.timeout > 0 {
		ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
		defer cancel()
		httpReq = httpReq.WithContext(ctx)
	} else if r.ctx != context.Background() {
		httpReq = httpReq.WithContext(r.ctx)
	}

	// 执行请求
	resp, err := r.session.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request execution failed: %w", err)
	}

	// 执行AfterResponse中间件
	for _, middleware := range r.middlewares {
		if err := middleware.AfterResponse(resp); err != nil {
			return nil, fmt.Errorf("middleware AfterResponse failed: %w", err)
		}
	}

	// 转换为我们的Response类型
	myResponse, err := FromHTTPResponse(resp, r.session.Is.isDecompressNoAccept)
	if err != nil {
		return nil, fmt.Errorf("failed to process response: %w", err)
	}

	myResponse.readResponse = resp
	return myResponse, nil
}

// WithMiddlewares 为Request添加中间件支持
func (r *Request) WithMiddlewares(middlewares ...Middleware) *RequestWithMiddleware {
	return &RequestWithMiddleware{
		Request:     r,
		middlewares: middlewares,
	}
}

// MetricsMiddleware 指标收集中间件
type MetricsMiddleware struct {
	RequestCounter  func(method, url string)
	ResponseCounter func(statusCode int, method, url string)
	DurationTracker func(duration time.Duration, method, url string)
	startTime       time.Time
}

func (m *MetricsMiddleware) BeforeRequest(req *http.Request) error {
	m.startTime = time.Now()
	if m.RequestCounter != nil {
		m.RequestCounter(req.Method, req.URL.String())
	}
	return nil
}

func (m *MetricsMiddleware) AfterResponse(resp *http.Response) error {
	duration := time.Since(m.startTime)

	if m.ResponseCounter != nil {
		m.ResponseCounter(resp.StatusCode, resp.Request.Method, resp.Request.URL.String())
	}

	if m.DurationTracker != nil {
		m.DurationTracker(duration, resp.Request.Method, resp.Request.URL.String())
	}

	return nil
}

// CacheMiddleware 缓存中间件
type CacheMiddleware struct {
	Cache map[string]*CacheEntry
}

type CacheEntry struct {
	Response  *http.Response
	ExpiresAt time.Time
}

func NewCacheMiddleware() *CacheMiddleware {
	return &CacheMiddleware{
		Cache: make(map[string]*CacheEntry),
	}
}

func (m *CacheMiddleware) BeforeRequest(req *http.Request) error {
	// 只缓存GET请求
	if req.Method != http.MethodGet {
		return nil
	}

	key := req.URL.String()
	if entry, exists := m.Cache[key]; exists {
		if time.Now().Before(entry.ExpiresAt) {
			// 缓存命中且未过期
			// 这里需要一个机制来返回缓存的响应
			// 实际实现中可能需要修改接口设计
		}
	}

	return nil
}

func (m *CacheMiddleware) AfterResponse(resp *http.Response) error {
	// 只缓存GET请求的200响应
	if resp.Request.Method != http.MethodGet || resp.StatusCode != http.StatusOK {
		return nil
	}

	// 检查Cache-Control头
	cacheControl := resp.Header.Get("Cache-Control")
	if cacheControl == "no-cache" || cacheControl == "no-store" {
		return nil
	}

	key := resp.Request.URL.String()
	m.Cache[key] = &CacheEntry{
		Response:  resp,
		ExpiresAt: time.Now().Add(5 * time.Minute), // 默认5分钟过期
	}

	return nil
}

// RequestIDMiddleware 请求ID中间件
type RequestIDMiddleware struct {
	Generator func() string
}

func (m *RequestIDMiddleware) BeforeRequest(req *http.Request) error {
	if m.Generator != nil {
		requestID := m.Generator()
		req.Header.Set("X-Request-ID", requestID)
	}
	return nil
}

func (m *RequestIDMiddleware) AfterResponse(resp *http.Response) error {
	return nil
}

// TimeoutMiddleware 超时中间件（基于Context）
type TimeoutMiddleware struct {
	Timeout time.Duration
}

func (m *TimeoutMiddleware) BeforeRequest(req *http.Request) error {
	if m.Timeout > 0 {
		ctx, cancel := context.WithTimeout(req.Context(), m.Timeout)
		// 注意：这里需要确保cancel被调用，但在当前架构下比较困难
		// 实际使用中可能需要调整接口设计
		_ = cancel
		*req = *req.WithContext(ctx)
	}
	return nil
}

func (m *TimeoutMiddleware) AfterResponse(resp *http.Response) error {
	return nil
}

// CircuitBreakerMiddleware 熔断器中间件
type CircuitBreakerMiddleware struct {
	FailureThreshold int
	ResetTimeout     time.Duration

	failures    int
	lastFailure time.Time
	state       CircuitState
}

type CircuitState int

const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

func NewCircuitBreakerMiddleware(failureThreshold int, resetTimeout time.Duration) *CircuitBreakerMiddleware {
	return &CircuitBreakerMiddleware{
		FailureThreshold: failureThreshold,
		ResetTimeout:     resetTimeout,
		state:            CircuitClosed,
	}
}

func (m *CircuitBreakerMiddleware) BeforeRequest(req *http.Request) error {
	switch m.state {
	case CircuitOpen:
		if time.Since(m.lastFailure) > m.ResetTimeout {
			m.state = CircuitHalfOpen
		} else {
			return fmt.Errorf("circuit breaker is open")
		}
	case CircuitHalfOpen:
		// 允许一个请求通过
	case CircuitClosed:
		// 正常状态，允许请求
	}

	return nil
}

func (m *CircuitBreakerMiddleware) AfterResponse(resp *http.Response) error {
	if resp.StatusCode >= 500 {
		m.failures++
		m.lastFailure = time.Now()

		if m.failures >= m.FailureThreshold {
			m.state = CircuitOpen
		}
	} else {
		// 成功响应，重置失败计数
		m.failures = 0
		if m.state == CircuitHalfOpen {
			m.state = CircuitClosed
		}
	}

	return nil
}

// UserAgentRotationMiddleware 用户代理轮换中间件
type UserAgentRotationMiddleware struct {
	UserAgents []string
	current    int
}

func NewUserAgentRotationMiddleware(userAgents []string) *UserAgentRotationMiddleware {
	return &UserAgentRotationMiddleware{
		UserAgents: userAgents,
	}
}

func (m *UserAgentRotationMiddleware) BeforeRequest(req *http.Request) error {
	if len(m.UserAgents) > 0 {
		req.Header.Set("User-Agent", m.UserAgents[m.current])
		m.current = (m.current + 1) % len(m.UserAgents)
	}
	return nil
}

func (m *UserAgentRotationMiddleware) AfterResponse(resp *http.Response) error {
	return nil
}
