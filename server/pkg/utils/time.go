package utils

import (
	"time"

	"github.com/dromara/carbon/v2"
)

// Now 获取当前时间
func Now() time.Time {
	return carbon.Now().StdTime()
}

// NowRFC3339 获取当前时间的 RFC3339 格式字符串
func NowRFC3339() string {
	return carbon.Now().ToRfc3339String()
}

// ParseRFC3339 解析 RFC3339 格式时间字符串
func ParseRFC3339(s string) (time.Time, error) {
	c := carbon.Parse(s)
	if c.Error != nil {
		return time.Time{}, c.Error
	}
	return c.StdTime(), nil
}

// FormatRFC3339 将时间格式化为 RFC3339 字符串
func FormatRFC3339(t time.Time) string {
	return carbon.CreateFromStdTime(t).ToRfc3339String()
}
