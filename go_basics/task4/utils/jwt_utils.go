package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// 自定义声明结构，可添加需要的用户信息
type CustomClaims struct {
	UserID   int64  `json:"userId"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// 密钥（务必保密，可通过环境变量配置）
var jwtSecret = []byte("123456") // 生产环境应使用强密钥并从安全配置中读取

// GenerateToken 生成 JWT Token
func GenerateToken(userID int64, username string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(72 * time.Hour) // 令牌过期时间

	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "task4",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// ParseToken 解析和验证 JWT Token
func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
