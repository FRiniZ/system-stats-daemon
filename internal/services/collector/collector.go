package collector

/*
func Run(ctx context.Context) <-chan sensors.Interface {
	wg := &sync.WaitGroup{}
	out := make(chan sensors.Interface, 1)

	sIOstat := iostat.Sensor{}
	sLA := loadavg.Sensor{}

	wg.Add(2)
	go sLA.Run(ctx, wg, out)
	go sIOstat.Run(ctx, wg, out)

	go func() {
		defer close(out)
		defer fmt.Println("Close channel collector")
		wg.Wait()
	}()

	return out
}
*/
