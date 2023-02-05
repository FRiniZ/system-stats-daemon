package grpcserver

import (
	"context"
	"log"
	"sync"
	"time"

	api "github.com/FRiniZ/system-stats-daemon/api/stub"
	"github.com/FRiniZ/system-stats-daemon/internal/services/sensors/common"
	"github.com/FRiniZ/system-stats-daemon/internal/services/sensors/iostatcpu"
	"github.com/FRiniZ/system-stats-daemon/internal/services/sensors/iostatdisk"
	loadavg "github.com/FRiniZ/system-stats-daemon/internal/services/sensors/load-average"
	"google.golang.org/grpc/peer"
)

type grpcserver struct {
	api.UnimplementedSSDServer
	wg *sync.WaitGroup
}

func New(wg *sync.WaitGroup) api.SSDServer {
	return grpcserver{wg: wg}
}

func (s grpcserver) Subsribe(req *api.Request, stream api.SSD_SubsribeServer) error {
	s.wg.Add(1)
	defer s.wg.Done()

	IPaddr := "unknown"
	p, _ := peer.FromContext(stream.Context())
	IPaddr = p.Addr.String()

	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	var controllers []common.Controller
	wg := &sync.WaitGroup{}
	N := time.Duration(req.GetN())
	M := time.Duration(req.GetM())
	sensors := req.GetBitmask()

	if sensors&api.STATS_ALL == api.STATS_ALL {
		log.Printf("[%s]Request ALL stats\n", IPaddr)
		controllers = append(controllers, loadavg.New(int(M)))
		controllers = append(controllers, iostatcpu.New(int(M)))
		controllers = append(controllers, iostatdisk.New(int(M)))
	}

	if sensors&api.STATS_CPU == api.STATS_CPU {
		log.Printf("[%s]Request CPU stats\n", IPaddr)
		controllers = append(controllers, iostatcpu.New(int(M)))
	}

	if sensors&api.STATS_LOADAVERAGE == api.STATS_LOADAVERAGE {
		log.Printf("[%s]Request LoadAverage stats\n", IPaddr)
		controllers = append(controllers, loadavg.New(int(M)))
	}
	if sensors&api.STATS_LOADDISK == api.STATS_LOADDISK {
		log.Printf("[%s]Request LoadDisk stats\n", IPaddr)
		controllers = append(controllers, iostatdisk.New(int(M)))
	}

	for _, ctrl := range controllers {
		ctrl.Run(ctx, wg)
	}

	sendRequst := func(in <-chan common.Sensor) {
		for v := range in {
			stream.Send(v.MakeResponse())
		}
	}

	ticker := time.NewTicker(N * time.Second)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				for _, ctrl := range controllers {
					go sendRequst(ctrl.GetAverageAfter(time.Now().Add(M * time.Second * -1)))
				}
			}
		}
	}()

	log.Printf("[%s]GRPC-Stream has started\n", IPaddr)
	wg.Wait()
	log.Printf("[%s]GRPC-Stream has closed\n", IPaddr)
	return nil
}

func (s grpcserver) mustEmbedUnimplementedSSDServer() {
	panic("not implemented") // TODO: Implement
}
