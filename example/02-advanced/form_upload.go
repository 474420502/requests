package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/474420502/requests"
)

// demonstrateFormUpload 展示表单和文件上传功能
func demonstrateFormUpload() {
	fmt.Println("=== 表单和文件上传演示 ===")

	session := requests.NewSession()

	// 1. 基础表单提交
	fmt.Println("1. 基础表单提交:")
	resp, err := session.Post("https://httpbin.org/post").
		SetFormFields(map[string]string{
			"username": "testuser",
			"email":    "test@example.com",
			"message":  "Hello World",
		}).
		Execute()

	if err != nil {
		log.Printf("基础表单请求失败: %v", err)
	} else {
		fmt.Printf("✓ 基础表单提交成功，状态码: %d\n", resp.GetStatusCode())
	}

	// 2. 类型安全的表单字段
	fmt.Println("\n2. 类型安全的表单字段:")
	resp, err = session.Post("https://httpbin.org/post").
		AddFormFieldInt("user_id", 12345).
		AddFormFieldFloat("score", 98.5).
		AddFormFieldBool("verified", true).
		AddFormField("name", "张三").
		Execute()

	if err != nil {
		log.Printf("类型安全表单请求失败: %v", err)
	} else {
		fmt.Printf("✓ 类型安全表单提交成功，状态码: %d\n", resp.GetStatusCode())
	}

	// 3. 混合类型表单
	fmt.Println("\n3. 混合类型表单:")
	formData := map[string]interface{}{
		"name":        "李四",
		"age":         28,
		"height":      175.5,
		"verified":    true,
		"description": "这是一个混合类型的表单",
	}

	resp, err = session.Post("https://httpbin.org/post").
		SetFormFieldsTyped(formData).
		Execute()

	if err != nil {
		log.Printf("混合类型表单请求失败: %v", err)
	} else {
		fmt.Printf("✓ 混合类型表单提交成功，状态码: %d\n", resp.GetStatusCode())
	}

	// 4. 文件上传
	fmt.Println("\n4. 文件上传:")

	// 模拟文件内容
	fileContent := strings.NewReader("这是一个测试文件的内容\n包含多行文本")

	resp, err = session.Post("https://httpbin.org/post").
		SetFormFields(map[string]string{
			"description": "文件上传测试",
			"category":    "test",
		}).
		AddFormFile("upload", "test.txt", fileContent).
		Execute()

	if err != nil {
		log.Printf("文件上传请求失败: %v", err)
	} else {
		fmt.Printf("✓ 文件上传成功，状态码: %d\n", resp.GetStatusCode())
	}

	// 5. 多文件上传
	fmt.Println("\n5. 多文件上传:")

	file1 := strings.NewReader("第一个文件的内容")
	file2 := strings.NewReader("第二个文件的内容")

	resp, err = session.Post("https://httpbin.org/post").
		AddFormField("description", "多文件上传测试").
		AddFormFile("file1", "document1.txt", file1).
		AddFormFile("file2", "document2.txt", file2).
		Execute()

	if err != nil {
		log.Printf("多文件上传请求失败: %v", err)
	} else {
		fmt.Printf("✓ 多文件上传成功，状态码: %d\n", resp.GetStatusCode())
	}

	// 6. URL编码表单数据
	fmt.Println("\n6. URL编码表单数据:")
	resp, err = session.Post("https://httpbin.org/post").
		SetBodyUrlencoded(map[string]string{
			"username": "urlencoded_user",
			"data":     "some data with spaces",
			"special":  "chars&=?#",
		}).
		Execute()

	if err != nil {
		log.Printf("URL编码表单请求失败: %v", err)
	} else {
		fmt.Printf("✓ URL编码表单提交成功，状态码: %d\n", resp.GetStatusCode())
	}

	// 7. 复杂的multipart表单（现代方式）
	fmt.Println("\n7. 复杂的multipart表单:")

	// 使用现代API设置表单字段和文件
	fileBytes := []byte("复杂表单中的文件内容")

	resp, err = session.Post("https://httpbin.org/post").
		SetFormFields(map[string]string{
			"user_name":  "复杂表单用户",
			"user_email": "complex@example.com",
		}).
		AddFormFile("attachment", "complex.txt", bytes.NewReader(fileBytes)).
		Execute()

	if err != nil {
		log.Printf("复杂multipart表单请求失败: %v", err)
	} else {
		fmt.Printf("✓ 复杂multipart表单提交成功，状态码: %d\n", resp.GetStatusCode())
	}

	fmt.Println("\n✅ 表单和文件上传演示完成")
}
