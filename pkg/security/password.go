package security

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// PasswordConfig 密码配置
type PasswordConfig struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// DefaultPasswordConfig 默认密码配置
var DefaultPasswordConfig = &PasswordConfig{
	Memory:      64 * 1024,
	Iterations:  3,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

// HashPassword 哈希密码 (兼容aq3cmsCMS的MD5格式)
func HashPassword(password string) string {
	// 为了兼容aq3cmsCMS数据库字段长度限制，使用MD5哈希
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

// HashPasswordArgon2 使用Argon2哈希密码 (更安全但字符串更长)
func HashPasswordArgon2(password string) string {
	return HashPasswordWithConfig(password, DefaultPasswordConfig)
}

// HashPasswordWithConfig 使用指定配置哈希密码
func HashPasswordWithConfig(password string, config *PasswordConfig) string {
	// 生成随机盐
	salt := make([]byte, config.SaltLength)
	_, err := rand.Read(salt)
	if err != nil {
		panic(err)
	}

	// 使用Argon2id哈希密码
	hash := argon2.IDKey([]byte(password), salt, config.Iterations, config.Memory, config.Parallelism, config.KeyLength)

	// 编码哈希和盐
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// 格式化哈希字符串
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, config.Memory, config.Iterations, config.Parallelism, b64Salt, b64Hash)

	return encodedHash
}

// CheckPassword 检查密码 (支持MD5和Argon2)
func CheckPassword(password, encodedHash string) bool {
	// 检查是否是MD5哈希 (32个十六进制字符)
	if len(encodedHash) == 32 && isHexString(encodedHash) {
		return CheckPasswordMD5(password, encodedHash)
	}

	// 检查是否是Argon2哈希
	if strings.HasPrefix(encodedHash, "$argon2id$") {
		return CheckPasswordArgon2(password, encodedHash)
	}

	// 如果都不匹配，尝试直接MD5比较（兼容性处理）
	expectedHash := HashPassword(password)
	return strings.ToLower(encodedHash) == strings.ToLower(expectedHash)
}

// CheckPasswordMD5 检查MD5密码
func CheckPasswordMD5(password, hash string) bool {
	expectedHash := HashPassword(password)
	return subtle.ConstantTimeCompare([]byte(strings.ToLower(hash)), []byte(strings.ToLower(expectedHash))) == 1
}

// isHexString 检查字符串是否为十六进制
func isHexString(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

// CheckPasswordArgon2 检查Argon2密码
func CheckPasswordArgon2(password, encodedHash string) bool {
	// 解析哈希字符串
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false
	}

	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return false
	}

	var memory, iterations uint32
	var parallelism uint8
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}

	// 使用相同的参数哈希密码
	keyLength := uint32(len(hash))
	newHash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLength)

	// 比较哈希
	return subtle.ConstantTimeCompare(hash, newHash) == 1
}

// UpgradePasswordHash 升级密码哈希
func UpgradePasswordHash(password, encodedHash string) (string, bool) {
	// 检查密码是否正确
	if !CheckPassword(password, encodedHash) {
		return "", false
	}

	// 解析哈希字符串
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return "", false
	}

	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return "", false
	}

	var memory, iterations uint32
	var parallelism uint8
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return "", false
	}

	// 检查是否需要升级
	if version == argon2.Version &&
		memory >= DefaultPasswordConfig.Memory &&
		iterations >= DefaultPasswordConfig.Iterations &&
		parallelism >= DefaultPasswordConfig.Parallelism {
		return encodedHash, false
	}

	// 使用新配置哈希密码
	newHash := HashPassword(password)
	return newHash, true
}
