package auth

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// 用户名检查, 只支持用户名，不支持邮箱
func CheckUsername(username string) bool {
	// 只支持用户名，不支持邮箱
	if strings.Contains(username, "@") {
		return false
	}
	// 用户名长度大于0小于20
	if len(username) > 20 {
		return false
	}
	return true
}

// 加密密码
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// 检查密码
func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

