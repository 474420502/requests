package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/474420502/requests"
)

func demonstratePhase3Improvements() {
	fmt.Println("=== 第三阶段改进：架构完善与开发者体验 ===")

	// 1. 演示强化的Session构建器
	fmt.Println("1. 强化的Session构建器:")

	// 使用预定义的Session配置
	_, err := requests.NewSessionForAPI()
	if err != nil {
		log.Printf("创建API Session失败: %v", err)
		return
	}
	fmt.Println("✓ 创建了专用于API调用的Session")

	_, err = requests.NewSessionForScraping()
	if err != nil {
		log.Printf("创建爬虫Session失败: %v", err)
		return
	}
	fmt.Println("✓ 创建了专用于网页抓取的Session")

	// 使用自定义选项创建Session
	customSession, err := requests.NewSessionWithOptions(
		requests.WithTimeout(15*time.Second),
		requests.WithUserAgent("MyApp/1.0"),
		requests.WithKeepAlives(true),
		requests.WithCompression(true),
		requests.WithMaxIdleConnsPerHost(5),
		requests.WithRetry(3, time.Second),
	)
	if err != nil {
		log.Printf("创建自定义Session失败: %v", err)
		return
	}
	fmt.Println("✓ 创建了自定义配置的Session（带重试功能）")

	// 2. 演示完善的中间件系统
	fmt.Println("\n2. 完善的中间件系统:")

	// 创建日志中间件
	logger := log.New(os.Stdout, "[HTTP] ", log.LstdFlags)
	loggingMiddleware := &requests.LoggingMiddleware{Logger: logger}

	// 创建指标收集中间件
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

	// 创建请求ID中间件
	requestIDMiddleware := &requests.RequestIDMiddleware{
		Generator: func() string {
			return fmt.Sprintf("req-%d", time.Now().UnixNano())
		},
	}

	// 添加中间件到Session
	customSession.AddMiddleware(loggingMiddleware)
	customSession.AddMiddleware(metricsMiddleware)
	customSession.AddMiddleware(requestIDMiddleware)

	resp, err := customSession.Get("https://httpbin.org/get").
		AddQuery("middleware", "demo").
		Execute()

	if err != nil {
		fmt.Printf("✗ 中间件请求失败: %v\n", err)
	} else {
		fmt.Printf("✓ 中间件请求成功，状态码: %d\n", resp.GetStatusCode())
	}

	// 3. 演示熔断器中间件
	fmt.Println("\n3. 熔断器中间件:")

	circuitBreaker := requests.NewCircuitBreakerMiddleware(2, 5*time.Second)
	testSession := requests.NewSession()
	testSession.AddMiddleware(circuitBreaker)

	// 模拟几次失败请求
	for i := 0; i < 3; i++ {
		_, err := testSession.Get("https://httpbin.org/status/500").Execute()
		if err != nil {
			fmt.Printf("  请求 %d 失败（预期）: %v\n", i+1, err)
		}
	}

	// 现在熔断器应该是打开状态
	_, err = testSession.Get("https://httpbin.org/get").Execute()
	if err != nil {
		fmt.Printf("✓ 熔断器生效，阻止了请求: %v\n", err)
	}

	// 4. 演示用户代理轮换中间件
	fmt.Println("\n4. 用户代理轮换中间件:")

	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
		"Mozilla/5.0 (Linux; Ubuntu 20.04) AppleWebKit/537.36",
	}

	uaMiddleware := requests.NewUserAgentRotationMiddleware(userAgents)
	rotationSession := requests.NewSession()
	rotationSession.AddMiddleware(uaMiddleware)

	for i := 0; i < 3; i++ {
		resp, err := rotationSession.Get("https://httpbin.org/headers").Execute()
		if err != nil {
			fmt.Printf("✗ 轮换请求 %d 失败: %v\n", i+1, err)
		} else {
			fmt.Printf("✓ 轮换请求 %d 成功，状态码: %d\n", i+1, resp.GetStatusCode())
		}
	}

	// 5. 演示上下文和取消功能
	fmt.Println("\n5. 上下文和取消功能:")

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 创建带默认上下文的Session
	contextSession, err := requests.NewSessionWithOptions(
		requests.WithContext(ctx),
		requests.WithTimeout(5*time.Second),
	)
	if err != nil {
		log.Printf("创建上下文Session失败: %v", err)
		return
	}

	start := time.Now()
	_, err = contextSession.Get("https://httpbin.org/delay/3").Execute()
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("✓ 上下文超时正确工作，耗时: %v, 错误: %v\n", duration, err)
	} else {
		fmt.Printf("✗ 上下文超时未生效，耗时: %v\n", duration)
	}

	// 6. 演示高性能Session
	fmt.Println("\n6. 高性能Session:")

	highPerfSession, err := requests.NewHighPerformanceSession()
	if err != nil {
		log.Printf("创建高性能Session失败: %v", err)
		return
	}

	// 并发请求测试
	start = time.Now()
	concurrentRequests := 5
	results := make(chan error, concurrentRequests)

	for i := 0; i < concurrentRequests; i++ {
		go func(id int) {
			_, err := highPerfSession.Get("https://httpbin.org/get").
				AddQueryInt("request_id", id).
				Execute()
			results <- err
		}(i)
	}

	// 等待所有请求完成
	successCount := 0
	for i := 0; i < concurrentRequests; i++ {
		if err := <-results; err == nil {
			successCount++
		}
	}

	totalTime := time.Since(start)
	fmt.Printf("✓ 高性能Session完成 %d/%d 个并发请求，总耗时: %v\n",
		successCount, concurrentRequests, totalTime)

	fmt.Println("\n=== 第三阶段改进总结 ===")
	fmt.Println("✓ 强化了Session构建器，提供多种预定义配置")
	fmt.Println("✓ 完善了中间件系统，支持日志、指标、缓存、熔断等功能")
	fmt.Println("✓ 统一了内部架构，改进了上下文传递")
	fmt.Println("✓ 提供了丰富的配置选项和性能优化")
	fmt.Println("✓ 增强了开发者体验，提供了类型安全和现代化的API")
	fmt.Println("✓ 保持了向后兼容性，现有代码可以无缝升级")
}
