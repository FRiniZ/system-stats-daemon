package collector

import (
	"context"
	"fmt"
	"sync"

	"github.com/FRiniZ/system-stats-daemon/internal/services/sensors"
	"github.com/FRiniZ/system-stats-daemon/internal/services/sensors/cpu"
	loadavg "github.com/FRiniZ/system-stats-daemon/internal/services/sensors/load-average"
)

func Run(ctx context.Context) <-chan sensors.Interface {
	wg := &sync.WaitGroup{}
	out := make(chan sensors.Interface, 1)

	sCPU := cpu.Sensor{}
	sLA := loadavg.Sensor{}

	wg.Add(1)
	sCPU.Run2(ctx, out, wg)

	wg.Add(1)
	sLA.Run(ctx, out, wg)

	go func() {
		defer close(out)
		defer fmt.Println("Close channel collector")
		wg.Wait()
	}()

	return out
}
