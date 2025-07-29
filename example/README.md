# Requests 库示例目录

本目录包含了requests库的各种使用示例，展示了从基础用法到高级功能的完整演示。

## 📁 目录结构

```
example/
├── README.md                    # 本文件
├── 01-basic/                    # 基础用法示例
│   ├── simple_requests.go       # 简单的HTTP请求
│   ├── session_management.go    # Session管理
│   └── main.go                  # 运行入口
├── 02-advanced/                 # 高级功能示例
│   ├── middleware_usage.go      # 中间件使用
│   ├── form_upload.go           # 表单和文件上传
│   ├── json_handling.go         # JSON处理
│   └── main.go                  # 运行入口
├── 03-phase-demos/              # 重构阶段演示
│   ├── phase1_unified_api.go    # 第一阶段：API统一
│   ├── phase2_type_safety.go    # 第二阶段：类型安全
│   ├── phase3_architecture.go   # 第三阶段：架构完善
│   └── main.go                  # 运行入口
├── 04-performance/              # 性能优化示例
│   ├── concurrent_requests.go   # 并发请求
│   ├── connection_pooling.go    # 连接池配置
│   └── main.go                  # 运行入口
├── 05-real-world/               # 真实场景示例
│   ├── api_client.go            # API客户端
│   ├── web_scraper.go           # 网页爬虫
│   ├── file_downloader.go       # 文件下载器
│   └── main.go                  # 运行入口
└── tests/                       # 示例测试
    └── example_test.go          # 示例代码测试
```

## 🚀 快速开始

```bash
# 运行基础示例
cd 01-basic && go run .

# 运行高级功能示例
cd 02-advanced && go run .

# 运行重构演示
cd 03-phase-demos && go run .

# 运行所有测试
cd tests && go test -v
```

## 📖 示例说明

### 01-basic 基础用法
- **simple_requests.go**: 展示GET、POST等基本HTTP请求
- **session_management.go**: 展示Session的创建和配置

### 02-advanced 高级功能
- **middleware_usage.go**: 展示各种中间件的使用
- **form_upload.go**: 展示表单提交和文件上传
- **json_handling.go**: 展示JSON数据的处理

### 03-phase-demos 重构演示
- **phase1_unified_api.go**: 展示第一阶段的API统一改进
- **phase2_type_safety.go**: 展示第二阶段的类型安全改进  
- **phase3_architecture.go**: 展示第三阶段的架构完善

### 04-performance 性能优化
- **concurrent_requests.go**: 展示如何进行高效的并发请求
- **connection_pooling.go**: 展示连接池的最佳实践

### 05-real-world 真实场景
- **api_client.go**: 构建生产级API客户端
- **web_scraper.go**: 构建网页爬虫
- **file_downloader.go**: 实现文件下载功能

## 🔧 开发指南

每个示例都是独立的，可以单独运行。建议按照编号顺序学习：

1. 先学习基础用法（01-basic）
2. 再学习高级功能（02-advanced）
3. 了解重构历程（03-phase-demos）
4. 掌握性能优化（04-performance）
5. 应用到实际项目（05-real-world）

## 🧪 测试

所有示例都有对应的测试，确保代码质量和正确性。运行测试：

```bash
cd tests
go test -v ./...
```

## 📝 贡献

欢迎为示例目录贡献新的示例或改进现有示例！
