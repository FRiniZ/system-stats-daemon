package grpcserver

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	api "github.com/FRiniZ/system-stats-daemon/api/stub"
	"github.com/FRiniZ/system-stats-daemon/internal/services/sensors/iostat"
	loadavg "github.com/FRiniZ/system-stats-daemon/internal/services/sensors/load-average"
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

	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	wg := &sync.WaitGroup{}
	N := time.Duration(req.GetN())
	M := time.Duration(req.GetM())
	sensors := req.GetBitmask()

	log.Println("Request stats:", sensors)

	if sensors == api.STATS_ALL {
		log.Println("Request to ALL Sensors")
	}
	if sensors&api.STATS_CPU == api.STATS_CPU {
		log.Println("Request to CPU sensor")
	}

	if sensors&api.STATS_LOADAVERAGE == api.STATS_LOADAVERAGE {
		log.Println("Request to LoadAverage sensor")
	}
	if sensors&api.STATS_LOADDISK == api.STATS_LOADDISK {
		log.Println("Request to Disks sensor")
	}

	cLA := loadavg.New(int(M))
	cIOS := iostat.New(int(M))

	cLA.Run(ctx, wg)
	cIOS.Run(ctx, wg)

	sendRequst := func(in <-chan interface{}) {
		for v := range in {
			switch v.(type) {
			case interface{}:
				s, ok := v.(loadavg.Sensor)
				if ok {
					log.Println("LA Sensor", s)
				} else {
					log.Println("Interface")
				}
			case loadavg.Sensor:
				fmt.Println("LA Sensor")
			}
			//stream.Send()
		}
	}

	ticker := time.NewTicker(N * time.Second)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				go sendRequst(cLA.GetAverageAfter(time.Now().Add(M * time.Second * -1)))
				go sendRequst(cIOS.GetAverageAfter(time.Now().Add(M * time.Second * -1)))
			}
		}
	}()

	log.Println("GRPC-Stream has stared")
	wg.Wait()
	log.Println("GRPC-Stream has closed")
	return nil
}

func (s grpcserver) mustEmbedUnimplementedSSDServer() {
	panic("not implemented") // TODO: Implement
}
