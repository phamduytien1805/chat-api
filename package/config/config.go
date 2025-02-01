package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env    string        `mapstructure:"env"`
	Auth   *AuthConfig   `mapstructure:"authsvc"`
	User   *UserConfig   `mapstructure:"usersvc"`
	DB     *DBConfig     `mapstructure:"db"`
	Hash   *HashConfig   `mapstructure:"hash"`
	Token  *TokenConfig  `mapstructure:"token"`
	Kafka  *KafkaConfig  `mapstructure:"kafka"`
	Redis  *RedisConfig  `mapstructure:"redis"`
	Scylla *ScyllaConfig `mapstructure:"scylla"`
	Mail   *MailConfig   `mapstructure:"mail"`
}

type AuthConfig struct {
	Http struct {
		Server struct {
			Port string
		}
		WS struct {
			Port string
		}
	}
}

type UserConfig struct {
	Http struct {
		Server struct {
			Port string
		}
	}
	Grpc struct {
		Server struct {
			Port string
			Host string
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

type ScyllaConfig struct {
	Hosts             []string
	Keyspace          string
	Class             string
	ReplicationFactor int
}

type MailConfig struct {
	Host           string
	Port           int
	Username       string
	Password       string
	Origin         string
	Expired        time.Duration
	VerifyEmailUrl string
}

func setDefault() {
	viper.SetDefault("auth.http.server.port", 5001)
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

	viper.SetDefault("mail.host", "smtp.gmail.com")
	viper.SetDefault("mail.port", 587)
	viper.SetDefault("mail.username", "")
	viper.SetDefault("mail.password", "")
	viper.SetDefault("mail.expired", 15*time.Minute)

	viper.SetDefault("scylla.hosts", []string{"localhost:9042"})
	viper.SetDefault("scylla.keyspace", "chatcore")
	viper.SetDefault("scylla.class", "SimpleStrategy")
	viper.SetDefault("scylla.replicationFactor", 2)

	viper.SetDefault("user.grpc.server.port", 5002)
	viper.SetDefault("user.grpc.server.host", "localhost")

}

func NewConfig() (*Config, error) {
	setDefault()

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
