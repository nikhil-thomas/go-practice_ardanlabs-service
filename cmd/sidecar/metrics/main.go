package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nikhil-thomas/go-practice_ardanlabs-service/cmd/sidecar/metrics/collector"
	"github.com/nikhil-thomas/go-practice_ardanlabs-service/cmd/sidecar/metrics/publisher"
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

	publishTo, err := c.String("PUBLISHER")
	if err != nil {
		publishTo = "CONSOLE"
	}

	log.Printf("%s=%s", "API_HOST", apiHost)
	log.Printf("%s=%s", "INTERVAL", interval)
	log.Printf("%s=%s", "PUBLISHER", publishTo)

	// Start collectors and publishers

	expvar, err := collector.New(apiHost)
	if err != nil {
		log.Fatalf("startup : Starting collector : %v", err)
	}

	f := publisher.Console

	switch publishTo {
	case publisher.TypeDatadog:
		f = publisher.Datadog
	}

	publish, err := publisher.New(expvar, f, interval)
	if err != nil {
		log.Fatalf("startup : Starting publisher : %v", err)
	}

	// make channel to listen for an interupt or terminate signal form the OS
	// Use a buffered channel because the signal package requires it
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
	<-osSignals

	defer log.Println("main : Completed")

	publish.Stop()

}
