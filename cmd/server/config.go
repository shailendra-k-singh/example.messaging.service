package main

import (
	"flag"
	"time"
)

type config struct {
	logLevel     string
	logFormat    string
	reqPort      string
	tracingLib   string
	tracingAddr  string
	tracingHost  string
	charLimit    int
	readTimeout  time.Duration
	writeTimeout time.Duration
}

var conf config

func loadConfig() {
	flag.StringVar(&conf.logLevel, "log-level", "debug", "logging level for application")
	flag.StringVar(&conf.logFormat, "log-format", "json", "logging format for application")
	flag.StringVar(&conf.reqPort, "req-port", ":8090", "request port for incoming route queries")
	flag.StringVar(&conf.tracingLib, "tracing-lib", "jaeger", "tracing library")
	flag.StringVar(&conf.tracingHost, "tracing-host", "localhost", "tracing host (localhost/service-name)")
	flag.StringVar(&conf.tracingAddr, "tracing-addr", ":6831", "tracing reporter address")
	flag.IntVar(&conf.charLimit, "char-limit", 280, "character limit in the input message")
	flag.DurationVar(&conf.readTimeout, "read-timeout", 3*time.Second, "read timeout for HTTP server")
	flag.DurationVar(&conf.writeTimeout, "write-timeout", 5*time.Second, "write timeout for HTTP server")
	flag.Parse()
}
