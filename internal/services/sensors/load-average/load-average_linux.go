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

	"github.com/FRiniZ/system-stats-daemon/internal/services/sensors"
)

type Sensor struct {
	L1 float32
	L2 float32
	L3 float32
}

func (s *Sensor) Run(ctx context.Context, out chan<- sensors.Interface) {
	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", "while true; do cat /proc/loadavg; sleep 1; done")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(cmdReader)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		//state := FSM_CPU_HEADER
		for scanner.Scan() {
			text := scanner.Text()
			fmt.Println("LA:", text)
		}
		fmt.Println("Closed Run3")
	}()

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	wg.Wait()

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}

func (s *Sensor) Read() sensors.Interface {
	return &Sensor{s.L1, s.L2, s.L3}
}
