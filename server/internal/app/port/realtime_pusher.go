package port

// RealtimePusher 实时推送出站端口
// 由 app 层定义，infra 层实现（如 ws.Hub）
type RealtimePusher interface {
	// SendToUser 向指定用户推送消息
	SendToUser(userID string, msg []byte)
}
