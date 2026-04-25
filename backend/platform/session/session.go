package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrSessionNotFound  = errors.New("session not found")
	ErrInvalidSessionID = errors.New("session id must not be empty")
	ErrInvalidTTL       = errors.New("ttl must be greater than zero")
)

type SessionData struct {
	UserID    string    `json:"user_id"`
	Role      string    `json:"role,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Store interface {
	Create(ctx context.Context, sessionID string, data SessionData, ttl time.Duration) error
	Get(ctx context.Context, sessionID string) (SessionData, error)
	Exists(ctx context.Context, sessionID string) (bool, error)
	Delete(ctx context.Context, sessionID string) error
}

type RedisStore struct {
	client redis.Cmdable
	prefix string
}

func NewRedisStore(client redis.Cmdable) *RedisStore {
	return &RedisStore{
		client: client,
		prefix: "session:",
	}
}

func (s *RedisStore) Create(ctx context.Context, sessionID string, data SessionData, ttl time.Duration) error {
	if sessionID == "" {
		return ErrInvalidSessionID
	}
	if ttl <= 0 {
		return ErrInvalidTTL
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal session data: %w", err)
	}

	if err := s.client.Set(ctx, s.key(sessionID), payload, ttl).Err(); err != nil {
		return fmt.Errorf("set session in redis: %w", err)
	}

	return nil
}

func (s *RedisStore) Get(ctx context.Context, sessionID string) (SessionData, error) {
	if sessionID == "" {
		return SessionData{}, ErrInvalidSessionID
	}

	raw, err := s.client.Get(ctx, s.key(sessionID)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return SessionData{}, ErrSessionNotFound
		}
		return SessionData{}, fmt.Errorf("get session from redis: %w", err)
	}

	var data SessionData
	if err := json.Unmarshal(raw, &data); err != nil {
		return SessionData{}, fmt.Errorf("unmarshal session data: %w", err)
	}

	return data, nil
}

func (s *RedisStore) Exists(ctx context.Context, sessionID string) (bool, error) {
	if sessionID == "" {
		return false, ErrInvalidSessionID
	}

	n, err := s.client.Exists(ctx, s.key(sessionID)).Result()
	if err != nil {
		return false, fmt.Errorf("check session existence in redis: %w", err)
	}

	return n > 0, nil
}

func (s *RedisStore) Delete(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return ErrInvalidSessionID
	}

	if err := s.client.Del(ctx, s.key(sessionID)).Err(); err != nil {
		return fmt.Errorf("delete session from redis: %w", err)
	}

	return nil
}

func (s *RedisStore) key(sessionID string) string {
	return s.prefix + sessionID
}
