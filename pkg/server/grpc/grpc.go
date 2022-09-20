package grpc

// import (
// 	"sync"

// 	"github.com/t1mon-ggg/gophkeeper/pkg/logging"
// 	"github.com/t1mon-ggg/gophkeeper/pkg/logging/zerolog"
// 	"github.com/t1mon-ggg/gophkeeper/pkg/server/config"
// )

// type Server struct {
// 	bind   string
// 	logger logging.Logger
// 	wg     *sync.WaitGroup
// }

// func New() *Server {
// 	s := new(Server)
// 	s.bind = config.New().GRPCBind
// 	s.logger = zerolog.New().WithPrefix("grpc")
// 	return s
// }

// func (s *Server) log() logging.Logger {
// 	return s.logger
// }

// func (s *Server) Start(wg *sync.WaitGroup) {
// 	s.wg = wg

// }

// func (s *Server) Stop() {
// 	s.log().Info(nil, "Graceful shutdown in progress...")

// 	s.log().Info(nil, "gRPC server stopped")
// 	s.wg.Done()
// }
