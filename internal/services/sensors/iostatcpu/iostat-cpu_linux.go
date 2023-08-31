//go:build linux

package iostatcpu

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

const (
	Name string = "CPU"
)

type Sensor struct {
	User   float32
	Nice   float32
	System float32
	IOWait float32
	Steal  float32
	Idle   float32
}

func (s *Sensor) Add(a *Sensor) {
	s.User += a.User
	s.Nice += a.Nice
	s.System += a.System
	s.IOWait += a.IOWait
	s.Steal += a.Steal
	s.Idle += a.Idle
}

func (s *Sensor) Div(n int32) {
	s.User /= float32(n)
	s.Nice /= float32(n)
	s.System /= float32(n)
	s.IOWait /= float32(n)
	s.Steal /= float32(n)
	s.Idle /= float32(n)
}

func (s *Sensor) MakeResponse() *api.Responce {
	return &api.Responce{
		CPU: &api.Cpu{
			User:   s.User,
			System: s.System,
			Idle:   s.Idle,
		},
	}
}

const (
	FsmCPUHeader = iota
	FsmCPUBody
)

type Controller struct {
	queue *storage.Queue
}

func New() *Controller {
	return &Controller{queue: storage.New(0)}
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
	var s *Sensor
	cmd := exec.CommandContext(ctx, "iostat", "-c", "1")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(cmdReader)

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	state := FsmCPUHeader

	wg.Add(1)
	go func() {
		for scanner.Scan() {
			text := scanner.Text()
			text = strings.ReplaceAll(text, "  ", " ")
			switch state {
			case FsmCPUHeader:
				if strings.Contains(text, "avg-cpu:") {
					s = &Sensor{}
					state = FsmCPUBody
				}
			case FsmCPUBody:
				text = strings.ReplaceAll(text, ",", ".")
				fmt.Sscanf(text, "%f %f %f %f %f %f", &s.User, &s.Nice, &s.System, &s.IOWait, &s.Steal, &s.Idle)
				state = FsmCPUHeader
				c.queue.Push(s, time.Now())
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

func (c *Controller) GetName() string {
	return Name
}

func (c *Controller) SetMaxM(m int32) {
	c.queue.SetSize(c.GetName(), m)
}
