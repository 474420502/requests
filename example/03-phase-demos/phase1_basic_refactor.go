package main

import (
	"fmt"
	"log"

	"github.com/474420502/requests"
)

// demonstratePhase1Features 演示第一阶段的基础重构功能
func demonstratePhase1Features() {
	fmt.Println("=== Phase 1: 基础重构演示 ===")

	// 1. 清理后的Session创建
	fmt.Println("1. 简洁的Session创建:")
	session := requests.NewSession()
	fmt.Println("✓ Session创建成功")

	// 2. 统一的请求接口
	fmt.Println("\n2. 统一的HTTP方法支持:")

	// GET请求
	resp, err := session.Get("https://httpbin.org/get").Execute()
	if err != nil {
		log.Printf("GET请求失败: %v", err)
	} else {
		fmt.Printf("✓ GET请求成功，状态码: %d\n", resp.GetStatusCode())
	}

	// POST请求
	resp, err = session.Post("https://httpbin.org/post").
		SetBodyJson(map[string]string{"message": "Hello Phase 1"}).
		Execute()
	if err != nil {
		log.Printf("POST请求失败: %v", err)
	} else {
		fmt.Printf("✓ POST请求成功，状态码: %d\n", resp.GetStatusCode())
	}

	// PUT请求
	resp, err = session.Put("https://httpbin.org/put").
		SetBodyJson(map[string]string{"update": "data"}).
		Execute()
	if err != nil {
		log.Printf("PUT请求失败: %v", err)
	} else {
		fmt.Printf("✓ PUT请求成功，状态码: %d\n", resp.GetStatusCode())
	}

	// DELETE请求
	resp, err = session.Delete("https://httpbin.org/delete").Execute()
	if err != nil {
		log.Printf("DELETE请求失败: %v", err)
	} else {
		fmt.Printf("✓ DELETE请求成功，状态码: %d\n", resp.GetStatusCode())
	}

	// 3. 改进的响应处理
	fmt.Println("\n3. 改进的响应处理:")

	resp, err = session.Get("https://httpbin.org/json").Execute()
	if err != nil {
		log.Printf("JSON请求失败: %v", err)
	} else {
		fmt.Printf("✓ JSON响应处理成功\n")

		// 字符串内容
		content := resp.ContentString()
		fmt.Printf("   响应长度: %d 字符\n", len(content))

		// 字节内容
		bytes := resp.Content()
		fmt.Printf("   响应字节数: %d\n", len(bytes))

		// 状态信息
		fmt.Printf("   状态: %s\n", resp.GetStatus())
		fmt.Printf("   状态码: %d\n", resp.GetStatusCode())
	}

	// 4. 链式调用的流畅性
	fmt.Println("\n4. 流畅的链式调用:")

	resp, err = session.Post("https://httpbin.org/post").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "Phase1-Demo/1.0").
		SetHeader("Accept", "application/json").
		SetBodyJson(map[string]interface{}{
			"phase":       1,
			"feature":     "基础重构",
			"description": "展示链式调用的优雅性",
		}).
		Execute()

	if err != nil {
		log.Printf("链式调用请求失败: %v", err)
	} else {
		fmt.Printf("✓ 链式调用成功，状态码: %d\n", resp.GetStatusCode())
	}

	// 5. 简化的参数设置
	fmt.Println("\n5. 简化的参数设置:")

	resp, err = session.Get("https://httpbin.org/get").
		AddParam("phase", "1").
		AddParam("feature", "简化参数").
		AddParam("demo", "true").
		Execute()

	if err != nil {
		log.Printf("参数设置请求失败: %v", err)
	} else {
		fmt.Printf("✓ 参数设置成功，状态码: %d\n", resp.GetStatusCode())
	}

	// 6. 错误处理的改进
	fmt.Println("\n6. 改进的错误处理:")

	// 尝试访问不存在的地址
	_, err = session.Get("https://nonexistent-domain-12345.com").Execute()
	if err != nil {
		fmt.Printf("✓ 错误处理正常工作: %v\n", err)
	} else {
		fmt.Printf("✗ 错误处理异常\n")
	}

	// 7. 基础配置的清理
	fmt.Println("\n7. 清理后的基础配置:")

	// 创建新session并设置基础配置
	configuredSession := requests.NewSession()

	resp, err = configuredSession.Get("https://httpbin.org/headers").
		SetHeader("Custom-Header", "Phase1-Value").
		Execute()

	if err != nil {
		log.Printf("配置请求失败: %v", err)
	} else {
		fmt.Printf("✓ 基础配置成功，状态码: %d\n", resp.GetStatusCode())
	}

	fmt.Println("\n✅ Phase 1 基础重构演示完成")
	fmt.Println("主要改进:")
	fmt.Println("• 统一了HTTP方法接口")
	fmt.Println("• 改进了响应处理机制")
	fmt.Println("• 优化了链式调用体验")
	fmt.Println("• 简化了参数和头部设置")
	fmt.Println("• 加强了错误处理")
}
