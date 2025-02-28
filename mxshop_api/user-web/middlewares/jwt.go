package middlewares

import (
	"fmt"
	"mxshop_api/user-web/global"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims 包含自定义的JWT声明信息。
// 你可以将此结构体放入model包中以便更好地组织代码。
type CustomClaims struct {
	ID                   int
	NickName             string
	Role                 int //1是普通用户，2是管理员
	jwt.RegisteredClaims     // JWT标准声明（如过期时间、签发者等）
}

// JWTManager 负责生成和验证JWT token。
type JWTManager struct {
	secretKey     string        // 用于签名token的密钥
	tokenDuration time.Duration // token的有效期
}

// NewJWTManager 创建一个新的JWTManager实例。
func NewJWTManager() *JWTManager {
	return &JWTManager{
		secretKey:     global.ServerConfig.JWTConfig.SigningKey,
		tokenDuration: 30 * 24 * time.Hour,
	}
}

// GenerateToken 使用指定的用户ID和用户名生成一个JWT token。
func (jm *JWTManager) GenerateToken(userID int, nickName string, role int) (string, error) {
	claims := CustomClaims{
		ID:       userID,
		NickName: nickName,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jm.tokenDuration)), // 设置token的过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                       // 设置token的签发时间
			Issuer:    "microservice-system",                                // 签发者
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jm.secretKey)) // 返回签名后的token字符串
}

// VerifyToken 验证给定的token字符串是否有效，并返回解析出的声明信息。
func (jm *JWTManager) VerifyToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jm.secretKey), nil // 返回用于验证token的密钥
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// GinMiddleware 提供了一个Gin中间件来保护路由，确保请求携带有效的JWT token。
func (jm *JWTManager) JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从HTTP头部获取Authorization字段
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			return
		}

		// 检查并分离Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
			})
			return
		}

		tokenString := parts[1]
		claims, err := jm.VerifyToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token: " + err.Error(),
			})
			return
		}

		// 将解析出的声明存储到Gin上下文中，以便后续处理逻辑使用
		c.Set("jwt_claims", claims)
		c.Next()
	}
}

// GetClaimsFromGinContext 是一个辅助函数，用于从Gin上下文中提取JWT声明。
func GetClaimsFromGinContext(c *gin.Context) (*CustomClaims, error) {
	claims, exists := c.Get("jwt_claims") // 从上下文中获取之前存储的claims
	if !exists {
		return nil, fmt.Errorf("claims not found in context")
	}

	customClaims, ok := claims.(*CustomClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims type")
	}

	return customClaims, nil
}
