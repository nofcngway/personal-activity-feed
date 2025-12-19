package sessions

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	rdb *redis.Client
}

type Session struct {
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func NewRedisStore(addr, password string, db int) *RedisStore {
	return &RedisStore{
		rdb: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
	}
}

func (s *RedisStore) Close() error { return s.rdb.Close() }

func (s *RedisStore) Set(ctx context.Context, token string, session Session, ttl time.Duration) error {
	key := sessionKey(token)
	b, err := json.Marshal(session)
	if err != nil {
		return err
	}
	return s.rdb.Set(ctx, key, b, ttl).Err()
}

func (s *RedisStore) Get(ctx context.Context, token string) (*Session, error) {
	key := sessionKey(token)
	val, err := s.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	var sess Session
	if err := json.Unmarshal(val, &sess); err != nil {
		return nil, err
	}
	return &sess, nil
}

func (s *RedisStore) Del(ctx context.Context, token string) error {
	return s.rdb.Del(ctx, sessionKey(token)).Err()
}

func sessionKey(token string) string {
	return fmt.Sprintf("session:%s", token)
}


