package main

import (
	"github.com/FRiniZ/system-stats-daemon/internal/app"
)

func main() {
	config := NewConfig().Config
	app := app.Application{Conf: config}
	app.Run()
}
