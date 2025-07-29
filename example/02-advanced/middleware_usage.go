package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/474420502/requests"
)

// demonstrateMiddlewareUsage 展示中间件的使用
func demonstrateMiddlewareUsage() {
	fmt.Println("=== 中间件使用演示 ===")

	// 创建Session
	session, err := requests.NewSessionWithOptions(
		requests.WithTimeout(10*time.Second),
		requests.WithUserAgent("MiddlewareDemo/1.0"),
	)
	if err != nil {
		log.Fatal("创建Session失败:", err)
	}

	// 1. 日志中间件
	fmt.Println("1. 日志中间件:")
	logger := log.New(os.Stdout, "[HTTP] ", log.LstdFlags)
	loggingMiddleware := &requests.LoggingMiddleware{Logger: logger}
	session.AddMiddleware(loggingMiddleware)

	resp, err := session.Get("https://httpbin.org/get").
		AddQuery("middleware", "logging").
		Execute()

	if err != nil {
		fmt.Printf("✗ 请求失败: %v\n", err)
	} else {
		fmt.Printf("✓ 日志中间件请求成功，状态码: %d\n", resp.GetStatusCode())
	}

	// 2. 指标收集中间件
	fmt.Println("\n2. 指标收集中间件:")
	metricsMiddleware := &requests.MetricsMiddleware{
		RequestCounter: func(method, url string) {
			fmt.Printf("📊 发起请求: %s %s\n", method, url)
		},
		ResponseCounter: func(statusCode int, method, url string) {
			fmt.Printf("📊 收到响应: %d %s %s\n", statusCode, method, url)
		},
		DurationTracker: func(duration time.Duration, method, url string) {
			fmt.Printf("📊 请求耗时: %v %s %s\n", duration, method, url)
		},
	}

	// 创建新Session用于指标演示
	metricsSession := requests.NewSession()
	metricsSession.AddMiddleware(metricsMiddleware)

	resp, err = metricsSession.Get("https://httpbin.org/delay/1").Execute()
	if err != nil {
		fmt.Printf("✗ 指标请求失败: %v\n", err)
	} else {
		fmt.Printf("✓ 指标中间件请求成功，状态码: %d\n", resp.GetStatusCode())
	}

	// 3. 请求ID中间件
	fmt.Println("\n3. 请求ID中间件:")
	requestIDMiddleware := &requests.RequestIDMiddleware{
		Generator: func() string {
			return fmt.Sprintf("req-%d", time.Now().UnixNano())
		},
	}

	idSession := requests.NewSession()
	idSession.AddMiddleware(requestIDMiddleware)

	resp, err = idSession.Get("https://httpbin.org/headers").Execute()
	if err != nil {
		fmt.Printf("✗ 请求ID请求失败: %v\n", err)
	} else {
		fmt.Printf("✓ 请求ID中间件请求成功，状态码: %d\n", resp.GetStatusCode())
	}

	// 4. 用户代理轮换中间件
	fmt.Println("\n4. 用户代理轮换中间件:")
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
		"Mozilla/5.0 (Linux; Ubuntu 20.04) AppleWebKit/537.36",
	}

	uaMiddleware := requests.NewUserAgentRotationMiddleware(userAgents)
	uaSession := requests.NewSession()
	uaSession.AddMiddleware(uaMiddleware)

	for i := 0; i < 3; i++ {
		resp, err := uaSession.Get("https://httpbin.org/headers").Execute()
		if err != nil {
			fmt.Printf("✗ 轮换请求 %d 失败: %v\n", i+1, err)
		} else {
			fmt.Printf("✓ 轮换请求 %d 成功，状态码: %d\n", i+1, resp.GetStatusCode())
		}
	}

	// 5. 熔断器中间件
	fmt.Println("\n5. 熔断器中间件:")
	circuitBreaker := requests.NewCircuitBreakerMiddleware(2, 5*time.Second)
	cbSession := requests.NewSession()
	cbSession.AddMiddleware(circuitBreaker)

	// 模拟几次失败请求
	fmt.Println("  模拟失败请求:")
	for i := 0; i < 3; i++ {
		_, err := cbSession.Get("https://httpbin.org/status/500").Execute()
		if err != nil {
			fmt.Printf("    请求 %d: %v\n", i+1, err)
		}
	}

	// 现在熔断器应该是打开状态
	fmt.Println("  测试熔断器状态:")
	_, err = cbSession.Get("https://httpbin.org/get").Execute()
	if err != nil {
		fmt.Printf("✓ 熔断器正确阻止了请求: %v\n", err)
	} else {
		fmt.Println("✗ 熔断器没有生效")
	}

	fmt.Println("\n✅ 中间件使用演示完成")
}

func main() {
	demonstrateAsyncPatterns()
	demonstrateFormUpload()
	demonstrateMiddlewareUsage()
}
