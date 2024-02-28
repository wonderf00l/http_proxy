package app

import (
	"crypto/tls"
	"errors"
	"net/http"
	"strconv"
	"sync"

	"github.com/wonderf00l/http_proxy/internal/configs"
	"go.uber.org/zap"
)

// proxy просто как отдельная ручка - а так  прокси и апи к одному серверу

type Server struct {
	http.Server
	logger  *zap.SugaredLogger
	proto   string
	crtFile string
	keyFile string
}

func NewServer(cfg *configs.SrvConfig, mux http.Handler, log *zap.SugaredLogger) *Server {
	srv := &Server{
		Server: http.Server{
			Addr:    cfg.Host + ":" + strconv.Itoa(cfg.Port),
			Handler: mux,
		},
		logger:  log,
		proto:   cfg.Proto,
		crtFile: cfg.CrtFile,
		keyFile: cfg.KeyFile,
	}

	if cfg.DisableHTTP2 {
		srv.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler))
	}

	return srv
}

func (s *Server) Run(wg *sync.WaitGroup) error {
	switch s.proto {
	case "http":
		go func() {
			defer wg.Done()

			s.logger.Infoln("Staring http server at", s.Addr)
			if err := s.Server.ListenAndServe(); err != http.ErrServerClosed {
				s.logger.Errorln("Listen and server: ", err)
			}
		}()
	case "https":
		go func() {
			defer wg.Done()

			s.logger.Infoln("Staring https server at", s.Addr)
			if err := s.Server.ListenAndServeTLS(s.crtFile, s.keyFile); err != http.ErrServerClosed {
				s.logger.Errorln("Listen and server tls: ", err)
			}
		}()
	default:
		return errors.New("run server: invalid proto")
	}
	return nil
}

// run(mux)
