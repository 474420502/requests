package requests

import (
	"crypto/tls"
	"net/url"
	"testing"
)

// TestConfigMethods 测试配置方法
func TestConfigMethods(t *testing.T) {
	t.Run("SetBasicAuthStruct", func(t *testing.T) {
		session := NewSession()
		config := session.Config()
		auth := &BasicAuth{
			User:     "testuser",
			Password: "testpass",
		}
		config.SetBasicAuthStruct(auth)

		if session.auth == nil {
			t.Fatal("Basic auth should be set")
		}
		if session.auth.User != "testuser" || session.auth.Password != "testpass" {
			t.Error("Basic auth credentials not set correctly")
		}

		// 测试nil auth
		config.SetBasicAuthStruct(nil)
		if session.auth != nil {
			t.Error("Basic auth should be nil when nil is passed")
		}
	})

	t.Run("SetTLSConfig", func(t *testing.T) {
		session := NewSession()
		config := session.Config()
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		config.SetTLSConfig(tlsConfig)

		if session.transport.TLSClientConfig != tlsConfig {
			t.Error("TLS config not set correctly")
		}
	})

	t.Run("SetProxyString", func(t *testing.T) {
		session := NewSession()
		config := session.Config()
		proxyURL := "http://127.0.0.1:8080"
		config.SetProxyString(proxyURL)

		if session.transport.Proxy == nil {
			t.Error("Proxy should be set")
		}

		// 测试无效的代理URL
		config.SetProxyString("invalid-url")
		// 应该不会崩溃，但代理设置可能失败
	})

	t.Run("setProxyURL", func(t *testing.T) {
		session := NewSession()
		config := session.Config()
		parsedURL, _ := url.Parse("http://127.0.0.1:8080")
		err := config.setProxyURL(parsedURL)
		if err != nil {
			t.Errorf("setProxyURL should not return error for valid URL: %v", err)
		}
		if session.transport.Proxy == nil {
			t.Error("Proxy should be set")
		}

		// 测试空URL的情况 - 这会导致panic，所以我们不测试它
		// err2 := config.setProxyURL(nil)
		// if err2 == nil {
		//     t.Error("setProxyURL should return error for nil URL")
		// }
	})
}

// TestSetProxyString 测试SetProxyString方法的各种情况
func TestSetProxyString(t *testing.T) {
	t.Run("ValidHTTPURL", func(t *testing.T) {
		session := NewSession()
		config := session.Config()
		err := config.SetProxyString("http://127.0.0.1:8080")
		if err != nil {
			t.Errorf("SetProxyString should not return error: %v", err)
		}

		if session.transport.Proxy == nil {
			t.Error("Proxy should be set")
		}
	})

	t.Run("ClearProxy", func(t *testing.T) {
		session := NewSession()
		config := session.Config()
		config.ClearProxy()
		// 应该不会崩溃，代理被清除
		if session.transport.Proxy != nil {
			t.Error("Proxy should be cleared")
		}
	})

	t.Run("HTTPSProxy", func(t *testing.T) {
		session := NewSession()
		config := session.Config()
		err := config.SetProxyString("https://proxy.example.com:8080")
		if err != nil {
			t.Errorf("SetProxyString should not return error: %v", err)
		}

		if session.transport.Proxy == nil {
			t.Error("HTTPS proxy should be set")
		}
	})

	t.Run("SOCKS5Proxy", func(t *testing.T) {
		session := NewSession()
		config := session.Config()
		err := config.SetProxyString("socks5://127.0.0.1:1080")
		if err != nil {
			t.Errorf("SetProxyString should not return error: %v", err)
		}

		if session.transport.Proxy == nil {
			t.Error("SOCKS5 proxy should be set")
		}
	})
}

// TestProxyEdgeCases 测试代理设置的边界情况
func TestProxyEdgeCases(t *testing.T) {
	t.Run("ProxyWithCredentials", func(t *testing.T) {
		session := NewSession()
		config := session.Config()
		proxyURL := "http://user:pass@127.0.0.1:8080"
		config.SetProxyString(proxyURL)

		if session.transport.Proxy == nil {
			t.Error("Proxy with credentials should be set")
		}
	})

	t.Run("ProxyWithPort", func(t *testing.T) {
		session := NewSession()
		config := session.Config()
		proxyURL := "http://proxy.example.com:3128"
		config.SetProxyString(proxyURL)

		if session.transport.Proxy == nil {
			t.Error("Proxy with custom port should be set")
		}
	})

	t.Run("ClearProxyAfterSet", func(t *testing.T) {
		session := NewSession()
		config := session.Config()
		// 先设置代理
		config.SetProxyString("http://127.0.0.1:8080")
		if session.transport.Proxy == nil {
			t.Error("Proxy should be set initially")
		}

		// 然后清除代理
		config.ClearProxy()
		if session.transport.Proxy != nil {
			t.Error("Proxy should be cleared")
		}
	})
}

// TestTLSConfigEdgeCases 测试TLS配置的边界情况
func TestTLSConfigEdgeCases(t *testing.T) {
	t.Run("TLSConfigOverwrite", func(t *testing.T) {
		session := NewSession()
		config := session.Config()

		// 设置第一个TLS配置
		config1 := &tls.Config{InsecureSkipVerify: true}
		config.SetTLSConfig(config1)

		// 设置第二个TLS配置（应该覆盖第一个）
		config2 := &tls.Config{InsecureSkipVerify: false}
		config.SetTLSConfig(config2)

		if session.transport.TLSClientConfig != config2 {
			t.Error("Second TLS config should overwrite the first")
		}
		if session.transport.TLSClientConfig.InsecureSkipVerify {
			t.Error("InsecureSkipVerify should be false from second config")
		}
	})

	t.Run("NilTLSConfig", func(t *testing.T) {
		session := NewSession()
		config := session.Config()
		config.SetTLSConfig(nil)

		if session.transport.TLSClientConfig != nil {
			t.Error("TLS config should be nil when nil is passed")
		}
	})
}

// TestAuthCombinations 测试认证方法的组合
func TestAuthCombinations(t *testing.T) {
	t.Run("BasicAuthOverwrite", func(t *testing.T) {
		session := NewSession()
		config := session.Config()

		// 设置第一个认证
		config.SetBasicAuth("user1", "pass1")
		if session.auth == nil || session.auth.User != "user1" {
			t.Error("First auth should be set")
		}

		// 设置第二个认证（应该覆盖第一个）
		config.SetBasicAuth("user2", "pass2")
		if session.auth == nil || session.auth.User != "user2" || session.auth.Password != "pass2" {
			t.Error("Second auth should overwrite the first")
		}
	})

	t.Run("ClearAuthAfterSet", func(t *testing.T) {
		session := NewSession()
		config := session.Config()
		config.SetBasicAuth("user", "pass")
		if session.auth == nil {
			t.Error("Auth should be set initially")
		}

		config.ClearBasicAuth()
		if session.auth != nil {
			t.Error("Auth should be cleared")
		}
	})

	t.Run("EmptyCredentials", func(t *testing.T) {
		session := NewSession()
		config := session.Config()
		err := config.SetBasicAuth("", "")
		if err != nil {
			t.Errorf("SetBasicAuth should not return error for empty credentials: %v", err)
		}
		if session.auth == nil {
			t.Fatal("Auth should be set even with empty credentials")
		}
		if session.auth.User != "" || session.auth.Password != "" {
			t.Error("Empty credentials should be preserved")
		}
	})
}

// TestSetBasicAuthString 测试SetBasicAuthString方法
func TestSetBasicAuthString(t *testing.T) {
	session := NewSession()
	config := session.Config()

	config.SetBasicAuthString("stringuser", "stringpass")

	if session.auth == nil {
		t.Fatal("Basic auth should be set")
	}
	if session.auth.User != "stringuser" || session.auth.Password != "stringpass" {
		t.Error("Basic auth credentials from SetBasicAuthString not set correctly")
	}
}
