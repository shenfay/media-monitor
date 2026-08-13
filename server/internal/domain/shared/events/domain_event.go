package events

import (
	"context"
	"time"
)

// DomainEvent 领域事件接口
// 所有领域事件需实现此接口
type DomainEvent interface {
	// EventName 返回事件名称（小写点分格式，如 "user.registered"）
	EventName() string

	// OccurredAt 返回事件发生时间
	OccurredAt() time.Time
}

// Handler 领域事件处理器
type Handler func(ctx context.Context, event DomainEvent) error

// Bus 进程内事件总线接口
// 用于在同一个进程中发布和订阅领域事件
type Bus interface {
	// Publish 发布领域事件（同步调用所有已注册的 Handler）
	Publish(ctx context.Context, event DomainEvent) error

	// Subscribe 订阅指定类型的领域事件
	Subscribe(eventName string, handler Handler)
}
