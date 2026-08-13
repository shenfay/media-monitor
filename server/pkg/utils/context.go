package utils

import "context"

// contextKey 自定义 Context 键类型
type contextKey string

const (
	// OperatorUserIDKey 操作人用户 ID 在 context 中的键
	OperatorUserIDKey contextKey = "operator_user_id"
	// OperatorEmailKey 操作人邮箱在 context 中的键
	OperatorEmailKey contextKey = "operator_email"
	// RequestIPKey 请求 IP 在 context 中的键
	RequestIPKey contextKey = "request_ip"
	// RequestUserAgentKey 请求 User-Agent 在 context 中的键
	RequestUserAgentKey contextKey = "request_user_agent"
	// RequestDeviceKey 请求设备类型在 context 中的键
	RequestDeviceKey contextKey = "request_device"
	// RequestBrowserKey 请求浏览器在 context 中的键
	RequestBrowserKey contextKey = "request_browser"
	// RequestOSKey 请求操作系统在 context 中的键
	RequestOSKey contextKey = "request_os"
)

// WithOperator 将操作人信息注入到 context
func WithOperator(ctx context.Context, userID, email string) context.Context {
	ctx = context.WithValue(ctx, OperatorUserIDKey, userID)
	ctx = context.WithValue(ctx, OperatorEmailKey, email)
	return ctx
}

// WithRequestInfo 将请求元数据注入到 context
func WithRequestInfo(ctx context.Context, ip, userAgent, device, browser, os string) context.Context {
	ctx = context.WithValue(ctx, RequestIPKey, ip)
	ctx = context.WithValue(ctx, RequestUserAgentKey, userAgent)
	ctx = context.WithValue(ctx, RequestDeviceKey, device)
	ctx = context.WithValue(ctx, RequestBrowserKey, browser)
	ctx = context.WithValue(ctx, RequestOSKey, os)
	return ctx
}

// GetOperatorUserID 从 context 获取操作人用户 ID（供 Service 层使用）
func GetOperatorUserID(ctx context.Context) string {
	if id, ok := ctx.Value(OperatorUserIDKey).(string); ok {
		return id
	}
	return ""
}

// GetOperatorEmail 从 context 获取操作人邮箱（供 Service 层使用）
func GetOperatorEmail(ctx context.Context) string {
	if email, ok := ctx.Value(OperatorEmailKey).(string); ok {
		return email
	}
	return ""
}

// GetRequestIP 从 context 获取请求 IP
func GetRequestIP(ctx context.Context) string {
	if v, ok := ctx.Value(RequestIPKey).(string); ok {
		return v
	}
	return ""
}

// GetRequestUserAgent 从 context 获取请求 User-Agent
func GetRequestUserAgent(ctx context.Context) string {
	if v, ok := ctx.Value(RequestUserAgentKey).(string); ok {
		return v
	}
	return ""
}

// GetRequestDevice 从 context 获取请求设备类型
func GetRequestDevice(ctx context.Context) string {
	if v, ok := ctx.Value(RequestDeviceKey).(string); ok {
		return v
	}
	return ""
}

// GetRequestBrowser 从 context 获取请求浏览器
func GetRequestBrowser(ctx context.Context) string {
	if v, ok := ctx.Value(RequestBrowserKey).(string); ok {
		return v
	}
	return ""
}

// GetRequestOS 从 context 获取请求操作系统
func GetRequestOS(ctx context.Context) string {
	if v, ok := ctx.Value(RequestOSKey).(string); ok {
		return v
	}
	return ""
}
