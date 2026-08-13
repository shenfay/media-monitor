package utils

import (
	"sync"

	"github.com/ua-parser/uap-go/uaparser"
)

var (
	parser     *uaparser.Parser
	parserOnce sync.Once
)

// UserAgentInfo User-Agent 解析结果
type UserAgentInfo struct {
	Browser string
	OS      string
	Device  string
}

// initParser 初始化 User-Agent 解析器（单例）
func initParser() {
	parserOnce.Do(func() {
		// 使用内置的正则表达式规则
		var err error
		parser, err = uaparser.New()
		if err != nil {
			// 如果初始化失败，使用 nil（会在 ParseUserAgent 中处理）
			parser = nil
		}
	})
}

// ParseUserAgent 解析 User-Agent 字符串
func ParseUserAgent(userAgent string) UserAgentInfo {
	initParser()

	info := UserAgentInfo{
		Browser: "Unknown",
		OS:      "Unknown",
		Device:  "Unknown",
	}

	if userAgent == "" {
		return info
	}

	// 如果解析器初始化失败，使用简单的启发式规则
	if parser == nil {
		return parseUserAgentSimple(userAgent)
	}

	client := parser.Parse(userAgent)

	// 提取浏览器信息
	if client.UserAgent.Family != "" {
		info.Browser = client.UserAgent.Family
	}

	// 提取操作系统信息
	if client.Os.Family != "" {
		info.OS = client.Os.Family
	}

	// 提取设备类型
	info.Device = detectDeviceType(userAgent)

	return info
}

// parseUserAgentSimple 简单的 User-Agent 解析（降级方案）
func parseUserAgentSimple(ua string) UserAgentInfo {
	info := UserAgentInfo{
		Browser: "Unknown",
		OS:      "Unknown",
		Device:  detectDeviceType(ua),
	}

	// 简单的浏览器检测
	if containsStr(ua, "Chrome") && !containsStr(ua, "Edg") {
		info.Browser = "Chrome"
	} else if containsStr(ua, "Firefox") {
		info.Browser = "Firefox"
	} else if containsStr(ua, "Safari") && !containsStr(ua, "Chrome") {
		info.Browser = "Safari"
	} else if containsStr(ua, "Edg") {
		info.Browser = "Edge"
	} else if containsStr(ua, "MSIE") || containsStr(ua, "Trident") {
		info.Browser = "IE"
	} else if containsStr(ua, "curl") {
		info.Browser = "curl"
	}

	// 简单的操作系统检测
	if containsStr(ua, "Windows") {
		info.OS = "Windows"
	} else if containsStr(ua, "Mac OS X") || containsStr(ua, "Macintosh") {
		info.OS = "macOS"
	} else if containsStr(ua, "Linux") && !containsStr(ua, "Android") {
		info.OS = "Linux"
	} else if containsStr(ua, "Android") {
		info.OS = "Android"
	} else if containsStr(ua, "iPhone") || containsStr(ua, "iPad") {
		info.OS = "iOS"
	}

	return info
}

// detectDeviceType 检测设备类型
func detectDeviceType(ua string) string {
	if ua == "" {
		return "Unknown"
	}

	if containsStr(ua, "iPad") || containsStr(ua, "Tablet") {
		return "tablet"
	}

	if containsStr(ua, "Mobile") || containsStr(ua, "Android") || containsStr(ua, "iPhone") {
		return "mobile"
	}

	return "desktop"
}

// containsStr 检查字符串是否包含子串
func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && indexOf(s, substr) >= 0
}

// indexOf 查找子串位置
func indexOf(s, substr string) int {
	if len(substr) == 0 {
		return 0
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
