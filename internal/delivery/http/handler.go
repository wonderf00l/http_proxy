package http

import (
	"net/http"

	"github.com/wonderf00l/http_proxy/internal/repository"
	"go.uber.org/zap"
)

// proxy included

type HandlerHTTP struct {
	logger *zap.SugaredLogger
	repo   repository.Repository
}

func NewHandler(log *zap.SugaredLogger, repo repository.Repository) *HandlerHTTP {
	return &HandlerHTTP{
		logger: log,
		repo:   repo,
	}
}

func (h *HandlerHTTP) DoProxy(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodConnect {
		// ...
	} else {
		h.ProxyHTTP(w, req)
	}
}
