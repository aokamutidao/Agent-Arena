package user

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-jwt/jwt/v5"
)

// JWT 配置
var jwtSecret []byte

func init() {
	// 从环境变量读取 JWT secret，默认使用硬编码值（开发环境）
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "agent-arena-dev-secret-key-change-in-production"
	}
	jwtSecret = []byte(secret)
}

// Claims JWT payload
type Claims struct {
	UserID  string `json:"user_id"`
	Address string `json:"address"`
	jwt.RegisteredClaims
}

// LoginRequest 登录请求
type LoginRequest struct {
	Address   string `json:"address"`
	Message   string `json:"message"`
	Signature string `json:"signature"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

// VerifySignature 验证钱包签名
// message: 签名的原始消息
// signature: 签名（hex 字符串，0x 开头）
// expectedAddress: 期望的地址
func VerifySignature(message, signature, expectedAddress string) (bool, error) {
	// 1. 解码签名
	sig := common.FromHex(signature)
	if len(sig) != 65 {
		return false, fmt.Errorf("invalid signature length: %d", len(sig))
	}

	// 2. 恢复 v 值（以太坊签名 v = 27 or 28）
	if sig[64] >= 27 {
		sig[64] -= 27
	}

	// 3. 创建消息哈希（以太坊签名格式）
	prefixedMsg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	hash := crypto.Keccak256Hash([]byte(prefixedMsg))

	// 4. 恢复公钥
	pubKey, err := crypto.SigToPub(hash.Bytes(), sig)
	if err != nil {
		return false, fmt.Errorf("failed to recover public key: %w", err)
	}

	// 5. 从公钥派生地址
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	// 6. 比较地址
	expectedAddr := common.HexToAddress(expectedAddress)
	return recoveredAddr == expectedAddr, nil
}

// GenerateJWT 生成 JWT token
func GenerateJWT(userID, address string) (string, error) {
	claims := &Claims{
		UserID:  userID,
		Address: strings.ToLower(address),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 天
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "agent-arena",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateJWT 验证 JWT token
func ValidateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// RecoverAddressFromSignature 从签名恢复地址（不验证期望地址）
func RecoverAddressFromSignature(message, signature string) (string, error) {
	// 1. 解码签名
	sig := common.FromHex(signature)
	if len(sig) != 65 {
		return "", fmt.Errorf("invalid signature length: %d", len(sig))
	}

	// 2. 恢复 v 值
	if sig[64] >= 27 {
		sig[64] -= 27
	}

	// 3. 创建消息哈希
	prefixedMsg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	hash := crypto.Keccak256Hash([]byte(prefixedMsg))

	// 4. 恢复公钥
	pubKey, err := crypto.SigToPub(hash.Bytes(), sig)
	if err != nil {
		return "", fmt.Errorf("failed to recover public key: %w", err)
	}

	// 5. 派生地址
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	return recoveredAddr.Hex(), nil
}

// GenerateLoginMessage 生成登录消息
func GenerateLoginMessage(address string, timestamp int64) string {
	return fmt.Sprintf("Login to Agent Arena: %s at %d", address, timestamp)
}

// PublicKeyToAddress 从公钥派生地址（辅助函数）
func PublicKeyToAddress(pubKey *ecdsa.PublicKey) string {
	return crypto.PubkeyToAddress(*pubKey).Hex()
}
