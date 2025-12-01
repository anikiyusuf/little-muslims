package cache

import (
	"context"
	"fmt"
	"time"
)


const (
	VerificationCodeKey = "email_verification_code:"
	TTL 				= time.Minute * 10
)


func (r *RedisCache) SetVerificationCode(ctx context.Context, key, value string) error {
	key = fmt.Sprintf("%s:%s", VerificationCodeKey, key)
	return r.Set(ctx, key, value, TTL).Err()
} 

func (r *RedisCache) GetVerificationCode(ctx context.Context, key string) (string, error) {
	key = fmt.Sprintf("%s:%s", VerificationCodeKey, key)
	return r.Get(ctx, key).Result()
}