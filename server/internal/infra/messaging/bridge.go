package messaging

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/shenfay/go-react-admin/internal/domain/shared/events"
	"github.com/shenfay/go-react-admin/pkg/constants"
)

// DomainToIntegrationBridge 领域事件到集成事件的桥接器
// 订阅 InProcessBus 中的领域事件，转换为 Asynq 集成任务入队
type DomainToIntegrationBridge struct {
	client   *asynq.Client
	queueMap map[constants.EventName]constants.QueueName
}

// NewBridge 创建桥接器
func NewBridge(client *asynq.Client) *DomainToIntegrationBridge {
	return &DomainToIntegrationBridge{
		client:   client,
		queueMap: logEventQueueMap,
	}
}

// SubscribeTo 订阅 InProcessBus 中的所有领域事件
func (b *DomainToIntegrationBridge) SubscribeTo(bus events.Bus) {
	if b.client == nil {
		return // nil client 模式仅用于事件类型发现，不执行订阅
	}
	for eventName := range b.queueMap {
		name := eventName // capture for closure
		bus.Subscribe(string(name), func(ctx context.Context, evt events.DomainEvent) error {
			return b.enqueue(ctx, name, evt)
		})
	}
}

// LogEventTypes 返回路由到 logs 队列的所有事件类型
// 作为事件注册表的单一真相来源，供 Worker 进程统一注册
func (b *DomainToIntegrationBridge) LogEventTypes() []constants.EventName {
	return LogEventTypes()
}

// LogEventTypes 包级函数：返回路由到 logs 队列的所有事件类型
// 无需 Bridge 实例即可获取事件注册表，供 Worker 进程独立使用
func LogEventTypes() []constants.EventName {
	types := make([]constants.EventName, 0, len(logEventQueueMap))
	for eventName, queue := range logEventQueueMap {
		if queue == constants.QueueLogs {
			types = append(types, eventName)
		}
	}
	return types
}

// logEventQueueMap 事件到队列的路由映射表（包级变量，作为单一真相来源）
var logEventQueueMap = map[constants.EventName]constants.QueueName{
	// 统一操作日志事件路由到 logs 队列
	constants.EventOperationLog: constants.QueueLogs,
	// 发送邮件事件路由到 notification 队列
	constants.EventSendEmail: constants.QueueNotification,
}

// enqueue 将领域事件入队到 Asynq
func (b *DomainToIntegrationBridge) enqueue(ctx context.Context, eventName constants.EventName, evt events.DomainEvent) error {
	if b.client == nil {
		return nil
	}
	payload, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	queue := b.queueMap[eventName]
	maxRetry := 3
	if queue == constants.QueueCritical {
		maxRetry = 5
	}

	_, err = b.client.EnqueueContext(ctx,
		asynq.NewTask(string(eventName), payload),
		asynq.Queue(string(queue)),
		asynq.MaxRetry(maxRetry),
	)
	return err
}
