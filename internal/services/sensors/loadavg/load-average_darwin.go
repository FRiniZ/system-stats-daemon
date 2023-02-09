//go:build darwin

package loadavg

/*
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
	L1 float32
	L2 float32
	L3 float32
}

func (s *Sensor) Run(ctx context.Context, out chan<- sensors.Interface, wg *sync.WaitGroup) {
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		defer wg.Done()
		defer fmt.Println("Stop load-average.sensor")
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

// sysctl -n vm.loadavg
func (s *Sensor) Read() sensors.Interface {
	var err error
	var f1, f2, f3 float32

	c1 := exec.Command("top", "-F", "-l", "1")
	c2 := exec.Command("grep", "-E", "^Load")

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

	fmt.Sscanf(outb.String(), "Load Avg: %f, %f, %f", &f1, &f2, &f3)

	return &Sensor{L1: f1, L2: f2, L3: f3}
}
*/
