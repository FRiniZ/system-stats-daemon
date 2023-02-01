package sensors

/*
type SensorType int

const (

	SENSOR_LOAD_AVERAGE SensorType = iota
	SENSOR_CPU

)
*/
type Interface interface {
	Read() Interface
}
