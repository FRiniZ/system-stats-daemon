package common

import (
	"context"
	"sync"
	"time"

	api "github.com/FRiniZ/system-stats-daemon/api/stub"
)

type Sensor interface {
	MakeResponse() *api.Responce
}

type Controller interface {
	Run(context.Context, *sync.WaitGroup)
	GetAverageAfter(time.Time) <-chan Sensor
	CheckSensor(interface{}) bool
}
