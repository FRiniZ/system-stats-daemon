package ssclient

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os/signal"
	"syscall"
	"time"

	api "github.com/FRiniZ/system-stats-daemon/api/stub"
	"github.com/FRiniZ/system-stats-daemon/internal/services/sensors/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/durationpb"
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

	ID string `toml:"ClientID"`
}

type Application struct {
	Conf Config
}

func (app *Application) Run() {
	var stats api.STATS
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
		log.Printf("Can't connect to GRPC-Server[%s:%s]: %v\n", app.Conf.GRPC.Host, app.Conf.GRPC.Port, err)
		return
	}

	defer conn.Close()

	grpcClient := api.NewSSDClient(conn)

	log.Println("Request sensors:")
stoploop:
	for _, name := range app.Conf.Core.Sensors {
		log.Println("               : ", name)
		switch name {
		case common.ALL:
			stats = 0
			stats |= api.STATS_ALL
			break stoploop
		case common.CPU:
			stats |= api.STATS_CPU
		case common.LOADAVERAGE:
			stats |= api.STATS_LOADAVERAGE
		case common.LOADDISK:
			stats |= api.STATS_LOADDISK
		case common.SIZEDISK:
			stats |= api.STATS_SIZEDISK
		case common.INODEDISK:
			stats |= api.STATS_INODEDISK
		}
	}

	stream, err := grpcClient.Subsribe(ctx, &api.Request{
		N:       durationpb.New(app.Conf.Core.N),
		M:       durationpb.New(app.Conf.Core.M),
		Bitmask: stats,
	})

	if err != nil {
		fmt.Printf("Can't call remote function:Subscribe:%v\n", err)
		return
	}

	for {
		recv, err := stream.Recv()

		if errors.Is(err, io.EOF) || stream.Context().Err() == context.Canceled {
			break
		}

		if err != nil {
			log.Printf("stream.Recv: %v, %T\n", err, err)
			break
		}

		if recv.GetCPU() != nil {
			log.Printf("%-18s[User:%5.2f%% System:%5.2f%% Idle:%5.2f%%]", "CPU", recv.GetCPU().User, recv.GetCPU().System, recv.GetCPU().Idle)
		}

		if recv.GetLoadAvg() != nil {
			log.Printf("%-18s[1m:%5.2f 5m:%5.2f 15m:%5.2f]", "LoadAverage", recv.GetLoadAvg().L1, recv.GetLoadAvg().L2, recv.GetLoadAvg().L3)
		}

		if recv.GetDisks() != nil {
			log.Printf("%-18s%8s%14s%14s\n", "Device", "tps", "Read KB/s", "Write KB/s")
			for _, d := range recv.GetDisks() {
				log.Printf("%-18s%8.2f%14.2f%14.2f\n", d.Name, d.TPS, d.ReadKBs, d.WriteKBs)
			}
		}

		if recv.GetDfsize() != nil {
			log.Printf("%-18s%15s%14s%%\n", "FileSystem", "Used", "Use")
			for _, d := range recv.GetDfsize() {
				log.Printf("%-18s%14dM%14d%%\n", d.Name, d.Used, d.Use)
			}
		}
		if recv.GetDfinode() != nil {
			log.Printf("%-18s%15s%14s%%\n", "FileSystem", "IUsed", "IUse")
			for _, d := range recv.GetDfinode() {
				log.Printf("%-18s%14dM%14d%%\n", d.Name, d.IUsed, d.IUse)
			}
		}
	}
}
