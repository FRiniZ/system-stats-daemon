//go:build darwin

package iostatdisk

import (
	"context"
	"sync"
	"time"

	api "github.com/FRiniZ/system-stats-daemon/api/stub"
	"github.com/FRiniZ/system-stats-daemon/internal/services/sensors/common"
	"github.com/FRiniZ/system-stats-daemon/internal/storage"
)

const (
	Name string = "Disk load"
)

type Disk struct {
	Name     string
	TPS      float32
	KBsRead  float32
	KBsWrite float32
}

func (d *Disk) Add(a Disk) {
	d.KBsRead += a.KBsRead
	d.KBsWrite += a.KBsWrite
	d.TPS += a.TPS
}

func (d *Disk) Div(n int32) {
	d.KBsRead /= float32(n)
	d.KBsWrite /= float32(n)
	d.TPS /= float32(n)
}

type MDisks map[string]Disk

func (d *MDisks) Add(a *MDisks) {
	for k, v := range *a {
		v2, ok := (*d)[k]
		if ok {
			v2.Add(v)
		} else {
			(*d)[k] = v
		}
	}
}

func (d *MDisks) Div(n int32) {
	for _, v := range *d {
		v.Div(n)
	}
}

type Sensor struct {
	Disks MDisks
}

func (s *Sensor) Add(a *Sensor) {
	s.Disks.Add(&a.Disks)
}

func (s *Sensor) Div(n int32) {
	s.Disks.Div(n)
}

func (s *Sensor) MakeResponse() *api.Responce {
	res := &api.Responce{
		Disks: make([]*api.Loaddisk, 0, len(s.Disks)),
	}

	res.Disks = append(res.Disks, &api.Loaddisk{ErrorMsg: "Not implemented"})

	return res
}

const (
	FsmDeviceHeader = iota
	FsmDeviceBody
)

type Controller struct {
	queue *storage.Queue
}

func New() *Controller {
	return &Controller{queue: storage.New(0)}
}

func (c *Controller) GetAverageAfter(t time.Time) <-chan common.Sensor {
	out := make(chan common.Sensor)
	avg := Sensor{
		Disks: map[string]Disk{},
	}

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
	wg.Add(1)
	go func() {
		tiker := time.NewTicker(1 * time.Second)
		defer tiker.Stop()
		defer wg.Done()
		for {
			select {
			case <-tiker.C:
				c.queue.Push(&Sensor{}, time.Now())
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
