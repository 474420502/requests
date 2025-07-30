package requests

import (
	"testing"
)

// TestConfigErrorPaths 测试Config模块的错误路径和边界情况
func TestConfigErrorPaths(t *testing.T) {
	t.Run("SetProxy_UnsupportedTypes", func(t *testing.T) {
		session := NewSession()
		cfg := session.Config()

		// 测试不支持的类型
		unsupportedTypes := []interface{}{
			123,
			true,
			struct{}{},
			[]string{"invalid"},
			make(chan int),
		}

		for _, unsupportedType := range unsupportedTypes {
			err := cfg.SetProxy(unsupportedType)
			if err == nil {
				t.Errorf("Expected error for unsupported proxy type: %T", unsupportedType)
			}

			// 验证错误信息包含类型信息
			if !contains(err.Error(), "unsupported proxy type") {
				t.Errorf("Error message should mention unsupported proxy type, got: %v", err)
			}
		}
	})

	t.Run("SetTimeout_UnsupportedTypes", func(t *testing.T) {
		session := NewSession()
		cfg := session.Config()

		// 测试不支持的类型
		unsupportedTypes := []interface{}{
			"string-timeout",
			true,
			struct{}{},
			[]int{30},
			make(chan int),
		}

		for _, unsupportedType := range unsupportedTypes {
			err := cfg.SetTimeout(unsupportedType)
			if err == nil {
				t.Errorf("Expected error for unsupported timeout type: %T", unsupportedType)
			}

			// 验证错误信息包含类型信息
			if !contains(err.Error(), "unsupported timeout type") {
				t.Errorf("Error message should mention unsupported timeout type, got: %v", err)
			}
		}
	})

	t.Run("SetBasicAuthLegacy_InvalidArguments", func(t *testing.T) {
		session := NewSession()
		cfg := session.Config()

		// 测试错误数量的参数
		err := cfg.SetBasicAuthLegacy()
		if err == nil {
			t.Error("Expected error for no arguments")
		}

		err = cfg.SetBasicAuthLegacy("user", "pass", "extra")
		if err == nil {
			t.Error("Expected error for too many arguments")
		}

		// 测试错误类型的参数
		err = cfg.SetBasicAuthLegacy(123)
		if err == nil {
			t.Error("Expected error for unsupported type")
		}

		err = cfg.SetBasicAuthLegacy(123, "pass")
		if err == nil {
			t.Error("Expected error for non-string first argument")
		}

		err = cfg.SetBasicAuthLegacy("user", 456)
		if err == nil {
			t.Error("Expected error for non-string second argument")
		}
	})

	t.Run("SetProxyString_InvalidURLs", func(t *testing.T) {
		session := NewSession()
		cfg := session.Config()

		// 测试各种无效URL格式
		invalidURLs := []string{
			"://missing-scheme",
			"http://[::1:invalid",
		}

		for _, invalidURL := range invalidURLs {
			err := cfg.SetProxyString(invalidURL)
			if err == nil {
				t.Errorf("Expected error for invalid proxy URL: %s", invalidURL)
			}
		}

		// 特殊情况：某些字符串虽然不是有效URL，但url.Parse不会返回错误
		// 这种情况下我们只是记录，不算测试失败
		err := cfg.SetProxyString("not-a-url-at-all")
		t.Logf("SetProxyString with 'not-a-url-at-all' returned error: %v", err)
	})
}

// TestDeprecatedMethodWarnings 测试deprecated方法的功能仍然工作
func TestDeprecatedMethodWarnings(t *testing.T) {
	t.Run("SetBasicAuthLegacy_StillWorks", func(t *testing.T) {
		session := NewSession()
		cfg := session.Config()

		// 尽管是deprecated，功能应该仍然工作
		err := cfg.SetBasicAuthLegacy("testuser", "testpass")
		if err != nil {
			t.Errorf("Deprecated method should still work: %v", err)
		}

		if session.auth == nil || session.auth.User != "testuser" {
			t.Error("Deprecated method should still set auth correctly")
		}
	})

	t.Run("SetProxy_DeprecatedStillWorks", func(t *testing.T) {
		session := NewSession()
		cfg := session.Config()

		// 尽管是deprecated，功能应该仍然工作
		err := cfg.SetProxy("http://deprecated-proxy:8080")
		if err != nil {
			t.Errorf("Deprecated method should still work: %v", err)
		}

		if session.transport.Proxy == nil {
			t.Error("Deprecated method should still set proxy")
		}
	})

	t.Run("SetTimeout_DeprecatedStillWorks", func(t *testing.T) {
		session := NewSession()
		cfg := session.Config()

		// 尽管是deprecated，功能应该仍然工作
		err := cfg.SetTimeout(30)
		if err != nil {
			t.Errorf("Deprecated method should still work: %v", err)
		}

		if session.client.Timeout.Seconds() != 30 {
			t.Error("Deprecated method should still set timeout correctly")
		}
	})
}

// contains 检查字符串是否包含子字符串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (substr == "" || indexString(s, substr) >= 0)
}

// indexString  返回子字符串在字符串中的索引
func indexString(s, substr string) int {
	n := len(substr)
	if n == 0 {
		return 0
	}
	for i := 0; i <= len(s)-n; i++ {
		if s[i:i+n] == substr {
			return i
		}
	}
	return -1
}
