package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nikhil-thomas/go-practice_ardanlabs-service/cmd/sidecar/metrics/collectors/expvar"
	"github.com/nikhil-thomas/go-practice_ardanlabs-service/cmd/sidecar/metrics/publishers/console"
	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/platform/cfg"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
}

func main() {
	c, err := cfg.New(cfg.EnvProvider{Namespace: "METRICS"})
	if err != nil {
		fmt.Printf("%s. All config defaults in use.", err)
	}

	apiHost, err := c.String("API_HOST")
	if err != nil {
		apiHost = "http://localhost:4000/debug/vars"
	}

	interval, err := c.Duration("INTERVAL")
	if err != nil {
		interval = 5 * time.Second
	}

	log.Printf("%s=%s", "API_HOST", apiHost)
	log.Printf("%s=%s", "INTERVAL", interval)

	// Start collectors and publishers

	expvar, err := expvar.New(apiHost)
	if err != nil {
		log.Fatalf("startup : Starting expavar collector : %v", err)
	}

	console, err := console.New(expvar, interval)
	if err != nil {
		log.Fatalf("startup : Starting console publisher : %v", err)
	}

	// make channel to listen for an interupt or terminate signal form the OS
	// Use a buffered channel because the signal package requires it
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
	<-osSignals

	log.Println("metric main : shutdown...")

	console.Stop()
	log.Println("main : Completed")
}
