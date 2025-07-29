package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/474420502/requests"
)

func demonstrateMiddleware() {
	fmt.Println("=== 中间件系统演示 ===")

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

	resp, err := session.Get("http://httpbin.org/get").
		WithMiddlewares(loggingMiddleware).
		ExecuteWithMiddleware()

	if err != nil {
		fmt.Printf("✗ 请求失败: %v\n", err)
	} else {
		fmt.Printf("✓ 请求成功，状态码: %d\n", resp.GetStatusCode())
	}

	// 2. 认证中间件
	fmt.Println("\n2. 认证中间件:")
	authMiddleware := &requests.AuthMiddleware{
		TokenProvider: func() (string, error) {
			// 模拟获取JWT token
			return "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...", nil
		},
	}

	resp, err = session.Get("http://httpbin.org/bearer").
		WithMiddlewares(authMiddleware).
		ExecuteWithMiddleware()

	if err != nil {
		fmt.Printf("✗ 认证请求失败: %v\n", err)
	} else {
		fmt.Printf("✓ 认证请求成功，状态码: %d\n", resp.GetStatusCode())
	}

	// 3. 限流中间件
	fmt.Println("\n3. 限流中间件 (2 requests/second):")
	rateLimitMiddleware := requests.NewRateLimitMiddleware(2) // 每秒2个请求
	defer rateLimitMiddleware.Close()

	start := time.Now()
	for i := 0; i < 3; i++ {
		resp, err := session.Get("http://httpbin.org/get").
			WithMiddlewares(rateLimitMiddleware, loggingMiddleware).
			ExecuteWithMiddleware()

		if err != nil {
			fmt.Printf("✗ 限流请求 #%d 失败: %v\n", i+1, err)
		} else {
			fmt.Printf("✓ 限流请求 #%d 成功，状态码: %d, 耗时: %v\n",
				i+1, resp.GetStatusCode(), time.Since(start))
		}
	}

	// 4. 组合多个中间件
	fmt.Println("\n4. 组合多个中间件:")
	resp, err = session.Post("http://httpbin.org/post").
		SetBodyJSON(map[string]string{"message": "Hello with middleware"}).
		WithMiddlewares(loggingMiddleware, authMiddleware).
		ExecuteWithMiddleware()

	if err != nil {
		fmt.Printf("✗ 组合中间件请求失败: %v\n", err)
	} else {
		fmt.Printf("✓ 组合中间件请求成功，状态码: %d\n", resp.GetStatusCode())
	}

	fmt.Println("\n=== 中间件系统优势 ===")
	fmt.Println("✓ 关注点分离：认证、日志、限流等逻辑独立")
	fmt.Println("✓ 可组合性：可以自由组合不同的中间件")
	fmt.Println("✓ 可扩展性：容易添加新的中间件类型")
	fmt.Println("✓ 一致性：所有请求都能受益于中间件逻辑")
}
