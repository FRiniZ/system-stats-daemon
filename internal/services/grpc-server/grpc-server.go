package grpcserver

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	api "github.com/FRiniZ/system-stats-daemon/api/stub"
	"github.com/FRiniZ/system-stats-daemon/internal/services/sensors/common"
	"github.com/FRiniZ/system-stats-daemon/internal/services/sensors/dfinode"
	"github.com/FRiniZ/system-stats-daemon/internal/services/sensors/dfsize"
	"github.com/FRiniZ/system-stats-daemon/internal/services/sensors/dummy"
	"github.com/FRiniZ/system-stats-daemon/internal/services/sensors/iostatcpu"
	"github.com/FRiniZ/system-stats-daemon/internal/services/sensors/iostatdisk"
	"github.com/FRiniZ/system-stats-daemon/internal/services/sensors/loadavg"
	"google.golang.org/grpc/peer"
)

type GRPCServer struct {
	api.UnimplementedSSDServer
	wg          *sync.WaitGroup
	lock        *sync.Mutex
	ctx         context.Context
	stop        context.CancelFunc
	controllers map[string]common.Controller
}

func New(wg *sync.WaitGroup) *GRPCServer {
	return &GRPCServer{wg: wg, lock: &sync.Mutex{}, controllers: make(map[string]common.Controller)}
}

func (s *GRPCServer) Stop(ctx context.Context) error {
	s.stop()
	return nil
}

func (s *GRPCServer) Start(ctx context.Context) error {
	s.ctx, s.stop = context.WithCancel(ctx)
	s.controllers[common.LOADAVERAGE] = loadavg.New()
	s.controllers[common.CPU] = iostatcpu.New()
	s.controllers[common.LOADDISK] = iostatdisk.New()
	s.controllers[common.SIZEDISK] = dfsize.New()
	s.controllers[common.INODEDISK] = dfinode.New()
	s.controllers[common.DUMMY] = dummy.New()

	for _, v := range s.controllers {
		log.Printf("Starting[%s]....\n", v.GetName())
		v.Run(ctx, s.wg)
	}
	return nil
}

func (s *GRPCServer) GetController(name string, m int64) common.Controller {
	if v, ok := s.controllers[name]; ok {
		v.SetMaxM(int32(m))
		return v
	}
	return nil
}

func (s *GRPCServer) Subsribe(req *api.Request, stream api.SSD_SubsribeServer) error {
	var once sync.Once
	var IPaddr string

	s.wg.Add(1)
	defer s.wg.Done()

	p, _ := peer.FromContext(stream.Context())
	IPaddr = p.Addr.String()

	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	var controllers []common.Controller
	wg := &sync.WaitGroup{}
	N := req.GetN()
	M := req.GetM()
	sensors := req.GetBitmask()

	if sensors == 0 {
		return fmt.Errorf("request without sensors")
	}

	if sensors&api.STATS_ALL == api.STATS_ALL {
		log.Printf("[%s]Request ALL stats\n", IPaddr)
		controllers = append(controllers, s.GetController(common.LOADAVERAGE, M.GetSeconds()))
		controllers = append(controllers, s.GetController(common.CPU, M.GetSeconds()))
		controllers = append(controllers, s.GetController(common.LOADDISK, M.GetSeconds()))
		controllers = append(controllers, s.GetController(common.SIZEDISK, M.GetSeconds()))
		controllers = append(controllers, s.GetController(common.INODEDISK, M.GetSeconds()))
	}

	if sensors&api.STATS_CPU == api.STATS_CPU {
		log.Printf("[%s]Request CPU stats\n", IPaddr)
		controllers = append(controllers, s.GetController(common.CPU, M.GetSeconds()))
	}

	if sensors&api.STATS_LOADAVERAGE == api.STATS_LOADAVERAGE {
		log.Printf("[%s]Request LoadAverage stats\n", IPaddr)
		controllers = append(controllers, s.GetController(common.LOADAVERAGE, M.GetSeconds()))
	}

	if sensors&api.STATS_LOADDISK == api.STATS_LOADDISK {
		log.Printf("[%s]Request LoadDisk stats\n", IPaddr)
		controllers = append(controllers, s.GetController(common.LOADDISK, M.GetSeconds()))
	}
	if sensors&api.STATS_SIZEDISK == api.STATS_SIZEDISK {
		log.Printf("[%s]Request SizeDisk stats\n", IPaddr)
		controllers = append(controllers, s.GetController(common.SIZEDISK, M.GetSeconds()))
	}
	if sensors&api.STATS_SIZEDISK == api.STATS_INODEDISK {
		log.Printf("[%s]Request InodeDisk stats\n", IPaddr)
		controllers = append(controllers, s.GetController(common.INODEDISK, M.GetSeconds()))
	}

	if sensors&api.STATS_DUMMY == api.STATS_DUMMY {
		log.Printf("[%s]Request Dummy stats\n", IPaddr)
		controllers = append(controllers, s.GetController(common.DUMMY, M.GetSeconds()))
	}

	sendRequst := func(in <-chan common.Sensor, wg *sync.WaitGroup) {
		defer wg.Done()
		for v := range in {
			s.lock.Lock()
			stream.Send(v.MakeResponse())
			s.lock.Unlock()
		}
	}

	tickerN := time.NewTicker(M.AsDuration())

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer tickerN.Stop()

		for {
			select {
			case <-s.ctx.Done():
				return
			case <-ctx.Done():
				return
			case <-tickerN.C:
				once.Do(func() { tickerN = time.NewTicker(N.AsDuration()) })
				for _, ctrl := range controllers {
					wg.Add(1)
					go sendRequst(ctrl.GetAverageAfter(time.Now().Add(M.AsDuration()*-1)), wg)
				}
			}
		}
	}()

	log.Printf("[%s]GRPC-Stream has started\n", IPaddr)
	wg.Wait()
	log.Printf("[%s]GRPC-Stream has closed\n", IPaddr)
	return nil
}
