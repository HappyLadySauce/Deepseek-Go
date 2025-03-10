package email

import (
	"regexp"
	"strings"
)

// 检查邮箱格式
func CheckEmailFormat(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// 允许的邮箱：QQ,163,gmail，并限制长度
func AllowEmailFormat(email string) bool {
	// 基本格式检查
	if !CheckEmailFormat(email) {
		return false
	}

	// 分割邮箱获取用户名和域名部分
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	username := parts[0]
	domain := parts[1]

	// 根据不同域名检查用户名长度
	switch domain {
	case "qq.com":
		// QQ邮箱用户名为纯数字，长度为5-11位
		qqRegex := regexp.MustCompile(`^[1-9][0-9]{4,10}$`)
		return qqRegex.MatchString(username)
	case "163.com":
		// 163邮箱用户名长度为6-18位
		return len(username) >= 6 && len(username) <= 18
	case "gmail.com":
		// Gmail用户名部分最多30个字符
		return len(username) <= 30
	default:
		return false
	}
}
