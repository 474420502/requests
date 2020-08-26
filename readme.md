# requests is easy used by the program of crawl(spider)

eg1:
```go
	ses := requests.NewSession()
	resp, err := ses.Get("http://httpbin.org/get").Execute()
	if err != nil {
		log.Panic(err)
	}
	log.Println(string(resp.Content()))
``` 

eg2:
```go
    ses := requests.NewSession() 
	tp := ses.Get("http://httpbin.org/anything")
	tp.SetBodyAuto(`{"a": 1, "b": 2}`)
	resp, _ := tp.Execute()
	log.Println(string(resp.Content()))
	// {
	// 	"args": {},
	// 	"data": "{\"a\": 1, \"b\": 2}",
	// 	"files": {},
	// 	"form": {},
	// 	"headers": {
	// 	  "Connection": "close",
	// 	  "Content-Length": "16",
	// 	  "Content-Type": "application/json",
	// 	  "Host": "httpbin.org",
	// 	  "User-Agent": "Go-http-client/1.1"
	// 	},
	// 	"json": {
	// 	  "a": 1,
	// 	  "b": 2
	// 	},
	// 	"method": "GET",
	// 	"origin": "172.17.0.1",
	// 	"url": "http://httpbin.org/anything"
	//   }

	tp = ses.Get("http://httpbin.org/anything")
	tp.SetBodyAuto(map[string]interface{}{"a": "1", "b": 2})
	resp, _ = tp.Execute()
	log.Println(string(resp.Content()))
	// {
	// 	"args": {},
	// 	"data": "{\"a\":\"1\",\"b\":2}",
	// 	"files": {},
	// 	"form": {},
	// 	"headers": {
	// 	  "Connection": "close",
	// 	  "Content-Length": "15",
	// 	  "Content-Type": "application/json",
	// 	  "Host": "httpbin.org",
	// 	  "User-Agent": "Go-http-client/1.1"
	// 	},
	// 	"json": {
	// 	  "a": "1",
	// 	  "b": 2
	// 	},
	// 	"method": "GET",
	// 	"origin": "172.17.0.1",
	// 	"url": "http://httpbin.org/anything"
	//   }

	tp = ses.Post("http://httpbin.org/anything")
	tp.SetBodyAuto("./tests/learn.js", TypeFormData)
	resp, _ = tp.Execute()
	log.Println(string(resp.Content()))
	// {
	// 	"args": {},
	// 	"data": "",
	// 	"files": {
	// 	  "file0": "learn.js\nfdsfsdavxlearnlearnlearnlearn"
	// 	},
	// 	"form": {},
	// 	"headers": {
	// 	  "Connection": "close",
	// 	  "Content-Length": "279",
	// 	  "Content-Type": "multipart/form-data; boundary=1b8ffe52a1241b6caa93af8d5d2c3b6172eb650224ad959c69ea8df7c04d",
	// 	  "Host": "httpbin.org",
	// 	  "User-Agent": "Go-http-client/1.1"
	// 	},
	// 	"json": null,
	// 	"method": "POST",
	// 	"origin": "172.17.0.1",
	// 	"url": "http://httpbin.org/anything"
    //   }
    
    ses = NewSession()
    tp = ses.Post("http://httpbin.org/post")
    ufile = NewUploadFile()
    ufile.SetFileName("MyFile")
    ufile.SetFieldName("MyField")
    ufile.SetFileFromPath("tests/json.file")
    tp.SetBodyAuto(ufile)
    resp, _ = tp.Execute()
```