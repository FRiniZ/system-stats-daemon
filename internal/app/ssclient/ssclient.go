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

type Application struct{}

func (app *Application) Process(resp *api.Responce) {
	if resp.GetCPU() != nil {
		log.Printf("%-18s[User:%5.2f%% System:%5.2f%% Idle:%5.2f%%]", "CPU",
			resp.GetCPU().User, resp.GetCPU().System, resp.GetCPU().Idle)
	}

	if resp.GetLoadAvg() != nil {
		log.Printf("%-18s[1m:%5.2f 5m:%5.2f 15m:%5.2f]", "LoadAverage",
			resp.GetLoadAvg().L1, resp.GetLoadAvg().L2, resp.GetLoadAvg().L3)
	}

	if resp.GetDisks() != nil {
		log.Printf("%-18s%8s%14s%14s\n", "Device", "tps", "Read KB/s", "Write KB/s")
		for _, d := range resp.GetDisks() {
			log.Printf("%-18s%8.2f%14.2f%14.2f\n", d.Name, d.TPS, d.ReadKBs, d.WriteKBs)
		}
	}

	if resp.GetDfsize() != nil {
		log.Printf("%-18s%15s%14s%%\n", "FileSystem", "Used", "Use")
		for _, d := range resp.GetDfsize() {
			log.Printf("%-18s%14dM%14d%%\n", d.Name, d.Used, d.Use)
		}
	}
	if resp.GetDfinode() != nil {
		log.Printf("%-18s%15s%14s%%\n", "FileSystem", "IUsed", "IUse")
		for _, d := range resp.GetDfinode() {
			log.Printf("%-18s%14dM%14d%%\n", d.Name, d.IUsed, d.IUse)
		}
	}
}

func (app *Application) MakeSTATS(sensors []string) api.STATS {
	var stats api.STATS
	log.Println("Request sensors:")
stoploop:
	for _, name := range sensors {
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
	return stats
}

func (app *Application) Run(conf Config) {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer cancel()

	ctxDial, cancelDial := context.WithTimeout(ctx, 5*time.Second)
	conn, err := grpc.DialContext(ctxDial, net.JoinHostPort(conf.GRPC.Host, conf.GRPC.Port),
		grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer cancelDial()

	if err != nil {
		log.Printf("Can't connect to GRPC-Server[%s:%s]: %v\n", conf.GRPC.Host, conf.GRPC.Port, err)
		return
	}

	defer conn.Close()

	grpcClient := api.NewSSDClient(conn)
	stats := app.MakeSTATS(conf.Core.Sensors)

	stream, err := grpcClient.Subsribe(ctx, &api.Request{
		N:       durationpb.New(conf.Core.N),
		M:       durationpb.New(conf.Core.M),
		Bitmask: stats,
	})
	if err != nil {
		fmt.Printf("Can't call remote function:Subscribe:%v\n", err)
		return
	}

	for {
		resp, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			log.Println("Server closed connection")
			break
		}

		if errors.Is(stream.Context().Err(), context.Canceled) {
			log.Println(err)
			log.Println("Canceled")
			break
		}

		if err != nil {
			log.Printf("stream.Recv: %v, %T\n", err, err)
			break
		}

		app.Process(resp)
	}
}
