package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/474420502/requests"
)

// ConcurrentTest 并发性能测试结构
type ConcurrentTest struct {
	concurrency int
	requests    int
	session     *requests.Session
}

// NewConcurrentTest 创建并发测试实例
func NewConcurrentTest(concurrency, totalRequests int) *ConcurrentTest {
	return &ConcurrentTest{
		concurrency: concurrency,
		requests:    totalRequests,
		session:     requests.NewSession(),
	}
}

// Result 测试结果结构
type Result struct {
	TotalRequests   int
	SuccessRequests int
	FailedRequests  int
	TotalTime       time.Duration
	MinTime         time.Duration
	MaxTime         time.Duration
	AvgTime         time.Duration
	RequestsPerSec  float64
}

// RunTest 执行并发测试
func (ct *ConcurrentTest) RunTest() *Result {
	fmt.Printf("开始并发测试: %d 并发数, %d 总请求数\n", ct.concurrency, ct.requests)

	var wg sync.WaitGroup
	results := make(chan time.Duration, ct.requests)
	errors := make(chan error, ct.requests)

	// 控制并发数的通道
	semaphore := make(chan struct{}, ct.concurrency)

	startTime := time.Now()

	// 启动请求
	for i := 0; i < ct.requests; i++ {
		wg.Add(1)
		go func(requestID int) {
			defer wg.Done()

			// 获取信号量
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			reqStart := time.Now()
			_, err := ct.session.Get("https://httpbin.org/get").
				AddParam("request_id", fmt.Sprintf("%d", requestID)).
				AddParam("concurrent_test", "true").
				Execute()
			reqDuration := time.Since(reqStart)

			results <- reqDuration
			errors <- err
		}(i)
	}

	// 等待所有请求完成
	wg.Wait()
	close(results)
	close(errors)

	totalTime := time.Since(startTime)

	// 统计结果
	var successCount, failedCount int
	var minTime, maxTime, totalReqTime time.Duration
	minTime = time.Hour // 初始化为一个大值

	for duration := range results {
		totalReqTime += duration
		if duration < minTime {
			minTime = duration
		}
		if duration > maxTime {
			maxTime = duration
		}
	}

	for err := range errors {
		if err == nil {
			successCount++
		} else {
			failedCount++
		}
	}

	avgTime := totalReqTime / time.Duration(ct.requests)
	requestsPerSec := float64(successCount) / totalTime.Seconds()

	return &Result{
		TotalRequests:   ct.requests,
		SuccessRequests: successCount,
		FailedRequests:  failedCount,
		TotalTime:       totalTime,
		MinTime:         minTime,
		MaxTime:         maxTime,
		AvgTime:         avgTime,
		RequestsPerSec:  requestsPerSec,
	}
}

// PrintResult 打印测试结果
func (r *Result) PrintResult() {
	fmt.Println("\n=== 并发测试结果 ===")
	fmt.Printf("总请求数: %d\n", r.TotalRequests)
	fmt.Printf("成功请求: %d\n", r.SuccessRequests)
	fmt.Printf("失败请求: %d\n", r.FailedRequests)
	fmt.Printf("成功率: %.2f%%\n", float64(r.SuccessRequests)/float64(r.TotalRequests)*100)
	fmt.Printf("总耗时: %v\n", r.TotalTime)
	fmt.Printf("最小响应时间: %v\n", r.MinTime)
	fmt.Printf("最大响应时间: %v\n", r.MaxTime)
	fmt.Printf("平均响应时间: %v\n", r.AvgTime)
	fmt.Printf("QPS (每秒请求数): %.2f\n", r.RequestsPerSec)
}

// MemoryStats 内存统计
func printMemoryStats(label string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%s - 内存使用:\n", label)
	fmt.Printf("  分配内存: %d KB\n", m.Alloc/1024)
	fmt.Printf("  系统内存: %d KB\n", m.Sys/1024)
	fmt.Printf("  GC次数: %d\n", m.NumGC)
	fmt.Println()
}

func main() {
	fmt.Println("=== Requests库并发性能测试 ===")

	// 测试前的内存状态
	printMemoryStats("测试开始前")

	// 测试场景1: 低并发，少量请求
	fmt.Println("场景1: 低并发测试 (5并发, 20请求)")
	test1 := NewConcurrentTest(5, 20)
	result1 := test1.RunTest()
	result1.PrintResult()

	printMemoryStats("场景1完成后")

	// 测试场景2: 中等并发
	fmt.Println("\n场景2: 中等并发测试 (10并发, 50请求)")
	test2 := NewConcurrentTest(10, 50)
	result2 := test2.RunTest()
	result2.PrintResult()

	printMemoryStats("场景2完成后")

	// 测试场景3: 高并发
	fmt.Println("\n场景3: 高并发测试 (20并发, 100请求)")
	test3 := NewConcurrentTest(20, 100)
	result3 := test3.RunTest()
	result3.PrintResult()

	printMemoryStats("场景3完成后")

	// 比较分析
	fmt.Println("\n=== 性能对比分析 ===")
	fmt.Printf("场景1 QPS: %.2f, 平均响应时间: %v\n", result1.RequestsPerSec, result1.AvgTime)
	fmt.Printf("场景2 QPS: %.2f, 平均响应时间: %v\n", result2.RequestsPerSec, result2.AvgTime)
	fmt.Printf("场景3 QPS: %.2f, 平均响应时间: %v\n", result3.RequestsPerSec, result3.AvgTime)

	// 强制GC并显示最终内存状态
	runtime.GC()
	printMemoryStats("测试完成后(GC)")

	fmt.Println("✅ 并发性能测试完成")
}
