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

	"github.com/FRiniZ/system-stats-daemon/internal/services/sensors"
)

type CPU struct {
	User   float32
	Nice   float32
	System float32
	IOWait float32
	Steal  float32
	Idle   float32
}

type Disk struct {
	Name     string
	TPS      float32
	KBsRead  float32
	KBsWrite float32
}

type Sensor struct {
	CPU
	Disks []Disk
}

const (
	FSM_CPU_HEADER = iota
	FSM_CPU_BODY
	FSM_DEVICE_HEADER
	FSM_DEVICE_BODY
)

func (s *Sensor) Run(ctx context.Context, out chan<- sensors.Interface) {
	cmd := exec.CommandContext(ctx, "iostat", "-k", "1")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(cmdReader)
	/*
	   avg-cpu:  %user   %nice %system %iowait  %steal   %idle
	              7,00    0,00    1,00    0,00    0,00   92,00

	   Device             tps    kB_read/s    kB_wrtn/s    kB_dscd/s    kB_read    kB_wrtn    kB_dscd
	   loop0             0,00         0,00         0,00         0,00          0          0          0
	   loop1             0,00         0,00         0,00         0,00          0          0          0
	   loop10            0,00         0,00         0,00         0,00          0          0          0
	   loop2             0,00         0,00         0,00         0,00          0          0          0
	   loop3             0,00         0,00         0,00         0,00          0          0          0
	   loop4             0,00         0,00         0,00         0,00          0          0          0
	   loop5             0,00         0,00         0,00         0,00          0          0          0
	   loop6             0,00         0,00         0,00         0,00          0          0          0
	   loop7             0,00         0,00         0,00         0,00          0          0          0
	   loop8             0,00         0,00         0,00         0,00          0          0          0
	   loop9             0,00         0,00         0,00         0,00          0          0          0
	   sda               0,00         0,00         0,00         0,00          0          0          0
	   sr0               0,00         0,00         0,00         0,00          0          0          0
	*/

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		state := FSM_CPU_HEADER
		for scanner.Scan() {
			text := scanner.Text()
			text = strings.ReplaceAll(text, "  ", " ")
			switch state {
			case FSM_CPU_HEADER:
				if strings.Contains(text, "avg-cpu:") {
					state = FSM_CPU_BODY
				}
			case FSM_CPU_BODY:
				text = strings.ReplaceAll(text, ",", ".")
				fmt.Sscanf(text, "%f %f %f %f %f %f", &s.CPU.User, &s.CPU.Nice, &s.CPU.System, &s.CPU.IOWait, &s.CPU.Steal, &s.CPU.Idle)
				state = FSM_DEVICE_HEADER
			case FSM_DEVICE_HEADER:
				if strings.Contains(text, "Device") {
					s.Disks = make([]Disk, 0)
					state = FSM_DEVICE_BODY
				}
			case FSM_DEVICE_BODY:
				if text == "" {
					fmt.Println("sensor filled:", s)
					state = FSM_CPU_HEADER
					continue
				}

				text = strings.ReplaceAll(text, ",", ".")
				disk := Disk{}
				fmt.Sscanf(text, "%s %f %f %f", &disk.Name, &disk.TPS, &disk.KBsRead, &disk.KBsWrite)
				s.Disks = append(s.Disks, disk)
			}
		}
		fmt.Println("Closed Run2")
	}()

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	wg.Wait()

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}
