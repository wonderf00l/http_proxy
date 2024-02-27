package app

import (
	"context"
	"fmt"
	"sync"

	"github.com/wonderf00l/http_proxy/internal/configs"
	delivery "github.com/wonderf00l/http_proxy/internal/delivery/http"
	"go.uber.org/zap"
)

// prepreq() -- create yaml cfg and other cfgs(srv, db)

// run() - get ctx, cfg -- create db conn, create repo,create handler, create router, create srvs, run srvs, ctx.Done
// run srv
// shoutdown func

type runErr struct {
	inner error
}

func (e *runErr) Error() string {
	return fmt.Sprintf("Run error: %s\n", e.inner.Error())
}

func Run(ctx context.Context, logger *zap.SugaredLogger, cfgs *configs.Configs) error {
	h := delivery.NewHandler(logger, nil)

	router, proxyRouter := NewRouter(h), NewProxyRouter(h)
	server, proxy := NewServer(&cfgs.SrvCfg, router, logger), NewServer(&cfgs.ProxyCfg, router, logger)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	server.Run(router, wg)
	wg.Add(1)
	proxy.Run(proxyRouter, wg)

	go func() {
		if err := proxy.Shutdown(ctx); err != nil {
			logger.Errorln("Shutdown proxy: ", err.Error())
		}
	}()
	if err := server.Shutdown(ctx); err != nil {
		return &runErr{inner: fmt.Errorf("Shutdown api: %w", err)}
	}

	wg.Wait()

	return nil
}
