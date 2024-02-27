package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	delivery "github.com/wonderf00l/http_proxy/internal/delivery/http"
)

// set api route
// set proxy route -- '/*', proxyHandler (check proto - maybe in mw)

// 2 инстанса сервера, старуем их в горутинах

func NewRouter(h *delivery.HandlerHTTP) http.Handler {
	mux := chi.NewMux()
	// route http handler methods
	return mux
}

func NewProxyRouter(h *delivery.HandlerHTTP) http.Handler {
	return http.HandlerFunc(h.ProxyHTTP)
}
