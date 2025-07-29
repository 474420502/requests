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
