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
		fmt.Printf("config : %s. All config defaults in use.", err)
	}

	apiHost, err := c.String("API_HOST")
	if err != nil {
		apiHost = "http://crud:4000/debug/vars"
	}

	interval, err := c.Duration("INTERVAL")
	if err != nil {
		interval = 5 * time.Second
	}

	publishTo, err := c.String("PUBLISHER")
	if err != nil {
		publishTo = "CONSOLE"
	}

	dataDogAPIKey, err := c.String("DATADOG_APIKEY")
	if err != nil {
		dataDogAPIKey = ""
	}

	dataDogHost, err := c.String("DATADOG_HOST")
	if err != nil {
		dataDogAPIKey = "https://app.datadoghq.com/api/v1/series"
	}

	log.Printf("config : %s=%s", "API_HOST", apiHost)
	log.Printf("config : %s=%s", "INTERVAL", interval)
	log.Printf("config : %s=%s", "PUBLISHER", publishTo)
	log.Printf("config : %s=%s", "DATADOG_APIKEY", dataDogAPIKey)
	log.Printf("config : %s=%s", "DATADOG_HOST", dataDogHost)

	// Start collectors and publishers

	expvar, err := collector.New(apiHost)
	if err != nil {
		log.Fatalf("main : Starting collector : %v", err)
	}

	f := publisher.Console

	switch publishTo {
	case publisher.TypeConsole:
		log.Println("config : PUB_TYPE=Console")
	case publisher.TypeDatadog:
		log.Println("config : PUB_TYPE=Datadog")
		d := publisher.NewDatadog(dataDogAPIKey, dataDogHost)
		f = d.Publish
	default:
		log.Fatalln("main : No publisher provided, using Console.")
	}

	publish, err := publisher.New(expvar, f, interval)
	if err != nil {
		log.Fatalf("main : Starting publisher : %v", err)
	}

	// make channel to listen for an interupt or terminate signal form the OS
	// Use a buffered channel because the signal package requires it
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
	<-osSignals

	defer log.Println("main : Completed")

	publish.Stop()

}
