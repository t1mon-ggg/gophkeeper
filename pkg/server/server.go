package server

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	logger "github.com/t1mon-ggg/gophkeeper/pkg/logging"
	"github.com/t1mon-ggg/gophkeeper/pkg/logging/zerolog"
	cfg "github.com/t1mon-ggg/gophkeeper/pkg/server/config"
	db "github.com/t1mon-ggg/gophkeeper/pkg/server/storage"
	_ "github.com/t1mon-ggg/gophkeeper/pkg/server/tls"
	web "github.com/t1mon-ggg/gophkeeper/pkg/server/web"
)

type KeeperServer struct {
	sig    chan os.Signal
	wg     sync.WaitGroup
	logger logger.Logger
	web    *web.Server
	// grpc   *mygrpc.Server
	db db.Storage
}

// New - returns struct of KeeperServer
func New() *KeeperServer {
	ks := new(KeeperServer)
	ks.sig = make(chan os.Signal, 1)
	ks.wg = sync.WaitGroup{}
	cfg.New()
	ks.logger = zerolog.New().WithPrefix("server")
	db, err := db.New()
	if err != nil {
		ks.log().Fatal(err, "fatal error on storage initialization")
	}
	ks.db = db
	ks.web = web.New()
	// ks.grpc = mygrpc.New()
	signal.Notify(ks.sig, os.Interrupt, syscall.SIGTERM)
	return ks
}

func (ks *KeeperServer) Start() error {
	ks.WG().Add(1)
	go ks.web.Start(ks.WG())
	// go ks.grpc.Start(ks.WG())
	<-ks.Signal()
	ks.Stop()
	ks.log().Info(nil, "Keeper stopped")
	return nil
}

func (ks *KeeperServer) Stop() {
	go ks.web.Stop()
	// go ks.grpc.Stop()
	ks.WG().Wait()
	ks.db.Close()
}

// WG - return application waitgroup
func (ks *KeeperServer) WG() *sync.WaitGroup {
	return &ks.wg
}

// Done - return chan of struct{}. Used for terminating.
func (ks *KeeperServer) Signal() chan os.Signal {
	return ks.sig
}

// log - returns server logger
func (ks *KeeperServer) log() logger.Logger {
	return ks.logger.WithPrefix("main")
}
