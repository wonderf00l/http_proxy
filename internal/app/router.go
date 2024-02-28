package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	delivery "github.com/wonderf00l/http_proxy/internal/delivery/http"
	"github.com/wonderf00l/http_proxy/internal/middleware"
	"go.uber.org/zap"
)

// set api route
// set proxy route -- '/*', proxyHandler (check proto - maybe in mw)

// 2 инстанса сервера, старуем их в горутинах

func NewRouter(h *delivery.HandlerHTTP) http.Handler {
	mux := chi.NewMux()
	// route http handler methods
	return mux
}

func NewProxyRouter(h *delivery.HandlerHTTP, log *zap.SugaredLogger) http.Handler {
	r := chi.NewMux()
	r.Use(middleware.LoggingMW(log))
	r.HandleFunc("/", http.HandlerFunc(h.DoProxy))
	return r
}
