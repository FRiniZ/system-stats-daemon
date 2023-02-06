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
		log.Println(err)
		log.Fatalf("Can't connect to GRPC-Server[%s:%s]\n", app.Conf.GRPC.Host, app.Conf.GRPC.Port)
	}

	defer conn.Close()

	grpcClient := api.NewSSDClient(conn)

	log.Println("Request sensors:")
	for _, name := range app.Conf.Core.Sensors {
		switch name {
		case "ALL":
			log.Println("               : ALL")
			stats = 0
			stats |= api.STATS_ALL
			goto exitfor
		case "CPU":
			stats |= api.STATS_CPU
			log.Println("               : CPU")
		case "LOADAVERAGE":
			stats |= api.STATS_LOADAVERAGE
			log.Println("               : LOADAVERAGE")
		case "LOADDISK":
			stats |= api.STATS_LOADDISK
			log.Println("               : LOADDISK")
		case "SIZEDISK":
			stats |= api.STATS_SIZEDISK
			log.Println("               : SIZEDISK")
		case "INODEDISK":
			stats |= api.STATS_INODEDISK
			log.Println("               : INODEDISK")
		}
	}

exitfor:

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
	}
}
