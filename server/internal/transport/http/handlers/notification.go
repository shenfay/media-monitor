package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/shenfay/go-react-admin/internal/app/notification"
	domain "github.com/shenfay/go-react-admin/internal/domain/notification"
	"github.com/shenfay/go-react-admin/internal/transport/http/response"
	validationErr "github.com/shenfay/go-react-admin/pkg/errors/validation"
	"github.com/shenfay/go-react-admin/pkg/utils"
)

// NotificationHandler 消息 HTTP 处理器
type NotificationHandler struct {
	service *notification.AppService
}

// NewNotificationHandler 创建消息处理器
func NewNotificationHandler(service *notification.AppService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

// RegisterRoutes 注册用户消息路由
func (h *NotificationHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.ListMessages)
	rg.GET("/unread-count", h.GetUnreadCount)
	rg.POST("/:id/read", h.MarkAsRead)
	rg.POST("/read-all", h.MarkAllAsRead)
}

// RegisterAdminRoutes 注册管理员消息路由
func (h *NotificationHandler) RegisterAdminRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.ListAllMessages)
	rg.POST("/system", h.SendSystemNotification)
}

// ListMessages 获取当前用户消息列表
// GET /api/v1/messages?type=system&category=verification&is_read=false&limit=20&offset=0
// @Summary 获取当前用户消息列表
// @Tags Messages
// @Produce json
// @Security BearerAuth
// @Param type query string false "消息类型" Enums(system, companion)
// @Param category query string false "消息分类"
// @Param is_read query bool false "已读状态"
// @Param limit query int false "每页条数" default(20)
// @Param offset query int false "偏移量" default(0)
// @Success 200 {object} response.SuccessResponse "消息列表"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /messages [get]
func (h *NotificationHandler) ListMessages(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Error(c, domain.ErrMessageAccessDenied)
		return
	}

	msgType := domain.MessageType(c.Query("type"))
	category := domain.MessageCategory(c.Query("category"))
	limit := utils.ToInt(c.DefaultQuery("limit", "20"))
	offset := utils.ToInt(c.DefaultQuery("offset", "0"))

	var isRead *bool
	if isReadStr := c.Query("is_read"); isReadStr != "" {
		val := isReadStr == "true"
		isRead = &val
	}

	result, err := h.service.GetMessages(c.Request.Context(), userID, msgType, category, isRead, limit, offset)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

// GetUnreadCount 获取未读消息计数
// GET /api/v1/messages/unread-count
// @Summary 获取未读消息计数
// @Tags Messages
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse "未读计数"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /messages/unread-count [get]
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Error(c, domain.ErrMessageAccessDenied)
		return
	}

	result, err := h.service.GetUnreadCount(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

// MarkAsRead 标记单条消息已读
// POST /api/v1/messages/:id/read
// @Summary 标记消息已读
// @Tags Messages
// @Produce json
// @Security BearerAuth
// @Param id path string true "消息ID"
// @Success 200 {object} response.SuccessResponse "标记成功"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "无权操作"
// @Failure 404 {object} response.ErrorResponse "消息不存在"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /messages/{id}/read [post]
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Error(c, domain.ErrMessageAccessDenied)
		return
	}

	id := c.Param("id")
	if err := h.service.MarkAsRead(c.Request.Context(), id, userID); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "marked as read"})
}

// MarkAllAsRead 标记全部已读
// POST /api/v1/messages/read-all
// @Summary 标记全部已读
// @Tags Messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body notification.ReadAllCmd false "可选类型筛选"
// @Success 200 {object} response.SuccessResponse "标记成功"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /messages/read-all [post]
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Error(c, domain.ErrMessageAccessDenied)
		return
	}

	var cmd notification.ReadAllCmd
	_ = c.ShouldBindJSON(&cmd)

	msgType := domain.MessageType(cmd.Type)
	if err := h.service.MarkAllAsRead(c.Request.Context(), userID, msgType); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "all marked as read"})
}

// SendSystemNotificationRequest 发送系统通知请求
type SendSystemNotificationRequest struct {
	RecipientID string `json:"recipient_id" binding:"required"`
	Category    string `json:"category" binding:"required"`
	Title       string `json:"title" binding:"required,max=255"`
	Content     string `json:"content" binding:"required,max=2000"`
}

// SendSystemNotification 发送系统通知
// @Summary 管理员发送系统通知
// @Tags Messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body SendSystemNotificationRequest true "通知内容"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /admin/messages/system [post]
func (h *NotificationHandler) SendSystemNotification(c *gin.Context) {
	var req SendSystemNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, validationErr.FromGinError(err))
		return
	}

	category := domain.MessageCategory(req.Category)
	if err := h.service.SendSystemNotification(c.Request.Context(), req.RecipientID, category, req.Title, req.Content); err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, gin.H{"message": "系统通知已发送"})
}

// ListAllMessages 管理员查询所有消息
// GET /api/v1/admin/messages?type=system&category=verification&limit=20&offset=0
// @Summary 管理员查询所有消息
// @Tags Messages
// @Produce json
// @Security BearerAuth
// @Param type query string false "消息类型"
// @Param category query string false "消息分类"
// @Param limit query int false "每页条数" default(20)
// @Param offset query int false "偏移量" default(0)
// @Success 200 {object} response.SuccessResponse "消息列表"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /admin/messages [get]
func (h *NotificationHandler) ListAllMessages(c *gin.Context) {
	msgType := domain.MessageType(c.Query("type"))
	category := domain.MessageCategory(c.Query("category"))
	limit := utils.ToInt(c.DefaultQuery("limit", "20"))
	offset := utils.ToInt(c.DefaultQuery("offset", "0"))

	result, err := h.service.GetAllMessages(c.Request.Context(), msgType, category, limit, offset)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}
