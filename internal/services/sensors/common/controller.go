package common

import (
	"context"
	"sync"
	"time"

	api "github.com/FRiniZ/system-stats-daemon/api/stub"
)

const (
	ALL         string = "ALL"
	CPU         string = "CPU"
	LOADAVERAGE string = "LOADAVERAGE"
	LOADDISK    string = "LOADDISK"
	SIZEDISK    string = "SIZEDISK"
	INODEDISK   string = "INODEDISK"
)

type Sensor interface {
	MakeResponse() *api.Responce
}

type Controller interface {
	Run(context.Context, *sync.WaitGroup)
	GetAverageAfter(time.Time) <-chan Sensor
	GetName() string
	SetMaxM(int32)
}
