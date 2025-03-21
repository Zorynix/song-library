package config

import (
	"fmt"
	"path"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App        `yaml:"app"`
		Server     `yaml:"server"`
		Log        `yaml:"log"`
		PG         `yaml:"postgres"`
		Prometheus `yaml:"prometheus"`
		MusicAPI   `yaml:"music_api"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	Server struct {
		Port int `env-required:"true" yaml:"port" env:"SERVER_PORT"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	}

	PG struct {
		MaxPoolSize  int           `env-required:"true" yaml:"max_pool_size" env:"PG_MAX_POOL_SIZE"`
		URL          string        `env-required:"true" env:"PG_URL"`
		ConnAttempts int           `env-required:"true" yaml:"conn_attempts" env:"PG_CONN_ATTEMPTS"`
		ConnTimeout  time.Duration `env-required:"true" yaml:"conn_timeout" env:"PG_CONN_TIMEOUT"`
	}

	Prometheus struct {
		MetricsPort int `env-required:"true" yaml:"metrics_port" env:"PROMETHEUS_METRICS_PORT"`
	}

	MusicAPI struct {
		URL string `env-required:"true" yaml:"url" env:"MUSIC_API_URL"`
	}
)

func NewConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(path.Join("./", configPath), cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	err = cleanenv.UpdateEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("error updating env: %w", err)
	}

	return cfg, nil
}
