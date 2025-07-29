package main

import (
	"fmt"
	"time"

	"github.com/474420502/requests"
)

// demonstrateAsyncPatterns 展示异步请求和并发模式
func demonstrateAsyncPatterns() {
	fmt.Println("=== 异步请求和并发模式演示 ===")

	session := requests.NewSession()

	// 1. 并发请求 - 基本模式
	fmt.Println("1. 并发请求 - 基本模式:")
	urls := []string{
		"https://httpbin.org/delay/1",
		"https://httpbin.org/delay/2",
		"https://httpbin.org/delay/1",
	}

	start := time.Now()
	results := make(chan struct {
		url      string
		status   int
		duration time.Duration
		err      error
	}, len(urls))

	// 启动并发请求
	for _, url := range urls {
		go func(u string) {
			reqStart := time.Now()
			resp, err := session.Get(u).Execute()
			var status int
			if resp != nil {
				status = resp.GetStatusCode()
			}
			results <- struct {
				url      string
				status   int
				duration time.Duration
				err      error
			}{u, status, time.Since(reqStart), err}
		}(url)
	}

	// 收集结果
	for i := 0; i < len(urls); i++ {
		result := <-results
		if result.err != nil {
			fmt.Printf("✗ %s 失败: %v\n", result.url, result.err)
		} else {
			fmt.Printf("✓ %s 成功，状态: %d，耗时: %v\n",
				result.url, result.status, result.duration)
		}
	}
	fmt.Printf("总耗时: %v\n", time.Since(start))

	// 2. 有限并发控制
	fmt.Println("\n2. 有限并发控制:")

	urls = []string{
		"https://httpbin.org/uuid",
		"https://httpbin.org/uuid",
		"https://httpbin.org/uuid",
		"https://httpbin.org/uuid",
		"https://httpbin.org/uuid",
	}

	maxConcurrent := 2 // 最大并发数
	semaphore := make(chan struct{}, maxConcurrent)
	results = make(chan struct {
		url      string
		status   int
		duration time.Duration
		err      error
	}, len(urls))

	start = time.Now()

	// 启动控制并发的请求
	for _, url := range urls {
		go func(u string) {
			semaphore <- struct{}{}        // 获取信号量
			defer func() { <-semaphore }() // 释放信号量

			reqStart := time.Now()
			resp, err := session.Get(u).Execute()
			var status int
			if resp != nil {
				status = resp.GetStatusCode()
			}
			results <- struct {
				url      string
				status   int
				duration time.Duration
				err      error
			}{u, status, time.Since(reqStart), err}
		}(url)
	}

	// 收集结果
	successCount := 0
	for i := 0; i < len(urls); i++ {
		result := <-results
		if result.err != nil {
			fmt.Printf("✗ 请求 %d 失败: %v\n", i+1, result.err)
		} else {
			fmt.Printf("✓ 请求 %d 成功，状态: %d，耗时: %v\n",
				i+1, result.status, result.duration)
			successCount++
		}
	}
	fmt.Printf("成功率: %d/%d，总耗时: %v\n",
		successCount, len(urls), time.Since(start))

	// 3. 流水线模式 - 请求处理分离
	fmt.Println("\n3. 流水线模式 - 请求处理分离:")

	requestChan := make(chan string, 10)
	responseChan := make(chan struct {
		url    string
		data   string
		status int
		err    error
	}, 10)

	// 请求生成器
	go func() {
		defer close(requestChan)
		urls := []string{
			"https://httpbin.org/json",
			"https://httpbin.org/headers",
			"https://httpbin.org/user-agent",
		}
		for _, url := range urls {
			requestChan <- url
		}
	}()

	// 请求处理器（多个worker）
	workerCount := 2
	for i := 0; i < workerCount; i++ {
		go func(workerID int) {
			for url := range requestChan {
				resp, err := session.Get(url).Execute()
				var data string
				var status int
				if resp != nil {
					status = resp.GetStatusCode()
					bodyData := resp.Content()
					if len(bodyData) > 0 {
						data = string(bodyData[:min(100, len(bodyData))])
					}
				}
				responseChan <- struct {
					url    string
					data   string
					status int
					err    error
				}{url, data, status, err}
			}
		}(i)
	}

	// 结果处理器
	go func() {
		defer close(responseChan)
		time.Sleep(time.Second * 3) // 等待所有请求完成
	}()

	// 收集流水线结果
	pipelineResults := 0
	for response := range responseChan {
		if response.err != nil {
			fmt.Printf("✗ 流水线请求失败 %s: %v\n", response.url, response.err)
		} else {
			fmt.Printf("✓ 流水线请求成功 %s，状态: %d\n",
				response.url, response.status)
			if len(response.data) > 0 {
				fmt.Printf("   数据预览: %s...\n", response.data)
			}
		}
		pipelineResults++
	}

	// 4. 超时和取消控制
	fmt.Println("\n4. 超时和取消控制:")

	// 创建带超时的请求
	timeoutSession := requests.NewSession()

	start = time.Now()
	resp, err := timeoutSession.Get("https://httpbin.org/delay/3").
		SetTimeout(2 * time.Second).
		Execute()

	if err != nil {
		fmt.Printf("✓ 请求按预期超时: %v，耗时: %v\n", err, time.Since(start))
	} else {
		fmt.Printf("✗ 请求未按预期超时，状态: %d\n", resp.GetStatusCode())
	}

	// 5. 批处理模式
	fmt.Println("\n5. 批处理模式:")

	type BatchRequest struct {
		URL    string
		Method string
		Data   interface{}
	}

	requests := []BatchRequest{
		{"https://httpbin.org/get", "GET", nil},
		{"https://httpbin.org/post", "POST", map[string]string{"key": "value"}},
		{"https://httpbin.org/put", "PUT", map[string]string{"update": "data"}},
	}

	batchResults := make([]struct {
		Request BatchRequest
		Status  int
		Error   error
	}, len(requests))

	batchChan := make(chan struct {
		index  int
		status int
		err    error
	}, len(requests))

	// 批量处理
	for i, req := range requests {
		go func(index int, batchReq BatchRequest) {
			var err error

			switch batchReq.Method {
			case "GET":
				_, err = session.Get(batchReq.URL).Execute()
			case "POST":
				_, err = session.Post(batchReq.URL).SetBodyJson(batchReq.Data).Execute()
			case "PUT":
				_, err = session.Put(batchReq.URL).SetBodyJson(batchReq.Data).Execute()
			}

			var status int
			if err == nil {
				status = 200 // 假设成功为200
			}

			batchChan <- struct {
				index  int
				status int
				err    error
			}{index, status, err}
		}(i, req)
	}

	// 收集批处理结果
	for i := 0; i < len(requests); i++ {
		result := <-batchChan
		batchResults[result.index] = struct {
			Request BatchRequest
			Status  int
			Error   error
		}{requests[result.index], result.status, result.err}
	}

	// 输出批处理结果
	for i, result := range batchResults {
		if result.Error != nil {
			fmt.Printf("✗ 批处理 %d [%s %s] 失败: %v\n",
				i+1, result.Request.Method, result.Request.URL, result.Error)
		} else {
			fmt.Printf("✓ 批处理 %d [%s %s] 成功，状态: %d\n",
				i+1, result.Request.Method, result.Request.URL, result.Status)
		}
	}

	fmt.Println("\n✅ 异步请求和并发模式演示完成")
}

// 辅助函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
