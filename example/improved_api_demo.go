package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/474420502/requests"
)

func demonstrateImprovedAPI() {
	fmt.Println("=== 改进后的 Go Requests 库 API 演示 ===")

	// 1. 函数式选项模式创建Session
	fmt.Println("1. 使用函数式选项创建Session:")
	session, err := requests.NewSessionWithOptions(
		requests.WithTimeout(30*time.Second),
		requests.WithUserAgent("MyApp/1.0"),
		requests.WithHeaders(map[string]string{
			"Accept": "application/json",
		}),
		requests.WithKeepAlives(true),
		requests.WithCompression(true),
	)
	if err != nil {
		log.Fatal("创建Session失败:", err)
	}
	fmt.Println("✓ Session创建成功，配置了超时、User-Agent、默认头部等")

	// 2. 健壮的错误处理
	fmt.Println("\n2. 健壮的错误处理:")
	resp, err := session.Get("invalid-url").Execute()
	if err != nil {
		fmt.Printf("✓ 正确捕获了URL错误: %v\n", err)
	} else {
		fmt.Printf("✗ 应该捕获URL错误，但得到了响应: %v\n", resp)
	}

	// 3. Context支持
	fmt.Println("\n3. Context支持（超时和取消）:")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond) // 极短超时用于演示
	defer cancel()

	start := time.Now()
	_, err = session.Get("http://httpbin.org/delay/1").
		WithContext(ctx).
		Execute()
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("✓ 请求被正确取消，耗时: %v, 错误: %v\n", duration, err)
	}

	// 4. 链式调用与类型安全
	fmt.Println("\n4. 链式调用与类型安全:")

	// JSON请求
	jsonData := map[string]interface{}{
		"name": "张三",
		"age":  25,
		"city": "北京",
	}

	resp, err = session.Post("http://httpbin.org/post").
		SetHeader("Authorization", "Bearer token123").
		SetCookieValue("session_id", "abc123").
		SetBodyJSON(jsonData).
		WithTimeout(10 * time.Second).
		Execute()

	if err != nil {
		fmt.Printf("✗ JSON请求失败: %v\n", err)
	} else {
		fmt.Printf("✓ JSON请求成功，状态码: %d\n", resp.GetStatusCode())
	}

	// 5. 表单请求
	fmt.Println("\n5. 表单请求:")
	formData := map[string]string{
		"username": "testuser",
		"password": "testpass",
	}

	resp, err = session.Post("http://httpbin.org/post").
		SetBodyFormValues(formData).
		Execute()

	if err != nil {
		fmt.Printf("✗ 表单请求失败: %v\n", err)
	} else {
		fmt.Printf("✓ 表单请求成功，状态码: %d\n", resp.GetStatusCode())
	}

	// 6. 查询参数
	fmt.Println("\n6. 查询参数:")
	resp, err = session.Get("http://httpbin.org/get").
		AddQuery("page", "1").
		AddQuery("size", "10").
		AddQuery("category", "tech").
		Execute()

	if err != nil {
		fmt.Printf("✗ 查询参数请求失败: %v\n", err)
	} else {
		fmt.Printf("✓ 查询参数请求成功，状态码: %d\n", resp.GetStatusCode())
	}

	// 7. 预定义Session配置
	fmt.Println("\n7. 预定义Session配置:")

	_, err = requests.NewSessionForAPI()
	if err != nil {
		log.Printf("创建API Session失败: %v", err)
	} else {
		fmt.Println("✓ API Session创建成功（30s超时，压缩开启，Keep-Alive开启）")
	}

	_, err = requests.NewSessionForScraping()
	if err != nil {
		log.Printf("创建爬虫Session失败: %v", err)
	} else {
		fmt.Println("✓ 爬虫Session创建成功（浏览器User-Agent，10s超时）")
	}

	fmt.Println("\n=== API改进总结 ===")
	fmt.Println("✓ 彻底消除了panic，改为返回error")
	fmt.Println("✓ 引入了context.Context支持超时和取消")
	fmt.Println("✓ 采用函数式选项模式，配置更灵活")
	fmt.Println("✓ 保持了流式API的便利性")
	fmt.Println("✓ 提供了类型安全的方法替代interface{}")
	fmt.Println("✓ 重命名Temporary为Request，更符合直觉")
	fmt.Println("✓ 添加了预定义的Session配置")
}
