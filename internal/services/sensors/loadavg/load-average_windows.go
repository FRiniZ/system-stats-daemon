//go:build windows

package loadavg


type Sensor struct {
	list list.List
}

func Read() string {

	return "not implemented"
}
