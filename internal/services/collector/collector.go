package collector

import (
	"context"

	"github.com/FRiniZ/system-stats-daemon/internal/services/sensors"
	loadavg "github.com/FRiniZ/system-stats-daemon/internal/services/sensors/load-average"
)

func Run(ctx context.Context) <-chan sensors.Interface {
	//wg := &sync.WaitGroup{}
	out := make(chan sensors.Interface, 1)

	sLA := loadavg.Sensor{}

	sLA.Run(ctx, out)
	//IOstat := iostat.Sensor{}

	//IOstat.Run2(ctx, out)

	/*
		sCPU := cpu.Sensor{}
		sLA := loadavg.Sensor{}

		wg.Add(1)
		sCPU.Run(ctx, out, wg)

		wg.Add(1)
		sLA.Run(ctx, out, wg)

		go func() {
			defer close(out)
			defer fmt.Println("Close channel collector")
			wg.Wait()
		}()
	*/

	return out
}
