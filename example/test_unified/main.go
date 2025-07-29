package main

import (
	"fmt"
	"log"
	"time"

	"github.com/474420502/requests"
	// 相对导入父目录
)

func main() {
	// 创建Session
	session := requests.NewSession()

	// 测试基本的GET请求（返回Request而不是Temporary）
	fmt.Println("=== 测试基本GET请求 ===")
	req := session.Get("http://httpbin.org/get")
	req.AddParam("test", "value")
	req.SetTimeout(10 * time.Second)

	// 测试便利方法
	text, err := req.Text()
	if err != nil {
		log.Printf("GET请求失败: %v", err)
	} else {
		fmt.Printf("响应长度: %d 字符\n", len(text))
	}

	// 测试中间件
	fmt.Println("\n=== 测试中间件 ===")
	session.AddMiddleware(&requests.LoggingMiddleware{})

	req2 := session.Post("http://httpbin.org/post")
	req2.SetBodyJSON(map[string]string{"message": "hello"})

	resp, err := req2.Execute()
	if err != nil {
		log.Printf("POST请求失败: %v", err)
	} else {
		fmt.Printf("状态码: %d\n", resp.GetStatusCode())
	}

	// 测试错误处理
	fmt.Println("\n=== 测试错误处理 ===")
	req3 := session.Get("invalid-url")
	if err := req3.Error(); err != nil {
		fmt.Printf("预期的错误: %v\n", err)
	}

	// 测试表单数据
	fmt.Println("\n=== 测试表单数据 ===")
	req4 := session.Post("http://httpbin.org/post")
	req4.SetBodyFormData(map[string]string{
		"field1": "value1",
		"field2": "value2",
	})

	bytes, err := req4.Bytes()
	if err != nil {
		log.Printf("表单请求失败: %v", err)
	} else {
		fmt.Printf("响应字节数: %d\n", len(bytes))
	}

	fmt.Println("\n=== 架构统一测试完成 ===")
}
