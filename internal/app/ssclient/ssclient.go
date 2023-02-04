package ssclient

import (
	"context"
	"io"
	"log"
	"net"
	"os/signal"
	"syscall"
	"time"

	api "github.com/FRiniZ/system-stats-daemon/api/stub"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	GRPC struct {
		Host string `toml:"host"`
		Port string `toml:"port"`
	} `toml:"GRPC"`

	Core struct {
		N       time.Duration `toml:"N"`
		M       time.Duration `toml:"M"`
		Sensors []string      `toml:"Sensors"`
	} `toml:"core"`
}

type Application struct {
	Conf Config
}

func (app *Application) Run() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer cancel()

	ctxDial, cancelDial := context.WithTimeout(ctx, 5*time.Second)
	conn, err := grpc.DialContext(ctxDial, net.JoinHostPort(app.Conf.GRPC.Host, app.Conf.GRPC.Port),
		grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer cancelDial()

	if err != nil {
		log.Println(err)
		log.Fatalf("Can't connect to GRPC-Server[%s:%s]\n", app.Conf.GRPC.Host, app.Conf.GRPC.Port)
	}

	defer conn.Close()

	grpcClient := api.NewSSDClient(conn)
	var stats api.STATS
	for _, name := range app.Conf.Core.Sensors {
		switch name {
		case "CPU":
			stats |= api.STATS_CPU
		case "LOADAVERAGE":
			stats |= api.STATS_LOADAVERAGE
		case "LOADDISK":
			stats |= api.STATS_LOADDISK
		case "ALL":
			stats |= api.STATS_CPU
			stats |= api.STATS_LOADAVERAGE
			stats |= api.STATS_LOADDISK
		}
	}
	log.Println("Stats:", stats)

	stream, err := grpcClient.Subsribe(ctx, &api.Request{
		N:       int32(app.Conf.Core.N.Seconds()),
		M:       int32(app.Conf.Core.M.Seconds()),
		Bitmask: stats,
	})

	for {
		recv, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalln(err)
		}
		log.Println(recv)
	}
}
