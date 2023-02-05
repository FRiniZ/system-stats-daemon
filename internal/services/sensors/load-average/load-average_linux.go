//go:build linux

package loadavg

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os/exec"
	"sync"
	"syscall"
	"time"

	api "github.com/FRiniZ/system-stats-daemon/api/stub"
	"github.com/FRiniZ/system-stats-daemon/internal/services/sensors/common"
	"github.com/FRiniZ/system-stats-daemon/internal/storage"
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
			L1: s.L1,
			L2: s.L2,
			L3: s.L3,
		},
	}
}

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

	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", "while true; do cat /proc/loadavg; sleep 1; done")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(cmdReader)

	if err := cmd.Start(); err != nil {
		log.Println(err)
		return
	}

	wg.Add(1)
	go func() {
		for scanner.Scan() {
			text := scanner.Text()
			fmt.Sscanf(text, "%f %f %f", &f1, &f2, &f3)
			c.queue.Push(&Sensor{L1: f1, L2: f2, L3: f3})
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
