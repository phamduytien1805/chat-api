package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env    string        `mapstructure:"env"`
	Web    *WebConfig    `mapstructure:"web"`
	DB     *DBConfig     `mapstructure:"db"`
	Hash   *HashConfig   `mapstructure:"hash"`
	Token  *TokenConfig  `mapstructure:"token"`
	Kafka  *KafkaConfig  `mapstructure:"kafka"`
	Redis  *RedisConfig  `mapstructure:"redis"`
	Common *CommonConfig `mapstructure:"common"`
}

type CommonConfig struct {
	Mail struct {
		Expired time.Duration
	}
}

type WebConfig struct {
	Http struct {
		Server struct {
			Port string
		}
		WS struct {
			Port string
		}
	}
}
type DBConfig struct {
	Source string
}

type HashConfig struct {
	Time uint32
	// cpu memory to be used.
	Memory uint32
	// threads for parallelism aspect
	// of the algorithm.
	Threads uint8
	// keyLen of the generate hash key.
	KeyLen uint32
	// saltLen the length of the salt used.
	SaltLen uint32
}

type TokenConfig struct {
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	SecretKey            string
}

type KafkaConfig struct {
	Brokers []string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func setDefault() {

	viper.SetDefault("common.mail.expired", 15*time.Minute)

	viper.SetDefault("web.http.server.port", 5001)
	viper.SetDefault("web.http.ws.port", 5002)
	viper.SetDefault("env", "development")
	viper.SetDefault("db.source", "postgresql://root:secret@localhost:5432/core?sslmode=disable")

	viper.SetDefault("hash.time", 1)
	viper.SetDefault("hash.memory", 64*1024)
	viper.SetDefault("hash.threads", 32)
	viper.SetDefault("hash.keyLen", 256)
	viper.SetDefault("hash.saltLen", 10)

	viper.SetDefault("token.accessTokenDuration", "5m")
	viper.SetDefault("token.refreshTokenDuration", "24h")
	viper.SetDefault("token.secretKey", "secret_secret_secret_secret_secret_secret_secret_secret")

	viper.SetDefault("kafka.brokers", []string{"localhost:9092"})

	viper.SetDefault("redis.addr", "localhost:6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
}

func NewConfig() (*Config, error) {
	setDefault()

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
