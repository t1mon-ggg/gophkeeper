package server

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	logger "github.com/t1mon-ggg/gophkeeper/pkg/logging"
	log "github.com/t1mon-ggg/gophkeeper/pkg/logging/zerolog"
	cfg "github.com/t1mon-ggg/gophkeeper/pkg/server/config"
	mygrpc "github.com/t1mon-ggg/gophkeeper/pkg/server/grpc"
	db "github.com/t1mon-ggg/gophkeeper/pkg/server/storage"
	_ "github.com/t1mon-ggg/gophkeeper/pkg/server/tls"
	web "github.com/t1mon-ggg/gophkeeper/pkg/server/web"
)

type KeeperServer struct {
	sig     chan os.Signal
	wg      sync.WaitGroup
	logger  logger.Logger
	config  *cfg.Config
	http    *web.Server
	grpc    *mygrpc.Server
	storage db.Storage
}

// New - returns struct of KeeperServer
func New() *KeeperServer {
	ks := new(KeeperServer)
	ks.sig = make(chan os.Signal, 1)
	ks.wg = sync.WaitGroup{}
	ks.logger = log.New()
	ks.config = cfg.New()
	ks.logger.SetLevel(ks.config.Level())
	storage, err := db.New(ks.config.DSN, ks.logger)
	if err != nil {
		ks.log().Fatal(err, "fatal error on storage initialization")
	}
	ks.storage = storage
	ks.http = web.New(ks.config.WebBind, ks.storage, ks.logger)
	ks.grpc = mygrpc.New(ks.config.GRPCBind, ks.storage, ks.logger)
	signal.Notify(ks.sig, os.Interrupt, syscall.SIGTERM)
	return ks
}

func (ks *KeeperServer) Start() error {
	ks.WG().Add(2)
	go ks.http.Start(ks.WG())
	go ks.grpc.Start(ks.WG())
	<-ks.Signal()
	ks.Stop()
	ks.log().Info(nil, "Keeper stopped")
	return nil
}

func (ks *KeeperServer) Stop() {
	go ks.http.Stop()
	go ks.grpc.Stop()
	ks.WG().Wait()
	ks.storage.Close()
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
