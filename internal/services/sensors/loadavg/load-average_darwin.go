//go:build darwin

package loadavg

import (
	"context"
	"sync"
	"time"

	api "github.com/FRiniZ/system-stats-daemon/api/stub"
	"github.com/FRiniZ/system-stats-daemon/internal/services/sensors/common"
	"github.com/FRiniZ/system-stats-daemon/internal/storage"
)

const (
	Name string = "Load average"
)

type Sensor struct {
	L1 float32
	L2 float32
	L3 float32
}

func (s *Sensor) Add(a *Sensor) {
	s.L1 += a.L1
	s.L2 += a.L2
	s.L3 += a.L3
}

func (s *Sensor) Div(d int32) {
	s.L1 /= float32(d)
	s.L2 /= float32(d)
	s.L3 /= float32(d)
}

func (s *Sensor) MakeResponse() *api.Responce {
	return &api.Responce{
		LoadAvg: &api.Loadaverage{
			ErrorMsg: "Not implemented",
		},
	}
}

type Controller struct {
	queue storage.Queue
}

func New() *Controller {
	return &Controller{queue: *storage.New(1)}
}

func (c *Controller) GetAverageAfter(t time.Time) <-chan common.Sensor {
	out := make(chan common.Sensor)
	avg := Sensor{}

	go func() {
		in := c.queue.GetElementsAfter(t)
		count := int32(0)
		for s := range in {
			count++
			avg.Add(s.(*Sensor))
		}
		if count > 0 {
			avg.Div(count)
			out <- &avg
		}
		close(out)
	}()

	return out
}

func (c *Controller) Run(ctx context.Context, wg *sync.WaitGroup) {
	var f1, f2, f3 float32

	wg.Add(1)
	go func() {
		tiker := time.NewTicker(1 * time.Second)
		defer tiker.Stop()
		defer wg.Done()
		for {
			select {
			case <-tiker.C:
				c.queue.Push(&Sensor{L1: f1, L2: f2, L3: f3}, time.Now())
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (c *Controller) GetName() string {
	return Name
}

func (c *Controller) SetMaxM(m int32) {
	c.queue.SetSize(c.GetName(), m)
}
