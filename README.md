# Go Requests - HTTP Client Library

A powerful and easy-to-use Go HTTP client library designed for web scraping, API calls, and HTTP request handling. Features modern design with method chaining, middleware support, and type-safe configuration.

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.18-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## ✨ Key Features

- 🔗 **Method Chaining** - Intuitive API design with fluent method chaining
- 🛡️ **Type Safety** - Strong typing with compile-time error checking
- 🔧 **Functional Options** - Flexible Session configuration pattern
- 🚀 **Middleware Support** - Extensible request/response processing
- 🌐 **Comprehensive Proxy Support** - HTTP/HTTPS/SOCKS5 proxy support
- 🔐 **Authentication** - Basic Auth, Bearer Token, and more
- 📁 **File Upload** - Multi-file upload and form data support
- ⏱️ **Context Support** - Timeout control and request cancellation
- 🍪 **Cookie Management** - Automatic cookie handling and session management
- 🗜️ **Compression** - Gzip/Deflate automatic compression
- 🔄 **Retry Mechanism** - Built-in retry logic
- 📋 **Rich Configuration** - TLS, connection pool, buffer settings

## 📦 Installation

```bash
go get github.com/474420502/requests
```

## 🚀 Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/474420502/requests"
)

func main() {
    // Create a Session
    session := requests.NewSession()
    
    // Send GET request
    resp, err := session.Get("http://httpbin.org/get").Execute()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Status Code:", resp.GetStatusCode())
    fmt.Println("Response:", string(resp.Content()))
}
```

### Creating Session with Functional Options

```go
session, err := requests.NewSessionWithOptions(
    requests.WithTimeout(30*time.Second),               // Set timeout
    requests.WithUserAgent("MyApp/1.0"),                // Set user agent
    requests.WithProxy("http://proxy.example.com:8080"), // Set proxy
    requests.WithBasicAuth("username", "password"),      // Set basic auth
    requests.WithHeaders(map[string]string{              // Set default headers
        "Accept": "application/json",
        "X-Custom-Header": "custom-value",
    }),
    requests.WithKeepAlives(true),    // Enable keep-alive
    requests.WithCompression(true),   // Enable compression
)
if err != nil {
    log.Fatal("Failed to create session:", err)
}
```

## 📖 Detailed Usage Guide

### 1. HTTP Methods

```go
session := requests.NewSession()

// GET request
resp, _ := session.Get("https://api.example.com/users").Execute()

// POST request with JSON
resp, _ := session.Post("https://api.example.com/users").
    SetBodyJSON(map[string]string{"name": "John"}).
    Execute()

// PUT request - update data
resp, _ := session.Put("https://api.example.com/users/1").
    SetBodyJSON(map[string]string{"name": "Jane"}).
    Execute()

// DELETE request
resp, _ := session.Delete("https://api.example.com/users/1").Execute()

// Other methods: HEAD, PATCH, OPTIONS, CONNECT, TRACE
```

### 2. Request Body Settings

```go
// JSON request body
session.Post("https://api.example.com/data").
    SetBodyJSON(map[string]interface{}{
        "name": "John Doe",
        "age":  30,
        "tags": []string{"developer", "golang"},
    }).
    Execute()

// Form data
session.Post("https://api.example.com/login").
    SetBodyFormValues(map[string]string{
        "username": "testuser",
        "password": "secret123",
    }).
    Execute()

// Raw string
session.Post("https://api.example.com/data").
    SetBodyString("raw text data").
    Execute()

// Byte data
session.Post("https://api.example.com/upload").
    SetBodyBytes([]byte("binary data")).
    Execute()
```

### 3. Headers and Cookies

```go
// Set request headers
session.Get("https://api.example.com/data").
    SetHeader("Authorization", "Bearer your-token").
    SetHeader("Content-Type", "application/json").
    AddHeader("X-Custom", "value1").
    AddHeader("X-Custom", "value2"). // Add multiple values
    Execute()

// Set cookies
session.Get("https://api.example.com/data").
    SetCookieValue("session_id", "abc123").
    SetCookie(&http.Cookie{
        Name:  "user_pref",
        Value: "dark_mode",
    }).
    Execute()
```

### 4. Query Parameters

```go
// Add query parameters
resp, _ := session.Get("https://api.example.com/search").
    AddQuery("q", "golang").
    AddQuery("page", "1").
    AddQuery("limit", "10").
    Execute()

// Batch set query parameters
params := url.Values{}
params.Add("category", "tech")
params.Add("sort", "date")
session.Get("https://api.example.com/articles").
    SetQuery(params).
    Execute()
```

### 5. File Upload

```go
// Single file upload
resp, _ := session.Post("https://api.example.com/upload").
    SetBodyFile("./document.pdf", "file").
    Execute()

// Multiple file upload
resp, _ := session.Post("https://api.example.com/upload").
    SetBodyFiles(map[string]string{
        "document": "./document.pdf",
        "image":    "./photo.jpg",
    }).
    Execute()

// Form data + File upload
formData := map[string]string{
    "title":       "My Document",
    "description": "Important file",
}
resp, _ := session.Post("https://api.example.com/upload").
    SetBodyFormData("./document.pdf"). // Upload file
    SetBodyFormValues(formData).       // Form fields
    Execute()
```

### 6. Proxy Configuration

```go
// HTTP proxy
session.Config().SetProxyString("http://proxy.example.com:8080")

// HTTPS proxy
session.Config().SetProxyString("https://proxy.example.com:8080")

// SOCKS5 proxy
session.Config().SetProxyString("socks5://127.0.0.1:1080")

// Proxy with authentication
session.Config().SetProxyString("http://user:pass@proxy.example.com:8080")

// Clear proxy
session.Config().ClearProxy()
```

### 7. Authentication Configuration

```go
// Basic authentication
session.Config().SetBasicAuth("username", "password")

// Or use type-safe method
session.Config().SetBasicAuthString("username", "password")

// Using struct
auth := &requests.BasicAuth{
    User:     "username",
    Password: "password",
}
session.Config().SetBasicAuthStruct(auth)

// Bearer Token
session.SetHeader("Authorization", "Bearer your-jwt-token")

// Clear authentication
session.Config().ClearBasicAuth()
```

### 8. TLS and Security Configuration

```go
// Skip TLS certificate verification (for testing only)
session.Config().SetInsecure(true)

// Custom TLS configuration
tlsConfig := &tls.Config{
    InsecureSkipVerify: false,
    MinVersion:         tls.VersionTLS12,
}
session.Config().SetTLSConfig(tlsConfig)
```

### 9. Timeout and Context Control

```go
// Set timeout
session.Config().SetTimeoutDuration(30 * time.Second)
// Or by seconds
session.Config().SetTimeoutSeconds(30)

// Using Context
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

resp, err := session.Get("https://api.example.com/slow").
    WithContext(ctx).
    Execute()
```

### 10. Dynamic Parameter Handling

Go Requests supports dynamic URL parameter modification, especially useful for web scraping and API traversal:

```go
// Dynamic URL query parameter modification
session := requests.NewSession()
req := session.Get("http://api.example.com/search?page=1&category=tech")

// Get and modify page parameter
pageParam := req.QueryParam("page")
pageParam.IntAdd(1) // page becomes 2
resp, _ := req.Execute()

// String setting
pageParam.StringSet("5") // page becomes 5
resp, _ = req.Execute()

// Regex path parameter modification
url := "http://api.example.com/articles/page-1-20/item-100"
req = session.Get(url)
pathParam := req.PathParam(`page-(\d+)-(\d+)`)
pathParam.IntAdd(1) // becomes page-2-20
resp, _ = req.Execute()
```

### 11. Middleware System

```go
// Logging middleware
logger := log.New(os.Stdout, "[HTTP] ", log.LstdFlags)
loggingMiddleware := &requests.LoggingMiddleware{Logger: logger}

// Authentication middleware
authMiddleware := &requests.AuthMiddleware{
    TokenProvider: func() (string, error) {
        return "your-jwt-token", nil
    },
}

// Use middleware
resp, err := session.Get("https://api.example.com/data").
    WithMiddlewares(loggingMiddleware, authMiddleware).
    ExecuteWithMiddleware()
```

### 12. Response Handling

```go
resp, err := session.Get("https://api.example.com/users").Execute()
if err != nil {
    log.Fatal(err)
}

// Status code
fmt.Println("Status Code:", resp.GetStatusCode())
fmt.Println("Status:", resp.GetStatus())

// Response headers
headers := resp.GetHeader()
contentType := resp.GetHeader().Get("Content-Type")

// Response body
body := resp.Content() // []byte
text := string(resp.Content())

// JSON parsing
var users []User
err = resp.BindJSON(&users)
if err != nil {
    log.Fatal("JSON parsing failed:", err)
}

// Or use UnmarshalJSON
err = resp.UnmarshalJSON(&users)
```

## 🏭 Pre-configured Sessions

The library provides several pre-configured Sessions for different use cases:

```go
// API-specific Session
session, _ := requests.NewSessionForAPI()

// Web scraping Session
session, _ := requests.NewSessionForScraping()

// Testing Session
session, _ := requests.NewSessionForTesting()

// High-performance Session
session, _ := requests.NewHighPerformanceSession()

// Security-enhanced Session
session, _ := requests.NewSecureSession()

// Session with retry capability
session, _ := requests.NewSessionWithRetry(3, time.Second*2)

// Session with proxy
session, _ := requests.NewSessionWithProxy("http://proxy.example.com:8080")
```

## ⚙️ Advanced Configuration

### Connection Pool and Performance Optimization

```go
session, _ := requests.NewSessionWithOptions(
    // Connection pool configuration
    requests.WithMaxIdleConns(100),           // Max idle connections
    requests.WithMaxIdleConnsPerHost(10),     // Max idle connections per host
    requests.WithMaxConnsPerHost(50),         // Max connections per host
    
    // Buffer configuration
    requests.WithReadBufferSize(64 * 1024),   // Read buffer size
    requests.WithWriteBufferSize(64 * 1024),  // Write buffer size
    
    // Connection timeout
    requests.WithDialTimeout(10 * time.Second), // Dial timeout
    
    // Enable Keep-Alive and compression
    requests.WithKeepAlives(true),    // Enable keep-alive
    requests.WithCompression(true),   // Enable compression
)
```

### Cookie Management

```go
// Enable Cookie jar
session.Config().SetWithCookiejar(true)

// Set cookies
u, _ := url.Parse("https://example.com")
cookies := []*http.Cookie{
    {Name: "session_id", Value: "abc123"},
    {Name: "user_pref", Value: "dark_mode"},
}
session.SetCookies(u, cookies)

// Get cookies
cookies = session.GetCookies(u)

// Delete specific cookie
session.DelCookies(u, "session_id")

// Clear all cookies
session.ClearCookies()
```

### Compression Configuration

```go
// Add accepted compression types
session.Config().AddAcceptEncoding(requests.AcceptEncodingGzip)
session.Config().AddAcceptEncoding(requests.AcceptEncodingDeflate)

// Set compression type for sending data
session.Config().SetContentEncoding(requests.ContentEncodingGzip)

// Set decompression behavior when no Accept-Encoding header
session.Config().SetDecompressNoAccept(true)
```

## 🔍 Error Handling

```go
resp, err := session.Get("https://api.example.com/data").Execute()
if err != nil {
    // Network errors, timeouts, etc.
    log.Printf("Request failed: %v", err)
    return
}

// Check HTTP status code
if resp.GetStatusCode() >= 400 {
    log.Printf("HTTP error: %d %s", resp.GetStatusCode(), resp.GetStatus())
    return
}

// Handle specific status codes
switch resp.GetStatusCode() {
case 200:
    // Success
case 401:
    // Unauthorized
case 404:
    // Not found
case 500:
    // Server error
}
```

## 🧪 Testing Examples

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

## 💡 Use Case Examples

### Web Scraping

```go
// Scraping-specific configuration
session, _ := requests.NewSessionForScraping()

// Simulate browser behavior
resp, _ := session.Get("https://example.com/page/1").
    SetHeader("Referer", "https://example.com").
    SetCookieValue("visited", "true").
    Execute()

// Handle pagination
for page := 1; page <= 10; page++ {
    url := fmt.Sprintf("https://example.com/page/%d", page)
    resp, err := session.Get(url).Execute()
    if err != nil {
        log.Printf("Failed to scrape page %d: %v", page, err)
        continue
    }
    
    // Parse page content
    parseHTML(resp.Content())
    
    // Add delay to avoid being blocked
    time.Sleep(time.Second)
}
```

### API Client

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
        return fmt.Errorf("failed to create user: %s", resp.GetStatus())
    }
    
    return nil
}
```

## 📈 Performance Recommendations

1. **Reuse Sessions**: Reuse Session instances for the same domain or API service
2. **Set Appropriate Timeouts**: Configure suitable timeout values based on network conditions
3. **Use Connection Pools**: Configure appropriate connection pool sizes for concurrent performance
4. **Enable Compression**: Enable Gzip compression for large data transfers
5. **Use Context**: Use Context for potentially long-running requests
6. **Avoid Over-concurrency**: Make moderate concurrent requests to avoid rate limiting
7. **Use Keep-Alive Wisely**: Enable Keep-Alive for scenarios requiring multiple requests to the same server

## 🛠️ Troubleshooting

### Common Issues

**Q: Request timeout?**
```go
// Increase timeout
session.Config().SetTimeoutDuration(60 * time.Second)

// Or use Context for precise control
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
resp, err := session.Get(url).WithContext(ctx).Execute()
```

**Q: HTTPS certificate errors?**
```go
// For testing environments only
session.Config().SetInsecure(true)

// For production, configure proper certificates
tlsConfig := &tls.Config{
    InsecureSkipVerify: false,
    // Add custom certificates, etc.
}
session.Config().SetTLSConfig(tlsConfig)
```

**Q: Proxy connection failed?**
```go
// Check proxy URL format
err := session.Config().SetProxyString("http://proxy.example.com:8080")
if err != nil {
    log.Printf("Proxy configuration error: %v", err)
}

// For authenticated proxy
session.Config().SetProxyString("http://username:password@proxy.example.com:8080")
```

## 🤝 Contributing

We welcome all forms of contributions!

### How to Contribute

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/new-feature`)
3. Commit your changes (`git commit -am 'Add new feature'`)
4. Push to the branch (`git push origin feature/new-feature`)
5. Create a Pull Request

### Code Standards

- Follow official Go coding standards
- Add comprehensive test cases
- Update relevant documentation
- Ensure all tests pass

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## 🔗 Related Links

- [GitHub Repository](https://github.com/474420502/requests)
- [Go Package Documentation](https://pkg.go.dev/github.com/474420502/requests)
- [Example Code](./example)
- [Chinese Documentation](./README_zh.md)

## ⭐ Support the Project

If this library helps you, please give it a ⭐️! Your support motivates us to keep improving.

---

**Made with ❤️ by Go developers, for Go developers**
