package requests

import (
	"net"
	"net/http"
	"runtime"
	"sync"
	"time"
)

// ConnectionPoolConfig 连接池配置
type ConnectionPoolConfig struct {
	// 基础连接池配置
	MaxIdleConns        int           // 全局最大空闲连接数
	MaxIdleConnsPerHost int           // 每个主机最大空闲连接数
	MaxConnsPerHost     int           // 每个主机最大连接数
	IdleConnTimeout     time.Duration // 空闲连接超时时间

	// 高级配置
	DialTimeout         time.Duration // 连接建立超时
	KeepAliveTimeout    time.Duration // Keep-Alive 超时
	TLSHandshakeTimeout time.Duration // TLS 握手超时

	// 自适应配置
	EnableAdaptive     bool    // 启用自适应调整
	MinConnsPerHost    int     // 每个主机最小连接数
	ScaleUpThreshold   float64 // 扩容阈值（连接使用率）
	ScaleDownThreshold float64 // 缩容阈值（连接使用率）
}

// DefaultConnectionPoolConfig 返回默认的连接池配置
func DefaultConnectionPoolConfig() *ConnectionPoolConfig {
	return &ConnectionPoolConfig{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     0, // 0表示无限制
		IdleConnTimeout:     90 * time.Second,

		DialTimeout:         30 * time.Second,
		KeepAliveTimeout:    30 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,

		EnableAdaptive:     false,
		MinConnsPerHost:    2,
		ScaleUpThreshold:   0.8,
		ScaleDownThreshold: 0.3,
	}
}

// HighPerformanceConnectionPoolConfig 返回高性能连接池配置
func HighPerformanceConnectionPoolConfig() *ConnectionPoolConfig {
	config := DefaultConnectionPoolConfig()

	// 根据CPU核心数调整连接池大小
	cpuCount := runtime.NumCPU()

	config.MaxIdleConns = cpuCount * 20
	config.MaxIdleConnsPerHost = cpuCount * 4
	config.MaxConnsPerHost = cpuCount * 10
	config.IdleConnTimeout = 120 * time.Second

	// 更激进的超时设置
	config.DialTimeout = 15 * time.Second
	config.KeepAliveTimeout = 60 * time.Second
	config.TLSHandshakeTimeout = 5 * time.Second

	// 启用自适应调整
	config.EnableAdaptive = true
	config.MinConnsPerHost = 2
	config.ScaleUpThreshold = 0.7
	config.ScaleDownThreshold = 0.2

	return config
}

// LowLatencyConnectionPoolConfig 返回低延迟连接池配置
func LowLatencyConnectionPoolConfig() *ConnectionPoolConfig {
	config := DefaultConnectionPoolConfig()

	// 优化延迟，增加连接数
	config.MaxIdleConns = 200
	config.MaxIdleConnsPerHost = 20
	config.MaxConnsPerHost = 50
	config.IdleConnTimeout = 60 * time.Second

	// 更短的超时时间
	config.DialTimeout = 10 * time.Second
	config.KeepAliveTimeout = 15 * time.Second
	config.TLSHandshakeTimeout = 3 * time.Second

	return config
}

// ResourceConstrainedConnectionPoolConfig 返回资源受限环境的连接池配置
func ResourceConstrainedConnectionPoolConfig() *ConnectionPoolConfig {
	config := DefaultConnectionPoolConfig()

	// 减少资源占用
	config.MaxIdleConns = 20
	config.MaxIdleConnsPerHost = 2
	config.MaxConnsPerHost = 5
	config.IdleConnTimeout = 30 * time.Second

	// 适中的超时设置
	config.DialTimeout = 20 * time.Second
	config.KeepAliveTimeout = 15 * time.Second
	config.TLSHandshakeTimeout = 8 * time.Second

	return config
}

// AdaptiveConnectionPool 自适应连接池管理器
type AdaptiveConnectionPool struct {
	config    *ConnectionPoolConfig
	transport *http.Transport
	stats     *ConnectionStats
	mutex     sync.RWMutex
	stopCh    chan struct{}
	running   bool
}

// ConnectionStats 连接统计信息
type ConnectionStats struct {
	ActiveConnections int64
	IdleConnections   int64
	TotalRequests     int64
	FailedConnections int64
	AverageConnTime   time.Duration
	LastAdjustment    time.Time
	mutex             sync.RWMutex
}

// NewAdaptiveConnectionPool 创建自适应连接池
func NewAdaptiveConnectionPool(config *ConnectionPoolConfig) *AdaptiveConnectionPool {
	if config == nil {
		config = DefaultConnectionPoolConfig()
	}

	pool := &AdaptiveConnectionPool{
		config:    config,
		transport: createOptimizedTransport(config),
		stats:     &ConnectionStats{},
		stopCh:    make(chan struct{}),
	}

	if config.EnableAdaptive {
		pool.startAdaptiveMonitoring()
	}

	return pool
}

// createOptimizedTransport 创建优化的HTTP传输层
func createOptimizedTransport(config *ConnectionPoolConfig) *http.Transport {
	dialer := &net.Dialer{
		Timeout:   config.DialTimeout,
		KeepAlive: config.KeepAliveTimeout,
		DualStack: true, // 支持IPv4和IPv6
	}

	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true, // 强制尝试HTTP/2
		MaxIdleConns:          config.MaxIdleConns,
		MaxIdleConnsPerHost:   config.MaxIdleConnsPerHost,
		MaxConnsPerHost:       config.MaxConnsPerHost,
		IdleConnTimeout:       config.IdleConnTimeout,
		TLSHandshakeTimeout:   config.TLSHandshakeTimeout,
		ExpectContinueTimeout: 1 * time.Second,

		// 优化TCP设置
		DisableKeepAlives:      false,
		DisableCompression:     false, // 启用压缩以减少传输数据
		MaxResponseHeaderBytes: 4096,  // 限制响应头大小
	}

	return transport
}

// GetTransport 获取HTTP传输层
func (p *AdaptiveConnectionPool) GetTransport() *http.Transport {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.transport
}

// GetStats 获取连接统计信息
func (p *AdaptiveConnectionPool) GetStats() *ConnectionStats {
	p.stats.mutex.RLock()
	defer p.stats.mutex.RUnlock()

	// 返回统计信息的副本
	return &ConnectionStats{
		ActiveConnections: p.stats.ActiveConnections,
		IdleConnections:   p.stats.IdleConnections,
		TotalRequests:     p.stats.TotalRequests,
		FailedConnections: p.stats.FailedConnections,
		AverageConnTime:   p.stats.AverageConnTime,
		LastAdjustment:    p.stats.LastAdjustment,
	}
}

// startAdaptiveMonitoring 启动自适应监控
func (p *AdaptiveConnectionPool) startAdaptiveMonitoring() {
	if p.running {
		return
	}

	p.running = true
	go func() {
		ticker := time.NewTicker(30 * time.Second) // 每30秒检查一次
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				p.adjustConnectionPool()
			case <-p.stopCh:
				return
			}
		}
	}()
}

// adjustConnectionPool 调整连接池参数
func (p *AdaptiveConnectionPool) adjustConnectionPool() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	stats := p.GetStats()

	// 计算连接使用率
	if stats.ActiveConnections > 0 {
		utilizationRate := float64(stats.ActiveConnections) /
			float64(p.config.MaxIdleConnsPerHost)

		// 需要扩容
		if utilizationRate > p.config.ScaleUpThreshold {
			newMax := int(float64(p.config.MaxIdleConnsPerHost) * 1.5)
			if newMax <= p.config.MaxConnsPerHost || p.config.MaxConnsPerHost == 0 {
				p.config.MaxIdleConnsPerHost = newMax
				p.updateTransportSettings()

				p.stats.mutex.Lock()
				p.stats.LastAdjustment = time.Now()
				p.stats.mutex.Unlock()
			}
		}

		// 需要缩容
		if utilizationRate < p.config.ScaleDownThreshold {
			newMax := int(float64(p.config.MaxIdleConnsPerHost) * 0.8)
			if newMax >= p.config.MinConnsPerHost {
				p.config.MaxIdleConnsPerHost = newMax
				p.updateTransportSettings()

				p.stats.mutex.Lock()
				p.stats.LastAdjustment = time.Now()
				p.stats.mutex.Unlock()
			}
		}
	}
}

// updateTransportSettings 更新传输层设置
func (p *AdaptiveConnectionPool) updateTransportSettings() {
	p.transport.MaxIdleConns = p.config.MaxIdleConns
	p.transport.MaxIdleConnsPerHost = p.config.MaxIdleConnsPerHost
	p.transport.MaxConnsPerHost = p.config.MaxConnsPerHost
	p.transport.IdleConnTimeout = p.config.IdleConnTimeout
}

// Close 关闭连接池
func (p *AdaptiveConnectionPool) Close() {
	if p.running {
		close(p.stopCh)
		p.running = false
	}

	if p.transport != nil {
		p.transport.CloseIdleConnections()
	}
}

// RecordRequest 记录请求统计
func (p *AdaptiveConnectionPool) RecordRequest(success bool, connTime time.Duration) {
	p.stats.mutex.Lock()
	defer p.stats.mutex.Unlock()

	p.stats.TotalRequests++
	if !success {
		p.stats.FailedConnections++
	}

	// 更新平均连接时间（简单移动平均）
	if p.stats.TotalRequests == 1 {
		p.stats.AverageConnTime = connTime
	} else {
		alpha := 0.1 // 平滑因子
		p.stats.AverageConnTime = time.Duration(
			(1-alpha)*float64(p.stats.AverageConnTime) +
				alpha*float64(connTime),
		)
	}
}

// WithConnectionPool 使用指定的连接池配置
func WithConnectionPool(config *ConnectionPoolConfig) SessionOption {
	return func(s *Session) error {
		pool := NewAdaptiveConnectionPool(config)
		s.transport = pool.GetTransport()

		// 如果Session有client，更新其Transport
		if s.client != nil {
			s.client.Transport = s.transport
		}

		return nil
	}
}

// WithHighPerformanceConnectionPool 使用高性能连接池配置
func WithHighPerformanceConnectionPool() SessionOption {
	return WithConnectionPool(HighPerformanceConnectionPoolConfig())
}

// WithLowLatencyConnectionPool 使用低延迟连接池配置
func WithLowLatencyConnectionPool() SessionOption {
	return WithConnectionPool(LowLatencyConnectionPoolConfig())
}

// WithResourceConstrainedConnectionPool 使用资源受限连接池配置
func WithResourceConstrainedConnectionPool() SessionOption {
	return WithConnectionPool(ResourceConstrainedConnectionPoolConfig())
}
