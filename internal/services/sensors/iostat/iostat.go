package iostat

type Sensor struct {
	LA1m   float32
	LA5m   float32
	LA15m  float32
	User   float32
	System float32
	Idle   float32
}
