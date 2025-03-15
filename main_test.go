package main

import (
	"net"
	"testing"
)

// TestCustomResolverWithoutServer 测试不指定DNS服务器时的解析
func TestCustomResolverWithoutServer(t *testing.T) {
	resolver := &customResolver{server: ""}

	// 测试lookupCNAME
	_, err := resolver.lookupCNAME("example.com")
	if err != nil {
		t.Errorf("lookupCNAME failed with default resolver: %v", err)
	}

	// 测试lookupIP
	ips, err := resolver.lookupIP("example.com")
	if err != nil {
		t.Errorf("lookupIP failed with default resolver: %v", err)
	}

	if len(ips) == 0 {
		t.Error("lookupIP returned empty result for example.com")
	}
}

// TestCustomResolverWithServer 测试指定DNS服务器时的解析
func TestCustomResolverWithServer(t *testing.T) {
	// 使用谷歌DNS服务器
	resolver := &customResolver{server: "8.8.8.8:53"}

	// 测试lookupCNAME
	_, err := resolver.lookupCNAME("example.com")
	if err != nil {
		t.Errorf("lookupCNAME failed with custom resolver: %v", err)
	}

	// 测试lookupIP
	ips, err := resolver.lookupIP("example.com")
	if err != nil {
		t.Errorf("lookupIP failed with custom resolver: %v", err)
	}

	if len(ips) == 0 {
		t.Error("lookupIP returned empty result for example.com")
	}
}

// TestContains 测试contains函数
func TestContains(t *testing.T) {
	testCases := []struct {
		slice    []string
		item     string
		expected bool
	}{
		{[]string{"a", "b", "c"}, "a", true},
		{[]string{"a", "b", "c"}, "d", false},
		{[]string{"a.", "b.", "c."}, "a", true},
		{[]string{}, "a", false},
	}

	for _, tc := range testCases {
		result := contains(tc.slice, tc.item)
		if result != tc.expected {
			t.Errorf("contains(%v, %s) = %v; expected %v",
				tc.slice, tc.item, result, tc.expected)
		}
	}
}

// TestInvalidDomain 测试无效域名
func TestInvalidDomain(t *testing.T) {
	resolver := &customResolver{server: ""}

	_, err := resolver.lookupIP("invalid-domain-that-does-not-exist.xyz")
	if err == nil {
		t.Error("Expected error for invalid domain, but got nil")
	}
}

// MockResolver 模拟解析器
type MockResolver struct {
	CNAMEResult string
	CNAMEError  error
	IPResult    []net.IP
	IPError     error
}
