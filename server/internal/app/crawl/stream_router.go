package crawl

import (
	"sort"
	"strings"
)

// ComputeStreamName 根据 tags 计算目标 Stream 名称。
//
//	tags 为空             → "crawl:task:dispatch"
//	tags=["overseas"]     → "crawl:task:dispatch:overseas"
//	tags=["overseas","a"] → "crawl:task:dispatch:a:overseas"（字母排序拼接）
func ComputeStreamName(baseStream string, tags []string) string {
	if len(tags) == 0 {
		return baseStream
	}
	// 去重 + 排序
	seen := make(map[string]struct{}, len(tags))
	unique := make([]string, 0, len(tags))
	for _, t := range tags {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		if _, ok := seen[t]; !ok {
			seen[t] = struct{}{}
			unique = append(unique, t)
		}
	}
	if len(unique) == 0 {
		return baseStream
	}
	sort.Strings(unique)
	return baseStream + ":" + strings.Join(unique, ":")
}
