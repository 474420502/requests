package requests

import (
	"context"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/elazarl/goproxy"
)

const ProxyAddress = "localhost:58080"

var ProxyServer *http.Server
var ProxyCloseChan chan int = make(chan int, 2)

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

	time.Sleep(time.Millisecond * 500)
}
