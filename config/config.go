package config

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/joho/godotenv"
)

type Config struct {
	// DB connection string
	DSN string `envconfig:"DB_DSN"`

	BcryptConfig
	ServerConfig
}

type BcryptConfig struct {
	Secret string `envconfig:"SECRET"`
	Cost   int    `envconfig:"BCRYPT_COST"`
}

type ServerConfig struct {
	Addr    string `envconfig:"ADDR" default:"0.0.0.0:8000"`
	Timeout int    `envconfig:"TIMEOUT" default:"10"` // in seconds
}

func (c *Config) MaxTimeout() time.Duration {
	return time.Duration(c.Timeout) * time.Second
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
		_ = godotenv.Load("../../.env")

		var config Config
		if err := envconfig.Process("", &config); err != nil {
			panic(err)
		}

		cfg = &config
	})

	return *cfg
}
