package utils

import (
	"crypto/rand"
	"io"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

var (
	entropy io.Reader
	once    sync.Once
)

// init 初始化熵源
func init() {
	once.Do(func() {
		entropy = &lockedReader{r: rand.Reader}
	})
}

// lockedReader 线程安全的读取器
type lockedReader struct {
	r  io.Reader
	mu sync.Mutex
}

func (l *lockedReader) Read(p []byte) (n int, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.r.Read(p)
}

// GenerateID 生成 ULID ID
// 所有实体统一使用此方法生成唯一标识
func GenerateID() string {
	t := time.Now()
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id.String()
}

// ParseULID 解析 ULID 字符串
func ParseULID(id string) (ulid.ULID, error) {
	return ulid.Parse(id)
}

// GetTimestampFromID 从 ULID 中提取时间戳
func GetTimestampFromID(id string) (time.Time, error) {
	parsed, err := ParseULID(id)
	if err != nil {
		return time.Time{}, err
	}

	return time.UnixMilli(int64(parsed.Time())), nil
}
