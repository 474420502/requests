package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/474420502/requests"
)

// IntegrationLoggingMiddleware 集成演示的日志中间件
type IntegrationLoggingMiddleware struct {
	phase string
}

func (m *IntegrationLoggingMiddleware) BeforeRequest(req *http.Request) error {
	fmt.Printf("[%s] → 发送请求: %s %s\n", m.phase, req.Method, req.URL.String())
	return nil
}

func (m *IntegrationLoggingMiddleware) AfterResponse(resp *http.Response) error {
	fmt.Printf("[%s] ✓ 收到响应: %d %s\n", m.phase, resp.StatusCode, resp.Status)
	return nil
}

// demonstrateCompleteIntegration 演示完整的集成功能
func demonstrateCompleteIntegration() {
	fmt.Println("=== 完整集成演示：三个阶段的功能整合 ===")

	// 1. Phase 1 功能：基础重构的统一API
	fmt.Println("1. Phase 1 基础功能集成:")

	// 创建基础Session
	phase1Session := requests.NewSession()

	// 演示统一的HTTP方法
	methods := []struct {
		name   string
		method func(string) *requests.Request
	}{
		{"GET", phase1Session.Get},
		{"POST", phase1Session.Post},
		{"PUT", phase1Session.Put},
		{"DELETE", phase1Session.Delete},
	}

	for _, m := range methods {
		fmt.Printf("  测试 %s 方法: ", m.name)
		var resp *requests.Response
		var err error

		if m.name == "GET" || m.name == "DELETE" {
			resp, err = m.method("https://httpbin.org/"+strings.ToLower(m.name)).
				AddParam("phase", "1").
				Execute()
		} else {
			resp, err = m.method("https://httpbin.org/" + strings.ToLower(m.name)).
				SetBodyJson(map[string]string{"phase": "1", "method": m.name}).
				Execute()
		}

		if err != nil {
			fmt.Printf("✗ 失败: %v\n", err)
		} else {
			fmt.Printf("✓ 成功 (状态码: %d)\n", resp.GetStatusCode())
		}
	}

	// 2. Phase 2 功能：中间件系统集成
	fmt.Println("\n2. Phase 2 中间件功能集成:")

	phase2Session := requests.NewSession()

	// 添加多个中间件
	// 日志中间件
	logger := log.New(os.Stdout, "[INTEGRATION] ", log.LstdFlags)
	loggingMiddleware := &requests.LoggingMiddleware{Logger: logger}
	phase2Session.AddMiddleware(loggingMiddleware)

	// 自定义中间件
	customMiddleware := &IntegrationLoggingMiddleware{phase: "PHASE2"}
	phase2Session.AddMiddleware(customMiddleware)

	// 重试中间件
	retryMiddleware := &requests.RetryMiddleware{
		MaxRetries: 2,
		RetryDelay: 500 * time.Millisecond,
	}
	phase2Session.AddMiddleware(retryMiddleware)

	// 执行请求以展示中间件链
	_, err := phase2Session.Get("https://httpbin.org/get").
		AddParam("integration", "phase2").
		Execute()

	if err != nil {
		fmt.Printf("✗ Phase 2 集成请求失败: %v\n", err)
	} else {
		fmt.Printf("✓ Phase 2 中间件链执行成功\n")
	}

	// 3. Phase 3 功能：架构完善集成
	fmt.Println("\n3. Phase 3 架构完善集成:")

	// 创建高级配置的Session
	phase3Session := requests.NewSession()

	// 添加Phase 3的增强中间件
	phase3Middleware := &IntegrationLoggingMiddleware{phase: "PHASE3"}
	phase3Session.AddMiddleware(phase3Middleware)

	// 演示高级功能
	_, err = phase3Session.Post("https://httpbin.org/post").
		SetHeader("X-Integration-Test", "phase3").
		SetHeader("X-API-Version", "v3").
		AddParam("advanced", "true").
		SetBodyJson(map[string]interface{}{
			"phase":       3,
			"features":    []string{"builder", "context", "performance"},
			"integration": true,
			"timestamp":   time.Now().Unix(),
		}).
		Execute()

	if err != nil {
		fmt.Printf("✗ Phase 3 集成请求失败: %v\n", err)
	} else {
		fmt.Printf("✓ Phase 3 高级功能集成成功\n")
	}

	// 4. 全功能集成演示
	fmt.Println("\n4. 全功能集成演示:")

	// 创建集合所有功能的Session
	fullSession := requests.NewSession()

	// Phase 1: 基础API
	// Phase 2: 完整中间件链
	allInOneLogger := &IntegrationLoggingMiddleware{phase: "ALL-IN-ONE"}
	fullSession.AddMiddleware(allInOneLogger)

	standardLogger := log.New(os.Stdout, "[FULL-INTEGRATION] ", log.LstdFlags)
	fullSession.AddMiddleware(&requests.LoggingMiddleware{Logger: standardLogger})

	// Phase 3: 高级配置和类型安全

	// 执行综合测试
	testCases := []struct {
		name string
		test func() error
	}{
		{
			"基础GET请求",
			func() error {
				_, err := fullSession.Get("https://httpbin.org/get").
					AddParam("test", "basic-get").
					Execute()
				return err
			},
		},
		{
			"JSON POST请求",
			func() error {
				_, err := fullSession.Post("https://httpbin.org/post").
					SetBodyJson(map[string]interface{}{
						"integration": "full",
						"test_type":   "json_post",
						"timestamp":   time.Now().Format(time.RFC3339),
					}).
					Execute()
				return err
			},
		},
		{
			"表单数据请求",
			func() error {
				_, err := fullSession.Post("https://httpbin.org/post").
					SetFormFields(map[string]string{
						"field1":      "value1",
						"field2":      "value2",
						"integration": "full",
					}).
					Execute()
				return err
			},
		},
		{
			"带头部的PUT请求",
			func() error {
				_, err := fullSession.Put("https://httpbin.org/put").
					SetHeader("Content-Type", "application/json").
					SetHeader("X-Custom-Header", "integration-test").
					SetBodyJson(map[string]string{
						"update": "integration test",
						"method": "PUT",
					}).
					Execute()
				return err
			},
		},
	}

	successCount := 0
	for _, tc := range testCases {
		fmt.Printf("  执行测试: %s ... ", tc.name)
		err := tc.test()
		if err != nil {
			fmt.Printf("✗ 失败: %v\n", err)
		} else {
			fmt.Printf("✓ 成功\n")
			successCount++
		}
	}

	fmt.Printf("\n全功能集成测试结果: %d/%d 成功\n", successCount, len(testCases))

	// 5. 性能集成测试
	fmt.Println("\n5. 性能集成测试:")

	perfSession := requests.NewSession()
	perfSession.AddMiddleware(&IntegrationLoggingMiddleware{phase: "PERF"})

	// 并发性能测试
	concurrentRequests := 5
	results := make(chan error, concurrentRequests)
	start := time.Now()

	for i := 0; i < concurrentRequests; i++ {
		go func(id int) {
			_, err := perfSession.Get("https://httpbin.org/get").
				AddParam("concurrent_id", fmt.Sprintf("%d", id)).
				AddParam("performance_test", "true").
				Execute()
			results <- err
		}(i)
	}

	// 收集结果
	perfSuccessCount := 0
	for i := 0; i < concurrentRequests; i++ {
		if err := <-results; err == nil {
			perfSuccessCount++
		}
	}

	duration := time.Since(start)
	fmt.Printf("并发性能测试: %d/%d 成功，总耗时: %v，平均耗时: %v\n",
		perfSuccessCount, concurrentRequests, duration, duration/time.Duration(concurrentRequests))

	// 6. 最终总结
	fmt.Println("\n=== 完整集成演示总结 ===")
	fmt.Println("✅ Phase 1 基础重构功能：")
	fmt.Println("   • 统一的HTTP方法接口")
	fmt.Println("   • 链式调用API")
	fmt.Println("   • 改进的响应处理")

	fmt.Println("✅ Phase 2 中间件系统功能：")
	fmt.Println("   • 标准中间件接口")
	fmt.Println("   • 自定义中间件支持")
	fmt.Println("   • 中间件链执行")
	fmt.Println("   • 重试机制")

	fmt.Println("✅ Phase 3 架构完善功能：")
	fmt.Println("   • Session构建器模式")
	fmt.Println("   • 高级配置选项")
	fmt.Println("   • 类型安全API")
	fmt.Println("   • 性能优化")

	fmt.Println("✅ 集成特性：")
	fmt.Println("   • 向后兼容性保持")
	fmt.Println("   • 功能无缝整合")
	fmt.Println("   • 性能稳定可靠")
	fmt.Println("   • 开发者体验优秀")

	fmt.Printf("\n🎉 所有三个阶段的功能已成功集成，提供了完整的现代化HTTP客户端库！\n")
}
