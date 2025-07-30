package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/474420502/requests"
)

// APIClient REST API客户端封装
type APIClient struct {
	baseURL string
	apiKey  string
	session *requests.Session
}

// NewAPIClient 创建新的API客户端
func NewAPIClient(baseURL, apiKey string) *APIClient {
	session := requests.NewSession()

	// 添加认证中间件
	authMiddleware := &AuthMiddleware{apiKey: apiKey}
	session.AddMiddleware(authMiddleware)

	// 添加日志中间件
	logger := log.New(os.Stdout, "[API-CLIENT] ", log.LstdFlags)
	logMiddleware := &requests.LoggingMiddleware{Logger: logger}
	session.AddMiddleware(logMiddleware)

	// 添加重试中间件
	retryMiddleware := &requests.RetryMiddleware{
		MaxRetries: 3,
		RetryDelay: time.Second,
	}
	session.AddMiddleware(retryMiddleware)

	return &APIClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		session: session,
	}
}

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	apiKey string
}

func (m *AuthMiddleware) BeforeRequest(req *http.Request) error {
	req.Header.Set("Authorization", "Bearer "+m.apiKey)
	req.Header.Set("Content-Type", "application/json")
	return nil
}

func (m *AuthMiddleware) AfterResponse(resp *http.Response) error {
	if resp.StatusCode == 401 {
		return fmt.Errorf("认证失败，请检查API密钥")
	}
	return nil
}

// User 用户数据结构
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

// APIResponse API响应结构
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Error   string      `json:"error,omitempty"`
}

// GetUser 获取用户信息
func (c *APIClient) GetUser(userID int) (*User, error) {
	url := fmt.Sprintf("%s/users/%d", c.baseURL, userID)

	resp, err := c.session.Get(url).Execute()
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}

	if resp.GetStatusCode() != 200 {
		return nil, fmt.Errorf("API错误: %d %s", resp.GetStatusCode(), resp.GetStatus())
	}

	var user User
	err = resp.DecodeJSON(&user)
	if err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &user, nil
}

// CreateUser 创建新用户
func (c *APIClient) CreateUser(req CreateUserRequest) (*User, error) {
	url := fmt.Sprintf("%s/users", c.baseURL)

	resp, err := c.session.Post(url).
		SetBodyJson(req).
		Execute()

	if err != nil {
		return nil, fmt.Errorf("创建用户请求失败: %v", err)
	}

	if resp.GetStatusCode() != 201 {
		return nil, fmt.Errorf("创建用户失败: %d %s", resp.GetStatusCode(), resp.GetStatus())
	}

	var user User
	err = resp.DecodeJSON(&user)
	if err != nil {
		return nil, fmt.Errorf("解析创建用户响应失败: %v", err)
	}

	return &user, nil
}

// UpdateUser 更新用户信息
func (c *APIClient) UpdateUser(userID int, updates map[string]interface{}) (*User, error) {
	url := fmt.Sprintf("%s/users/%d", c.baseURL, userID)

	resp, err := c.session.Put(url).
		SetBodyJson(updates).
		Execute()

	if err != nil {
		return nil, fmt.Errorf("更新用户请求失败: %v", err)
	}

	if resp.GetStatusCode() != 200 {
		return nil, fmt.Errorf("更新用户失败: %d %s", resp.GetStatusCode(), resp.GetStatus())
	}

	var user User
	err = resp.DecodeJSON(&user)
	if err != nil {
		return nil, fmt.Errorf("解析更新用户响应失败: %v", err)
	}

	return &user, nil
}

// DeleteUser 删除用户
func (c *APIClient) DeleteUser(userID int) error {
	url := fmt.Sprintf("%s/users/%d", c.baseURL, userID)

	resp, err := c.session.Delete(url).Execute()
	if err != nil {
		return fmt.Errorf("删除用户请求失败: %v", err)
	}

	if resp.GetStatusCode() != 204 && resp.GetStatusCode() != 200 {
		return fmt.Errorf("删除用户失败: %d %s", resp.GetStatusCode(), resp.GetStatus())
	}

	return nil
}

// ListUsers 获取用户列表
func (c *APIClient) ListUsers(page, limit int) ([]User, error) {
	url := fmt.Sprintf("%s/users", c.baseURL)

	resp, err := c.session.Get(url).
		AddParam("page", fmt.Sprintf("%d", page)).
		AddParam("limit", fmt.Sprintf("%d", limit)).
		Execute()

	if err != nil {
		return nil, fmt.Errorf("获取用户列表请求失败: %v", err)
	}

	if resp.GetStatusCode() != 200 {
		return nil, fmt.Errorf("获取用户列表失败: %d %s", resp.GetStatusCode(), resp.GetStatus())
	}

	var apiResp APIResponse
	err = resp.DecodeJSON(&apiResp)
	if err != nil {
		return nil, fmt.Errorf("解析用户列表响应失败: %v", err)
	}

	// 将interface{}转换为[]User
	dataBytes, err := json.Marshal(apiResp.Data)
	if err != nil {
		return nil, fmt.Errorf("序列化用户数据失败: %v", err)
	}

	var users []User
	err = json.Unmarshal(dataBytes, &users)
	if err != nil {
		return nil, fmt.Errorf("反序列化用户数据失败: %v", err)
	}

	return users, nil
}

// demonstrateAPIClient 演示API客户端的使用
func demonstrateAPIClient() {
	fmt.Println("=== REST API客户端演示 ===")

	// 创建API客户端
	// 注意：这里使用JSONPlaceholder作为演示API
	client := NewAPIClient("https://jsonplaceholder.typicode.com", "demo-api-key")

	// 1. 获取用户信息
	fmt.Println("1. 获取用户信息:")
	user, err := client.GetUser(1)
	if err != nil {
		fmt.Printf("✗ 获取用户失败: %v\n", err)
	} else {
		fmt.Printf("✓ 获取用户成功: %s (%s)\n", user.Name, user.Email)
	}

	// 2. 创建新用户
	fmt.Println("\n2. 创建新用户:")
	newUserReq := CreateUserRequest{
		Name:     "示例用户",
		Email:    "example@test.com",
		Username: "exampleuser",
	}

	newUser, err := client.CreateUser(newUserReq)
	if err != nil {
		fmt.Printf("✗ 创建用户失败: %v\n", err)
	} else {
		fmt.Printf("✓ 创建用户成功: ID=%d, Name=%s\n", newUser.ID, newUser.Name)
	}

	// 3. 更新用户信息
	fmt.Println("\n3. 更新用户信息:")
	updates := map[string]interface{}{
		"name":  "更新后的用户名",
		"email": "updated@test.com",
	}

	updatedUser, err := client.UpdateUser(1, updates)
	if err != nil {
		fmt.Printf("✗ 更新用户失败: %v\n", err)
	} else {
		fmt.Printf("✓ 更新用户成功: %s (%s)\n", updatedUser.Name, updatedUser.Email)
	}

	// 4. 获取用户列表
	fmt.Println("\n4. 获取用户列表:")
	users, err := client.ListUsers(1, 5)
	if err != nil {
		fmt.Printf("✗ 获取用户列表失败: %v\n", err)
	} else {
		fmt.Printf("✓ 获取用户列表成功，共 %d 个用户:\n", len(users))
		for i, u := range users {
			if i < 3 { // 只显示前3个
				fmt.Printf("   - %s (%s)\n", u.Name, u.Email)
			}
		}
		if len(users) > 3 {
			fmt.Printf("   ... 还有 %d 个用户\n", len(users)-3)
		}
	}

	// 5. 删除用户（演示，不会真正删除）
	fmt.Println("\n5. 删除用户:")
	err = client.DeleteUser(1)
	if err != nil {
		fmt.Printf("✗ 删除用户失败: %v\n", err)
	} else {
		fmt.Printf("✓ 删除用户成功\n")
	}

	fmt.Println("\n✅ API客户端演示完成")
	fmt.Println("主要特性:")
	fmt.Println("• 自动认证处理")
	fmt.Println("• 请求重试机制")
	fmt.Println("• 统一错误处理")
	fmt.Println("• JSON响应自动解析")
	fmt.Println("• 完整的CRUD操作支持")
}

func main() {
	demonstrateAPIClient()
}
