# 性能测试

本目录包含requests库的性能测试和基准测试代码。

## 目录结构

- `benchmark_test.go` - 基准测试
- `concurrent_test.go` - 并发性能测试  
- `memory_usage.go` - 内存使用分析
- `load_test.go` - 负载测试

## 运行方式

### 基准测试
```bash
go test -bench=.
go test -bench=. -benchmem
```

### 并发测试
```bash
go run concurrent_test.go
```

### 内存分析
```bash
go run memory_usage.go
```

### 负载测试
```bash
go run load_test.go
```

## 性能指标

- 请求延迟
- 吞吐量
- 内存使用
- 连接池效率
- 中间件开销
