package redis

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache - redis cache wrap.
type Cache struct {
	conn   *redis.Client
	logger *slog.Logger
}

var (
	// ErrConnection - error with the connection for any reason
	ErrConnection = errors.New("connection error")
)

// New - create a new Cache object.
func New(ctx context.Context, config Config, opts ...Option) (*Cache, error) {
	var cache Cache

	for _, opt := range opts {
		if e := opt(&cache); e != nil {
			return nil, cache.wrapError(e)
		}
	}

	c := redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		DB:           config.DB,
		Username:     config.User,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.Timeout,
		WriteTimeout: config.Timeout,
	})

	if cache.logger != nil {
		cache.logger.LogAttrs(
			nil,
			slog.LevelInfo,
			"[redis/New]",
			slog.String("Connected to", config.Addr),
			slog.String("Connected as", config.User),
			slog.Int("Connected DB", config.DB),
		)
	}

	if err := c.Ping(ctx).Err(); err != nil {
		if cache.logger != nil {
			cache.logger.LogAttrs(
				nil,
				slog.LevelError,
				"[redis/New]",
				slog.String("Ping Error", err.Error()),
			)
		}
		return nil, fmt.Errorf("%w: redis is unavailable", ErrConnection)
	}

	cache.conn = c
	return &cache, nil
}

// Set - write a pair key-value in cache 'c'.
// ttl same as in redis.
func (c *Cache) Set(
	ctx context.Context,
	key string,
	value string,
	ttl time.Duration,
) error {
	return c.wrapError(c.conn.Set(ctx, key, value, ttl).Err())
}

// Get - get stored value from cache 'c'.
func (c *Cache) Get(
	ctx context.Context,
	key string,
) (
	string,
	error,
) {
	r := c.conn.Get(ctx, key)
	res, err := r.Result()

	if err != nil {
		return "", c.wrapError(err)
	}

	return res, nil
}

// Keys - return all keys stored in cache 'c'.
func (c *Cache) Keys(
	ctx context.Context,
) (
	[]string,
	error,
) {
	// TODO: add limit for get N keys
	var (
		keys   []string
		cursor uint64
	)

	for {
		ks, nextCursor, err := c.conn.Scan(ctx, cursor, "*", 100).Result()

		if err != nil {
			return nil, err
		}

		keys = append(keys, ks...)
		cursor = nextCursor

		if cursor == 0 {
			break
		}
	}

	return keys, nil
}

// Values - return all values stored in cache 'c'.
func (c *Cache) Values(
	ctx context.Context,
) (
	[]any,
	error,
) {
	keys, err := c.Keys(ctx)

	if err != nil {
		return nil, err
	}

	values, err := c.conn.MGet(ctx, keys...).Result()

	if err != nil {
		return nil, err
	}

	return values, nil
}

// wrapError - log error if it not nil and call wrapError function.
func (c *Cache) wrapError(err error) error {
	if err == nil {
		return nil
	}

	if c.logger != nil {
		c.logger.LogAttrs(
			nil,
			slog.LevelError,
			"[redis/wrapError]",
			slog.String("Got error", err.Error()),
		)
	}

	return wrapError(err)
}
