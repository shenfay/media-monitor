package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	apperrors "github.com/shenfay/go-react-admin/pkg/errors"
	"github.com/shenfay/go-react-admin/internal/app/crawl"
	"github.com/shenfay/go-react-admin/internal/transport/http/response"
)

// CrawlHandler 抓取模块 HTTP 处理器
type CrawlHandler struct {
	svc *crawl.Service
}

// NewCrawlHandler 创建抓取处理器
func NewCrawlHandler(svc *crawl.Service) *CrawlHandler {
	return &CrawlHandler{svc: svc}
}

// RegisterAdminRoutes 注册需要 JWT + Casbin 权限的管理路由
func (h *CrawlHandler) RegisterAdminRoutes(rg *gin.RouterGroup) {
	rg.GET("/sources", h.ListSources)
	rg.POST("/sources", h.CreateSource)
	rg.PUT("/sources/:id", h.UpdateSource)
	rg.DELETE("/sources/:id", h.DeleteSource)
	rg.POST("/sources/:id/run", h.RunSource)
	rg.POST("/tasks", h.CreateTask)
	rg.GET("/tasks", h.ListTasks)
	rg.GET("/tasks/:id", h.GetTask)
}

// ListSources GET /api/v1/crawl/sources
// @Summary 列出数据源
// @Tags Crawl
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
func (h *CrawlHandler) ListSources(c *gin.Context) {
	sources, err := h.svc.ListSources(c.Request.Context(), false)
	if err != nil {
		response.Error(c, err)
		return
	}
	// 脱敏 auth（明文仅机器接口返回）
	for _, s := range sources {
		s.Auth = nil
	}
	response.Success(c, sources)
}

// CreateSource POST /api/v1/crawl/sources
// @Summary 创建数据源
// @Tags Crawl
// @Security BearerAuth
// @Success 201 {object} response.SuccessResponse
func (h *CrawlHandler) CreateSource(c *gin.Context) {
	var req crawl.CreateSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}
	owner := c.GetString("user_id")
	src, err := h.svc.CreateSource(c.Request.Context(), owner, req)
	if err != nil {
		response.Error(c, err)
		return
	}
	src.Auth = nil
	response.Created(c, src)
}

// UpdateSource PUT /api/v1/crawl/sources/:id
// @Summary 更新数据源
// @Tags Crawl
// @Security BearerAuth
func (h *CrawlHandler) UpdateSource(c *gin.Context) {
	id := c.Param("id")
	var req crawl.UpdateSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}
	src, err := h.svc.UpdateSource(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, crawl.ErrSourceNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, apperrors.NewAppError("CRAWL.NOT_FOUND", "数据源不存在", http.StatusNotFound).WithError(err))
			return
		}
		response.Error(c, err)
		return
	}
	src.Auth = nil
	response.Success(c, src)
}

// DeleteSource DELETE /api/v1/crawl/sources/:id
// @Summary 删除数据源
// @Tags Crawl
// @Security BearerAuth
func (h *CrawlHandler) DeleteSource(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteSource(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}
	response.NoContent(c)
}

// RunSource POST /api/v1/crawl/sources/:id/run
// @Summary 手动触发一次抓取
// @Tags Crawl
// @Security BearerAuth
func (h *CrawlHandler) RunSource(c *gin.Context) {
	id := c.Param("id")
	taskID, err := h.svc.EnqueueForSource(c.Request.Context(), id, "manual")
	if err != nil {
		if errors.Is(err, crawl.ErrSourceNotFound) {
			response.Error(c, apperrors.NewAppError("CRAWL.NOT_FOUND", "数据源不存在", http.StatusNotFound).WithError(err))
			return
		}
		if errors.Is(err, crawl.ErrSourceDisabled) {
			response.Error(c, apperrors.NewAppError("CRAWL.DISABLED", "数据源已禁用", http.StatusBadRequest).WithError(err))
			return
		}
		if errors.Is(err, crawl.ErrTaskAlreadyRun) {
			response.Error(c, apperrors.NewAppError("CRAWL.ALREADY_RUNNING", "该数据源已有任务在运行中", http.StatusConflict).WithError(err))
			return
		}
		response.Error(c, err)
		return
	}
	response.Success(c, gin.H{"task_id": taskID})
}

// CreateTask POST /api/v1/crawl/tasks
// @Summary 创建抓取任务
// @Tags Crawl
// @Security BearerAuth
func (h *CrawlHandler) CreateTask(c *gin.Context) {
	var req crawl.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}
	taskID, err := h.svc.CreateTask(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, crawl.ErrSourceNotFound) {
			response.Error(c, apperrors.NewAppError("CRAWL.NOT_FOUND", "数据源不存在", http.StatusNotFound).WithError(err))
			return
		}
		if errors.Is(err, crawl.ErrTaskAlreadyRun) {
			response.Error(c, apperrors.NewAppError("CRAWL.ALREADY_RUNNING", "该数据源已有任务在运行中", http.StatusConflict).WithError(err))
			return
		}
		response.Error(c, err)
		return
	}
	response.Created(c, gin.H{"task_id": taskID})
}

// GetTask GET /api/v1/crawl/tasks/:id
// @Summary 获取任务详情
// @Tags Crawl
// @Security BearerAuth
func (h *CrawlHandler) GetTask(c *gin.Context) {
	id := c.Param("id")
	t, err := h.svc.GetTask(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, crawl.ErrTaskNotFound) {
			response.Error(c, apperrors.NewAppError("CRAWL.NOT_FOUND", "任务不存在", http.StatusNotFound).WithError(err))
			return
		}
		response.Error(c, err)
		return
	}
	response.Success(c, t)
}

// ListTasks GET /api/v1/crawl/tasks
// @Summary 列出任务（可按 source_id 过滤）
// @Tags Crawl
// @Security BearerAuth
func (h *CrawlHandler) ListTasks(c *gin.Context) {
	sourceID := c.Query("source_id")
	tasks, err := h.svc.ListTasks(c.Request.Context(), sourceID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, tasks)
}
