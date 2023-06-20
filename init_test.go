package requests

import (
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/elazarl/goproxy"
)

const ProxyAddress = "localhost:58080"

func init() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	go func() {
		proxy := goproxy.NewProxyHttpServer()
		proxy.Verbose = true
		http.ListenAndServe(ProxyAddress, proxy)
	}()

	cmd := exec.Command("/bin/bash", "-c", "docker ps | grep httpbin")
	_, err := cmd.Output()
	if err != nil {
		log.Println("recommend:\n1. docker run --rm -p 80:80 kennethreitz/httpbin\n2. echo \"127.0.0.1 httpbin.org\" >> /etc/hosts")
	}

	time.Sleep(time.Millisecond * 100)
}
