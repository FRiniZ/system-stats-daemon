package main

import (
	"github.com/FRiniZ/system-stats-daemon/internal/app"
)

/*
var (
	Release   string
	BuildDate string
	GitHash   string
)
*/

func main() {
	app := app.Application{}
	app.Run()
}
