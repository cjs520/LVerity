package service

import (
	"bytes"
	"encoding/base64"
	"errors"
	"github.com/dchest/captcha"
	"image/png"
	"log"
	"sync"
	"time"
)

// CaptchaService 定义验证码服务接口
type CaptchaService interface {
	Generate() (string, string, error)
	Verify(id string, answer string) bool
}

// captchaServiceImpl 实现验证码服务
type captchaServiceImpl struct {
	store     CaptchaStore
	width     int
	height    int
	length    int
	logPrefix string
}

// CaptchaStore 定义验证码存储接口
type CaptchaStore interface {
	Set(id string, digits []byte) bool
	Get(id string, clear bool) []byte
}

// memoryStore 内存存储实现
type memoryStore struct {
	sync.RWMutex
	data       map[string]captchaData
	capacity   int
	expiration time.Duration
}

type captchaData struct {
	digits    []byte
	timestamp time.Time
}

// NewCaptchaService 创建新的验证码服务实例
func NewCaptchaService() CaptchaService {
	store := newMemoryStore(100, 300*time.Second)
	return &captchaServiceImpl{
		store:     store,
		width:     240,
		height:    80,
		length:    6,
		logPrefix: "[验证码服务] ",
	}
}

// newMemoryStore 创建新的内存存储实例
func newMemoryStore(capacity int, expiration time.Duration) *memoryStore {
	return &memoryStore{
		data:       make(map[string]captchaData),
		capacity:   capacity,
		expiration: expiration,
	}
}

// Set 存储验证码
func (s *memoryStore) Set(id string, digits []byte) bool {
	s.Lock()
	defer s.Unlock()

	// 检查容量
	if len(s.data) >= s.capacity {
		// 清理过期数据
		s.cleanup()
		// 如果仍然超出容量，返回失败
		if len(s.data) >= s.capacity {
			return false
		}
	}

	s.data[id] = captchaData{
		digits:    digits,
		timestamp: time.Now(),
	}
	return true
}

// Get 获取验证码
func (s *memoryStore) Get(id string, clear bool) []byte {
	s.RLock()
	data, exists := s.data[id]
	s.RUnlock()

	if !exists || time.Since(data.timestamp) > s.expiration {
		return nil
	}

	if clear {
		s.Lock()
		delete(s.data, id)
		s.Unlock()
	}

	return data.digits
}

// cleanup 清理过期数据
func (s *memoryStore) cleanup() {
	now := time.Now()
	for id, data := range s.data {
		if now.Sub(data.timestamp) > s.expiration {
			delete(s.data, id)
		}
	}
}

// Generate 生成验证码
func (s *captchaServiceImpl) Generate() (string, string, error) {
	// 生成验证码ID和数字
	id := captcha.NewLen(s.length)
	digits := captcha.RandomDigits(s.length)

	// 生成图片
	var buf bytes.Buffer
	img := captcha.NewImage(id, digits, s.width, s.height)
	if err := png.Encode(&buf, img); err != nil {
		log.Printf(s.logPrefix+"生成验证码图片失败: %v", err)
		return "", "", errors.New("生成验证码图片失败")
	}

	// 存储验证码数字
	if !s.store.Set(id, digits) {
		log.Printf(s.logPrefix+"存储验证码失败，ID: %s", id)
		return "", "", errors.New("验证码存储失败")
	}

	// 记录生成的验证码信息（仅用于调试）
	expected := make([]byte, len(digits))
	for i := 0; i < len(digits); i++ {
		expected[i] = digits[i] + '0'
	}
	log.Printf(s.logPrefix+"生成验证码 - ID: %s, 值: %s", id, expected)

	// 转换为base64
	b64s := base64.StdEncoding.EncodeToString(buf.Bytes())
	return id, "data:image/png;base64," + b64s, nil
}

// Verify 验证验证码
func (s *captchaServiceImpl) Verify(id string, answer string) bool {
	// 预处理验证码答案
	answer = string(bytes.TrimSpace([]byte(answer)))

	// 获取存储的验证码
	digits := s.store.Get(id, true) // 验证后立即删除
	if digits == nil {
		log.Printf(s.logPrefix+"验证码不存在或已过期，ID: %s", id)
		return false
	}

	// 验证长度
	if len(digits) != len(answer) {
		log.Printf(s.logPrefix+"验证码长度不匹配: 期望 %d, 实际 %d", len(digits), len(answer))
		return false
	}

	// 验证字符
	for i := 0; i < len(answer); i++ {
		if answer[i] < '0' || answer[i] > '9' {
			log.Printf(s.logPrefix+"验证码包含非数字字符: %c", answer[i])
			return false
		}
		if digits[i] != answer[i]-'0' {
			log.Printf(s.logPrefix+"验证码不匹配: 位置 %d, 期望 %d, 实际 %d", i, digits[i], answer[i]-'0')
			return false
		}
	}

	log.Printf(s.logPrefix+"验证码验证成功，ID: %s", id)
	return true
}

// 全局验证码服务实例
var defaultCaptchaService CaptchaService

func init() {
	defaultCaptchaService = NewCaptchaService()
}

// GenerateCaptcha 使用默认服务生成验证码
func GenerateCaptcha() (string, string, error) {
	return defaultCaptchaService.Generate()
}

// VerifyCaptcha 使用默认服务验证验证码
func VerifyCaptcha(id string, answer string) bool {
	return defaultCaptchaService.Verify(id, answer)
}
