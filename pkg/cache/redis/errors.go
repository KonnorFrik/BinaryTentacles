package redis

import (
	"errors"

	"github.com/redis/go-redis/v9"
)

var (
	// ErrUnknow - any undocumented error from redis
	ErrUnknow = errors.New("cache error")
	// ErrTimeOut - from redis "timed out waiting to get a connection from the connection pool"
	ErrTimeOut = errors.New("time out")
	ErrNil     = errors.New("value is nil")
)

// wrapError - Map errors from other packages.
// Wrap errors from this package.
func wrapError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case err == redis.Nil:
		return ErrNil
	case errors.Is(err, redis.ErrPoolTimeout):
		return ErrTimeOut
	case errors.Is(err, ErrInvalidOption):
		return err
	default:
		return ErrUnknow
	}
}
