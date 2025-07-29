package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/474420502/requests"
)

// demonstrateSessionManagement 展示Session的管理和配置
func demonstrateSessionManagement() {
	fmt.Println("=== Session管理演示 ===")

	// 1. 基础Session创建
	fmt.Println("1. 基础Session:")
	basicSession := requests.NewSession()
	resp, err := basicSession.Get("https://httpbin.org/get").Execute()
	if err != nil {
		log.Printf("基础Session请求失败: %v", err)
	} else {
		fmt.Printf("✓ 基础Session请求成功，状态码: %d\n", resp.GetStatusCode())
	}

	// 2. 使用函数式选项创建Session
	fmt.Println("\n2. 函数式选项Session:")
	advancedSession, err := requests.NewSessionWithOptions(
		requests.WithTimeout(30*time.Second),
		requests.WithUserAgent("MyApp/1.0"),
		requests.WithHeaders(map[string]string{
			"Accept": "application/json",
		}),
		requests.WithKeepAlives(true),
		requests.WithCompression(true),
	)
	if err != nil {
		log.Printf("创建高级Session失败: %v", err)
	} else {
		fmt.Println("✓ 高级Session创建成功")
		resp, err := advancedSession.Get("https://httpbin.org/headers").Execute()
		if err != nil {
			log.Printf("高级Session请求失败: %v", err)
		} else {
			fmt.Printf("✓ 高级Session请求成功，状态码: %d\n", resp.GetStatusCode())
		}
	}

	// 3. 预定义Session配置
	fmt.Println("\n3. 预定义Session配置:")

	// API专用Session
	_, err = requests.NewSessionForAPI()
	if err != nil {
		log.Printf("创建API Session失败: %v", err)
	} else {
		fmt.Println("✓ API专用Session创建成功")
	}

	// 爬虫专用Session
	_, err = requests.NewSessionForScraping()
	if err != nil {
		log.Printf("创建爬虫Session失败: %v", err)
	} else {
		fmt.Println("✓ 爬虫专用Session创建成功")
	}

	// 高性能Session
	_, err = requests.NewHighPerformanceSession()
	if err != nil {
		log.Printf("创建高性能Session失败: %v", err)
	} else {
		fmt.Println("✓ 高性能Session创建成功")
	}

	// 4. Session配置管理
	fmt.Println("\n4. Session配置管理:")
	session := requests.NewSession()

	// 设置认证
	session.Config().SetBasicAuthString("username", "password")
	fmt.Println("✓ 设置了基础认证")

	// 设置代理（仅演示，不实际使用）
	err = session.Config().SetProxyString("http://proxy.example.com:8080")
	if err != nil {
		fmt.Printf("  代理设置失败（预期）: %v\n", err)
	}

	// 设置超时
	session.Config().SetTimeoutDuration(60 * time.Second)
	fmt.Println("✓ 设置了60秒超时")

	// 5. Session持久化设置
	fmt.Println("\n5. Session持久化设置:")
	persistentSession := requests.NewSession()

	// 设置持久化Header
	persistentSession.SetHeader(http.Header{
		"Authorization": []string{"Bearer token123"},
		"Accept":        []string{"application/json"},
	})

	// 设置持久化Query参数
	persistentSession.SetQuery(url.Values{
		"api_key": {"your-api-key"},
		"version": {"v1"},
	})

	fmt.Println("✓ 设置了持久化Header和Query参数")

	// 验证持久化设置
	resp, err = persistentSession.Get("https://httpbin.org/get").Execute()
	if err != nil {
		log.Printf("持久化Session请求失败: %v", err)
	} else {
		fmt.Printf("✓ 持久化Session请求成功，状态码: %d\n", resp.GetStatusCode())
	}

	fmt.Println("\n✅ Session管理演示完成")
}
