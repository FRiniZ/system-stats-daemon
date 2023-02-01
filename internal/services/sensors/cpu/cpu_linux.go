package cpu

import (
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

// top -l 1 | grep -E "^CPU"
func (s *Sensor) Read() sensors.Interface {
	var err error
	var f1, f2, f3, f4, f5, f6, f7, f8 float32

	c1 := exec.Command("top", "-b", "-n1")
	c2 := exec.Command("grep", "-E", "Cpu")

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

	fmt.Sscanf(outb.String(), "%%Cpu(s):  %f us,  %f sy,  %f ni,%f id,  %f wa,  %f hi,  %f si,  %f st", &f1, &f2, &f3, &f4, &f5, &f6, &f7, &f8)

	return &Sensor{User: f1, System: f2, Idle: f4}
}
