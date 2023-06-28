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

	time.Sleep(time.Millisecond * 500)
}
