package handlers

import (
	"encoding/json"

	"github.com/gin-gonic/gin"

	appsetting "github.com/shenfay/go-react-admin/internal/app/setting"
	"github.com/shenfay/go-react-admin/internal/domain/setting"
	"github.com/shenfay/go-react-admin/internal/transport/http/response"
)

// SettingHandler 系统设置 HTTP 处理器
type SettingHandler struct {
	service *appsetting.Service
}

// NewSettingHandler 创建系统设置处理器
func NewSettingHandler(service *appsetting.Service) *SettingHandler {
	return &SettingHandler{service: service}
}

// RegisterRoutes 注册系统设置路由
func (h *SettingHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.ListSettings)
	rg.GET("/:key", h.GetSetting)
	rg.PUT("", h.BatchUpdateSettings)
}

// ListSettings 获取设置列表（支持 category 过滤）
// GET /api/v1/settings?category=basic
// @Summary 获取系统设置列表
// @Tags Settings
// @Produce json
// @Security BearerAuth
// @Param category query string false "设置分类"
// @Success 200 {object} response.SuccessResponse "设置列表"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /settings [get]
func (h *SettingHandler) ListSettings(c *gin.Context) {
	category := c.Query("category")

	settings, err := h.service.GetAllSettings(c.Request.Context(), category)
	if err != nil {
		response.Error(c, err)
		return
	}

	// 对敏感字段进行脱敏处理（返回时不返回真实密码）
	for _, s := range settings {
		if setting.IsSensitiveKey(s.Key) {
			s.Value = maskSensitiveFields(s.Value)
		}
	}

	response.Success(c, settings)
}

// GetSetting 获取单个设置
// GET /api/v1/settings/:key
// @Summary 获取单个系统设置
// @Tags Settings
// @Produce json
// @Security BearerAuth
// @Param key path string true "设置项Key"
// @Success 200 {object} response.SuccessResponse "设置详情"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "设置项不存在"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /settings/{key} [get]
func (h *SettingHandler) GetSetting(c *gin.Context) {
	key := c.Param("key")

	s, err := h.service.GetSettingByKey(c.Request.Context(), key)
	if err != nil {
		response.Error(c, err)
		return
	}

	// 敏感字段脱敏
	if setting.IsSensitiveKey(s.Key) {
		s.Value = maskSensitiveFields(s.Value)
	}

	response.Success(c, s)
}

// batchUpdateRequest 批量更新请求体
type batchUpdateRequest struct {
	Settings []settingItem `json:"settings" binding:"required,dive"`
}

// settingItem 单条设置更新项
type settingItem struct {
	Key   string          `json:"key" binding:"required"`
	Value json.RawMessage `json:"value" binding:"required"`
}

// BatchUpdateSettings 批量更新设置
// PUT /api/v1/settings
// @Summary 批量更新系统设置
// @Tags Settings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body batchUpdateRequest true "批量更新数据"
// @Success 200 {object} response.SuccessResponse "更新成功"
// @Failure 400 {object} response.ErrorResponse "请求参数错误"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /settings [put]
func (h *SettingHandler) BatchUpdateSettings(c *gin.Context) {
	var req batchUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	// 转换为 service 层的更新结构
	updates := make([]setting.SettingUpdate, len(req.Settings))
	for i, item := range req.Settings {
		updates[i] = setting.SettingUpdate{
			Key:   item.Key,
			Value: item.Value,
		}
	}

	// 从 JWT context 中获取当前用户 ID 作为 updated_by
	updatedByStr := c.GetString("user_id")
	var updatedBy *string
	if updatedByStr != "" {
		updatedBy = &updatedByStr
	}

	if err := h.service.BatchUpdate(c.Request.Context(), updates, updatedBy); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"message": "设置已保存",
	})
}

// maskSensitiveFields 对渠道配置中的敏感字段进行脱敏
// 返回时将 password/secret 替换为占位符
func maskSensitiveFields(value json.RawMessage) json.RawMessage {
	var obj map[string]interface{}
	if err := json.Unmarshal(value, &obj); err != nil {
		return value
	}

	sensitiveFields := []string{"password", "secret"}
	masked := false
	for _, field := range sensitiveFields {
		if v, ok := obj[field]; ok {
			if str, isStr := v.(string); isStr && str != "" {
				obj[field] = "••••••"
				masked = true
			}
		}
	}

	if !masked {
		return value
	}

	result, err := json.Marshal(obj)
	if err != nil {
		return value
	}
	return result
}
