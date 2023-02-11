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

func (app *Application) ProcessDummy(resp *api.Dummy) {
	if resp.GetErrorMsg() == "" {
		log.Printf("%-18s", common.DUMMY)
	} else {
		log.Printf("%-18s:%s", common.DUMMY, resp.GetErrorMsg())
	}
}

func (app *Application) ProcessCPU(resp *api.Cpu) {
	if resp.GetErrorMsg() == "" {
		log.Printf("%-18s[User:%5.2f%% System:%5.2f%% Idle:%5.2f%%]", common.CPU,
			resp.User, resp.System, resp.Idle)
	} else {
		log.Printf("%-18s:%s", common.CPU, resp.GetErrorMsg())
	}
}

func (app *Application) ProcessLoadAvg(resp *api.Loadaverage) {
	if resp.GetErrorMsg() == "" {
		log.Printf("%-18s[1m:%5.2f 5m:%5.2f 15m:%5.2f]", common.LOADAVERAGE,
			resp.L1, resp.L2, resp.L3)
	} else {
		log.Printf("%-18s:%s", common.LOADAVERAGE, resp.GetErrorMsg())
	}
}

func (app *Application) ProcessDisks(resp []*api.Loaddisk) {
	log.Printf("%-18s%8s%14s%14s\n", common.LOADDISK, "tps", "Read KB/s", "Write KB/s")
	for _, d := range resp {
		if d.GetErrorMsg() == "" {
			log.Printf("%-18s%8.2f%14.2f%14.2f\n", d.Name, d.TPS, d.ReadKBs, d.WriteKBs)
		} else {
			log.Printf("%-18s:%s", common.LOADDISK, d.GetErrorMsg())
		}
	}
}

func (app *Application) ProcessDfsize(resp []*api.Dfsize) {
	log.Printf("%-18s%15s%14s%%\n", common.SIZEDISK, "Used", "Use")
	for _, d := range resp {
		if d.GetErrorMsg() == "" {
			log.Printf("%-18s%14dM%14d%%\n", d.Name, d.Used, d.Use)
		} else {
			log.Printf("%-18s:%s", common.SIZEDISK, d.GetErrorMsg())
		}
	}
}

func (app *Application) ProcessDfinode(resp []*api.Dfinode) {
	log.Printf("%-18s%15s%14s%%\n", common.INODEDISK, "IUsed", "IUse")
	for _, d := range resp {
		if d.GetErrorMsg() == "" {
			log.Printf("%-18s%14dM%14d%%\n", d.Name, d.IUsed, d.IUse)
		} else {
			log.Printf("%-18s:%s", common.INODEDISK, d.GetErrorMsg())
		}
	}
}

func (app *Application) Process(resp *api.Responce) {
	switch {
	case resp.GetDummy() != nil:
		app.ProcessDummy(resp.GetDummy())
	case resp.GetCPU() != nil:
		app.ProcessCPU(resp.GetCPU())
	case resp.GetLoadAvg() != nil:
		app.ProcessLoadAvg(resp.GetLoadAvg())
	case resp.GetDisks() != nil:
		app.ProcessDisks(resp.GetDisks())
	case resp.GetDfsize() != nil:
		app.ProcessDfsize(resp.GetDfsize())
	case resp.GetDfinode() != nil:
		app.ProcessDfinode(resp.GetDfinode())
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
		case common.DUMMY:
			stats |= api.STATS_DUMMY
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
