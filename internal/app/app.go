package app

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/FRiniZ/system-stats-daemon/internal/services/collector"
	"github.com/FRiniZ/system-stats-daemon/internal/services/sensors/cpu"
	loadavg "github.com/FRiniZ/system-stats-daemon/internal/services/sensors/load-average"
)

type Config struct {
	Core struct {
		Frequency time.Duration `toml:"frequency"`
	} `toml:"core"`
}

type Application struct {
	Conf Config
}

func (app *Application) Run() {

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer cancel()

	fmt.Println("Run app")
	fmt.Println("With:")
	fmt.Printf("	frequency = %v\n", app.Conf.Core.Frequency)

	sCPU := []*cpu.Sensor{}
	sLA := []*loadavg.Sensor{}

	ch := collector.Run(ctx)

	for v := range ch {
		switch v.(type) {
		case *cpu.Sensor:
			sCPU = append(sCPU, v.(*cpu.Sensor))
		case *loadavg.Sensor:
			sLA = append(sLA, v.(*loadavg.Sensor))
		}
	}

	fmt.Println("We collected")
	for _, v := range sCPU {
		fmt.Printf("sensorCPU:%v %T\n", v, v)
	}
	for _, v := range sLA {
		fmt.Printf("sensorLA:%v %T\n", v, v)
	}
}
