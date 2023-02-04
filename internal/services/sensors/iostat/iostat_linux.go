//go:build linux

package iostat

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

	"github.com/FRiniZ/system-stats-daemon/internal/storage"
)

type CPU struct {
	User   float32
	Nice   float32
	System float32
	IOWait float32
	Steal  float32
	Idle   float32
}

func (c *CPU) Add(a *CPU) {
	c.User += a.User
	c.Nice += a.Nice
	c.System += a.System
	c.IOWait += a.IOWait
	c.Steal += a.Steal
	c.Idle += a.Idle
}

func (c *CPU) Div(n int32) {
	c.User /= float32(n)
	c.Nice /= float32(n)
	c.System /= float32(n)
	c.IOWait /= float32(n)
	c.Steal /= float32(n)
	c.Idle /= float32(n)
}

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
	CPU
	Disks MDisks
}

func (s *Sensor) Add(a *Sensor) {
	s.CPU.Add(&a.CPU)
	s.Disks.Add(&a.Disks)
}

func (s *Sensor) Div(n int32) {
	s.CPU.Div(n)
	s.Disks.Div(n)
}

const (
	FSM_CPU_HEADER = iota
	FSM_CPU_BODY
	FSM_DEVICE_HEADER
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

func (c *Controller) GetAverageAfter(t time.Time) <-chan interface{} {
	out := make(chan interface{})
	avg := Sensor{
		CPU:   CPU{},
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
			out <- avg
		}
		close(out)
	}()

	return out
}

func (c *Controller) Run(ctx context.Context, wg *sync.WaitGroup) {
	var s *Sensor
	cmd := exec.CommandContext(ctx, "iostat", "-k", "1")
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
					s = &Sensor{
						CPU:   CPU{},
						Disks: map[string]Disk{},
					}
					state = FSM_CPU_BODY
				}
			case FSM_CPU_BODY:
				text = strings.ReplaceAll(text, ",", ".")
				fmt.Sscanf(text, "%f %f %f %f %f %f", &s.CPU.User, &s.CPU.Nice, &s.CPU.System, &s.CPU.IOWait, &s.CPU.Steal, &s.CPU.Idle)
				state = FSM_DEVICE_HEADER
			case FSM_DEVICE_HEADER:
				if strings.Contains(text, "Device") {
					s.Disks = make(MDisks, 0)
					state = FSM_DEVICE_BODY
				}
			case FSM_DEVICE_BODY:
				if text == "" {
					state = FSM_CPU_HEADER
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
			log.Println(err)
		}

		fmt.Println("Controller[IOStat] has closed")
		wg.Done()
	}()

	fmt.Println("Controller[IOStat] has started")
}
