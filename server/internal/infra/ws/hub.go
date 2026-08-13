package ws

import (
	"sync"
)

// Client WebSocket 客户端连接
type Client struct {
	UserID string
	Send   chan []byte
	done   chan struct{}
}

// NewClient 创建客户端
func NewClient(userID string) *Client {
	return &Client{
		UserID: userID,
		Send:   make(chan []byte, 64),
		done:   make(chan struct{}),
	}
}

// Close 关闭客户端
func (c *Client) Close() {
	select {
	case <-c.done:
	default:
		close(c.done)
	}
}

// Done 返回关闭通知 channel
func (c *Client) Done() <-chan struct{} {
	return c.done
}

// Hub 广播中心
// 管理按用户分组的 WebSocket 连接，支持向特定用户发送消息
type Hub struct {
	// rooms: roomID -> set of clients
	rooms map[string]map[*Client]bool
	mu    sync.RWMutex
}

// NewHub 创建广播中心
func NewHub() *Hub {
	return &Hub{
		rooms: make(map[string]map[*Client]bool),
	}
}

// JoinRoom 将客户端加入房间
func (h *Hub) JoinRoom(roomID string, c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[roomID] == nil {
		h.rooms[roomID] = make(map[*Client]bool)
	}
	h.rooms[roomID][c] = true
}

// LeaveRoom 将客户端移出房间
func (h *Hub) LeaveRoom(roomID string, c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.rooms[roomID]; ok {
		delete(clients, c)
		if len(clients) == 0 {
			delete(h.rooms, roomID)
		}
	}
}

// SendToUser 向指定用户房间发送消息（非阻塞）
func (h *Hub) SendToUser(userID string, msg []byte) {
	roomID := userRoomID(userID)

	h.mu.RLock()
	clients := h.rooms[roomID]
	h.mu.RUnlock()

	for c := range clients {
		select {
		case c.Send <- msg:
		default:
			// 发送缓冲区满，跳过
		}
	}
}

// UserRoomID 用户房间 ID
func userRoomID(userID string) string {
	return "user:" + userID
}

// Shutdown 关闭所有客户端连接，触发 HandleWS 退出
func (h *Hub) Shutdown() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, clients := range h.rooms {
		for c := range clients {
			c.Close()
		}
	}
}
