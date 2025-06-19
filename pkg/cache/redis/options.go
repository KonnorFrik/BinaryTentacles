package redis

import (
	"errors"
	"fmt"
	"log/slog"
)

// Option - type for customize cache object.
type Option func(*Cache) error

var (
	// ErrInvalidOption - invalid option for any reason.
	ErrInvalidOption = errors.New("invalid option")
)

// WithSlog - register given logger 'l'.
func WithSlog(l *slog.Logger) Option {
	return func(c *Cache) error {
		if l == nil {
			return fmt.Errorf("%w: logger cannot be nil", ErrInvalidOption)
		}

		c.logger = l
		return nil
	}
}
