package utils

import "github.com/spf13/cast"

// ToString 转换为字符串
func ToString(i interface{}) string {
	return cast.ToString(i)
}

// ToInt 转换为整数
func ToInt(i interface{}) int {
	return cast.ToInt(i)
}

// ToIntE 转换为整数（带错误返回）
func ToIntE(i interface{}) (int, error) {
	return cast.ToIntE(i)
}
