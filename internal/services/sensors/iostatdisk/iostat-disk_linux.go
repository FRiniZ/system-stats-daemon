//go:build linux

package iostatdisk

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	api "github.com/FRiniZ/system-stats-daemon/api/stub"
	"github.com/FRiniZ/system-stats-daemon/internal/services/sensors/common"
	"github.com/FRiniZ/system-stats-daemon/internal/storage"
)

type Disk struct {
	Name     string
	TPS      float32
	KBsRead  float32
	KBsWrite float32
}

func (d *Disk) Add(a *Disk) {
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
			v2.Add(&v)
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

	for _, v := range s.Disks {
		res.Disks = append(res.Disks, &api.Loaddisk{
			Name:     v.Name,
			TPS:      v.TPS,
			WriteKBs: v.KBsWrite,
			ReadKBs:  v.KBsRead,
		})
	}

	return res
}

const (
	FSM_DEVICE_HEADER = iota
	FSM_DEVICE_BODY
)

type Controller struct {
	queue storage.Queue
}

func New(size int) *Controller {
	return &Controller{
		queue: *storage.New(size),
	}
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
	var s *Sensor
	cmd := exec.CommandContext(ctx, "iostat", "-d", "-k", "1")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(cmdReader)

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	state := FSM_DEVICE_HEADER

	wg.Add(1)
	go func() {
		for scanner.Scan() {
			text := scanner.Text()
			text = strings.ReplaceAll(text, "  ", " ")
			switch state {
			case FSM_DEVICE_HEADER:
				if strings.Contains(text, "Device") {
					s = &Sensor{
						Disks: make(MDisks, 0),
					}
					state = FSM_DEVICE_BODY
				}
			case FSM_DEVICE_BODY:
				if text == "" {
					state = FSM_DEVICE_HEADER
					c.queue.Push(s)
					s = nil
					continue
				}

				text = strings.ReplaceAll(text, ",", ".")
				disk := Disk{}
				fmt.Sscanf(text, "%s %f %f %f", &disk.Name, &disk.TPS, &disk.KBsRead, &disk.KBsWrite)
				s.Disks[disk.Name] = disk
			}
		}

		if err := cmd.Wait(); err != nil {
			if err.Error() != "signal: killed" {
				log.Println(err)
			}
		}
		wg.Done()
	}()
}

func (c *Controller) CheckSensor(i interface{}) bool {
	_, ok := i.(Sensor)
	return ok
}
