package grpcserver

import (
	"context"
	"log"
	"sync"
	"time"

	api "github.com/FRiniZ/system-stats-daemon/api/stub"
	"github.com/FRiniZ/system-stats-daemon/internal/services/sensors/common"
	"github.com/FRiniZ/system-stats-daemon/internal/services/sensors/dfsize"
	"github.com/FRiniZ/system-stats-daemon/internal/services/sensors/iostatcpu"
	"github.com/FRiniZ/system-stats-daemon/internal/services/sensors/iostatdisk"
	loadavg "github.com/FRiniZ/system-stats-daemon/internal/services/sensors/load-average"
	"google.golang.org/grpc/peer"
)

type grpcserver struct {
	api.UnimplementedSSDServer
	wg          *sync.WaitGroup
	lock        *sync.Mutex
	controllers map[string]common.Controller
}

func New(wg *sync.WaitGroup) grpcserver {
	return grpcserver{wg: wg, lock: &sync.Mutex{}, controllers: make(map[string]common.Controller)}
}

func (s grpcserver) Start(ctx context.Context, M int) error {
	s.controllers["LA"] = loadavg.New(int(M))
	s.controllers["CPU"] = iostatcpu.New(int(M))
	s.controllers["DISKLOAD"] = iostatdisk.New(int(M))
	s.controllers["DISKSIZE"] = dfsize.New(int(M))

	for _, v := range s.controllers {
		log.Printf("Starting[%s]....\n", v.GetName())
		v.Run(ctx, s.wg)
	}
	return nil
}

func (s grpcserver) Stop(ctx context.Context) error {

	return nil
}

func (s grpcserver) GetController(name string, M int32) common.Controller {
	if v, ok := s.controllers[name]; ok {
		v.SetMaxM(M)
		return v
	}
	return nil
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
	M := req.GetM()
	sensors := req.GetBitmask()

	if sensors&api.STATS_ALL == api.STATS_ALL {
		log.Printf("[%s]Request ALL stats\n", IPaddr)
		controllers = append(controllers, s.GetController("LA", M))
		controllers = append(controllers, s.GetController("CPU", M))
		controllers = append(controllers, s.GetController("DISKLOAD", M))
		controllers = append(controllers, s.GetController("DISKSIZE", M))
	} else {
		if sensors&api.STATS_CPU == api.STATS_CPU {
			log.Printf("[%s]Request CPU stats\n", IPaddr)
			controllers = append(controllers, s.GetController("CPU", M))
		}

		if sensors&api.STATS_LOADAVERAGE == api.STATS_LOADAVERAGE {
			log.Printf("[%s]Request LoadAverage stats\n", IPaddr)
			controllers = append(controllers, s.GetController("LA", M))
		}
		if sensors&api.STATS_LOADDISK == api.STATS_LOADDISK {
			log.Printf("[%s]Request LoadDisk stats\n", IPaddr)
			controllers = append(controllers, s.GetController("DISKLOAD", M))
		}
		if sensors&api.STATS_SIZEDISK == api.STATS_SIZEDISK {
			log.Printf("[%s]Request SizeDisk stats\n", IPaddr)
			controllers = append(controllers, s.GetController("DISKSIZE", M))
		}
	}

	sendRequst := func(in <-chan common.Sensor) {
		for v := range in {
			s.lock.Lock()
			stream.Send(v.MakeResponse())
			s.lock.Unlock()
		}
	}

	ticker := time.NewTicker(N * time.Second)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				for _, ctrl := range controllers {
					go sendRequst(ctrl.GetAverageAfter(time.Now().Add(time.Duration(M) * time.Second * -1)))
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
