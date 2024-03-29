package requests

import (
	"log"
	"testing"
	"time"
)

func TestRequestPool(t *testing.T) {
	pool := NewRequestPool(3) // 设置最大并发数3

	ses := NewSession()
	ses.Config().SetTimeout(3 * time.Second)
	// 添加6个GET请求到httpbin
	urls := []string{"http://httpbin.org/get",
		"http://httpbin.org/get",
		"http://httpbin.org/get"}
	for _, url := range urls {
		tp := ses.Get(url)
		pool.Add(tp)
	}

	// 执行并验证正确的响应
	resps := pool.Execute()
	if len(resps) != 3 {
		t.Fatal("expected 3 responses, got", len(resps))
	}
	for i, url := range urls {
		if resps[i].Error != nil {
			log.Println(resps[i].Json())
			t.Fatal("invalid response for", url, resps[i].Error)
		}
	}

	// 添加两个故意超时的请求
	pool.Add(ses.Get("http://httpbin.org/delay/4"))
	pool.Add(ses.Get("http://httpbin.org/delay/5"))

	resps = pool.Execute()
	if resps[0].Error != nil || resps[1].Error != nil || resps[2].Error != nil {
		t.Fatal("expected timeout for last two requests")
	}

	if resps[3].Error == nil || resps[4].Error == nil {
		t.Fatal("expected timeout for last two requests")
	}
}
