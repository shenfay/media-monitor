package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/gin-gonic/gin"

	"github.com/shenfay/go-react-admin/internal/app/authentication"
	"github.com/shenfay/go-react-admin/internal/infra/ws"
	"github.com/shenfay/go-react-admin/pkg/logger"
)

const (
	wsReadTimeout  = 60 * time.Second
	wsWriteTimeout = 10 * time.Second
	wsPingInterval = 30 * time.Second
)

// WSHandler WebSocket 连接处理器
type WSHandler struct {
	hub          *ws.Hub
	tokenManager authentication.TokenManager
}

// NewWSHandler 创建 WebSocket 处理器
func NewWSHandler(hub *ws.Hub, tokenManager authentication.TokenManager) *WSHandler {
	return &WSHandler{hub: hub, tokenManager: tokenManager}
}

// HandleWS 处理 WebSocket 升级和连接生命周期
// userID 已由中间件鉴权提取
func (h *WSHandler) HandleWS(conn *websocket.Conn, userID string) {
	client := ws.NewClient(userID)

	// 加入用户房间
	roomID := "user:" + userID
	h.hub.JoinRoom(roomID, client)
	defer h.hub.LeaveRoom(roomID, client)
	defer conn.Close(websocket.StatusNormalClosure, "connection closed")

	logger.Info("WebSocket connected", "user_id", userID)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 读协程：处理 ping/pong 和客户端关闭
	go func() {
		defer cancel()
		for {
			_, _, err := conn.Read(ctx)
			if err != nil {
				return
			}
		}
	}()

	// 主协程：发送消息
	for {
		select {
		case <-ctx.Done():
			return
		case <-client.Done():
			return
		case msg, ok := <-client.Send:
			if !ok {
				return
			}
			writeCtx, writeCancel := context.WithTimeout(ctx, wsWriteTimeout)
			if err := conn.Write(writeCtx, websocket.MessageText, msg); err != nil {
				writeCancel()
				return
			}
			writeCancel()
		}
	}
}

// ServeWS Gin 处理器：从 query param 获取 JWT，升级为 WebSocket 连接
func (h *WSHandler) ServeWS(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	claims, err := h.tokenManager.ValidateAccessToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	conn, err := websocket.Accept(c.Writer, c.Request, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return
	}

	h.HandleWS(conn, claims.UserID)
}
