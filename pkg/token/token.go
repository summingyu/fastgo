package token

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
)

type Config struct {
	key         string
	identityKey string
	expiration  time.Duration
}

var (
	config = Config{"Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5", "identityKey", 2 * time.Hour}
	once   sync.Once
)

func Init(key string, identityKey string, expiration time.Duration) {
	once.Do(func() {
		if key != "" {
			config.key = key
		}
		if identityKey != "" {
			config.identityKey = identityKey
		}
		if expiration != 0 {
			config.expiration = expiration
		}
	})
}

func Parse(tokenString string, key string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			slog.Error("Parse token", "error", "unexpected signing method")
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(key), nil
	})
	if err != nil {
		return "", err
	}
	var identityKey string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if key, exists := claims[config.identityKey]; exists {
			if identity, valid := key.(string); valid {
				identityKey = identity
			}
		}
	}
	if identityKey == "" {
		return "", jwt.ErrSignatureInvalid
	}

	return identityKey, nil
}

// ParseRequest 从 Gin 上下文的 Authorization 请求头中解析出身份标识。
// 该函数会检查请求头是否存在且长度不为零。
// 参数 c 是 Gin 上下文对象，包含了请求的相关信息。
// 返回值为解析出的身份标识字符串和可能出现的错误。
func ParseRequest(c *gin.Context) (string, error) {
	// 从请求头中获取 Authorization 字段的值
	header := c.Request.Header.Get("Authorization")

	// 检查 Authorization 头的长度是否为 0
	if len(header) == 0 {
		// 若长度为 0，返回错误信息
		return "", errors.New("the length of the `Authorization` header is zero")
	}

	// 使用 fmt.Sscanf 从 Authorization 头中提取 token
	var token string
	fmt.Sscanf(header, "Bearer %s", &token)

	// 调用 Parse 函数解析 token
	return Parse(token, config.key)
}

// Sign 生成 JWT 令牌（签名）
// 参数 identityKey 是需要存储在令牌中的身份标识（如用户ID）
// 返回值依次为：生成的令牌字符串、令牌过期时间、可能出现的错误
func Sign(identityKey string) (string, time.Time, error) {
	// 计算令牌过期时间（当前时间 + 配置的过期时长）
	expireAt := time.Now().Add(config.expiration)

	// 创建 JWT 令牌并设置声明（Claims）
	// 使用 HS256 算法签名，包含以下标准和自定义声明：
	// - config.identityKey: 自定义身份标识字段（由配置指定字段名）
	// - nbf（Not Before）: 令牌生效时间（当前时间戳）
	// - iat（Issued At）: 令牌签发时间（当前时间戳）
	// - exp（Expires At）: 令牌过期时间（当前时间 + 配置时长的时间戳）
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		config.identityKey: identityKey,
		"nbf":              time.Now().Unix(),
		"iat":              time.Now().Unix(),
		"exp":              expireAt.Unix(),
	})

	// 检查配置的签名密钥是否为空（空密钥无法生成有效令牌）
	if config.key == "" {
		return "", time.Time{}, jwt.ErrInvalidKey
	}

	// 使用配置的密钥对令牌进行签名，生成最终的令牌字符串
	tokenString, err := token.SignedString([]byte(config.key))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expireAt, nil
}
