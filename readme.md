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

* Request Url With Change Query 
* eg3:
```go
	ses := requests.NewSession()
	tp := ses.Get("http://httpbin.org/get?page=1&name=xiaoming")
	p := tp.QueryParam("page") // get the param of page.
	p.IntAdd(1) // change page. => page += 1 
	resp, _ := tp.Execute()
	log.Println("\n", string(resp.Content()))
	// {
	//   "args": {
	//     "name": "xiaoming", 
	//     "page": "2"
	//   }, 
	//   "headers": {
	//     "Connection": "close", 
	//     "Host": "httpbin.org", 
	//     "User-Agent": "Go-http-client/1.1"
	//   }, 
	//   "origin": "172.17.0.1", 
	//   "url": "http://httpbin.org/get?name=xiaoming&page=2"
	// }

	p.StringSet("5") // Page String Set. equal to IntSet
	resp, _ := tp.Execute()
	log.Println("\n",string(resp.Content()))
	//{
	//  "args": {
	//    "name": "xiaoming", 
	//    "page": "5"
	//  }, 
	//  "headers": {
	//    "Connection": "close", 
	//    "Host": "httpbin.org", 
	//    "User-Agent": "Go-http-client/1.1"
	//  }, 
	//  "origin": "172.17.0.1", 
	//  "url": "http://httpbin.org/get?name=xiaoming&page=5"
	//}
``` 

* Request Url With Change the Path By regexp 
* eg4:
```go
	ses := requests.NewSession()
	surl := "http://httpbin.org/anything/Page-1-30/1028/1000"
	tp := ses.Get(surl)
	param := tp.PathParam(`.+Page-(\d+)-(\d+).+`)
	param.IntAdd(1) // equal to IntArraySet(0, 2)
	resp, _ := tp.Execute()
	log.Println("\n", string(resp.Content())) // Page-2-30
	// {
	//   "args": {}, 
	//   "data": "", 
	//   "files": {}, 
	//   "form": {}, 
	//   "headers": {
	//     "Connection": "close", 
	//     "Host": "httpbin.org", 
	//     "User-Agent": "Go-http-client/1.1"
	//   }, 
	//   "json": null, 
	//   "method": "GET", 
	//   "origin": "172.17.0.1", 
	//   "url": "http://httpbin.org/anything/Page-2-30/1028/1000"
	// }
``` 