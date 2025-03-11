package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 生成token
func GenerateToken(username string) (string, error) {
	// JWT令牌生成函数，传入用户名，生成其JWT令牌.
	// Header 为 HS256，JWT; Payload 为 username; Signature 为 "deepseek-chat"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("deepseek-chat"))
	if err != nil {
		return "", err
	}

	return "Bearer " + tokenString, nil
}

// 验证token
func ValidateToken(tokenString string) (string, error) {
	// 去除可能的引号
	tokenString = strings.Trim(tokenString, "\"")

	// 去掉 Token 的 "Bearer " 前缀。
	if len(tokenString) > 7 && strings.HasPrefix(strings.ToLower(tokenString), "bearer ") {
		tokenString = tokenString[7:]
	}

	// 去除令牌两端的空格
	tokenString = strings.TrimSpace(tokenString)

	// 验证令牌是否为空
	if tokenString == "" {
		return "", errors.New("empty token")
	}

	// 验证 JWT 的签名方法是否为 HMAC。
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证 JWT 的签名方法是否为 HMAC。
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("deepseek-chat"), nil
	})

	// 提取 Claims 中的用户名字段。
	// 如果解析失败或用户名字段无效，返回错误。
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", errors.New("invalid token")
	}

	return username, nil
}
