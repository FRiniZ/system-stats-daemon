package dummy

import (
	"context"
	"sync"
	"time"

	api "github.com/FRiniZ/system-stats-daemon/api/stub"
	"github.com/FRiniZ/system-stats-daemon/internal/services/sensors/common"
	"github.com/FRiniZ/system-stats-daemon/internal/storage"
)

const (
	Name string = "Dummy"
)

type Sensor struct{}

func (s *Sensor) Add(_ *Sensor) {}

func (s *Sensor) Div(_ int32) {}

func (s *Sensor) MakeResponse() *api.Responce {
	res := &api.Responce{
		Dummy: &api.Dummy{ErrorMsg: "Not implemented"},
	}
	return res
}

type Controller struct {
	queue *storage.Queue
}

func New() *Controller {
	return &Controller{queue: storage.New(1)}
}

// Dummy alwayes return only one record.
func (c *Controller) GetAverageAfter(_ time.Time) <-chan common.Sensor {
	out := make(chan common.Sensor)
	avg := Sensor{}

	go func() {
		out <- &avg
		close(out)
	}()

	return out
}

func (c *Controller) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		t := time.NewTicker(1 * time.Second)
		defer t.Stop()
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				c.queue.Push(&Sensor{}, time.Now())
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
