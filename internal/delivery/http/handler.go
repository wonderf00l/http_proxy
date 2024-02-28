package http

import (
	"net/http"

	"github.com/wonderf00l/http_proxy/internal/repository"
	"go.uber.org/zap"
)

type HandlerHTTP struct {
	proxyClient http.Client
	logger      *zap.SugaredLogger
	repo        repository.Repository
}

func NewHandler(log *zap.SugaredLogger, repo repository.Repository) *HandlerHTTP {
	return &HandlerHTTP{
		proxyClient: http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }},
		logger:      log,
		repo:        repo,
	}
}

func (h *HandlerHTTP) DoProxy(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodConnect {
		h.tunnelHTTP(w, req)
	} else {
		h.ProxyHTTP(w, req)
	}
}
