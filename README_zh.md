# Go Requests - Go HTTP客户端库

一个功能强大、易于使用的Go HTTP客户端库，专为爬虫、API调用和HTTP请求处理设计。具有链式调用、中间件支持、类型安全配置等现代化特性。

[![Go版本](https://img.shields.io/badge/Go-%3E%3D%201.18-blue.svg)](https://golang.org)
[![许可证](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## ✨ 核心特性

- 🔗 **链式调用** - 直观的API设计，支持方法链式调用
- 🛡️ **类型安全** - 强类型配置，编译时错误检查
- 🔧 **函数式选项** - 灵活的Session配置模式
- 🚀 **中间件支持** - 可扩展的请求/响应处理机制
- 🌐 **全面代理支持** - HTTP/HTTPS/SOCKS5代理
- 🔐 **认证机制** - 基础认证、Bearer Token等
- 📁 **文件上传** - 支持多文件上传和表单数据
- ⏱️ **Context支持** - 超时控制和请求取消
- 🍪 **Cookie管理** - 自动Cookie处理和会话管理
- 🗜️ **压缩支持** - Gzip/Deflate自动压缩
- 🔄 **重试机制** - 内置重试逻辑
- 📋 **丰富配置** - TLS、连接池、缓冲区等详细配置

## 📦 安装

```bash
go get github.com/474420502/requests
```

## 🚀 快速开始

### 基础用法

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/474420502/requests"
)

func main() {
    // 创建Session
    session := requests.NewSession()
    
    // 发送GET请求
    resp, err := session.Get("http://httpbin.org/get").Execute()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("状态码:", resp.GetStatusCode())
    fmt.Println("响应内容:", string(resp.Content()))
}
```

### 使用函数式选项创建Session

```go
session, err := requests.NewSessionWithOptions(
    requests.WithTimeout(30*time.Second),               // 设置超时时间
    requests.WithUserAgent("我的应用/1.0"),               // 设置用户代理
    requests.WithProxy("http://proxy.example.com:8080"), // 设置代理
    requests.WithBasicAuth("用户名", "密码"),               // 设置基础认证
    requests.WithHeaders(map[string]string{              // 设置默认请求头
        "Accept": "application/json",
        "自定义头部": "自定义值",
    }),
    requests.WithKeepAlives(true),    // 启用连接保持
    requests.WithCompression(true),   // 启用压缩
)
if err != nil {
    log.Fatal("创建Session失败:", err)
}
```

## 📖 详细使用指南

### 1. HTTP请求方法

```go
session := requests.NewSession()

// GET请求
resp, _ := session.Get("https://api.example.com/users").Execute()

// POST请求 - 发送JSON数据
resp, _ := session.Post("https://api.example.com/users").
    SetBodyJSON(map[string]string{"name": "张三"}).
    Execute()

// PUT请求 - 更新数据
resp, _ := session.Put("https://api.example.com/users/1").
    SetBodyJSON(map[string]string{"name": "李四"}).
    Execute()

// DELETE请求 - 删除数据
resp, _ := session.Delete("https://api.example.com/users/1").Execute()

// 其他方法：HEAD, PATCH, OPTIONS, CONNECT, TRACE
```

### 2. 请求体设置

```go
// JSON请求体
session.Post("https://api.example.com/data").
    SetBodyJSON(map[string]interface{}{
        "姓名": "张三",
        "年龄": 25,
        "标签": []string{"开发者", "golang"},
    }).
    Execute()

// 表单数据
session.Post("https://api.example.com/login").
    SetBodyFormValues(map[string]string{
        "username": "测试用户",
        "password": "密码123",
    }).
    Execute()

// 原始字符串
session.Post("https://api.example.com/data").
    SetBodyString("原始文本数据").
    Execute()

// 字节数据
session.Post("https://api.example.com/upload").
    SetBodyBytes([]byte("二进制数据")).
    Execute()
```

### 3. 请求头和Cookie设置

```go
// 设置请求头
session.Get("https://api.example.com/data").
    SetHeader("Authorization", "Bearer 你的令牌").
    SetHeader("Content-Type", "application/json").
    AddHeader("X-自定义", "值1").
    AddHeader("X-自定义", "值2"). // 添加多个值
    Execute()

// 设置Cookie
session.Get("https://api.example.com/data").
    SetCookieValue("session_id", "abc123").
    SetCookie(&http.Cookie{
        Name:  "用户偏好",
        Value: "深色模式",
    }).
    Execute()
```

### 4. 查询参数

```go
// 添加查询参数
resp, _ := session.Get("https://api.example.com/search").
    AddQuery("q", "golang").
    AddQuery("页码", "1").
    AddQuery("每页数量", "10").
    Execute()

// 批量设置查询参数
params := url.Values{}
params.Add("分类", "技术")
params.Add("排序", "时间")
session.Get("https://api.example.com/articles").
    SetQuery(params).
    Execute()
```

### 5. 文件上传

```go
// 单文件上传
resp, _ := session.Post("https://api.example.com/upload").
    SetBodyFile("./文档.pdf", "file").
    Execute()

// 多文件上传
resp, _ := session.Post("https://api.example.com/upload").
    SetBodyFiles(map[string]string{
        "文档": "./文档.pdf",
        "图片": "./照片.jpg",
    }).
    Execute()

// 表单数据 + 文件上传
formData := map[string]string{
    "标题":  "我的文档",
    "描述":  "重要文件",
}
resp, _ := session.Post("https://api.example.com/upload").
    SetBodyFormData("./文档.pdf").      // 上传文件
    SetBodyFormValues(formData).       // 表单字段
    Execute()
```

### 6. 代理配置

```go
// HTTP代理
session.Config().SetProxyString("http://proxy.example.com:8080")

// HTTPS代理
session.Config().SetProxyString("https://proxy.example.com:8080")

// SOCKS5代理
session.Config().SetProxyString("socks5://127.0.0.1:1080")

// 带认证的代理
session.Config().SetProxyString("http://用户名:密码@proxy.example.com:8080")

// 清除代理设置
session.Config().ClearProxy()
```

### 7. 身份认证配置

```go
// 基础认证
session.Config().SetBasicAuth("用户名", "密码")

// 或者使用类型安全的方法
session.Config().SetBasicAuthString("用户名", "密码")

// 使用认证结构体
auth := &requests.BasicAuth{
    User:     "用户名",
    Password: "密码",
}
session.Config().SetBasicAuthStruct(auth)

// Bearer Token认证
session.SetHeader("Authorization", "Bearer 你的JWT令牌")

// 清除认证设置
session.Config().ClearBasicAuth()
```

### 8. TLS和安全配置

```go
// 跳过TLS证书验证（仅用于测试环境）
session.Config().SetInsecure(true)

// 自定义TLS配置
tlsConfig := &tls.Config{
    InsecureSkipVerify: false,
    MinVersion:         tls.VersionTLS12,
}
session.Config().SetTLSConfig(tlsConfig)
```

### 9. 超时和Context控制

```go
// 设置超时时间
session.Config().SetTimeoutDuration(30 * time.Second)
// 或者按秒设置
session.Config().SetTimeoutSeconds(30)

// 使用Context控制
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

resp, err := session.Get("https://api.example.com/slow").
    WithContext(ctx).
    Execute()
```

### 10. 动态参数处理

Go Requests 支持动态修改URL参数，特别适用于爬虫和API遍历：

```go
// URL查询参数动态修改
session := requests.NewSession()
req := session.Get("http://api.example.com/search?page=1&category=tech")

// 获取并修改page参数
pageParam := req.QueryParam("page")
pageParam.IntAdd(1) // page变为2
resp, _ := req.Execute()

// 字符串方式设置
pageParam.StringSet("5") // page变为5
resp, _ = req.Execute()

// 正则表达式路径参数修改
url := "http://api.example.com/articles/page-1-20/item-100"
req = session.Get(url)
pathParam := req.PathParam(`page-(\d+)-(\d+)`)
pathParam.IntAdd(1) // 变为page-2-20
resp, _ = req.Execute()
```

### 11. 中间件系统

```go
// 日志中间件
logger := log.New(os.Stdout, "[HTTP] ", log.LstdFlags)
loggingMiddleware := &requests.LoggingMiddleware{Logger: logger}

// 认证中间件
authMiddleware := &requests.AuthMiddleware{
    TokenProvider: func() (string, error) {
        return "你的JWT令牌", nil
    },
}

// 使用中间件
resp, err := session.Get("https://api.example.com/data").
    WithMiddlewares(loggingMiddleware, authMiddleware).
    ExecuteWithMiddleware()
```

### 12. 响应处理

```go
resp, err := session.Get("https://api.example.com/users").Execute()
if err != nil {
    log.Fatal(err)
}

// 状态码
fmt.Println("状态码:", resp.GetStatusCode())
fmt.Println("状态信息:", resp.GetStatus())

// 响应头
headers := resp.GetHeader()
contentType := resp.GetHeader().Get("Content-Type")

// 响应体
body := resp.Content() // []byte
text := string(resp.Content())

// JSON解析
var users []User
err = resp.BindJSON(&users)
if err != nil {
    log.Fatal("JSON解析失败:", err)
}

// 或使用UnmarshalJSON
err = resp.UnmarshalJSON(&users)
```

## 🏭 预定义Session配置

库提供了几种预配置的Session，适用于不同使用场景：

```go
// API调用专用Session
session, _ := requests.NewSessionForAPI()

// 网页爬虫专用Session
session, _ := requests.NewSessionForScraping()

// 测试专用Session
session, _ := requests.NewSessionForTesting()

// 高性能Session
session, _ := requests.NewHighPerformanceSession()

// 安全增强Session
session, _ := requests.NewSecureSession()

// 带重试功能的Session
session, _ := requests.NewSessionWithRetry(3, time.Second*2)

// 带代理的Session
session, _ := requests.NewSessionWithProxy("http://proxy.example.com:8080")
```

## ⚙️ 高级配置

### 连接池和性能优化

```go
session, _ := requests.NewSessionWithOptions(
    // 连接池配置
    requests.WithMaxIdleConns(100),           // 最大空闲连接数
    requests.WithMaxIdleConnsPerHost(10),     // 每个主机最大空闲连接数
    requests.WithMaxConnsPerHost(50),         // 每个主机最大连接数
    
    // 缓冲区配置
    requests.WithReadBufferSize(64 * 1024),   // 读缓冲区大小
    requests.WithWriteBufferSize(64 * 1024),  // 写缓冲区大小
    
    // 连接超时
    requests.WithDialTimeout(10 * time.Second), // 连接超时时间
    
    // 启用Keep-Alive和压缩
    requests.WithKeepAlives(true),    // 启用连接保持
    requests.WithCompression(true),   // 启用压缩
)
```

### Cookie管理

```go
// 启用Cookie jar
session.Config().SetWithCookiejar(true)

// 设置Cookie
u, _ := url.Parse("https://example.com")
cookies := []*http.Cookie{
    {Name: "session_id", Value: "abc123"},
    {Name: "user_pref", Value: "dark_mode"},
}
session.SetCookies(u, cookies)

// 获取Cookie
cookies = session.GetCookies(u)

// 删除特定Cookie
session.DelCookies(u, "session_id")

// 清除所有Cookie
session.ClearCookies()
```

### 压缩配置

```go
// 添加接受的压缩类型
session.Config().AddAcceptEncoding(requests.AcceptEncodingGzip)
session.Config().AddAcceptEncoding(requests.AcceptEncodingDeflate)

// 设置发送数据的压缩类型
session.Config().SetContentEncoding(requests.ContentEncodingGzip)

// 设置无Accept-Encoding头时的解压行为
session.Config().SetDecompressNoAccept(true)
```

## 🔍 错误处理

```go
resp, err := session.Get("https://api.example.com/data").Execute()
if err != nil {
    // 网络错误、超时等
    log.Printf("请求失败: %v", err)
    return
}

// 检查HTTP状态码
if resp.GetStatusCode() >= 400 {
    log.Printf("HTTP错误: %d %s", resp.GetStatusCode(), resp.GetStatus())
    return
}

// 处理特定状态码
switch resp.GetStatusCode() {
case 200:
    // 请求成功
case 401:
    // 未授权，需要登录
case 404:
    // 资源未找到
case 500:
    // 服务器内部错误
}
```

## 🧪 测试示例

```go
func TestAPICall(t *testing.T) {
    session := requests.NewSession()
    
    resp, err := session.Get("https://httpbin.org/json").Execute()
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.GetStatusCode())
    
    var data map[string]interface{}
    err = resp.BindJSON(&data)
    assert.NoError(t, err)
    assert.NotEmpty(t, data)
}
```

## 💡 使用场景示例

### Web爬虫

```go
// 爬虫专用配置
session, _ := requests.NewSessionForScraping()

// 模拟浏览器行为
resp, _ := session.Get("https://example.com/page/1").
    SetHeader("Referer", "https://example.com").
    SetCookieValue("visited", "true").
    Execute()

// 处理分页
for page := 1; page <= 10; page++ {
    url := fmt.Sprintf("https://example.com/page/%d", page)
    resp, err := session.Get(url).Execute()
    if err != nil {
        log.Printf("爬取第%d页失败: %v", page, err)
        continue
    }
    
    // 解析页面内容
    parseHTML(resp.Content())
    
    // 添加延迟避免被封
    time.Sleep(time.Second)
}
```

### API客户端

```go
type APIClient struct {
    session *requests.Session
    baseURL string
}

func NewAPIClient(token string) *APIClient {
    session, _ := requests.NewSessionWithOptions(
        requests.WithTimeout(30*time.Second),
        requests.WithHeaders(map[string]string{
            "Authorization": "Bearer " + token,
            "Content-Type":  "application/json",
        }),
    )
    
    return &APIClient{
        session: session,
        baseURL: "https://api.example.com",
    }
}

func (c *APIClient) GetUser(userID int) (*User, error) {
    url := fmt.Sprintf("%s/users/%d", c.baseURL, userID)
    resp, err := c.session.Get(url).Execute()
    if err != nil {
        return nil, err
    }
    
    var user User
    err = resp.BindJSON(&user)
    return &user, err
}

func (c *APIClient) CreateUser(user *User) error {
    url := c.baseURL + "/users"
    resp, err := c.session.Post(url).
        SetBodyJSON(user).
        Execute()
    if err != nil {
        return err
    }
    
    if resp.GetStatusCode() != 201 {
        return fmt.Errorf("创建用户失败: %s", resp.GetStatus())
    }
    
    return nil
}
```

## 📈 性能优化建议

1. **复用Session**: 为同一个域名或API服务复用Session实例
2. **合理设置超时**: 根据实际网络环境设置合适的超时时间
3. **使用连接池**: 配置合适的连接池大小以提高并发性能
4. **启用压缩**: 对于大量数据传输启用Gzip压缩
5. **使用Context**: 对于可能长时间运行的请求使用Context进行控制
6. **避免过度并发**: 对目标服务器进行适度的并发请求，避免被限流
7. **合理使用Keep-Alive**: 对于需要多次请求同一服务器的场景启用Keep-Alive

## 🛠️ 故障排除

### 常见问题

**Q: 请求超时怎么办？**
```go
// 增加超时时间
session.Config().SetTimeoutDuration(60 * time.Second)

// 或使用Context精确控制
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
resp, err := session.Get(url).WithContext(ctx).Execute()
```

**Q: 如何处理HTTPS证书错误？**
```go
// 仅在测试环境使用
session.Config().SetInsecure(true)

// 生产环境应该配置正确的证书
tlsConfig := &tls.Config{
    InsecureSkipVerify: false,
    // 添加自定义证书等
}
session.Config().SetTLSConfig(tlsConfig)
```

**Q: 代理连接失败？**
```go
// 检查代理URL格式
err := session.Config().SetProxyString("http://proxy.example.com:8080")
if err != nil {
    log.Printf("代理配置错误: %v", err)
}

// 对于需要认证的代理
session.Config().SetProxyString("http://username:password@proxy.example.com:8080")
```

## 🤝 贡献指南

我们欢迎所有形式的贡献！

### 如何贡献

1. Fork 本仓库
2. 创建功能分支 (`git checkout -b feature/新功能`)
3. 提交更改 (`git commit -am '添加新功能'`)
4. 推送到分支 (`git push origin feature/新功能`)
5. 创建 Pull Request

### 代码规范

- 遵循Go语言官方代码规范
- 添加充分的测试用例
- 更新相关文档
- 确保所有测试通过

## 📄 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

## 🔗 相关链接

- [GitHub 仓库](https://github.com/474420502/requests)
- [Go Package 文档](https://pkg.go.dev/github.com/474420502/requests)
- [示例代码](./example)
- [更新日志](CHANGELOG.md)

## ⭐ 支持项目

如果这个库对你有帮助，请给个⭐️！你的支持是我们持续改进的动力。

---

**Made with ❤️ by Go开发者，为Go开发者服务**
