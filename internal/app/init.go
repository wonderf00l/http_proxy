package app

import (
	"fmt"

	"github.com/wonderf00l/http_proxy/internal/configs"
	"github.com/wonderf00l/http_proxy/internal/pkg/logger"
	"go.uber.org/zap"
)

type initErr struct {
	inner error
}

func (e *initErr) Error() string {
	return fmt.Sprintf("Init error: %s\n", e.inner.Error())
}

func Init() (*zap.SugaredLogger, *configs.Configs, error) {
	appCfg, err := configs.NewYAML(configs.AppCfgPath)
	if err != nil {
		return nil, nil, &initErr{inner: fmt.Errorf("app cfg: %w", err)}
	}

	sConfig, err := configs.NewServerConfig(appCfg, configs.ServerKey)
	if err != nil {
		return nil, nil, &initErr{inner: fmt.Errorf("srv cfg: %w", err)}
	}
	proxyConfig, err := configs.NewServerConfig(appCfg, configs.ProxyKey)
	if err != nil {
		return nil, nil, &initErr{inner: fmt.Errorf("proxy cfg: %w", err)}
	}

	// db cfg
	// env -- godotenv.Load()

	serviceLogger, err := logger.New(zap.NewProductionConfig())
	if err != nil {
		return nil, nil, &initErr{inner: fmt.Errorf("create logger: %w", err)}
	}

	return serviceLogger, &configs.Configs{
		SrvCfg:   *sConfig,
		ProxyCfg: *proxyConfig,
	}, nil
}
