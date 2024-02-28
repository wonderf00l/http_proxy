package configs

import (
	"fmt"

	"go.uber.org/config"
)

type SrvConfig struct {
	Proto        string `yaml:"proto"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	CrtFile      string `yaml:"crtFile"`
	KeyFile      string `yaml:"keyFile"`
	DisableHTTP2 bool   `yaml:"-"`
}

func NewServerConfig(appCfg *config.YAML, key string, isProxy bool) (*SrvConfig, error) {
	cfg := &SrvConfig{}
	if err := appCfg.Get(key).Populate(cfg); err != nil {
		return nil, fmt.Errorf("new srv config: %w", err)
	}
	cfg.DisableHTTP2 = isProxy
	return cfg, nil
}
