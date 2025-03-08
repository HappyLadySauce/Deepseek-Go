package email

import (
	"regexp"
)

// 检查邮箱格式
func CheckEmailFormat(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// 允许的邮箱：QQ,163,gmail
func AllowEmailFormat(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@(qq\.com|163\.com|gmail\.com)$`)
	return emailRegex.MatchString(email)
}