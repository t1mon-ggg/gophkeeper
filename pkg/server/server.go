package server

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	logger "github.com/t1mon-ggg/gophkeeper/pkg/logging"
	"github.com/t1mon-ggg/gophkeeper/pkg/logging/zerolog"
	apigrpc "github.com/t1mon-ggg/gophkeeper/pkg/server/api-grpc"
	cfg "github.com/t1mon-ggg/gophkeeper/pkg/server/config"
	db "github.com/t1mon-ggg/gophkeeper/pkg/server/storage"
	_ "github.com/t1mon-ggg/gophkeeper/pkg/server/tls"
	web "github.com/t1mon-ggg/gophkeeper/pkg/server/web"
)

type KeeperServer struct {
	sig    chan os.Signal
	Wg     *sync.WaitGroup
	logger logger.Logger
	ag     []apigrpc.APIandGRPC
	// grpc   *mygrpc.Server
	db  db.Storage
	mux *sync.Mutex
}

// New - returns struct of KeeperServer
func New() *KeeperServer {
	ks := new(KeeperServer)
	ks.sig = make(chan os.Signal, 1)
	ks.Wg = new(sync.WaitGroup)
	ks.mux = new(sync.Mutex)
	cfg.New()
	ks.logger = zerolog.New().WithPrefix("server")
	db, err := db.New()
	if err != nil {
		ks.log().Fatal(err, "fatal error on storage initialization")
	}
	ks.db = db
	ks.ag = append(ks.ag, web.New())
	// ks.grpc = mygrpc.New()
	signal.Notify(ks.sig, os.Interrupt, syscall.SIGTERM)
	return ks
}

func (ks *KeeperServer) Start() error {
	ks.Wg.Add(len(ks.ag))
	ks.mux.Lock()
	for _, v := range ks.ag {
		go func(v apigrpc.APIandGRPC) {
			err := v.Start(ks.Wg)
			if err != nil {
				ks.log().Error(err, "endpoint not started")
			}
		}(v)
	}
	ks.mux.Unlock()
	<-ks.Signal()
	ks.Stop()
	ks.log().Info(nil, "Keeper stopped")
	return nil
}

func (ks *KeeperServer) Stop() {
	ks.mux.Lock()
	for _, v := range ks.ag {
		go func(v apigrpc.APIandGRPC) {
			err := v.Stop()
			if err != nil {
				ks.log().Error(err, "failed to stop")
			}
		}(v)
	}
	ks.mux.Unlock()
	// go ks.grpc.Stop()
	ks.Wg.Wait()
	ks.db.Close()
}

// Done - return chan of struct{}. Used for terminating.
func (ks *KeeperServer) Signal() chan os.Signal {
	return ks.sig
}

// log - returns server logger
func (ks *KeeperServer) log() logger.Logger {
	return ks.logger.WithPrefix("main")
}
