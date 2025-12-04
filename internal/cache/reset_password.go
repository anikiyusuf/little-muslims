package cache

import(
	"context"
	"fmt"
	"time"
)

const (
	ResetPasswordKey = "reset_password_code"
	ResetPasswordTTL = time.Minute * 3
)

func (r *RedisCache) SetResetPasswordCode(ctx context.Context, key, value string) error {
	key = fmt.Sprintf("%s:%s", ResetPasswordKey, key)
	return r.Set(ctx, key, value, ResetPasswordTTL).Err()
}



func (r *RedisCache) GetResetPasswordCode(ctx context.Context, key string) (string, error) {
	key = fmt.Sprintf("%s:%s", ResetPasswordKey, key)
	return r.Get(ctx, key).Result()
}

func (r *RedisCache) IsTokenCorrect(ctx context.Context, token string) bool {
    storedToken, err := r.GetResetPasswordCode(ctx, token)
	if err != nil{
		return false
	}
	return token == storedToken
}

func (r *RedisCache) GetResetPasswordTTL(ctx context.Context, key string) (time.Duration, error) {
	key = fmt.Sprintf("%s:%s", ResetPasswordKey, key)
	return r.Client.TTL(ctx, key).Result()
}