package main

import (
	"net/http"

	"github.com/shailendra-k-singh/example.messaging.service/app"
	logger "github.com/shailendra-k-singh/example.messaging.service/pkg/log"
)

func main() {
	loadConfig()
	log := logger.NewLogger(conf.logLevel, conf.logFormat)
	log.Info("Starting Messaging Service")

	// Get new appRouter instance
	log.Info("Creating new appRouter instance")
	r := app.NewAppRouter(conf.charLimit)
	log.Info("Initializing tracing and routes")
	err := r.InitTracing(conf.tracingLib, conf.tracingHost,conf.tracingAddr)
	if err != nil {
		log.Fatal("Error while initializing tracing: ", err)
	}
	r.SetRoutes()

	srv := http.Server{
		Addr:         conf.reqPort,
		Handler:      r.GetRouter(),
		ReadTimeout:  conf.readTimeout,
		WriteTimeout: conf.writeTimeout,
	}

	log.Info("Starting HTTP server")
	err = srv.ListenAndServe()
	r.Close()
	if err != nil {
		log.Fatal("Server closed with error: ", err)
	}
	log.Info("Server exiting...")
}
