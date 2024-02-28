package app

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/wonderf00l/http_proxy/internal/configs"
	delivery "github.com/wonderf00l/http_proxy/internal/delivery/http"
	"go.uber.org/zap"
)

// prepreq() -- create yaml cfg and other cfgs(srv, db)

// run() - get ctx, cfg -- create db conn, create repo,create handler, create router, create srvs, run srvs, ctx.Done
// run srv
// shoutdown func

var (
	shutdownTimeout = 5 * time.Second
)

type runErr struct {
	inner error
}

func (e *runErr) Error() string {
	return fmt.Sprintf("Run error: %s\n", e.inner.Error())
}

func Run(ctx context.Context, logger *zap.SugaredLogger, cfgs *configs.Configs) error {
	h := delivery.NewHandler(logger, nil)

	router, proxyRouter := NewRouter(h), NewProxyRouter(h, logger)
	server, proxy := NewServer(&cfgs.SrvCfg, router, logger), NewServer(&cfgs.ProxyCfg, proxyRouter, logger)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	if err := server.Run(wg); err != nil {
		return err
	}
	wg.Add(1)
	if err := proxy.Run(wg); err != nil {
		return err
	}

	<-ctx.Done()
	logger.Infoln("Shutting down gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	go func() {
		if err := proxy.Shutdown(shutdownCtx); err != nil {
			logger.Errorln("shutdown proxy: ", err.Error())
		}
	}()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return &runErr{inner: fmt.Errorf("shutdown api: %w", err)}
	}

	wg.Wait()

	return nil
}
