package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// 生成JWT令牌
func GenerateJWT(username string) (string, error) {
	// 创建JWT令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})	

	// 生成JWT令牌
	return token.SignedString([]byte("secret"))
}	


// 验证JWT令牌
func VerifyJWT(tokenString string) (string, error) {
	// 解析JWT令牌
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	// 验证JWT令牌
	if err != nil {
		return "", err
	}

	// 获取JWT令牌中的用户名
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("无法获取JWT令牌中的用户名")
	}

	// 返回JWT令牌中的用户名
	return claims["username"].(string), nil
}
