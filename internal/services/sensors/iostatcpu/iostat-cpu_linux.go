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

func (c *Sensor) Add(a *Sensor) {
	c.User += a.User
	c.Nice += a.Nice
	c.System += a.System
	c.IOWait += a.IOWait
	c.Steal += a.Steal
	c.Idle += a.Idle
}

func (c *Sensor) Div(n int32) {
	c.User /= float32(n)
	c.Nice /= float32(n)
	c.System /= float32(n)
	c.IOWait /= float32(n)
	c.Steal /= float32(n)
	c.Idle /= float32(n)
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
	FSM_CPU_HEADER = iota
	FSM_CPU_BODY
)

type Controller struct {
	queue storage.Queue
}

func New(size int) *Controller {
	return &Controller{
		queue: *storage.New(size)}
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

	state := FSM_CPU_HEADER

	wg.Add(1)
	go func() {
		for scanner.Scan() {
			text := scanner.Text()
			text = strings.ReplaceAll(text, "  ", " ")
			switch state {
			case FSM_CPU_HEADER:
				if strings.Contains(text, "avg-cpu:") {
					s = &Sensor{}
					state = FSM_CPU_BODY
				}
			case FSM_CPU_BODY:
				text = strings.ReplaceAll(text, ",", ".")
				fmt.Sscanf(text, "%f %f %f %f %f %f", &s.User, &s.Nice, &s.System, &s.IOWait, &s.Steal, &s.Idle)
				state = FSM_CPU_HEADER
				c.queue.Push(s)
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

func (c *Controller) SetMaxM(M int32) {
	if c.queue.SetSize(M) {
		log.Printf("[%s] Changed size of queue to:%d", Name, M)
	}
}
