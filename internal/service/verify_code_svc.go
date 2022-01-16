package service

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/redis"
	"github.com/pkg/errors"
)

const (
	verifyCodeRedisKey = "user_svc:verify:code:%s" // 验证码key
	maxDurationTime    = 10 * time.Minute          // 验证码有效期
)

// phone white list
var phoneWhiteLit = []string{
	"13000001111",
}

// VerifyCodeService define a interface
type VerifyCodeService interface {
	GenVerifyCode(ctx context.Context, phone string) (int, error)
	CheckVerifyCode(ctx context.Context, phone string, vCode string) bool
	GetVerifyCode(ctx context.Context, phone string) (vCode string, err error)
}

type verifyCodeService struct {
}

var _ VerifyCodeService = (*verifyCodeService)(nil)

func newVerifyCodeService() *verifyCodeService {
	return &verifyCodeService{}
}

// GenLoginVCode 生成校验码
func (s *verifyCodeService) GenVerifyCode(ctx context.Context, phone string) (int, error) {
	// step1: 生成随机数
	vCodeStr := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))

	// step2: 写入到redis里
	// 使用set, key使用前缀+手机号 缓存10分钟）
	key := fmt.Sprintf(verifyCodeRedisKey, phone)
	err := redis.RedisClient.Set(ctx, key, vCodeStr, maxDurationTime).Err()
	if err != nil {
		return 0, errors.Wrap(err, "[vcode_svc] gen verify code from redis set err")
	}

	vCode, err := strconv.Atoi(vCodeStr)
	if err != nil {
		return 0, errors.Wrap(err, "[vcode_svc] string convert int err")
	}

	return vCode, nil
}

// isTestPhone add test phone number to here and avoid to verify
func isTestPhone(phone string) bool {
	for _, val := range phoneWhiteLit {
		if val == phone {
			return true
		}
	}
	return false
}

// CheckLoginVCode check phone if correct
func (s *verifyCodeService) CheckVerifyCode(ctx context.Context, phone string, vCode string) bool {
	if isTestPhone(phone) {
		return true
	}

	oldVCode, err := s.GetVerifyCode(ctx, phone)
	if err != nil {
		log.Warnf("[vcode_svc] get verify code err, %v", err)
		return false
	}

	if vCode != oldVCode {
		return false
	}

	return true
}

// GetLoginVCode get verify code
func (s *verifyCodeService) GetVerifyCode(ctx context.Context, phone string) (vCode string, err error) {
	key := fmt.Sprintf(verifyCodeRedisKey, phone)
	vCode, err = redis.RedisClient.Get(ctx, key).Result()
	if err == redis.ErrRedisNotFound {
		return
	} else if err != nil {
		return "", errors.Wrap(err, "[vcode_svc] redis get verify code err")
	}

	return vCode, nil
}
