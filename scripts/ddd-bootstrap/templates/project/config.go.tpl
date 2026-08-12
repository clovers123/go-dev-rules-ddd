package configs

import (
	"fmt"
	"os"

	"git.yugeeker.com/SHARED/go-lazy/config"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Basic    *config.BasicConfiguration                     `json:"-" yaml:"-"`
	Postgres map[config.Label]*config.PostgresConfiguration `json:"postgres" yaml:"postgres"`
}

func LoadConfig() (*config.BasicConfiguration, *Config) {
	cfg, err := Load("configs/config.yml")
	if err != nil {
		panic(err)
	}
	config.SetGlobalConfiguration(cfg.Basic)
	return config.GlobalConfiguration(), cfg
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	var basic config.BasicConfiguration
	if err := yaml.Unmarshal(data, &basic); err != nil {
		return nil, err
	}
	cfg.Basic = &basic
	if cfg.Postgres == nil {
		cfg.Postgres = basic.Postgres
	}
	return &cfg, nil
}

func (c *Config) PostgresDSN(name string) (string, error) {
	if name == "" {
		name = "default"
	}
	pg, ok := c.Postgres[config.Label(name)]
	if !ok || pg == nil {
		return "", fmt.Errorf("postgres config %q not found", name)
	}
	sslMode := "disable"
	if pg.SslMode != nil && *pg.SslMode != "" {
		sslMode = *pg.SslMode
	}
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		pg.Host,
		pg.Port,
		pg.User,
		pg.Password,
		pg.Db,
		sslMode,
	), nil
}
