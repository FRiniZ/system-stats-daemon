package main

import "github.com/FRiniZ/system-stats-daemon/internal/app/ssclient"

func main() {
	config := NewConfig().Config
	app := ssclient.Application{}
	app.Run(config)
}
