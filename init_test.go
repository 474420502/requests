package requests

import (
	"compress/flate"
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/elazarl/goproxy"
)

const ProxyAddress = "localhost:58080"

var ProxyServer *http.Server
var ProxyCloseChan chan int = make(chan int, 2)

var TestServer *http.ServeMux

func init() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	go func() {

		proxy := goproxy.NewProxyHttpServer()
		proxy.Verbose = true
		ProxyServer = &http.Server{Addr: ProxyAddress, Handler: proxy}

		go func() {
			for range ProxyCloseChan {

			}
			ProxyServer.Shutdown(context.TODO())
		}()

		ProxyServer.ListenAndServe()

	}()

	cmd := exec.Command("/bin/bash", "-c", "docker ps | grep httpbin")
	_, err := cmd.Output()
	if err != nil {
		log.Println("recommend:\n1. docker run --rm -p 80:80 kennethreitz/httpbin\n2. echo \"127.0.0.1 httpbin.org\" >> /etc/hosts")
	}

	TestServer = http.NewServeMux()
	TestServer.HandleFunc("/compress", func(w http.ResponseWriter, r *http.Request) {
		var writer io.Writer = w

		encodings := r.Header.Values("Accept-Encoding")
		for _, encoding := range encodings {
			if strings.Contains(encoding, "gzip") {
				w.Header().Set("Content-Encoding", "gzip")
				writer = gzip.NewWriter(writer)
				writer.Write([]byte("hello compress"))
				defer writer.(*gzip.Writer).Close()
				return
			} else if strings.Contains(encoding, "deflate") {
				w.Header().Set("Content-Encoding", "deflate")
				writer, err = flate.NewWriter(writer, flate.DefaultCompression)
				if err != nil {
					panic(err)
				}
				writer.Write([]byte("hello compress"))
				defer writer.(*flate.Writer).Close()
				return
			}
		}
		writer.Write([]byte("hello"))
	})

	TestServer.HandleFunc("/content-compress", func(w http.ResponseWriter, r *http.Request) {

		encodings := r.Header.Get("Content-Encoding")
		if strings.Contains(encodings, "gzip") {

			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				panic(err)
			}
			data, err := io.ReadAll(reader)
			if err != nil {
				panic(err)
			}
			var o map[string]interface{}
			err = json.Unmarshal(data, &o)
			if err != nil {
				panic(err)
			}

			if o["key"] != "hello compress" {
				panic(o)
			}

			w.Write([]byte(o["key"].(string)))

		} else if strings.Contains(encodings, "deflate") {

			reader := flate.NewReader(r.Body)
			data, err := io.ReadAll(reader)
			if err != nil {
				panic(err)
			}
			var o map[string]interface{}
			err = json.Unmarshal(data, &o)
			if err != nil {
				panic(err)
			}

			if o["key"] != "hello compress" {
				panic(o)
			}
			w.Write([]byte(o["key"].(string)))

		} else if strings.Contains(encodings, "br") {
			reader := brotli.NewReader(r.Body)
			data, err := io.ReadAll(reader)
			if err != nil {
				panic(err)
			}
			var o map[string]interface{}
			err = json.Unmarshal(data, &o)
			if err != nil {
				panic(err)
			}

			if o["key"] != "hello compress" {
				panic(o)
			}
			w.Write([]byte(o["key"].(string)))
		} else {
			w.Write([]byte("error compress"))
		}

	})

	// 添加 /anything 端点来模拟 httpbin.org/anything
	TestServer.HandleFunc("/anything", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// 读取请求体
		var bodyData interface{}
		var bodyString string
		if r.Body != nil {
			bodyBytes, err := io.ReadAll(r.Body)
			if err == nil && len(bodyBytes) > 0 {
				bodyString = string(bodyBytes)
				// 尝试解析JSON
				var jsonData interface{}
				if json.Unmarshal(bodyBytes, &jsonData) == nil {
					bodyData = jsonData
				}
			}
		}

		// 获取查询参数
		args := make(map[string]interface{})
		for key, values := range r.URL.Query() {
			if len(values) == 1 {
				args[key] = values[0]
			} else {
				args[key] = values
			}
		}

		// 获取表单数据
		form := make(map[string]interface{})
		if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
			r.ParseForm()
			for key, values := range r.PostForm {
				if len(values) == 1 {
					form[key] = values[0]
				} else {
					form[key] = values
				}
			}
		}

		// 获取请求头
		headers := make(map[string]string)
		for key, values := range r.Header {
			headers[key] = strings.Join(values, ", ")
		}

		// 构建响应
		response := map[string]interface{}{
			"args":    args,
			"data":    bodyString,
			"files":   map[string]interface{}{},
			"form":    form,
			"headers": headers,
			"json":    bodyData,
			"method":  r.Method,
			"origin":  "127.0.0.1", // 本地测试
			"url":     "http://localhost" + r.RequestURI,
		}

		responseBytes, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Write(responseBytes)
	})

	// 添加 /get 端点来模拟 httpbin.org/get
	TestServer.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// 获取查询参数
		args := make(map[string]interface{})
		for key, values := range r.URL.Query() {
			if len(values) == 1 {
				args[key] = values[0]
			} else {
				args[key] = values
			}
		}

		// 获取请求头
		headers := make(map[string]string)
		for key, values := range r.Header {
			headers[key] = strings.Join(values, ", ")
		}

		// 构建响应
		response := map[string]interface{}{
			"args":    args,
			"headers": headers,
			"origin":  "127.0.0.1",
			"url":     "http://localhost" + r.RequestURI,
		}

		responseBytes, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Write(responseBytes)
	})

	time.Sleep(time.Millisecond * 500)
}
