package main

import "github.com/FRiniZ/system-stats-daemon/internal/app/ssdaemon"

func main() {
	config := NewConfig().Config
	app := ssdaemon.Application{Conf: config}
	app.Run()
}
