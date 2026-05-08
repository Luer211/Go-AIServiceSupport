package config

type Config struct {
	Server    ServerConfig
	MySQL     MySQLConfig
	Redis     RedisConfig
	JWT       JWTConfig
	RateLimit RateLimitConfig
	Task      TaskConfig
	MQ        MQConfig
}

type ServerConfig struct {
	Port string
}

type MySQLConfig struct {
	DSN string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret        string
	ExpireSeconds int64
}

type RateLimitConfig struct {
	IPLimitPerMinute   int
	UserLimitPerMinute int
}

type TaskConfig struct {
	RedisStatusTTLSeconds int64
}

type MQConfig struct {
	MaxRetry int
}

func Default() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "8080",
		},
		MySQL: MySQLConfig{
			DSN: "user:password@tcp(127.0.0.1:3306)/ai_task?charset=utf8mb4&parseTime=True&loc=Local",
		},
		Redis: RedisConfig{
			Addr: "127.0.0.1:6379",
			DB:   0,
		},
		JWT: JWTConfig{
			Secret:        "dev-secret",
			ExpireSeconds: 7200,
		},
		RateLimit: RateLimitConfig{
			IPLimitPerMinute:   30,
			UserLimitPerMinute: 30,
		},
		Task: TaskConfig{
			RedisStatusTTLSeconds: 1800,
		},
		MQ: MQConfig{
			MaxRetry: 3,
		},
	}
}

func Load(path string) (*Config, error) {
	// TODO: 用 Viper 读取 path 中的 YAML，并覆盖 Default 返回的配置。
	// 也就是说，Default() 先给一套默认值，然后 Load(path) 用 Viper 去读取 YAML 配置文件；
	// 如果 YAML 里面有不同的配置，就用 YAML 的值覆盖默认值。
	return Default(), nil
}
