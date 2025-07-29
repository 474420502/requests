package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/474420502/requests"
)

func demonstrateImprovedAPIv2() {
	fmt.Println("=== Go Requests 库现代化改进演示 ===")

	// 1. 类型安全的Session创建（无panic）
	fmt.Println("1. 类型安全的Session创建:")
	session, err := requests.NewSessionWithOptions(
		requests.WithTimeout(30*time.Second),
		requests.WithUserAgent("ModernApp/2.0"),
	)
	if err != nil {
		log.Fatal("创建Session失败:", err)
	}

	// 新的类型安全配置方法
	session.Config().SetBasicAuth("username", "password")
	err = session.Config().SetProxy("http://proxy.example.com:8080")
	if err != nil {
		fmt.Printf("设置代理失败: %v\n", err)
	} else {
		fmt.Println("✓ Session创建成功，配置了超时、认证等")
	}

	// 2. 类型安全的查询参数
	fmt.Println("2. 类型安全的查询参数:")
	resp, err := session.Get("https://httpbin.org/get").
		AddQuery("name", "张三").
		AddQueryInt("age", 25).
		AddQueryBool("active", true).
		AddQueryFloat("score", 95.5).
		Execute()

	if err != nil {
		fmt.Printf("✗ 查询参数请求失败: %v\n", err)
	} else {
		fmt.Printf("✓ 类型安全查询参数请求成功，状态码: %d\n", resp.GetStatusCode())
	}

	// 3. 简单的路径参数替换
	fmt.Println("3. 简单的路径参数替换:")

	// 单个参数替换
	resp, err = session.Get("https://httpbin.org/status/{code}").
		SetPathParam("code", "200").
		Execute()

	if err != nil {
		fmt.Printf("✗ 路径参数请求失败: %v\n", err)
	} else {
		fmt.Printf("✓ 路径参数替换成功，状态码: %d\n", resp.GetStatusCode())
	}

	// 批量参数替换
	resp, err = session.Get("https://httpbin.org/delay/{seconds}").
		SetPathParams(map[string]string{
			"seconds": "1",
		}).
		Execute()

	if err != nil {
		fmt.Printf("✗ 批量路径参数请求失败: %v\n", err)
	} else {
		fmt.Printf("✓ 批量路径参数替换成功，状态码: %d\n", resp.GetStatusCode())
	}

	// 4. 类型安全的表单处理
	fmt.Println("4. 类型安全的表单处理:")

	// 使用新的SetFormFields方法
	resp, err = session.Post("https://httpbin.org/post").
		SetFormFields(map[string]string{
			"username": "testuser",
			"email":    "test@example.com",
		}).
		Execute()

	if err != nil {
		fmt.Printf("✗ 表单字段请求失败: %v\n", err)
	} else {
		fmt.Printf("✓ 类型安全表单字段请求成功，状态码: %d\n", resp.GetStatusCode())
	}

	// 添加文件字段
	fileContent := strings.NewReader("这是一个测试文件的内容")
	resp, err = session.Post("https://httpbin.org/post").
		SetFormFields(map[string]string{
			"description": "文件上传测试",
		}).
		AddFormFile("upload", "test.txt", fileContent).
		Execute()

	if err != nil {
		fmt.Printf("✗ 文件上传请求失败: %v\n", err)
	} else {
		fmt.Printf("✓ 类型安全文件上传请求成功，状态码: %d\n", resp.GetStatusCode())
	}

	// 5. 改进的JSON处理
	fmt.Println("5. 改进的JSON处理:")

	// 发送JSON
	jsonData := map[string]interface{}{
		"name":  "李四",
		"age":   28,
		"email": "lisi@example.com",
	}

	resp, err = session.Post("https://httpbin.org/post").
		SetBodyJSON(jsonData).
		Execute()

	if err != nil {
		fmt.Printf("✗ JSON请求失败: %v\n", err)
	} else {
		fmt.Printf("✓ JSON请求成功，状态码: %d\n", resp.GetStatusCode())

		// 使用新的结构体绑定方法
		type ResponseData struct {
			Json map[string]interface{} `json:"json"`
			URL  string                 `json:"url"`
		}

		var responseData ResponseData
		err = resp.BindJSON(&responseData)
		if err != nil {
			fmt.Printf("✗ JSON解析失败: %v\n", err)
		} else {
			fmt.Printf("✓ JSON解析成功，接收到的数据: %+v\n", responseData.Json)
		}
	}

	// 6. 错误处理改进演示
	fmt.Println("6. 错误处理改进:")

	// 无效URL的处理
	_, err = session.Get("invalid-url").Execute()
	if err != nil {
		fmt.Printf("✓ 正确捕获URL错误: %v\n", err)
	}

	// 无效代理配置的处理
	err = session.Config().SetProxy("invalid-proxy-url")
	if err != nil {
		fmt.Printf("✓ 正确捕获代理配置错误: %v\n", err)
	}

	fmt.Println("=== 现代化改进总结 ===")
	fmt.Println("✓ 彻底消除了panic，所有错误都通过error返回")
	fmt.Println("✓ 提供了类型安全的查询参数方法（AddQueryInt, AddQueryBool等）")
	fmt.Println("✓ 实现了简单的路径参数替换（SetPathParam, SetPathParams）")
	fmt.Println("✓ 重构了表单处理，提供SetFormFields和AddFormFile方法")
	fmt.Println("✓ 增强了JSON处理，添加了BindJSON方法")
	fmt.Println("✓ 所有配置方法都返回error，提高了健壮性")
	fmt.Println("✓ 向后兼容旧API，同时引导用户使用现代化的方法")
	fmt.Println("✓ 【新增】统一了内部架构，ParamQuery和ParamRegexp直接使用Request")
	fmt.Println("✓ 【新增】提供了类型安全的配置方法：SetBasicAuthString, SetProxyString, SetTimeoutDuration")
}

func demonstrateTypeSafeConfig() {
	fmt.Println("=== 类型安全的配置方法演示 ===")

	session := requests.NewSession()

	// 类型安全的基础认证
	session.Config().SetBasicAuthString("username", "password")
	fmt.Println("✓ 使用SetBasicAuthString设置认证，无需类型转换")

	// 类型安全的代理设置
	err := session.Config().SetProxyString("http://proxy.example.com:8080")
	if err != nil {
		fmt.Printf("代理设置错误: %v\n", err)
	} else {
		fmt.Println("✓ 使用SetProxyString设置代理，自动URL解析和错误处理")
	}

	// 清除代理
	session.Config().ClearProxy()
	fmt.Println("✓ 使用ClearProxy清除代理设置")

	// 类型安全的超时设置
	session.Config().SetTimeoutDuration(30 * time.Second)
	fmt.Println("✓ 使用SetTimeoutDuration设置超时，类型明确")

	session.Config().SetTimeoutSeconds(60)
	fmt.Println("✓ 使用SetTimeoutSeconds设置超时，更直观的秒数")

	fmt.Println("现在配置方法都是类型安全的，减少了运行时错误的可能性")
}

func main() {
	demonstrateImprovedAPIv2()
	fmt.Println()
	demonstrateTypeSafeConfig()
}
