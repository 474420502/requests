package main

import (
	"fmt"
	"log"

	"github.com/474420502/requests"
)

// demonstrateSimpleRequests 展示基础的HTTP请求方法
func demonstrateSimpleRequests() {
	fmt.Println("=== 基础HTTP请求演示 ===")

	// 创建Session
	session := requests.NewSession()

	// 1. GET请求
	fmt.Println("1. GET请求:")
	resp, err := session.Get("https://httpbin.org/get").
		AddQuery("name", "张三").
		AddQuery("age", "25").
		Execute()

	if err != nil {
		log.Printf("GET请求失败: %v", err)
	} else {
		fmt.Printf("✓ GET请求成功，状态码: %d\n", resp.GetStatusCode())
	}

	// 2. POST请求 - JSON数据
	fmt.Println("\n2. POST请求（JSON）:")
	jsonData := map[string]interface{}{
		"name":  "李四",
		"email": "lisi@example.com",
		"age":   28,
	}

	resp, err = session.Post("https://httpbin.org/post").
		SetBodyJSON(jsonData).
		Execute()

	if err != nil {
		log.Printf("POST请求失败: %v", err)
	} else {
		fmt.Printf("✓ POST请求成功，状态码: %d\n", resp.GetStatusCode())
	}

	// 3. POST请求 - 表单数据
	fmt.Println("\n3. POST请求（表单）:")
	resp, err = session.Post("https://httpbin.org/post").
		SetFormFields(map[string]string{
			"username": "testuser",
			"password": "testpass",
		}).
		Execute()

	if err != nil {
		log.Printf("表单请求失败: %v", err)
	} else {
		fmt.Printf("✓ 表单请求成功，状态码: %d\n", resp.GetStatusCode())
	}

	// 4. PUT请求
	fmt.Println("\n4. PUT请求:")
	resp, err = session.Put("https://httpbin.org/put").
		SetBodyString("这是PUT请求的内容").
		SetHeader("Content-Type", "text/plain").
		Execute()

	if err != nil {
		log.Printf("PUT请求失败: %v", err)
	} else {
		fmt.Printf("✓ PUT请求成功，状态码: %d\n", resp.GetStatusCode())
	}

	// 5. DELETE请求
	fmt.Println("\n5. DELETE请求:")
	resp, err = session.Delete("https://httpbin.org/delete").
		AddQuery("id", "123").
		Execute()

	if err != nil {
		log.Printf("DELETE请求失败: %v", err)
	} else {
		fmt.Printf("✓ DELETE请求成功，状态码: %d\n", resp.GetStatusCode())
	}

	// 6. HEAD请求
	fmt.Println("\n6. HEAD请求:")
	resp, err = session.Head("https://httpbin.org/get").Execute()

	if err != nil {
		log.Printf("HEAD请求失败: %v", err)
	} else {
		fmt.Printf("✓ HEAD请求成功，状态码: %d\n", resp.GetStatusCode())
		fmt.Printf("  Content-Length: %d\n", resp.GetContentLength())
	}

	// 7. 设置请求头
	fmt.Println("\n7. 自定义请求头:")
	resp, err = session.Get("https://httpbin.org/headers").
		SetHeader("X-Custom-Header", "MyValue").
		SetHeader("User-Agent", "MyApp/1.0").
		Execute()

	if err != nil {
		log.Printf("自定义头部请求失败: %v", err)
	} else {
		fmt.Printf("✓ 自定义头部请求成功，状态码: %d\n", resp.GetStatusCode())
	}

	fmt.Println("\n✅ 基础HTTP请求演示完成")
}

func main() {
	demonstrateSimpleRequests()
}
