//go:build darwin

package dfinode

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
	Name string = "Disk inodes"
)

type Disk struct {
	Name  string
	IUsed int32
	IUse  int32
}

func (d *Disk) Add(a Disk) {
	d.IUsed += a.IUsed
	d.IUse += a.IUse
}

func (d *Disk) Div(n int32) {
	d.IUsed /= n
	d.IUse /= n
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
		Dfinode: make([]*api.Dfinode, 0, len(s.Disks)),
	}
	for _, v := range s.Disks {
		res.Dfinode = append(res.Dfinode, &api.Dfinode{
			Name:  v.Name,
			IUsed: v.IUsed,
			IUse:  v.IUse,
		})
	}
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
	var s *Sensor
	var inodes int32
	var ifree int32

	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", "while true; do df -BM -i; echo; echo; sleep 1; done")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(cmdReader)

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	state := FsmDeviceHeader

	wg.Add(1)
	go func() {
		for scanner.Scan() {
			text := scanner.Text()
			text = strings.ReplaceAll(text, "  ", " ")
			switch state {
			case FsmDeviceHeader:
				if strings.Contains(text, "Filesystem") {
					s = &Sensor{
						Disks: make(MDisks, 0),
					}
					state = FsmDeviceBody
				}
			case FsmDeviceBody:
				if text == "" {
					state = FsmDeviceHeader
					c.queue.Push(s, time.Now())
					s = nil
					continue
				}
				disk := Disk{}
				fmt.Sscanf(text, "%s %dM %dM %dM %d%%", &disk.Name, &inodes, &disk.IUsed, &ifree, &disk.IUse)
				// Skip not phisical disk
				if disk.Name[0] == '/' {
					s.Disks[disk.Name] = disk
				}
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
