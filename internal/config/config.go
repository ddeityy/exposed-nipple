package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"time"

	_ "embed"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	RCON   RCON   `yaml:"rcon"`
	HTTP   HTTP   `yaml:"http"`
	Logger Logger `yaml:"logger"`
}

type RCON struct {
	Host     string `yaml:"host"                     env:"RCON_HOST"`
	Password string `yaml:"password"                 env:"RCON_PASSWORD"`
}

type HTTP struct {
	Port         string        `yaml:"port"          env:"HTTP_SERVER_PORT"`
	ReadTimeout  time.Duration `yaml:"read-timeout"`
	WriteTimeout time.Duration `yaml:"write-timeout"`
	IdleTimeout  time.Duration `yaml:"idle-timeout"`
}

type Logger struct {
	LogLevel        string `yaml:"level"`
	ReportTimestamp bool   `yaml:"report-timestamp"`
}

//go:embed config.yaml
var configBytes []byte

func Load(fileName string) (*Config, error) {
	cfg := Config{}

	err := cleanenv.ReadConfig(fileName, &cfg)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("cleanenv.ReadConfig: %w", err)
		}
		reader := bytes.NewReader(configBytes)
		err = cleanenv.ParseYAML(reader, &cfg)
		if err != nil {
			return nil, fmt.Errorf("cleanenv.ParseYAML: %w", err)
		}
	}

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("cleanenv.ReadEnv: %w", err)
	}

	return &cfg, nil
}
