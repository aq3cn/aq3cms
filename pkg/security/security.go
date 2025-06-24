package security

import (
	"bytes"
	"crypto/hmac"
	cryptorand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/big"
	"math/rand"
	"strings"
	"time"
)

// GenerateToken 生成JWT令牌
func GenerateToken(claims map[string]interface{}, secret string, expireSeconds int) (string, error) {
	// 设置过期时间
	if expireSeconds <= 0 {
		expireSeconds = 3600 // 默认1小时
	}
	claims["exp"] = time.Now().Add(time.Duration(expireSeconds) * time.Second).Unix()
	claims["iat"] = time.Now().Unix()

	// 构建头部
	header := map[string]interface{}{
		"alg": "HS256",
		"typ": "JWT",
	}

	// 序列化头部和载荷
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	// Base64编码
	headerBase64 := base64.RawURLEncoding.EncodeToString(headerJSON)
	claimsBase64 := base64.RawURLEncoding.EncodeToString(claimsJSON)

	// 构建签名
	signatureInput := headerBase64 + "." + claimsBase64
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(signatureInput))
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	// 构建令牌
	token := headerBase64 + "." + claimsBase64 + "." + signature

	return token, nil
}

// ParseToken 解析JWT令牌
func ParseToken(token string, secret string) (map[string]interface{}, error) {
	// 分割令牌
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	// 解码载荷
	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	// 解析载荷
	var claims map[string]interface{}
	err = json.Unmarshal(claimsJSON, &claims)
	if err != nil {
		return nil, err
	}

	// 验证签名
	signatureInput := parts[0] + "." + parts[1]
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(signatureInput))
	expectedSignature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	if expectedSignature != parts[2] {
		return nil, fmt.Errorf("invalid token signature")
	}

	// 验证过期时间
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, fmt.Errorf("token expired")
		}
	}

	return claims, nil
}

// URLEncode URL编码
func URLEncode(str string) string {
	return base64.URLEncoding.EncodeToString([]byte(str))
}

// URLDecode URL解码
func URLDecode(str string) (string, error) {
	decoded, err := base64.URLEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// FilterXSS 过滤XSS攻击
func FilterXSS(content string) string {
	// 替换特殊字符
	content = strings.ReplaceAll(content, "<", "&lt;")
	content = strings.ReplaceAll(content, ">", "&gt;")
	content = strings.ReplaceAll(content, "\"", "&quot;")
	content = strings.ReplaceAll(content, "'", "&#39;")
	content = strings.ReplaceAll(content, "(", "&#40;")
	content = strings.ReplaceAll(content, ")", "&#41;")

	return content
}

// RandomString 生成随机字符串
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	charsetLength := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		n, err := cryptorand.Int(cryptorand.Reader, charsetLength)
		if err != nil {
			// 如果出错，使用时间戳作为备选方案
			result[i] = charset[time.Now().Nanosecond()%len(charset)]
			continue
		}
		result[i] = charset[n.Int64()]
	}

	return string(result)
}

// GenerateCaptcha 生成验证码
func GenerateCaptcha() (string, []byte, error) {
	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())
	// 生成随机验证码
	code := RandomString(4)

	// 创建图片
	width := 120
	height := 40
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 填充背景色
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.RGBA{uint8(rand.Intn(50) + 200), uint8(rand.Intn(50) + 200), uint8(rand.Intn(50) + 200), 255})
		}
	}

	// 添加干扰线
	for i := 0; i < 5; i++ {
		x1 := rand.Intn(width)
		y1 := rand.Intn(height)
		x2 := rand.Intn(width)
		y2 := rand.Intn(height)

		r := uint8(rand.Intn(255))
		g := uint8(rand.Intn(255))
		b := uint8(rand.Intn(255))

		for t := 0.0; t < 1.0; t += 0.01 {
			x := int(float64(x1) * (1.0 - t) + float64(x2) * t)
			y := int(float64(y1) * (1.0 - t) + float64(y2) * t)
			if x >= 0 && x < width && y >= 0 && y < height {
				img.Set(x, y, color.RGBA{r, g, b, 255})
			}
		}
	}

	// 添加文字
	for i, char := range code {
		size := 20 + rand.Intn(10)
		x := 10 + i*30 + rand.Intn(5)
		y := 20 + rand.Intn(10)

		r := uint8(rand.Intn(100))
		g := uint8(rand.Intn(100))
		b := uint8(rand.Intn(100))

		drawChar(img, x, y, char, size, color.RGBA{r, g, b, 255})
	}

	// 编码为PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", nil, err
	}

	return code, buf.Bytes(), nil
}

// drawChar 在图片上绘制字符
func drawChar(img *image.RGBA, x, y int, char rune, size int, col color.Color) {
	// 简单实现，实际应使用字体库
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			// 简单的字符形状
			if (i == 0 || i == size-1 || j == 0 || j == size-1 || i == j || i == size-j-1) && rand.Intn(3) != 0 {
				img.Set(x+i, y+j-size/2, col)
			}
		}
	}
}



