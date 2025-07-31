package requests

import (
	"testing"
	"time"
)

// TestConfigErrorPaths 测试Config模块的错误路径和边界情况
func TestConfigErrorPaths(t *testing.T) {
	t.Run("SetProxyString_InvalidURL", func(t *testing.T) {
		session := NewSession()
		cfg := session.Config()

		// 测试无效的代理URL
		invalidURLs := []string{
			"invalid-url",
			"://missing-scheme",
		}

		for _, invalidURL := range invalidURLs {
			err := cfg.SetProxyString(invalidURL)
			if err == nil {
				t.Errorf("Expected error for invalid proxy URL: %s", invalidURL)
			}
		}
	})

	t.Run("SetTimeoutDuration_EdgeCases", func(t *testing.T) {
		session := NewSession()
		cfg := session.Config()

		// 测试边界情况
		testCases := []struct {
			name     string
			duration time.Duration
		}{
			{"Zero duration", 0},
			{"Negative duration", -time.Second}, // Go允许负超时
			{"Very large duration", time.Hour * 24 * 365},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				cfg.SetTimeoutDuration(tc.duration)
				// 验证设置成功
				if session.client.Timeout != tc.duration {
					t.Errorf("Expected timeout %v, got %v", tc.duration, session.client.Timeout)
				}
			})
		}
	})

	t.Run("SetBasicAuth_ErrorHandling", func(t *testing.T) {
		session := NewSession()
		cfg := session.Config()

		// SetBasicAuth现在总是成功，为了向后兼容
		err := cfg.SetBasicAuth("", "password")
		if err != nil {
			t.Errorf("SetBasicAuth should not return error for empty username: %v", err)
		}

		err = cfg.SetBasicAuth("username", "")
		if err != nil {
			t.Errorf("SetBasicAuth should not return error for empty password: %v", err)
		}

		// 测试正常情况
		err = cfg.SetBasicAuth("user", "pass")
		if err != nil {
			t.Errorf("SetBasicAuth should not return error for valid credentials: %v", err)
		}
	})
}
