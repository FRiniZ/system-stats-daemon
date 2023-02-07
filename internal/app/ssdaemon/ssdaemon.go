package ssdaemon

import (
	"context"
	"log"
	"net"
	"os/signal"
	"sync"
	"syscall"
	"time"

	api "github.com/FRiniZ/system-stats-daemon/api/stub"
	grpcserver "github.com/FRiniZ/system-stats-daemon/internal/services/grpc-server"
	"google.golang.org/grpc"
)

type Config struct {
	GRPC struct {
		Host string `toml:"host"`
		Port string `toml:"port"`
	} `toml:"GRPC"`

	Core struct {
		Frequency time.Duration `toml:"frequency"`
	} `toml:"core"`
}

type Application struct {
	Conf Config
}

func (app *Application) Run() {
	var opts []grpc.ServerOption

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer cancel()

	lis, err := net.Listen("tcp", net.JoinHostPort(app.Conf.GRPC.Host, app.Conf.GRPC.Port))
	if err != nil {
		log.Fatalln(err)
	}

	wg := &sync.WaitGroup{}
	grpcBase := grpc.NewServer(opts...)
	grpcSrv := grpcserver.New(wg)
	grpcSrv.Start(ctx, 0)
	api.RegisterSSDServer(grpcBase, grpcSrv)

	go func() {
		<-ctx.Done()
		grpcBase.Stop()
		log.Println("GRPC-server stopping...")
	}()

	log.Printf("GRPC-server listening:[%s:%s]\n", app.Conf.GRPC.Host, app.Conf.GRPC.Port)
	grpcBase.Serve(lis)

	wg.Wait()
	log.Println("App closed")
}
