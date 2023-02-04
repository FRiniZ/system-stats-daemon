package sensors

import (
	"context"
	"sync"
	"time"
)

/*
type SensorType int

const (

	SENSOR_LOAD_AVERAGE SensorType = iota
	SENSOR_CPU

)
*/
type Interface interface {
	Run(ctx context.Context, wg *sync.WaitGroup, out chan<- Interface)
	GetTimestamp() time.Time
	Sum(*Interface) *Interface
	Div(d int32) *Interface
}


type Controller interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
	GetAverageAfter (t time.Time) interface
}