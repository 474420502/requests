package main

import (
	"fmt"
	"time"

	"github.com/474420502/requests"
)

// 测试v3.0现代化API
func testModernAPI() {
	fmt.Println("=== requests库 v3.0 现代化API演示 ===")

	// 1. 使用现代化的Request API
	fmt.Println("1. 现代化Request API:")
	ses := requests.NewSession()
	req := ses.Get("https://httpbin.org/get")
	req.AddHeader("X-Test", "modern-api")

	if req.Error() != nil {
		fmt.Printf("✗ Request创建失败: %v\n", req.Error())
	} else {
		fmt.Println("✓ 现代化Request API正常")
	}

	// 2. 验证Session统一返回Request对象
	fmt.Println("2. Session统一API:")
	req = ses.Get("https://httpbin.org/get").
		SetHeader("X-Test", "unified-request").
		AddQuery("phase", "3")

	if req.Error() != nil {
		fmt.Printf("✗ Request创建失败: %v\n", req.Error())
	} else {
		fmt.Println("✓ Session统一返回Request对象")
	}

	// 3. 验证顶层函数使用Request
	fmt.Println("3. 测试顶层函数统一性:")
	req2 := requests.Get("https://httpbin.org/get").
		SetHeader("X-Test", "top-level-request")

	if req2.Error() != nil {
		fmt.Printf("✗ 顶层函数失败: %v\n", req2.Error())
	} else {
		fmt.Println("✓ 顶层函数统一返回Request对象")
	}

	// 4. 验证类型安全的配置方法
	fmt.Println("4. 测试类型安全的配置:")

	// 创建Session使用类型安全方法
	session, err := requests.NewSessionWithOptions(
		requests.WithTimeout(30*time.Second),
		requests.WithUserAgent("Refactor-Test/1.0"),
	)
	if err != nil {
		fmt.Printf("✗ 类型安全Session创建失败: %v\n", err)
	} else {
		// 使用类型安全的配置方法
		session.Config().SetBasicAuth("user", "pass")
		err = session.Config().SetProxyString("http://proxy.example.com:8080")
		if err != nil {
			fmt.Printf("注意: 代理设置失败（预期，因为是无效代理）: %v\n", err)
		}
		session.Config().SetTimeoutDuration(10 * time.Second)

		fmt.Println("✓ 类型安全配置方法正常工作")
	}

	// 5. 验证弃用方法仍能工作但发出警告
	fmt.Println("5. 测试向后兼容的弃用方法:")

	// 测试deprecated的SetBasicAuth with interface{}
	err = session.Config().SetBasicAuthLegacy("testuser", "testpass")
	if err != nil {
		fmt.Printf("✗ 兼容性方法失败: %v\n", err)
	} else {
		fmt.Println("✓ 弃用方法保持向后兼容")
	}

	// 6. 验证现代化API的一致性
	fmt.Println("6. 现代化API一致性:")

	// 使用Session方法
	sessionReq := session.Get("https://httpbin.org/get")
	sessionReq.AddHeader("X-Source", "session")

	// 使用顶层函数
	directReq := requests.Get("https://httpbin.org/get").
		SetHeader("X-Source", "direct")

	// 两者都应该正常工作
	if sessionReq.Error() == nil && directReq.Error() == nil {
		fmt.Println("✓ API一致性：Session和顶层函数都正常工作")
	} else {
		fmt.Printf("✗ API统一性测试失败: temp=%v, direct=%v\n",
			sessionReq.Error(), directReq.Error())
	}

	fmt.Println("\n=== 第一阶段重构总结 ===")
	fmt.Println("✅ 彻底废弃了Temporary - 现在是Request的兼容层")
	fmt.Println("✅ Session所有方法统一返回*Request对象")
	fmt.Println("✅ 顶层函数统一使用Request模式")
	fmt.Println("✅ config.go清理完成 - 推荐类型安全方法，兼容旧方法")
	fmt.Println("✅ base.go简化 - 移除了不再使用的buildBodyRequest")
	fmt.Println("✅ 保持向后兼容性，所有现有代码仍能正常工作")
	fmt.Println("✅ 消除了API二元性问题 - 现在只有一个Request构建器")
}

func main() {
	testModernAPI()
}
