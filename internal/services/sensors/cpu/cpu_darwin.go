package cpu

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/FRiniZ/system-stats-daemon/internal/services/sensors"
)

type Sensor struct {
	User   float32
	System float32
	Idle   float32
}

func (s *Sensor) Run(ctx context.Context, out chan<- sensors.Interface, wg *sync.WaitGroup) {
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		defer wg.Done()
		defer fmt.Println("Stop cpu.sensor")
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				out <- s.Read()
			}
		}
	}()
}

func (s *Sensor) Run2(ctx context.Context, out chan<- sensors.Interface, wg *sync.WaitGroup) error {

	cmd := exec.Command("iostat", "-w", "1")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmdOut, err := cmd.StdoutPipe()
	//cmdErr, err := cmd.StderrPipe()

	err = cmd.Start()
	if err != nil {
		log.Printf("Stderr:%v", err)
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(cmdOut)

	go func() {

		for scanner.Scan() {
			fmt.Println("cpu.sensor:", scanner.Text())
			out <- &Sensor{s.User, s.System, s.Idle}
		}
		cmd.Wait()
	}()

	return nil
}

// top -l 1 | grep -E "^CPU"
func (s *Sensor) Read() sensors.Interface {
	var err error

	c1 := exec.Command("top", "-F", "-l", "1")
	c2 := exec.Command("grep", "-E", "^CPU")

	c1.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	c2.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	var outb, errb bytes.Buffer
	r, w := io.Pipe()

	c1.Stdout = w
	c1.Stderr = &errb

	c2.Stdin = r
	c2.Stdout = &outb
	c2.Stderr = &errb

	err = c1.Start()
	if err != nil {
		log.Printf("Stderr:%v", errb.String())
		log.Fatal(err)
	}
	err = c2.Start()
	if err != nil {
		log.Printf("Stderr:%v", errb.String())
		log.Fatal(err)
	}
	c1.Wait()
	w.Close()
	c2.Wait()

	fmt.Sscanf(outb.String(), "CPU usage: %f%% user, %f%% sys, %f%% idle", &s.User, &s.System, &s.Idle)
	return &Sensor{s.User, s.System, s.Idle}
}
