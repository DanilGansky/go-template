package config

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/kelseyhightower/envconfig"

	"github.com/joho/godotenv"
)

type Environment string

const (
	LOCAL      Environment = "local"
	STAGING    Environment = "staging"
	PRODUCTION Environment = "production"
)

type Config struct {
	// DB connection string
	DSN string `envconfig:"DB_DSN"`

	BcryptConfig
	ServerConfig
}

type BcryptConfig struct {
	Secret string `required:"true" envconfig:"SECRET"`
	Cost   int    `required:"true" envconfig:"BCRYPT_COST"`
	Issuer string `required:"true" envconfig:"ISSUER"`
}

type ServerConfig struct {
	Addr        string      `required:"true" envconfig:"ADDR" default:"0.0.0.0:8000"`
	Environment Environment `required:"true" envconfig:"ENVIRONMENT" default:"local"`
	Timeout     int         `required:"true" envconfig:"TIMEOUT" default:"10"` // in seconds
}

func (c *Config) MaxTimeout() time.Duration {
	return time.Duration(c.Timeout) * time.Second
}

func (c *Config) LogLevel() logrus.Level {
	switch c.Environment {
	case LOCAL:
		return logrus.DebugLevel
	case STAGING:
		return logrus.WarnLevel
	case PRODUCTION:
		return logrus.ErrorLevel
	default:
		return logrus.InfoLevel
	}
}

func (c Config) String() string {
	result, _ := json.MarshalIndent(&c, "", "    ")
	return string(result)
}

var (
	once sync.Once
	cfg  *Config
)

func Get() Config {
	once.Do(func() {
		if err := godotenv.Load(".env"); err != nil {
			log.Println(err)
		}

		var config Config
		if err := envconfig.Process("", &config); err != nil {
			log.Fatalf("failed to load config: %v", err)
		}

		cfg = &config
	})

	return *cfg
}
