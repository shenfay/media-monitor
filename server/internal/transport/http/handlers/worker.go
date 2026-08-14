package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/shenfay/go-react-admin/internal/app/crawl"
	"github.com/shenfay/go-react-admin/internal/transport/http/response"
)

// WorkerHandler Worker 管理 HTTP 处理器
type WorkerHandler struct {
	registry *crawl.WorkerRegistry
}

// NewWorkerHandler 创建 Worker 管理器
func NewWorkerHandler(registry *crawl.WorkerRegistry) *WorkerHandler {
	return &WorkerHandler{registry: registry}
}

// RegisterRoutes 注册 Worker 管理路由
func (h *WorkerHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/workers", h.ListWorkers)
	rg.GET("/workers/adapters", h.ListAdapters)
	rg.POST("/workers/:id/pause", h.PauseWorker)
	rg.POST("/workers/:id/resume", h.ResumeWorker)
	rg.POST("/workers/:id/shutdown", h.ShutdownWorker)
}

// ListWorkers GET /api/v1/crawl/workers
// @Summary 列出所有 Worker 实例
// @Tags Worker
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
func (h *WorkerHandler) ListWorkers(c *gin.Context) {
	workers, err := h.registry.ListWorkers(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, workers)
}

// ListAdapters GET /api/v1/crawl/workers/adapters
// @Summary 列出所有已知适配器
// @Tags Worker
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
func (h *WorkerHandler) ListAdapters(c *gin.Context) {
	adapters, err := h.registry.ListAdapters(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, adapters)
}

// PauseWorker POST /api/v1/crawl/workers/:id/pause
// @Summary 暂停 Worker
// @Tags Worker
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
func (h *WorkerHandler) PauseWorker(c *gin.Context) {
	id := c.Param("id")
	if err := h.registry.SendCommand(c.Request.Context(), id, "pause"); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, gin.H{"message": "pause command sent"})
}

// ResumeWorker POST /api/v1/crawl/workers/:id/resume
// @Summary 恢复 Worker
// @Tags Worker
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
func (h *WorkerHandler) ResumeWorker(c *gin.Context) {
	id := c.Param("id")
	if err := h.registry.SendCommand(c.Request.Context(), id, "resume"); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, gin.H{"message": "resume command sent"})
}

// ShutdownWorker POST /api/v1/crawl/workers/:id/shutdown
// @Summary 下线 Worker（优雅停止）
// @Tags Worker
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
func (h *WorkerHandler) ShutdownWorker(c *gin.Context) {
	id := c.Param("id")
	if err := h.registry.SendCommand(c.Request.Context(), id, "shutdown"); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, gin.H{"message": "shutdown command sent"})
}

// RespondError 返回错误响应（供中间件使用）
func RespondError(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{"error": msg})
}
