package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	MySQL     MySQLConfig     `mapstructure:"mysql"`
	Redis     RedisConfig     `mapstructure:"redis"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
	Task      TaskConfig      `mapstructure:"task"`
	MQ        MQConfig        `mapstructure:"mq"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type MySQLConfig struct {
	DSN string `mapstructure:"dsn"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int 	`mapstructure:"db"`
}

type JWTConfig struct {
	Secret        string `mapstructure:"secret"`
	ExpireSeconds int64  `mapstructure:"expire_seconds"`
}

type RateLimitConfig struct {
	IPLimitPerMinute   int `mapstructure:"ip_limit_per_minute"`
	UserLimitPerMinute int `mapstructure:"user_limit_per_minute"`
}

type TaskConfig struct {
	RedisStatusTTLSeconds int64 `mapstructure:"redis_status_ttl_seconds"`
}

type MQConfig struct {
	URL        string `mapstructure:"url"`
	Exchange   string `mapstructure:"exchange"`
	Queue      string `mapstructure:"queue"`
	RoutingKey string `mapstructure:"routing_key"`
	MaxRetry   int    `mapstructure:"max_retry"`
}

// 使用 Viper 从 yaml 中读取配置
func Load(path string) (*Config, error) {
	cfg := &Config{}

	v := viper.New()
	// 指定配置文件路径
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	// 将 YAML 中的配置映射到 cfg 结构体中
	if err := v.Unmarshal(cfg); err != nil {
		return nil, err
	}

	// 校验cfg的各项配置是否合法
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Todo: 这里的错误要接入日志，我们目前只是打印出来
func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server.port is required")
	}
	if c.MySQL.DSN == "" {
		return fmt.Errorf("mysql.dsn is required")
	}
	if c.JWT.Secret == "" {
		return fmt.Errorf("jwt.secret is required")
	}
	if c.JWT.ExpireSeconds <= 0 {
		return fmt.Errorf("jwt.expire_seconds must be positive")
	}
	if c.Redis.Addr == "" {
		return fmt.Errorf("redis.addr is required")
	}
	if c.MQ.URL == "" {
	return fmt.Errorf("mq.url is required")
	}
	if c.MQ.Exchange == "" {
		return fmt.Errorf("mq.exchange is required")
	}
	if c.MQ.Queue == "" {
		return fmt.Errorf("mq.queue is required")
	}
	if c.MQ.RoutingKey == "" {
		return fmt.Errorf("mq.routing_key is required")
	}
	if c.MQ.MaxRetry < 0 {
		return fmt.Errorf("mq.max_retry must be non-negative")
	}

	return nil
}
