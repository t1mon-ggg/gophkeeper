package grpc

import (
	"sync"

	"github.com/t1mon-ggg/gophkeeper/pkg/logging"
	"github.com/t1mon-ggg/gophkeeper/pkg/server/storage"
)

type Server struct {
	bind   string
	logger logging.Logger
	wg     *sync.WaitGroup
}

func New(bind string, storage storage.Storage, logger logging.Logger) *Server {
	s := new(Server)
	s.bind = bind
	s.logger = logger.WithPrefix("grpc")
	return s
}

func (s *Server) log() logging.Logger {
	return s.logger
}

func (s *Server) Start(wg *sync.WaitGroup) {
	s.wg = wg

}

func (s *Server) Stop() {
	s.log().Info(nil, "Graceful shutdown in progress...")

	s.log().Info(nil, "gRPC server stopped")
	s.wg.Done()
}
