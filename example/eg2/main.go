package main

import (
	"log"

	"github.com/474420502/requests"
)

func main() {
	ses := requests.NewSession()
	tp := ses.Get("http://httpbin.org/anything")
	tp.SetBodyJson(`{"a": 1, "b": 2}`)
	resp, _ := tp.Execute()
	log.Println(string(resp.Content()))

	tp = ses.Get("http://httpbin.org/anything")
	tp.SetBodyJson(map[string]interface{}{"a": "1", "b": 2})
	resp, _ = tp.Execute()
	log.Println(string(resp.Content()))

	tp = ses.Post("http://httpbin.org/anything")
	tp.SetFormFileFromPath("file", "./tests/learn.js")
	resp, _ = tp.Execute()
	log.Println(string(resp.Content()))
}
