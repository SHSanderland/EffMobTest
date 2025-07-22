package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string `yaml:"env" env:"ENV" env-default:"local"`
	Server   `yaml:"server"`
	Database `yaml:"database"`
}

type Server struct {
	Addr         string        `yaml:"address" env:"ADDR" env-default:":8080"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env:"RT" env-default:"5s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env:"WT" env-default:"5s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env:"IT" env-default:"10s"`
}

type Database struct {
	DSN       string `yaml:"dsn" env:"DSN" env-required:"true"`
	SourceURL string `yaml:"sourceURL" env:"SURL" env-default:"./migrations"`
}

func InitConfig(flagPath string) *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		if flagPath != "" {
			configPath = flagPath
		} else {
			log.Fatal("CONFIG_PATH is not set")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("error check path: %v", err)
	}

	cfg := Config{}

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("error parce file: %v", err)
	}

	return &cfg
}
