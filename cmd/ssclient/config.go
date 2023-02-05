package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/FRiniZ/system-stats-daemon/internal/app"
	"github.com/FRiniZ/system-stats-daemon/internal/app/ssclient"
)

type Sensors []string

func (s *Sensors) String() string {
	return "not implemented"
}

func (s *Sensors) Set(v string) error {
	*s = append(*s, v)
	return nil
}

var (
	grpcHost   string
	grpcPort   string
	N          time.Duration
	M          time.Duration
	sensors    Sensors
	configFile string
)

func init() {
	flag.DurationVar(&N, "N", 0, "frequency of statistics output")
	flag.DurationVar(&M, "M", 0, "Duration of statistics")
	flag.StringVar(&grpcHost, "host", "", "GRPC Host")
	flag.StringVar(&grpcPort, "port", "", "GRPC Port")
	flag.Var(&sensors, "sensor", "[ALL|CPU|LoadAverage|LoadDisk]")
	flag.StringVar(&configFile, "config", "", "Path to configuration file")
	flag.Parse()

	if flag.Arg(0) == "version" {
		app.PrintVersion()
		return
	}
}

type Config struct {
	ssclient.Config
}

func NewConfig() Config {
	var config Config
	if configFile != "" {
		if err := config.LoadFileTOML(configFile); err != nil {
			fmt.Fprintf(os.Stderr, "Can't load config file:%v error: %v\n", configFile, err)
			os.Exit(1)
		}
	}

	if grpcHost != "" {
		config.GRPC.Host = grpcHost
	}

	if grpcPort != "" {
		config.GRPC.Port = grpcPort
	}

	if N != 0 {
		config.Core.N = N
	}

	if M != 0 {
		config.Core.M = M
	}

	if len(sensors) > 0 {
		config.Core.Sensors = sensors
	}

	log.Println("Config:", config)
	return config
}

func (c *Config) LoadFileTOML(filename string) error {
	filedata, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return toml.Unmarshal(filedata, c)
}
