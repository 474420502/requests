package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/474420502/requests"
)

// CustomLoggingMiddleware 自定义日志中间件
type CustomLoggingMiddleware struct {
	prefix string
}

func (m *CustomLoggingMiddleware) BeforeRequest(req *http.Request) error {
	fmt.Printf("%s → 发送请求: %s %s\n", m.prefix, req.Method, req.URL.String())
	return nil
}

func (m *CustomLoggingMiddleware) AfterResponse(resp *http.Response) error {
	fmt.Printf("%s ✓ 收到响应: %d %s\n", m.prefix, resp.StatusCode, resp.Status)
	return nil
}

// TimingMiddleware 计时中间件
type TimingMiddleware struct {
	startTime time.Time
}

func (m *TimingMiddleware) BeforeRequest(req *http.Request) error {
	m.startTime = time.Now()
	fmt.Printf("⏱️  开始计时: %s\n", req.URL.String())
	return nil
}

func (m *TimingMiddleware) AfterResponse(resp *http.Response) error {
	duration := time.Since(m.startTime)
	fmt.Printf("⏱️  请求耗时: %v\n", duration)
	return nil
}

// UserAgentMiddleware 用户代理中间件
type UserAgentMiddleware struct {
	userAgent string
}

func (m *UserAgentMiddleware) BeforeRequest(req *http.Request) error {
	req.Header.Set("User-Agent", m.userAgent)
	fmt.Printf("🏷️  设置User-Agent: %s\n", m.userAgent)
	return nil
}

func (m *UserAgentMiddleware) AfterResponse(resp *http.Response) error {
	return nil
}

// demonstratePhase2Features 演示第二阶段的中间件系统功能
func demonstratePhase2Features() {
	fmt.Println("=== Phase 2: 中间件系统演示 ===")

	// 1. 基础日志中间件演示
	fmt.Println("1. 基础日志中间件:")

	session := requests.NewSession()

	// 添加标准日志中间件
	logger := log.New(os.Stdout, "[HTTP] ", log.LstdFlags)
	loggingMiddleware := &requests.LoggingMiddleware{Logger: logger}
	session.AddMiddleware(loggingMiddleware)

	// 执行请求以演示中间件
	_, err := session.Get("https://httpbin.org/get").Execute()
	if err != nil {
		log.Printf("请求失败: %v", err)
	}

	// 2. 自定义中间件演示
	fmt.Println("\n2. 自定义中间件:")

	customSession := requests.NewSession()

	// 添加自定义日志中间件
	customLogging := &CustomLoggingMiddleware{prefix: "   [CUSTOM]"}
	customSession.AddMiddleware(customLogging)

	// 添加计时中间件
	timingMiddleware := &TimingMiddleware{}
	customSession.AddMiddleware(timingMiddleware)

	// 执行请求
	_, err = customSession.Get("https://httpbin.org/delay/1").Execute()
	if err != nil {
		log.Printf("自定义中间件请求失败: %v", err)
	}

	// 3. 多层中间件堆叠演示
	fmt.Println("\n3. 多层中间件堆叠:")

	multiSession := requests.NewSession()

	// 第一层：自定义日志
	layer1 := &CustomLoggingMiddleware{prefix: "   [层1]"}
	multiSession.AddMiddleware(layer1)

	// 第二层：用户代理设置
	layer2 := &UserAgentMiddleware{userAgent: "Phase2-Demo/1.0 (Middleware-System)"}
	multiSession.AddMiddleware(layer2)

	// 第三层：计时
	layer3 := &TimingMiddleware{}
	multiSession.AddMiddleware(layer3)

	// 执行请求以演示多层中间件
	_, err = multiSession.Get("https://httpbin.org/headers").Execute()
	if err != nil {
		log.Printf("多层中间件请求失败: %v", err)
	}

	// 4. 重试中间件演示
	fmt.Println("\n4. 重试机制中间件:")

	retrySession := requests.NewSession()

	// 添加重试中间件
	retryMiddleware := &requests.RetryMiddleware{
		MaxRetries: 3,
		RetryDelay: time.Second,
	}
	retrySession.AddMiddleware(retryMiddleware)

	// 添加日志以查看重试过程
	retryLogger := &CustomLoggingMiddleware{prefix: "   [RETRY]"}
	retrySession.AddMiddleware(retryLogger)

	// 尝试访问正常地址
	_, err = retrySession.Get("https://httpbin.org/get").Execute()
	if err != nil {
		fmt.Printf("重试后仍然失败: %v\n", err)
	} else {
		fmt.Printf("✓ 重试机制工作正常\n")
	}

	// 5. 中间件组合演示
	fmt.Println("\n5. 中间件组合和管理:")

	combinedSession := requests.NewSession()

	// 批量添加中间件
	middlewares := []requests.Middleware{
		&CustomLoggingMiddleware{prefix: "   [组合1]"},
		&UserAgentMiddleware{userAgent: "Combined-Demo/1.0"},
		&TimingMiddleware{},
	}

	// 使用SetMiddlewares批量设置
	combinedSession.SetMiddlewares(middlewares)

	// 执行请求
	_, err = combinedSession.Post("https://httpbin.org/post").
		SetBodyJson(map[string]string{
			"phase":   "2",
			"feature": "中间件组合",
			"test":    "middleware combination",
		}).
		Execute()

	if err != nil {
		log.Printf("组合中间件请求失败: %v", err)
	} else {
		fmt.Printf("✓ 中间件组合执行成功\n")
	}

	// 6. 中间件清理演示
	fmt.Println("\n6. 中间件管理:")

	// 清除所有中间件
	combinedSession.ClearMiddlewares()
	fmt.Printf("✓ 已清除所有中间件\n")

	// 添加单个简单中间件
	simpleMiddleware := &CustomLoggingMiddleware{prefix: "   [简单]"}
	combinedSession.AddMiddleware(simpleMiddleware)

	// 执行清理后的请求
	_, err = combinedSession.Get("https://httpbin.org/json").Execute()
	if err != nil {
		log.Printf("清理后请求失败: %v", err)
	} else {
		fmt.Printf("✓ 中间件清理和重新设置成功\n")
	}

	fmt.Println("\n✅ Phase 2 中间件系统演示完成")
	fmt.Println("主要功能:")
	fmt.Println("• 标准中间件接口（BeforeRequest/AfterResponse）")
	fmt.Println("• 自定义中间件开发")
	fmt.Println("• 多层中间件堆叠支持")
	fmt.Println("• 重试机制集成")
	fmt.Println("• 中间件批量管理")
	fmt.Println("• 动态中间件清理和重置")
}
