package repository

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shenfay/go-react-admin/pkg/utils"
)

// TimeNull 可空的时间类型
// 用于 GORM 持久化对象中表示可为 NULL 的时间字段
type TimeNull struct {
	Time  time.Time
	Valid bool
}

// Value 实现 driver.Valuer 接口
func (t TimeNull) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

// Scan 实现 sql.Scanner 接口
func (t *TimeNull) Scan(value interface{}) error {
	if value == nil {
		t.Time = time.Time{}
		t.Valid = false
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		t.Time = v
		t.Valid = true
		return nil
	default:
		return fmt.Errorf("failed to scan TimeNull: %v", value)
	}
}

// MarshalJSON 实现 JSON 序列化
func (t TimeNull) MarshalJSON() ([]byte, error) {
	if t.Valid {
		return json.Marshal(utils.FormatRFC3339(t.Time))
	}
	return json.Marshal(nil)
}

// UnmarshalJSON 实现 JSON 反序列化
func (t *TimeNull) UnmarshalJSON(data []byte) error {
	var s interface{}
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if s == nil {
		t.Valid = false
		return nil
	}

	str, ok := s.(string)
	if !ok {
		return fmt.Errorf("invalid time value")
	}

	parsed, err := utils.ParseRFC3339(str)
	if err != nil {
		return err
	}

	t.Time = parsed
	t.Valid = true
	return nil
}
