package app

import (
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
	logger *zap.SugaredLogger
	proto  string
}

// init srv
func NewServer(cfg *configs.SrvConfig, mux http.Handler, log *zap.SugaredLogger) *Server {
	return &Server{
		Server: http.Server{
			Addr:    cfg.Host + ":" + strconv.Itoa(cfg.Port),
			Handler: mux,
		},
		logger: log,
	}
}

func (s *Server) Run(mux http.Handler, wg *sync.WaitGroup) error {
	s.logger.Infoln("Staring server at ", s.Addr)
	switch s.proto {
	case "http":
		go func() {
			defer wg.Done()

			if err := s.Server.ListenAndServe(); err != http.ErrServerClosed {
				s.logger.Errorln("Listen and server: ", err)
			}
		}()
	default:
		return errors.New("run server: invalid proto")
	}
	return nil
}

// run(mux)
