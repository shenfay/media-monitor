package shared

import (
	"time"

	"github.com/shenfay/go-react-admin/pkg/utils"
)

// BaseEvent 领域事件基类（可选）
// Go语言中通常使用接口而非继承，这里提供作为参考
type BaseEvent struct {
	EventID   string    `json:"event_id"`
	Timestamp time.Time `json:"timestamp"`
}

// NewBaseEvent 创建基础事件
func NewBaseEvent() BaseEvent {
	return BaseEvent{
		Timestamp: utils.Now(),
	}
}
