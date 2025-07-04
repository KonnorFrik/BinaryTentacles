package redis

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	Addr        string        `env:"REDIS_ADDR" env-required`
	Password    string        `env:"REDIS_PASSWORD" env-required`
	User        string        `env:"REDIS_USER" env-required`
	DB          int           `env:"REDIS_DB" env-required`
	MaxRetries  int           `env:"REDIS_MAX_RETRIES" env-required`
	DialTimeout time.Duration `env:"REDIS_DIAL_TIMEOUT" env-required`
	Timeout     time.Duration `env:"REDIS_RW_TIMEOUT" env-required`
}

// ConfigOption - option for customize the Config object at creation.
type ConfigOption func(*Config) error

// NewConfig - create a new config for redis.
// Read values from env as default.
func NewConfig(opts ...ConfigOption) (Config, error) {
	var config Config
	err := cleanenv.ReadEnv(&config)

	if err == nil {
		for _, opt := range opts {
			if e := opt(&config); e != nil {
				return config, e
			}
		}
	}

	return config, err
}

// WithDB - ConfigOption for specify redis db number.
func WithDB(num int) ConfigOption {
	return func(c *Config) error {
		c.DB = num
		return nil
	}
}
